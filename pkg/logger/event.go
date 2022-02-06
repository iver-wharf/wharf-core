package logger

import (
	"fmt"
	"sync"
	"time"

	"github.com/iver-wharf/wharf-core/v2/internal/traceutil"
)

// DoneFunc is the signature of the function that is called at the end of a
// submitted log event.
type DoneFunc func(message string)

// Event is a single log message that's aimed to be submitted. It may hold
// multiple logging contexts created by log sinks used to finally submit to that
// range of sinks.
type Event interface {
	// Messagef submits this log event to the different sinks using a formatted
	// message. The formatting is the same applied from the fmt package.
	Messagef(format string, args ...interface{})

	// Message submits this log event to the different sinks using a message. To
	// submit without a message you may pass an empty string into this method, like
	// so:
	// 	ev.WithString("hello", "world").Message("")
	Message(message string)

	// WithFunc applies a function to the event and then forwards the return value.
	//
	// Useful for reusing "with statements" for multiple logs.
	WithFunc(f func(Event) Event) Event

	// WithCaller adds a caller field to the log contexts inside this log event.
	//
	// This method is called automatically by NewEvent and all Logger methods,
	// though you can override the value set there by calling it again manually.
	//
	// It's up to the logger sink to decide how this error is rendered in the log
	// message. Commonly, but not necessarily, this is rendered as fields with names
	// "caller" and "line".
	WithCaller(file string, line int) Event

	// WithString adds a string field to this logged message. Calling this method
	// multiple times with the same key may lead to unexpected behaviour.
	WithString(key string, value string) Event

	// WithStringf adds a formatted string field to this logged message. The
	// formatting is the same applied from the fmt package. Calling this method
	// multiple times with the same key may lead to unexpected behaviour.
	WithStringf(key string, format string, args ...interface{}) Event

	// WithStringer adds a string field to this logged message using the value
	// from fmt.Stringer.String(). Calling this method multiple times with the
	// same key may lead to unexpected behaviour.
	WithStringer(key string, value fmt.Stringer) Event

	// WithRune adds a rune field to this logged message. Calling this method
	// multiple times with the same key may lead to unexpected behaviour.
	WithRune(key string, value rune) Event

	// WithBool adds a boolean field to this logged message. Calling this method
	// multiple times with the same key may lead to unexpected behaviour.
	WithBool(key string, value bool) Event

	// WithInt adds an integer field to this logged message. Calling this method
	// multiple times with the same key may lead to unexpected behaviour.
	WithInt(key string, value int) Event

	// WithInt64 adds an integer field to this logged message. Calling this method
	// multiple times with the same key may lead to unexpected behaviour.
	WithInt64(key string, value int64) Event

	// WithInt32 adds an integer field to this logged message. Calling this method
	// multiple times with the same key may lead to unexpected behaviour.
	WithInt32(key string, value int32) Event

	// WithUint adds an integer field to this logged message. Calling this method
	// multiple times with the same key may lead to unexpected behaviour.
	WithUint(key string, value uint) Event

	// WithUint64 adds an integer field to this logged message. Calling this method
	// multiple times with the same key may lead to unexpected behaviour.
	WithUint64(key string, value uint64) Event

	// WithUint32 adds an integer field to this logged message. Calling this method
	// multiple times with the same key may lead to unexpected behaviour.
	WithUint32(key string, value uint32) Event

	// WithFloat32 adds a floating point number field to this logged message. Calling
	// this method multiple times with the same key may lead to unexpected behaviour.
	WithFloat32(key string, value float32) Event

	// WithFloat64 adds a floating point number field to this logged message. Calling
	// this method multiple times with the same key may lead to unexpected behaviour.
	WithFloat64(key string, value float64) Event

	// WithError adds an error field to this logged message. Calling this method
	// multiple times may lead to unexpected behaviour.
	//
	// It's up to the logger sink to decide how this error is rendered in the log
	// message. Commonly, but not necessarily, this is rendered as a field with name
	// "error".
	WithError(value error) Event

	// WithTime adds a timestamp field to this logged message. Calling
	// this method multiple times with the same key may lead to unexpected behaviour.
	//
	// It's up to the logger sink to decide how this error is rendered in the log
	// message, e.g. in UNIX timestamp integer form or string formatted datetime.
	WithTime(key string, value time.Time) Event

	// WithDuration adds a timestamp field to this logged message. Calling
	// this method multiple times with the same key may lead to unexpected behaviour.
	//
	// It's up to the logger sink to decide how this error is rendered in the log
	// message, e.g. in milliseconds integer form or string formatted duration.
	WithDuration(key string, value time.Duration) Event
}

