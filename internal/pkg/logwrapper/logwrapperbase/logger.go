package logwrapperbase

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/gstdl/crypto-arbitrage/internal/config"
)

type LoggerBase struct {
	*logrus.Logger
}

func NewLogger(packagePath, exchange, loggerType string) (logger *LoggerBase) {
	environment := config.GetEnvironmentValue()

	logger = &LoggerBase{logrus.New()}

	logger.AddEnvironmentHook(environment)
	logger.AddPackageHook(packagePath)
	logger.AddExchangeHook(exchange)
	logger.AddLogTypeHook(loggerType)

	switch environment {
	case "production":
		logger.SetOutput(io.Discard)
	default:
		logger.SetOutput(os.Stderr)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
			ForceQuote:    true,
		})
	}

	return
}

func (l *LoggerBase) LogToFile(logFilePath string) {
	l.LogToFileWithLevel(logFilePath, logrus.AllLevels)
}

func (l *LoggerBase) LogToFileWithLevel(logFilePath string, level []logrus.Level) {
	if reflect.DeepEqual(level, logrus.AllLevels) {
		l.AddHook(lfshook.NewHook(logFilePath, &logrus.JSONFormatter{}))
	} else {
		f := make(lfshook.PathMap)
		for _, l := range level {
			f[l] = logFilePath
		}
		l.AddHook(lfshook.NewHook(f, &logrus.JSONFormatter{}))
	}
}

func (l *LoggerBase) LogWithFields(e Event, fields map[string]interface{}, args ...interface{}) {
	log := l.WithFields(logrus.Fields{
		"log_category":     e.logCategory,
		"log_sub_category": e.logSubCategory,
	})
	for key, value := range fields {
		log = log.WithField(key, value)
	}
	switch e.eventLevel {
	case logrus.TraceLevel:
		log.Tracef(e.message, args...)
	case logrus.DebugLevel:
		log.Debugf(e.message, args...)
	case logrus.InfoLevel:
		log.Infof(e.message, args...)
	case logrus.WarnLevel:
		log.Warnf(e.message, args...)
	case logrus.FatalLevel:
		log.Fatalf(e.message, args...)
	case logrus.ErrorLevel:
		log.Errorf(e.message, args...)
	default:
		// log.Panicf("Invalid logging level")
		l.Panic(fmt.Errorf("invalid logging level"))
	}
}

func (l *LoggerBase) Log(e Event, args ...interface{}) {
	l.LogWithFields(e, nil, args...)
}

func (l *LoggerBase) Panicf(msg string, err error) {
	l.WithFields(logrus.Fields{
		"event_id":     500,
		"event_sub_id": 500,
	}).WithError(
		err,
	).Panic(msg)
}

func (l *LoggerBase) Panic(err error) {
	l.Panicf("", err)
}

func (l *LoggerBase) Fatalf(msg string, err error) {
	l.WithFields(logrus.Fields{
		"event_id":     500,
		"event_sub_id": 500,
	}).WithError(
		err,
	).Fatal(msg)
}

func (l *LoggerBase) Fatal(err error) {
	l.Fatalf("", err)
}

func (l *LoggerBase) FmtPrintf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func (l *LoggerBase) FmtPrintln(format string) {
	fmt.Println(format)
}
