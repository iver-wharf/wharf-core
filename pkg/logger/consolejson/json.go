package consolejson

import (
	"encoding/json"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/iver-wharf/wharf-core/v2/pkg/logger"
)

// TimeFormat specifies the formatting used when logging time.Time values.
//
// You may use a custom time format by casting a time-package compatible format
// into this type.
type TimeFormat string

const (
	// TimeRFC3339 will render a time.Time as a string with the format
	// 	2006-01-02T15:04:05Z07:00
	TimeRFC3339 TimeFormat = time.RFC3339
	// TimeUnix will render a time.Time as an integer of seconds since
	// January 1, 1970 UTC.
	TimeUnix TimeFormat = "wharf-core/Unix"
	// TimeUnixMs will render a time.Time as an integer of milliseconds since
	// January 1, 1970 UTC.
	TimeUnixMs TimeFormat = "wharf-core/UnixMs"
	// TimeUnixMicro will render a time.Time as an integer of microseconds since
	// January 1, 1970 UTC.
	TimeUnixMicro TimeFormat = "wharf-core/UnixMicro"
	// TimeUnixNano will render a time.Time as an integer of nanoseconds since
	// January 1, 1970 UTC.
	TimeUnixNano TimeFormat = "wharf-core/UnixNano"
)

// Config lets you gradually configure the output of the logger by disabling
// certain features or changing the format of certain field types.
type Config struct {
	// DisableDate removes the date field from the log when set to true.
	//
	// When set to false:
	// 	{"level":"info","date":"2006-01-02T15:04:05Z","caller":"example.go","line":20,"message":"Sample message."}
	// When set to true:
	// 	{"level":"info","caller":"example.go","line":20,"message":"Sample message."}
	DisableDate bool
	// DisableCaller removes the caller file name and line fields from the log
	// when set to true.
	//
	// When set to false:
	// 	{"level":"info","date":"2006-01-02T15:04:05Z","caller":"example.go","line":20,"message":"Sample message."}
	// When set to true:
	// 	{"level":"info","date":"2006-01-02T15:04:05Z","message":"Sample message."}
	DisableCaller bool
	// DisableCallerLine removes just the caller line field from the log
	// when set to true, but leaves the caller file name as-is.
	//
	// When set to false:
	// 	{"level":"info","caller":"example.go","line":20,"message":"Sample message."}
	// When set to true:
	// 	{"level":"info","caller":"example.go","message":"Sample message."}
	DisableCallerLine bool
	// CallerFileField sets the name of the JSON property used in the logs
	// caller file path. The value is automatically escaped.
	// Defaults to "caller".
	//
	// When set to "" (empty string):
	// 	{"level":"info","caller":"example.go","line":20,"message":"Sample message."}
	// When set to "foo":
	// 	{"level":"info","foo":"example.go","line":20,"message":"Sample message."}
	CallerFileField string
	// CallerLineField sets the name of the JSON property used in the logs
	// caller line number. The value is automatically escaped.
	// Defaults to "line".
	//
	// When set to "" (empty string):
	// 	{"level":"info","caller":"example.go","line":20,"message":"Sample message."}
	// When set to "foo":
	// 	{"level":"info","caller":"example.go","foo":20,"message":"Sample message."}
	CallerLineField string
	// ErrorField sets the name of the JSON property used in the logs error.
	// The value is automatically escaped.
	// Defaults to "error".
	//
	// When set to "" (empty string):
	// 	{"level":"info","message":"Sample message.","error":"strconv.Atoi: parsing \"bar\": invalid syntax"}
	// When set to "foo":
	// 	{"level":"info","message":"Sample message.","foo":"strconv.Atoi: parsing \"bar\": invalid syntax"}
	ErrorField string
	// LevelField sets the name of the JSON property used in the logs severity
	// level. The value is automatically escaped.
	// Defaults to "level".
	//
	// When set to "" (empty string):
	// 	{"level":"info","message":"Sample message."}
	// When set to "foo":
	// 	{"foo":"info","message":"Sample message."}
	LevelField string
	// MessageField sets the name of the JSON property used in the logs message.
	// The value is automatically escaped.
	// Defaults to "message".
	//
	// When set to "" (empty string):
	// 	{"level":"info","message":"Sample message."}
	// When set to "foo":
	// 	{"level":"info","foo":"Sample message."}
	MessageField string
	// ScopeField sets the name of the JSON property used in the logs scope.
	// The value is automatically escaped.
	// Defaults to "scope".
	//
	// When set to "" (empty string):
	// 	{"level":"info","scope":"GORM","message":"Sample message."}
	// When set to "foo":
	// 	{"level":"info","foo":"GORM","message":"Sample message."}
	ScopeField string
	// DateField sets the name of the JSON property used in the logs scope.
	// The value is automatically escaped.
	// Defaults to "date".
	//
	// When set to "" (empty string):
	// 	{"level":"info","date":"2006-01-02T15:04:05Z","message":"Sample message."}
	// When set to "foo":
	// 	{"level":"info","foo":"2006-01-02T15:04:05Z","message":"Sample message."}
	DateField string
	// TiemFormat defines how time.Time fields added via Event.WithTime is
	// rendered. Defaults to TimeRFC3339, which looks like so:
	// 	2006-01-02T15:04:05Z
	TimeFormat TimeFormat
	// TimeDurationUnit defines how time.Duration fields added via
	// Event.WithDuration is rendered. A value of time.Second will then show
	// the duration in whole seconds, whereas time.Minute will
	// show the duration in full minutes.
	//
	// Defaults to 0, which will show the time in nanoseconds.
	//
	// If not using floats for duration, then the values will be rounded down
	// to the nearest unit. Where adding a 30 seconds duration to the log Event
	// with a unit of time.Minute will result in the integer value 0.
	TimeDurationUnit time.Duration
	// TimeDurationUseFloat defines how time.Duration fields added via
	// Event.WithDuration is rendered.
	//
	// Usually combined with also specifying the time unit. Setting this to true
	// and the duration unit to time.Minute will result in showing a 30 seconds
	// duration as 0.5 instead of rounding down to the integer 0.
	//
	// When set to false (which is the default) the duration is formatted as an
	// integer.
	TimeDurationUseFloat bool
}

