package standardlogger

import (
	"github.com/gstdl/crypto-arbitrage/internal/pkg/logwrapper/logwrapperbase"
	"github.com/sirupsen/logrus"
)

func newEvent(logCategory int, logSubCategory int, eventLevel logrus.Level, message string) (e logwrapperbase.Event) {
	e = logwrapperbase.NewEvent(logCategory, logSubCategory, eventLevel, message, "standard")
	return
}

// Declare variables to store log messages as new Events
var (
	serviceStartMessage   = newEvent(200, 0, logrus.InfoLevel, "Starting Service!")
	serviceStopMessage    = newEvent(500, 0, logrus.ErrorLevel, "Service stopped! Error: %v")
	invalidRequestMessage = newEvent(500, 0, logrus.ErrorLevel, "Service stopped! Error: %v")
)

func (l *StandardLogger) LogServiceStart() {
	l.Log(serviceStartMessage)
}

func (l *StandardLogger) LogServiceStopped(err interface{}) {
	l.Log(serviceStopMessage, err)
}

func (l *StandardLogger) LogRequestErrors(req interface{}, err error) {
	fields := map[string]interface{}{"request": req}
	l.LogWithFields(invalidRequestMessage, fields, err)
}
