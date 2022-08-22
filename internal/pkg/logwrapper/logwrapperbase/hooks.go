package logwrapperbase

import (
	"github.com/sirupsen/logrus"
)

type fieldHook struct {
	key   string
	value interface{}
}

func (h *fieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *fieldHook) Fire(e *logrus.Entry) error {
	e.Data[h.key] = h.value
	return nil
}

func (l *LoggerBase) AddFieldHook(key string, value interface{}) {
	l.AddHook(&fieldHook{key: key, value: value})
}

func (l *LoggerBase) AddPackageHook(packagePath string) {
	l.AddFieldHook("package", packagePath)
}

func (l *LoggerBase) AddLogTypeHook(logType string) {
	l.AddFieldHook("log_type", logType)
}

func (l *LoggerBase) AddExchangeHook(exchangeName string) {
	l.AddFieldHook("exchange_name", exchangeName)
}

func (l *LoggerBase) AddEnvironmentHook(exchangeName string) {
	l.AddFieldHook("environment", exchangeName)
}