// Default is a logger Sink that outputs JSON-formatted logs to the console
// using its default settings.
var Default = New(Config{})

// New creates a new JSON-console logging Sink.
func New(conf Config) logger.Sink {
	conf.CallerFileField = prepareFieldName(conf.CallerFileField, "caller")
	conf.CallerLineField = prepareFieldName(conf.CallerLineField, "line")
	conf.ErrorField = prepareFieldName(conf.ErrorField, "error")
	conf.LevelField = prepareFieldName(conf.LevelField, "level")
	conf.MessageField = prepareFieldName(conf.MessageField, "message")
	conf.ScopeField = prepareFieldName(conf.ScopeField, "scope")
	conf.DateField = prepareFieldName(conf.DateField, "date")
	return sink{&conf}
}

type sink struct {
	config *Config
}

// NewContext creates a new JSON-console logging Context using the
// same configuration as the one given when creating the Sink.
func (s sink) NewContext(scope string) logger.Context {
	return context{
		Config: s.config,
		scope:  scope,
	}
}

type context struct {
	*Config
	fields     []byte
	caller     string
	callerLine int
	scope      string
	error      error
}

func (c context) WriteOut(level logger.Level, message string) {
	var buf []byte
	buf = append(buf, `{"`...)
	buf = append(buf, c.LevelField...)
	buf = append(buf, `":"`...)
	buf = append(buf, levelString(level)...)
	buf = append(buf, '"')

	if !c.DisableDate {
		buf = appendFieldNameRaw(buf, c.DateField)
		buf = appendTime(buf, time.Now(), c.TimeFormat)
	}

	if !c.DisableCaller {
		buf = appendFieldNameRaw(buf, c.CallerFileField)
		buf = appendEscapedString(buf, c.caller)
		if !c.DisableCallerLine {
			buf = appendFieldNameRaw(buf, c.CallerLineField)
			buf = strconv.AppendInt(buf, int64(c.callerLine), 10)
		}
	}

	if c.scope != "" {
		buf = appendFieldNameRaw(buf, c.ScopeField)
		buf = appendEscapedString(buf, c.scope)
	}

	if message != "" {
		buf = appendFieldNameRaw(buf, c.MessageField)
		buf = appendEscapedString(buf, message)
	}

	if c.error != nil {
		buf = appendFieldNameRaw(buf, c.ErrorField)
		buf = appendEscapedString(buf, c.error.Error())
	}

	buf = append(buf, c.fields...)
	buf = append(buf, "}\n"...)

	os.Stdout.Write(buf)
}

func (c context) SetCaller(file string, line int) logger.Context {
	c.caller, c.callerLine = file, line
	return c
}

func (c context) SetError(value error) logger.Context {
	c.error = value
	return c
}

func (c context) AppendString(key string, value string) logger.Context {
	c.fields = appendFieldName(c.fields, key)
	c.fields = appendEscapedString(c.fields, value)
	return c
}

func (c context) AppendRune(key string, value rune) logger.Context {
	return c.AppendString(key, string(value))
}

func (c context) AppendBool(key string, value bool) logger.Context {
	c.fields = appendFieldName(c.fields, key)
	c.fields = strconv.AppendBool(c.fields, value)
	return c
}