var contextPool = sync.Pool{
	New: func() interface{} {
		return []Context{}
	},
}

type event struct {
	level Level
	ctxs  []Context
	done  DoneFunc
}

// NewEvent creates a new event and prepares it to use a list of logging sinks
// based on the logging level fed into it using the globally registered sinks
// added using logger.AddOutput(...).
func NewEvent(level Level, scope string, done DoneFunc) Event {
	return newEventFromSinks(level, scope, done, registeredSinks)
}

func newEventFromSinks(level Level, scope string, done DoneFunc, sinks []registeredSink) Event {
	if level < getLevelScoped(scope) {
		return event{}
	}
	ctxs := contextPool.Get().([]Context)
	ctxs = ctxs[:0]
	for _, reg := range sinks {
		if level < reg.minLevel {
			continue
		}
		ctxs = append(ctxs, reg.sink.NewContext(scope))
	}
	ev := event{level, ctxs, done}
	if caller, line := traceutil.CallerFileWithLineNum(); caller != "" {
		return ev.WithCaller(caller, line)
	}
	return ev
}

// NewEventFromLogger creates an event using the logger itself based on the
// logging level. Useful in edge-cases and when testing with a slice of test
// cases.
func NewEventFromLogger(log Logger, level Level) Event {
	switch level {
	case LevelDebug:
		return log.Debug()
	case LevelInfo:
		return log.Info()
	case LevelWarn:
		return log.Warn()
	case LevelError:
		return log.Error()
	case LevelPanic:
		return log.Panic()
	default:
		panic(fmt.Sprintf("invalid log level: %s", level))
	}
}

func (ev event) Messagef(format string, args ...interface{}) {
	if len(ev.ctxs) > 0 {
		ev.Message(fmt.Sprintf(format, args...))
	} else {
		ev.returnPooledSlice()
		if ev.done != nil {
			ev.done(fmt.Sprintf(format, args...))
		}
	}
}

func (ev event) Message(message string) {
	for _, log := range ev.ctxs {
		log.WriteOut(ev.level, message)
	}
	ev.returnPooledSlice()
	if ev.done != nil {
		ev.done(message)
	}
}

func (ev event) returnPooledSlice() {
	if ev.ctxs != nil {
		contextPool.Put(ev.ctxs)
	}
}

func (ev event) WithFunc(f func(Event) Event) Event {
	return f(ev)
}

func (ev event) WithCaller(file string, line int) Event {
	return ev.with(func(ctx Context) Context { return ctx.SetCaller(file, line) })
}

func (ev event) WithString(key string, value string) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendString(key, value) })
}

func (ev event) WithStringf(key string, format string, args ...interface{}) Event {
	if len(ev.ctxs) > 0 {
		return ev.WithString(key, fmt.Sprintf(format, args...))
	}
	return ev
}

func (ev event) WithStringer(key string, value fmt.Stringer) Event {
	if len(ev.ctxs) > 0 {
		return ev.WithString(key, value.String())
	}
	return ev
}

func (ev event) WithRune(key string, value rune) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendRune(key, value) })
}

func (ev event) WithBool(key string, value bool) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendBool(key, value) })
}

func (ev event) WithInt(key string, value int) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendInt(key, value) })
}

func (ev event) WithInt64(key string, value int64) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendInt64(key, value) })
}

func (ev event) WithInt32(key string, value int32) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendInt32(key, value) })
}

func (ev event) WithUint(key string, value uint) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendUint(key, value) })
}

func (ev event) WithUint64(key string, value uint64) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendUint64(key, value) })
}

func (ev event) WithUint32(key string, value uint32) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendUint32(key, value) })
}

func (ev event) WithFloat32(key string, value float32) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendFloat32(key, value) })
}

func (ev event) WithFloat64(key string, value float64) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendFloat64(key, value) })
}

func (ev event) WithError(value error) Event {
	return ev.with(func(ctx Context) Context { return ctx.SetError(value) })
}

func (ev event) WithTime(key string, value time.Time) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendTime(key, value) })
}

func (ev event) WithDuration(key string, value time.Duration) Event {
	return ev.with(func(ctx Context) Context { return ctx.AppendDuration(key, value) })
}

func (ev event) with(f func(Context) Context) Event {
	for i, ctx := range ev.ctxs {
		ev.ctxs[i] = f(ctx)
	}
	return ev
}
