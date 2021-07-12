package logger

import "time"

// Context is data held about a certain logging event for a particular sink.
// The data can be stored in any way that seems suitable for efficiently
// composing a logging message for that sink.
//
// Most methods return the context itself. It is up to the implementation to
// take advantage of this or use pointers to the same object throughout.
//
// It is up to the user of this type to honor the specification of feeding the
// Context around to itself and calling the Append... methods on the return
// value from any other such method.
//
// Good:
// 	ctx.AppendString("hello", "world").WriteOut(logger.LevelDebug, "")
//
// Good:
// 	ctx = ctx.AppendString("hello", "world")
// 	ctx.WriteOut(logger.LevelDebug, "")
//
// Bad:
// 	ctx.AppendString("hello", "world") // undefined behaviour
// 	ctx.WriteOut(logger.LevelDebug, "")
type Context interface {
	// WriteOut sends the log message with the collected data from all the
	// Append... methods
	WriteOut(level Level, message string)

	// SetCaller sets the caller and its line value for this context.
	//
	// Calling this method multiple times shall override the previous value.
	// An empty string on the file name signifies to unset this field.
	//
	// In contrast to AppendString, the logging sink is allowed to render this
	// differently. E.g. some may render it as yet another fields named "caller"
	// and "line", others may render it as a specific HTTP header in a request.
	SetCaller(file string, line int) Context
	// SetError sets the error value for this context.
	//
	// Calling this method multiple times shall override the previous value.
	// An nil signifies to unset this field.
	//
	// In contrast to AppendString, the logging sink is allowed to render this
	// differently. E.g. some may render it as yet another field named "error",
	// others may render it as a specific HTTP header in a request.
	SetError(value error) Context
	// AppendString adds a string value for a specific key to this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendString(key string, value string) Context
	// AppendRune adds a rune value for a specific key to this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendRune(key string, value rune) Context
	// AppendBool adds a boolean value for a specific key to this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendBool(key string, value bool) Context
	// AppendInt adds an integer value for a specific key to this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendInt(key string, value int) Context
	// AppendInt32 adds an integer value for a specific key to this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendInt32(key string, value int32) Context
	// AppendInt64 adds an integer value for a specific key to this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendInt64(key string, value int64) Context
	// AppendUint adds an integer value for a specific key to this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendUint(key string, value uint) Context
	// AppendUint32 adds an integer value for a specific key to this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendUint32(key string, value uint32) Context
	// AppendUint64 adds an integer value for a specific key to this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendUint64(key string, value uint64) Context
	// AppendFloat32 adds a floating point number value for a specific key to
	// this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendFloat32(key string, value float32) Context
	// AppendFloat64 adds a floating point number value for a specific key to
	// this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendFloat64(key string, value float64) Context
	// AppendTime adds a timestamp value for a specific key to
	// this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendTime(key string, value time.Time) Context
	// AppendDuration adds a time duration value for a specific key to
	// this context.
	//
	// Calling this method multiple times with the same key may lead to
	// unexpected behaviour.
	AppendDuration(key string, value time.Duration) Context
}
