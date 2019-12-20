package logwrapper

import (
	log "github.com/sirupsen/logrus"
)

// type Logger interface {
// 	Log(v ...interface{})
// 	Logf(format string, v ...interface{})
// }

type Event struct {
	id      int
	message string
}

type StdLogger struct {
	*log.Logger
}

const (
	DebugLevel = log.DebugLevel
)

var (
	invalidArgMessage      = Event{1, "Invalid arg: %s"}
	invalidArgValueMessage = Event{2, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{3, "Missing arg: %s"}
	invalidConfigMessage   = Event{4, "Invalid config: %s"}
	missingConfigValue     = Event{5, "Missing config value: %s"}
	JSONFormatter          = log.JSONFormatter{}
)

func New() *StdLogger {
	var baseLogger = log.New()
	var stdLogger = &StdLogger{baseLogger}
	stdLogger.Formatter = &log.TextFormatter{}

	return stdLogger
}

func (l *StdLogger) InvalidArg(argumentName string) {
	l.Errorf(invalidArgMessage.message, argumentName)
}

func (l *StdLogger) InvalidArgValue(argumentName string, argumentValue string) {
	l.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

func (l *StdLogger) MissingArg(argumentName string) {
	l.Errorf(missingArgMessage.message, argumentName)
}

func (l *StdLogger) InvalidConfig(argumentName string) {
	l.Errorf(invalidArgMessage.message, argumentName)
}

func (l *StdLogger) Log(v ...interface{}) {
	l.Print(v...)
}

func (l *StdLogger) Logf(format string, v ...interface{}) {
	l.Printf(format, v...)
}
