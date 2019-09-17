package server

import (
	"CDcoding2333/scaffold/apps/user"
	"CDcoding2333/scaffold/apps/ws"
	"CDcoding2333/scaffold/constant"

	"github.com/gin-gonic/gin"
)

// userRegister ...
func (s *Server) userServiceRegister(router *gin.RouterGroup) {
	h := user.NewUserHandler(s.jwtAuth, s.c)
	s.c.RegistHandlers(constant.UserHandler, *h)
	u := router.Group("/users")
	{
		u.POST("/regist", h.UsersRegist)
		u.POST("/login", h.UsersLogin)
	}

	user := router.Group("/users", h.Auth.HandleAuth)
	{
		user.GET("/info", h.GetUserInfo)
		user.DELETE("/", h.DelUsers)
	}
}

func (s *Server) wsServiceRegister(router *gin.RouterGroup) {

	h := ws.NewWebsocketHandler(s.c)
	router.GET("/ws", s.jwtAuth.HandleAuth, h.WebsocketConnect)
}