func (c context) AppendInt(k string, v int) logger.Context {
	c.fields = appendInt64(c.fields, k, int64(v))
	return c
}
func (c context) AppendInt32(k string, v int32) logger.Context {
	c.fields = appendInt64(c.fields, k, int64(v))
	return c
}
func (c context) AppendInt64(k string, v int64) logger.Context {
	c.fields = appendInt64(c.fields, k, v)
	return c
}

func (c context) AppendUint(k string, v uint) logger.Context {
	c.fields = appendUint64(c.fields, k, uint64(v))
	return c
}
func (c context) AppendUint32(k string, v uint32) logger.Context {
	c.fields = appendUint64(c.fields, k, uint64(v))
	return c
}
func (c context) AppendUint64(k string, v uint64) logger.Context {
	c.fields = appendUint64(c.fields, k, v)
	return c
}

func (c context) AppendFloat32(key string, value float32) logger.Context {
	c.fields = appendFloat(c.fields, key, float64(value), 32)
	return c
}

func (c context) AppendFloat64(key string, value float64) logger.Context {
	c.fields = appendFloat(c.fields, key, value, 64)
	return c
}

func (c context) AppendTime(key string, value time.Time) logger.Context {
	c.fields = appendFieldName(c.fields, key)
	c.fields = appendTime(c.fields, value, c.TimeFormat)
	return c
}

func (c context) AppendDuration(key string, value time.Duration) logger.Context {
	switch {
	case c.TimeDurationUseFloat:
		valueFloat := float64(value)
		if c.TimeDurationUnit > 0 {
			valueFloat /= float64(c.TimeDurationUnit)
		}
		c.fields = appendFloat(c.fields, key, valueFloat, 64)
	default:
		valueInt := int64(value)
		if c.TimeDurationUnit > 0 {
			valueInt /= int64(c.TimeDurationUnit)
		}
		c.fields = appendInt64(c.fields, key, valueInt)
	}
	return c
}

func appendTime(b []byte, value time.Time, format TimeFormat) []byte {
	switch format {
	case TimeUnix:
		b = strconv.AppendInt(b, value.Unix(), 10)
	case TimeUnixMs:
		const nanoToMilliDivisor = 1000000
		b = strconv.AppendInt(b, value.UnixNano()/nanoToMilliDivisor, 10)
	case TimeUnixMicro:
		const nanoToMicroDivisor = 1000
		b = strconv.AppendInt(b, value.UnixNano()/nanoToMicroDivisor, 10)
	case TimeUnixNano:
		b = strconv.AppendInt(b, value.UnixNano(), 10)
	default:
		b = append(b, '"')
		b = value.AppendFormat(b, string(format))
		b = append(b, '"')
	}
	return b
}

func appendFloat(b []byte, key string, value float64, bitSize int) []byte {
	const (
		floatFormat    byte = 'f'
		floatPrecision int  = -1
	)
	b = appendFieldName(b, key)
	switch {
	case math.IsNaN(value):
		b = append(b, `"NaN"`...)
	case math.IsInf(value, 1):
		b = append(b, `"+Inf"`...)
	case math.IsInf(value, -1):
		b = append(b, `"-Inf"`...)
	default:
		b = strconv.AppendFloat(b, value, floatFormat, floatPrecision, bitSize)
	}
	return b
}

func appendUint64(b []byte, key string, value uint64) []byte {
	b = appendFieldName(b, key)
	b = strconv.AppendUint(b, value, 10)
	return b
}

func appendInt64(b []byte, key string, value int64) []byte {
	b = appendFieldName(b, key)
	b = strconv.AppendInt(b, value, 10)
	return b
}

func appendFieldName(b []byte, key string) []byte {
	b = append(b, ',')
	b = appendEscapedString(b, key)
	b = append(b, ':')
	return b
}

func appendFieldNameRaw(b []byte, preEscapedKey string) []byte {
	b = append(b, ',', '"')
	b = append(b, preEscapedKey...)
	b = append(b, '"', ':')
	return b
}

func appendEscapedString(b []byte, value string) []byte {
	if out, err := json.Marshal(value); err == nil {
		return append(b, out...)
	}
	return append(b, '"', '"')
}

func levelString(level logger.Level) string {
	switch level {
	case logger.LevelDebug:
		return "debug"
	case logger.LevelInfo:
		return "info"
	case logger.LevelWarn:
		return "warn"
	case logger.LevelError:
		return "error"
	case logger.LevelPanic:
		return "panic"
	default:
		return "unknown"
	}
}

func prepareFieldName(field, fallback string) string {
	if field == "" {
		return inefficientlyEscapeJSON(fallback)
	}
	return inefficientlyEscapeJSON(field)
}

func inefficientlyEscapeJSON(value string) string {
	b, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(b[1 : len(b)-1])
}
