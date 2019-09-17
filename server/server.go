package server

import (
	"CDcoding2333/scaffold/apps/rpc"
	pb "CDcoding2333/scaffold/apps/rpc/proto"
	"CDcoding2333/scaffold/conf"
	"CDcoding2333/scaffold/core"
	"CDcoding2333/scaffold/utils/middleware"
	"CDcoding2333/scaffold/utils/sign"
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	c          *core.Core
	httpServer *http.Server
	wsServer   *http.Server
	grpcServer *grpc.Server
	config     *conf.ServerConfig
	jwtAuth    *sign.Auth
}

// Run ...
func Run(c *core.Core, config *conf.ServerConfig) (*Server, error) {
	server := &Server{
		c:       c,
		config:  config,
		jwtAuth: sign.NewAuth(config.JwtISS, config.JwtScrect),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", server.config.Port))
	if err != nil {
		log.WithError(err).Fatalf("new listener error")
		return nil, err
	}

	m := cmux.New(listener)
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	wsl := m.Match(cmux.HTTP1HeaderField("Upgrade", "websocket"))
	httpL := m.Match(cmux.HTTP1Fast())

	if server.config.WebsocketServerEnabled {
		go server.serverWS(wsl)
		log.Infoln("WebsocketServer start")
	}

	if server.config.GRPCServerEnabled {
		go server.severGRPC(grpcL)
		log.Infoln("GRPCServer start")
	}

	if server.config.HTTPServerEnabled {
		go server.serverHTTP(httpL)
		log.Infoln("HTTPServer start")
	}

	go func() {
		if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
			panic(err)
		}
	}()
	log.Infof("ServerRun:%d", server.config.Port)
	return server, nil
}

func (s *Server) serverHTTP(l net.Listener) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Recovery())

	if s.config.TraceEnabled {
		router.Use(middleware.RequestLog())
	}

	if s.config.CorsEnabled {
		router.Use(middleware.Cors(s.config.AllowOrigins))
	}

	router.GET("/favicon.ico", func(c *gin.Context) {})
	router.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	v1 := router.Group("/scaffold/v1")
	{
		s.userServiceRegister(v1)
	}

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: router,
	}

	s.httpServer.Serve(l)
}

func (s *Server) severGRPC(l net.Listener) {
	var customFunc grpc_recovery.RecoveryHandlerFunc
	logrusEntry := log.NewEntry(log.New())
	// Shared options for the logger, with a custom duration to log field function.
	logrusopts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}

	s.grpcServer = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, logrusopts...),
			grpc_recovery.UnaryServerInterceptor(opts...),
			grpc_auth.UnaryServerInterceptor(grpc_auth.AuthFunc(s.jwtAuth.BuildGRPCAUTH())),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrusEntry, logrusopts...),
			grpc_recovery.StreamServerInterceptor(opts...),
			grpc_auth.StreamServerInterceptor(s.jwtAuth.BuildGRPCAUTH()),
		),
	)

	pb.RegisterContentServiceServer(s.grpcServer, &rpc.Handler{})
	s.grpcServer.Serve(l)
}

func (s *Server) serverWS(l net.Listener) {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Recovery())

	v1 := router.Group("/scaffold/v1")
	{
		s.wsServiceRegister(v1)
	}

	s.wsServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: router,
	}

	s.wsServer.Serve(l)
}

// Stop ...
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if s.config.WebsocketServerEnabled {
		if err := s.wsServer.Shutdown(ctx); err != nil {
			log.WithError(err).Fatal("WebsocketServerStop")
		}
		log.Infoln("WebsocketServer stopped")
	}

	if s.config.HTTPServerEnabled {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.WithError(err).Fatal("HTTPServerStop")
		}
		log.Infoln("HTTPServer stopped")
	}

	if s.config.GRPCServerEnabled {
		s.grpcServer.Stop()
		log.Infoln("GRPCServer stopped")
	}

	close(s.c.GetWg().ExitChan)
	s.c.GetWg().Wait()

	log.Infoln("ServerStop stopped")
}
