package logger

import (
	"io"
	"strings"
)

// Requirements:
// - [x] LogLevel
// - [x] Multiple writers/sinks (ex: for Kafka, for console)
// - [x] Scope per-module configs (ex: warn for Gin, info for API audits)
// - [x] Easy integration with Gin
// - [x] Easy integration with GORM
// - [ ] Easy integration with io.Writer's, ex gin.DefaultWriter & gin.DefaultErrorWriter
// - [x] Fields

// Sink is an interface that creates logging contexts. Each sink could be for
// different log collectors such as Kibana or Logstash, or simply a console
// logging sink that outputs all the logs to STDOUT.
type Sink interface {
	NewContext() Context
}

type registeredSink struct {
	sink     Sink
	minLevel Level
}

var registeredSinks []registeredSink

// ClearOutputs resets the outputs added by AddOutput. Should not be needed in
// production code, but is quite useful to be called at the beginning of an
// example test.
func ClearOutputs() {
	registeredSinks = nil
}

// AddOutput registers a logging sink globally. Multiple sinks can be added, and
// they will be used in the order of when they are added.
//
// To only use a particular sink for warning or higher logging levels, you pass
// in the warning log level:
//
// 	logger.AddOutput(logger.LevelWarn, myLogSink)
//
// To let a particular sink log all messages, use the "debug" logging level:
//
// 	logger.AddOutput(logger.LevelDebug, myLogSink)
func AddOutput(minLevel Level, sink Sink) {
	registeredSinks = append(registeredSinks, registeredSink{
		sink:     sink,
		minLevel: minLevel,
	})
}

// Logger is an interface that is used to initiate logging events of different
// log levels. This is done before populating the log messages with fields so
// that those calls can be ignored if no sink listens for that particular
// logging level.
type Logger interface {
	// Debug creates a new event using new contexts from all registered sinks of
	// "debugging" logging level or higher.
	Debug() Event
	// Info creates a new event using new contexts from all registered sinks of
	// "information" logging level or higher.
	Info() Event
	// Warn creates a new event using new contexts from all registered sinks of
	// "warning" logging level or higher.
	Warn() Event
	// Error creates a new event using new contexts from all registered sinks of
	// "error" logging level or higher.
	Error() Event
	// Panic creates a new event using new contexts from all registered sinks of
	// "panic" logging level or higher.
	//
	// Compared to the other logging events, after submitting the logged
	// messages this method calls panic with the final message string.
	Panic() Event
}

type logger struct {
	newEvent func(Level, DoneFunc) Event
}

// New creates a new basic Logger without a scope. Use NewScoped instead to add
// a "scope" field to each logged message.
func New() Logger {
	return logger{NewEvent}
}

func (log logger) Debug() Event { return log.newEvent(LevelDebug, nil) }
func (log logger) Info() Event  { return log.newEvent(LevelInfo, nil) }
func (log logger) Warn() Event  { return log.newEvent(LevelWarn, nil) }
func (log logger) Error() Event { return log.newEvent(LevelError, nil) }
func (log logger) Panic() Event { return log.newEvent(LevelPanic, panicString) }

func panicString(message string) {
	panic(message)
}

// NewScoped creates a new logger and assigns a scope to it. Useful when you
// want to group logs from different parts of the system on a string name.
//
// For example:
// 	logger.NewScoped("GORM") // use when registering logger to GORM
// 	logger.NewScoped("GIN") // use when registering logger to gin-gonic
// 	logger.New() // use in the apps top-level domain
func NewScoped(scope string) Logger {
	return logger{
		newEvent: func(level Level, done DoneFunc) Event {
			return NewEvent(level, done).WithScope(scope)
		},
	}
}

type loggerWriter struct {
	logger Logger
	level  Level
}

// NewWriter creates a logger that channels everything written to it via a
// wharf-core logger, using the given logging level for all of the logs.
func NewWriter(log Logger, level Level) io.Writer {
	return loggerWriter{log, level}
}

func (w loggerWriter) Write(p []byte) (n int, err error) {
	var message = strings.TrimRight(string(p), "\n")
	NewEventFromLogger(w.logger, w.level).Message(message)
	return len(p), nil
}
