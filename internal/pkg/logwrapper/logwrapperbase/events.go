package logwrapperbase

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Event stores messages to log later, from our logger interface
type Event struct {
	logCategory    int
	logSubCategory int
	eventLevel     logrus.Level
	message        string
}

func NewEvent(logCategory int, logSubCategory int, eventLevel logrus.Level, message, packageName string) (e Event) {
	if eventLevel < logrus.WarnLevel && logCategory < 300 || eventLevel >= logrus.WarnLevel && logCategory >= 300 {
		levelString := []string{"Panic", "Fatal", "Error", "Warn", "Info", "Debug", "Trace"}[eventLevel]
		panic(fmt.Sprintf("Assigned invalid event_level (%s) and log_id (%d)\npackage: %s", levelString, logCategory, packageName))
	}
	e = Event{
		logCategory:    logCategory,
		logSubCategory: logSubCategory,
		eventLevel:     eventLevel,
		message:        message,
	}
	return
}

func newEvent(logCategory int, logSubCategory int, eventLevel logrus.Level, message string) (e Event) {
	e = NewEvent(logCategory, logSubCategory, eventLevel, message, "logwrapper")
	return
}

// Declare variables to store log messages as new Events
var (
	errorMessage           = newEvent(500, 0, logrus.ErrorLevel, "Got an error: %v")
	fatalMessage           = newEvent(500, 0, logrus.FatalLevel, "Got an error: %v")
	invalidArgMessage      = newEvent(406, 1, logrus.ErrorLevel, "Invalid arg: %s")
	invalidArgValueMessage = newEvent(406, 2, logrus.ErrorLevel, "Invalid value for argument: %s: %v")
	missingArgMessage      = newEvent(400, 1, logrus.ErrorLevel, "Missing arg: %s")
)

// InvalidArg is a standard error message
func (l *LoggerBase) InvalidArg(argumentName string) {
	// l.Errorf(invalidArgMessage.message, argumentName)
	l.Log(invalidArgMessage, argumentName)
}

// InvalidArgValue is a standard error message
func (l *LoggerBase) InvalidArgValue(argumentName string, argumentValue string) {
	// l.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
	l.Log(invalidArgValueMessage, argumentName, argumentValue)
}

// MissingArg is a standard error message
func (l *LoggerBase) MissingArg(argumentName string) {
	// l.Errorf(missingArgMessage.message, argumentName)
	l.Log(missingArgMessage, argumentName)
}

// Error is a standard error message
func (l *LoggerBase) ErrorFunc(funcName string, err error) {
	fields := map[string]interface{}{"funcName": funcName}
	l.LogWithFields(errorMessage, fields, err)
}

func (l *LoggerBase) FatalFunc(funcName string, err error) {
	fields := map[string]interface{}{"funcName": funcName}
	l.LogWithFields(fatalMessage, fields, err)
}
