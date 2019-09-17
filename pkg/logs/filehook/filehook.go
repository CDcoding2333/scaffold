package filehook

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

// Config ...
type Config struct {
	Path         string
	RotationTime time.Duration
	MaxAge       time.Duration
}

// NewLfsHook ...
func NewLfsHook(conf *Config, format log.Formatter) (log.Hook, error) {
	writer, err := rotatelogs.New(
		conf.Path+".%Y%m%d%H",
		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(conf.Path),
		// 多长时间分割一次文件
		rotatelogs.WithRotationTime(conf.RotationTime),
		// 日志文件保留最长时间
		rotatelogs.WithMaxAge(conf.MaxAge),
	)

	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
		return nil, err
	}

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, format)

	return lfsHook, nil
}
