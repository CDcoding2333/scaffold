package main

import (
	"CDcoding2333/scaffold/conf"
	"CDcoding2333/scaffold/core"
	"CDcoding2333/scaffold/dao"
	"CDcoding2333/scaffold/pkg/cache"
	"CDcoding2333/scaffold/pkg/flagtools"
	"CDcoding2333/scaffold/pkg/logs"
	"CDcoding2333/scaffold/server"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	baseConf *conf.BaseConfig
	// globalConf *conf.GlobalConfig
)

func init() {
	baseConf = &conf.BaseConfig{RedisConf: new(conf.RedisConfig), DBConf: new(conf.DbConfig), LogConf: new(conf.LogConfig), ServerConf: new(conf.ServerConfig)}

	pflag.StringVar(&baseConf.RedisConf.Host, "redis_host", "127.0.0.1", "redis host")
	pflag.StringVar(&baseConf.RedisConf.Password, "redis_password", "", "redis password")
	pflag.IntVar(&baseConf.RedisConf.Port, "redis_port", 6379, "redis port")
	pflag.IntVar(&baseConf.RedisConf.DB, "redis_db", 0, "redis db")

	pflag.StringVar(&baseConf.DBConf.Driver, "db_driver", "mysql", "db driver")
	pflag.StringVar(&baseConf.DBConf.Source, "db_con", "root:root@tcp(127.0.0.1:3306)/scaffold?charset=utf8&parseTime=True&loc=Local", "driver connection string")
	pflag.BoolVar(&baseConf.DBConf.LogMode, "db_log", true, "db log")

	pflag.IntVar(&baseConf.LogConf.Level, "log_level", 5, "redis db")
	pflag.BoolVar(&baseConf.LogConf.EnableFile, "log_file_enabled", true, "logs to file enabled")
	pflag.StringVar(&baseConf.LogConf.LogPath, "log_file_path", "./logs/app.log", "logs file path")

	pflag.BoolVar(&baseConf.ServerConf.HTTPServerEnabled, "server_http_server_enabled", true, "open or close http server")
	pflag.BoolVar(&baseConf.ServerConf.CorsEnabled, "server_cors", false, "open cors")
	pflag.StringSliceVar(&baseConf.ServerConf.AllowOrigins, "server_origins", []string{"http://127.0.0.1:8080"}, "allow origins")
	pflag.BoolVar(&baseConf.ServerConf.TraceEnabled, "server_trace_enabled", true, "enabled trace to log")
	pflag.BoolVar(&baseConf.ServerConf.GRPCServerEnabled, "server_grpc_server_enabled", true, "open or close grpc server")
	pflag.BoolVar(&baseConf.ServerConf.WebsocketServerEnabled, "server_ws_server_enabled", true, "open or close ws server")
	pflag.IntVar(&baseConf.ServerConf.Port, "server_port", 8080, "redis db")
	pflag.StringVar(&baseConf.ServerConf.JwtISS, "server_jwt_iss", "9p6yjuoaVxn0VwwmSttIcm3XJwmcfRCk", "server auth jwt iss")
	pflag.StringVar(&baseConf.ServerConf.JwtScrect, "server_jwt_screct", "cd_scaffold", "server auth jwt secret")
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())

	flagtools.InitFlags()

	if err := logs.InitLogs(baseConf.LogConf); err != nil {
		return
	}

	if err := dao.InitDB(baseConf.DBConf); err != nil {
		log.WithError(err).Fatal("init db error")
		return
	}

	if err := cache.InitRedis(baseConf.RedisConf); err != nil {
		log.WithError(err).Fatal("init redis error")
		return
	}

	server, err := server.Run(core.New(), baseConf.ServerConf)
	if err != nil {
		log.WithError(err).Fatal("init server error")
		return
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.Stop()
	log.Infoln("process exit")
}
