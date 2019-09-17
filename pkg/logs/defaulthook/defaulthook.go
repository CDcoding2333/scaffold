package defaulthook

import (
	log "github.com/sirupsen/logrus"
)

// DefaultFieldHook ...
type DefaultFieldHook struct {
	AppName string
}

// Fire ...
func (hook *DefaultFieldHook) Fire(entry *log.Entry) error {
	entry.Data["appName"] = hook.AppName
	return nil
}

// Levels ...
func (hook *DefaultFieldHook) Levels() []log.Level {
	return log.AllLevels
}
