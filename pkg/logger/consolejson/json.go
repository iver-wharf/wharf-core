package consolejson

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/iver-wharf/wharf-core/pkg/logger"
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
	// DisableCaller removee the caller file name and line fields from the log
	// when set to true.
	//
	// When set to false:
	// 	{"level":"info","date":"2006-01-02T15:04:05Z","caller":"example.go","line":20,"message":"Sample message."}
	// When set to true:
	// 	{"level":"info","date":"2006-01-02T15:04:05Z","message":"Sample message."}
	DisableCaller bool
	// DisableCallerLine removee the just the caller line field from the log
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
func (s sink) NewContext() logger.Context {
	return context{
		Config: s.config,
	}
}

type context struct {
	*Config
	fields     []byte
	caller     string
	callerLine int
	scope      string
}

func (c context) WriteOut(level logger.Level, message string) {
	var buf bytes.Buffer
	buf.WriteString(`{"`)
	buf.WriteString(c.LevelField)
	buf.WriteString(`":"`)
	buf.WriteString(levelString(level))
	buf.WriteRune('"')

	if !c.DisableDate {
		buf.WriteString(`,"`)
		buf.WriteString(c.DateField)
		buf.WriteString(`":`)
		writeEscapedString(&buf, time.Now().Format(time.RFC3339))
	}

	if !c.DisableCaller {
		buf.WriteString(`,"`)
		buf.WriteString(c.CallerFileField)
		buf.WriteString(`":`)
		writeEscapedString(&buf, c.caller)
		if !c.DisableCallerLine {
			buf.WriteString(`,"`)
			buf.WriteString(c.CallerLineField)
			buf.WriteString(`":`)
			buf.WriteString(strconv.FormatInt(int64(c.callerLine), 10))
		}
	}

	if c.scope != "" {
		buf.WriteString(`,"`)
		buf.WriteString(c.ScopeField)
		buf.WriteString(`":`)
		writeEscapedString(&buf, c.scope)
	}

	if message != "" {
		buf.WriteString(`,"`)
		buf.WriteString(c.MessageField)
		buf.WriteString(`":`)
		writeEscapedString(&buf, message)
	}

	buf.Write(c.fields)
	buf.WriteString("}\n")

	buf.WriteTo(os.Stdout)
}

func (c context) AppendScope(value string) logger.Context {
	c.scope = value
	return c
}

func (c context) AppendCaller(file string, line int) logger.Context {
	c.caller, c.callerLine = file, line
	return c
}

func (c context) AppendError(value error) logger.Context {
	c = c.appendRawFieldName(c.ErrorField)
	c.fields = appendEscapedString(c.fields, value.Error())
	return c
}

func (c context) AppendString(key string, value string) logger.Context {
	c = c.appendFieldName(key)
	c.fields = appendEscapedString(c.fields, value)
	return c
}

func (c context) AppendRune(key string, value rune) logger.Context {
	return c.AppendString(key, string(value))
}

func (c context) AppendBool(key string, value bool) logger.Context {
	c = c.appendFieldName(key)
	c.fields = strconv.AppendBool(c.fields, value)
	return c
}

func (c context) AppendInt(k string, v int) logger.Context     { return c.appendInt64(k, int64(v)) }
func (c context) AppendInt32(k string, v int32) logger.Context { return c.appendInt64(k, int64(v)) }
func (c context) AppendInt64(k string, v int64) logger.Context { return c.appendInt64(k, v) }

func (c context) AppendUint(k string, v uint) logger.Context     { return c.appendUint64(k, uint64(v)) }
func (c context) AppendUint32(k string, v uint32) logger.Context { return c.appendUint64(k, uint64(v)) }
func (c context) AppendUint64(k string, v uint64) logger.Context { return c.appendUint64(k, v) }

func (c context) AppendFloat32(key string, value float32) logger.Context {
	return c.appendFloat(key, float64(value), 32)
}

func (c context) AppendFloat64(key string, value float64) logger.Context {
	return c.appendFloat(key, value, 64)
}

func (c context) AppendTime(key string, value time.Time) logger.Context {
	switch c.TimeFormat {
	case TimeUnix:
		return c.appendInt64(key, value.Unix())
	case TimeUnixMs:
		const nanoToMilliDivisor = 1000000
		return c.appendInt64(key, value.UnixNano()/nanoToMilliDivisor)
	case TimeUnixMicro:
		const nanoToMicroDivisor = 1000
		return c.appendInt64(key, value.UnixNano()/nanoToMicroDivisor)
	case TimeUnixNano:
		return c.appendInt64(key, value.UnixNano())
	default:
		c = c.appendFieldName(key)
		c.fields = appendEscapedString(c.fields, value.Format(string(c.TimeFormat)))
		return c
	}
}

func (c context) AppendDuration(key string, value time.Duration) logger.Context {
	switch {
	case c.TimeDurationUseFloat:
		valueFloat := float64(value)
		if c.TimeDurationUnit > 0 {
			valueFloat /= float64(c.TimeDurationUnit)
		}
		return c.appendFloat(key, valueFloat, 64)
	default:
		valueInt := int64(value)
		if c.TimeDurationUnit > 0 {
			valueInt /= int64(c.TimeDurationUnit)
		}
		return c.appendInt64(key, valueInt)
	}
}

func (c context) appendFloat(key string, value float64, bitSize int) context {
	const (
		floatFormat    byte = 'f'
		floatPrecision int  = -1
	)
	c = c.appendFieldName(key)
	switch {
	case math.IsNaN(value):
		c.fields = append(c.fields, `"NaN"`...)
	case math.IsInf(value, 1):
		c.fields = append(c.fields, `"+Inf"`...)
	case math.IsInf(value, -1):
		c.fields = append(c.fields, `"-Inf"`...)
	default:
		c.fields = strconv.AppendFloat(c.fields, value, floatFormat, floatPrecision, bitSize)
	}
	return c
}

func (c context) appendUint64(key string, value uint64) context {
	c = c.appendFieldName(key)
	c.fields = strconv.AppendUint(c.fields, value, 10)
	return c
}

func (c context) appendInt64(key string, value int64) context {
	c = c.appendFieldName(key)
	c.fields = strconv.AppendInt(c.fields, value, 10)
	return c
}

func (c context) appendFieldName(key string) context {
	c.fields = append(c.fields, ',')
	c.fields = appendEscapedString(c.fields, key)
	c.fields = append(c.fields, ':')
	return c
}

func (c context) appendRawFieldName(preEscapedKey string) context {
	c.fields = append(c.fields, ',', '"')
	c.fields = append(c.fields, preEscapedKey...)
	c.fields = append(c.fields, '"', ':')
	return c
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

func writeEscapedString(w io.Writer, value string) {
	// using json.NewEncoder here instead will always result in an additional
	// trailing newline, which we don't want.
	if out, err := json.Marshal(value); err == nil {
		w.Write(out)
	} else {
		w.Write([]byte(`""`))
	}
}

func appendEscapedString(b []byte, value string) []byte {
	if out, err := json.Marshal(value); err == nil {
		return append(b, out...)
	}
	return append(b, '"', '"')
}
