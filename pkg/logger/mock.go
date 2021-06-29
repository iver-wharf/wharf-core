package logger

import "time"

// MockLogger is a logger meant to be used in testing. It records all logs sent
// to it and provides useful fields for you to verify that your application
// logs as expected.
type MockLogger struct {
	// Logs is an array of logs recorded. Each logged event is stored as a
	// separate item in this array..
	Logs []MockLog
}

// MockLog is a single logged event with additional public fields containing
// data of the logged event.
type MockLog struct {
	// Level is the logging level that was used. If calling MockLogger.Debug(),
	// then this field holds the value of LevelDebug.
	Level Level
	// Message is the final message that was used to submit this logged event.
	Message string
	// Fields holds each field added to this logged event. Some logged data is
	// translated into this map array:
	//
	// 	Event.WithScope("foo")   => MockLog.Fields["scope"] = "foo"
	// 	Event.WithError(someErr) => MockLog.Fields["error"] = someErr
	// 	Event.WithCaller("foo", 42)
	// 		=> MockLog.Fields["caller"] = "foo"
	// 		=> MockLog.Fields["line"] = 42
	Fields map[string]interface{}
	// FieldNames is a slice of strings with all the keys added to the Fields
	// map. This includes the custom mapping of Event.WithScope,
	// Event.WithError, and Event.WithCaller as mentioned in the Fields docs.
	FieldNames []string
}

// NewMock creates a new Logger interface compatible type that holds additional
// fields of the logs that has been submitted.
func NewMock() *MockLogger {
	return &MockLogger{}
}

// Debug creates a new event using new contexts connected to this mock logger of
// "debugging" logging level or higher.
func (log *MockLogger) Debug() Event { return log.newEvent(LevelDebug) }

// Info creates a new event using new contexts connected to this mock logger of
// "information" logging level or higher.
func (log *MockLogger) Info() Event { return log.newEvent(LevelInfo) }

// Warn creates a new event using new contexts connected to this mock logger of
// "warning" logging level or higher.
func (log *MockLogger) Warn() Event { return log.newEvent(LevelWarn) }

// Error creates a new event using new contexts connected to this mock logger of
// "error" logging level or higher.
func (log *MockLogger) Error() Event { return log.newEvent(LevelError) }

// Panic creates a new event using new contexts connected to this mock logger of
// "panic" logging level or higher.
//
// Compared to the other logging events, after submitting the logged
// messages this method calls panic with the final message string.
func (log *MockLogger) Panic() Event { return log.newEvent(LevelPanic) }

func (log *MockLogger) newEvent(level Level) Event {
	var done DoneFunc
	if level == LevelPanic {
		done = panicString
	}
	return event{
		level: level,
		ctxs:  []Context{newMockContext(log)},
		done:  done,
	}
}

func newMockContext(logger *MockLogger) Context {
	return mockCtx{
		MockLog: MockLog{
			Fields: map[string]interface{}{},
		},
		logger: logger,
	}
}

type mockCtx struct {
	MockLog
	logger *MockLogger
}

func (c mockCtx) WriteOut(level Level, message string) {
	c.Level = level
	c.Message = message
	c.logger.Logs = append(c.logger.Logs, c.MockLog)
}

func (c mockCtx) SetCaller(file string, line int) Context {
	c.Fields["caller"] = file
	c.Fields["line"] = line
	c.FieldNames = append(c.FieldNames, "caller", "line")
	return c
}

func (c mockCtx) SetScope(v string) Context                        { return c.addField("scope", v) }
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

func (c mockCtx) addField(key string, value interface{}) Context {
	c.Fields[key] = value
	c.FieldNames = append(c.FieldNames, key)
	return c
}
