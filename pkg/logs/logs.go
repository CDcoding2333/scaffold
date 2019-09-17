package logs

import (
	"CDcoding2333/scaffold/conf"
	"CDcoding2333/scaffold/pkg/logs/defaulthook"
	"CDcoding2333/scaffold/pkg/logs/filehook"
	"CDcoding2333/scaffold/pkg/logs/filename"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// InitLogs ...
func InitLogs(config *conf.LogConfig) error {

	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
	log.SetLevel(log.Level(config.Level))
	log.SetOutput(os.Stdout)
	log.AddHook(&defaulthook.DefaultFieldHook{AppName: "app"})
	log.AddHook(filename.NewHook())

	if config.EnableFile {
		fileHook, err := filehook.NewLfsHook(&filehook.Config{Path: config.LogPath,
			RotationTime: time.Hour * 24,
			MaxAge:       time.Hour * 24 * 365}, &log.TextFormatter{FullTimestamp: true})
		if err != nil {
			return err
		}
		log.AddHook(fileHook)
	}

	return nil
}
