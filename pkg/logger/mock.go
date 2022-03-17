package logger

import "time"

// Mock is a logger and a sink meant to be used in testing. It records all
// logs sent to it and provides useful fields for you to verify that your
// application logs as expected.
type Mock struct {
	// Logs is an array of logs recorded. Each logged event is stored as a
	// separate item in this array.
	Logs []MockLog
	// LogMessages is an array of log messages recorded. This complements the
	// array of MockLog for easier assertion that the expected messages have been
	// logged. Empty messages are also stored in this array as empty strings.
	LogMessages []string
}

// NewMock creates a new Logger interface compatible type that holds additional
// fields of the logs that have been submitted.
func NewMock() *Mock {
	return &Mock{}
}

// NewContext creates a new log event context for this mock. The scope is added
// as a field unless it's an empty string.
func (log *Mock) NewContext(scope string) Context {
	ctx := mockCtx{
		MockLog: MockLog{
			Fields: make(map[string]any),
		},
		logger: log,
	}
	if scope != "" {
		ctx.addField("scope", scope)
	}
	return ctx
}

// MockLog is a single logged event with additional public fields containing
// data of the logged event.
type MockLog struct {
	// Level is the logging level that was used. If calling MockLogger.Debug(),
	// then this field holds the value of LevelDebug.
	Level Level
	// Message is the final message that was used to submit this logged event.
	Message string
	// Fields holds each field added to this logged event, in addition to the
	// following fields:
	//
	// 	Event.SetScope("foo")   => MockLog.Fields["scope"] = "foo"
	// 	Event.SetError(someErr) => MockLog.Fields["error"] = someErr
	// 	Event.SetCaller("foo", 42)
	// 		=> MockLog.Fields["caller"] = "foo"
	// 		=> MockLog.Fields["line"] = 42
	Fields map[string]any
	// FieldsAdded is a slice of strings with all the keys added to the Fields
	// map. This includes the custom mapping of Event.SetScope,
	// Event.SetError, and Event.SetCaller as mentioned in the Fields docs.
	//
	// If a field is added more than one time, then it will show up in this list
	// equally many times. Useful for checking if fields are misstakenly added
	// multiple times.
	FieldsAdded []string
}

// Debug creates a new event using new contexts connected to this mock logger of
// "debugging" logging level or higher.
func (log *Mock) Debug() Event { return log.newEvent(LevelDebug) }

// Info creates a new event using new contexts connected to this mock logger of
// "information" logging level or higher.
func (log *Mock) Info() Event { return log.newEvent(LevelInfo) }

// Warn creates a new event using new contexts connected to this mock logger of
// "warning" logging level or higher.
func (log *Mock) Warn() Event { return log.newEvent(LevelWarn) }

// Error creates a new event using new contexts connected to this mock logger of
// "error" logging level or higher.
func (log *Mock) Error() Event { return log.newEvent(LevelError) }

// Panic creates a new event using new contexts connected to this mock logger of
// "panic" logging level or higher.
//
// Compared to the other logging events, after submitting the logged
// messages this method calls panic with the final message string.
func (log *Mock) Panic() Event { return log.newEvent(LevelPanic) }

func (log *Mock) newEvent(level Level) Event {
	var done DoneFunc
	if level == LevelPanic {
		done = panicString
	}
	return newEventFromSinks(level, "", done, []registeredSink{
		{log, LevelDebug},
	})
}

type mockCtx struct {
	MockLog
	logger *Mock
}

func (c mockCtx) WriteOut(level Level, message string) {
	c.Level = level
	c.Message = message
	c.logger.Logs = append(c.logger.Logs, c.MockLog)
	c.logger.LogMessages = append(c.logger.LogMessages, message)
}

func (c mockCtx) SetCaller(file string, line int) Context {
	c.Fields["caller"] = file
	c.Fields["line"] = line
	c.FieldsAdded = append(c.FieldsAdded, "caller", "line")
	return c
}

func (c mockCtx) SetError(v error) Context                         { return c.addField("error", v) }
func (c mockCtx) AppendString(k string, v string) Context          { return c.addField(k, v) }
func (c mockCtx) AppendRune(k string, v rune) Context              { return c.addField(k, v) }
func (c mockCtx) AppendBool(k string, v bool) Context              { return c.addField(k, v) }
func (c mockCtx) AppendInt(k string, v int) Context                { return c.addField(k, v) }
func (c mockCtx) AppendInt32(k string, v int32) Context            { return c.addField(k, v) }
func (c mockCtx) AppendInt64(k string, v int64) Context            { return c.addField(k, v) }
func (c mockCtx) AppendUint(k string, v uint) Context              { return c.addField(k, v) }
func (c mockCtx) AppendUint32(k string, v uint32) Context          { return c.addField(k, v) }
func (c mockCtx) AppendUint64(k string, v uint64) Context          { return c.addField(k, v) }
func (c mockCtx) AppendFloat32(k string, v float32) Context        { return c.addField(k, v) }
func (c mockCtx) AppendFloat64(k string, v float64) Context        { return c.addField(k, v) }
func (c mockCtx) AppendTime(k string, v time.Time) Context         { return c.addField(k, v) }
func (c mockCtx) AppendDuration(k string, v time.Duration) Context { return c.addField(k, v) }

func (c mockCtx) addField(key string, value any) Context {
	c.Fields[key] = value
	c.FieldsAdded = append(c.FieldsAdded, key)
	return c
}
