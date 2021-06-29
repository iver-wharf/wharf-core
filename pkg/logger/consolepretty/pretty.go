package consolepretty

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/iver-wharf/wharf-core/pkg/logger"
	"github.com/mattn/go-colorable"
)

// ColorConfig lets you gradually configure the coloring of the logger.
type ColorConfig struct {
	// Date sets the color attributes for the timestamp of the logs.
	Date *color.Color
	// Scope sets the color attributes for the scope value of the logs.
	Scope *color.Color
	// CallerFile sets the color attributes for the caller file path of the logs.
	CallerFile *color.Color
	// CallerDelimiter sets the color attributes for the delimiter between the
	// caller file path and the caller line number of the logs.
	CallerDelimiter *color.Color
	// CallerLine sets the color attributes for the caller line number of the logs.
	CallerLine *color.Color
	// PreMessageDelimiter sets the color attributes for the delimiters between
	// the date timestamp, logging level, scope, and caller of the logs.
	PreMessageDelimiter *color.Color
	// MessageDebug sets the color attributes for the message on debug logs.
	MessageDebug *color.Color
	// MessageInfo sets the color attributes for the message on info logs.
	MessageInfo *color.Color
	// MessageWarn sets the color attributes for the message on warning logs.
	MessageWarn *color.Color
	// MessageError sets the color attributes for the message on error logs.
	MessageError *color.Color
	// MessagePanic sets the color attributes for the message on panic logs.
	MessagePanic *color.Color
	// LevelDebug sets the color attributes for the log level on debug logs.
	LevelDebug *color.Color
	// LevelInfo sets the color attributes for the log level on info logs.
	LevelInfo *color.Color
	// LevelWarn sets the color attributes for the log level on warning logs.
	LevelWarn *color.Color
	// LevelError sets the color attributes for the log level on error logs.
	LevelError *color.Color
	// LevelPanic sets the color attributes for the log level on panic logs.
	LevelPanic *color.Color
	// FieldKey sets the color attributes for the string key of each field added
	// via the Event.With* methods for the logs.
	FieldKey *color.Color
	// FieldDelimiter sets the color attributes for the delimiter between the
	// string key and the formatted value of each field added via the
	// Event.With* methods for the logs.
	FieldDelimiter *color.Color
	// FieldValue sets the color attributes for the formatted value of each
	// field added via the Event.With* methods for the logs for any non-zero
	// values.
	//
	// A zero-value here is more narrow than Go's definition. Here a zero-value
	// only refers to nil and empty strings.
	FieldValue *color.Color
	// FieldValueZero sets the color attributes for the formatted value of each
	// field added via the Event.With* methods for the logs for any zero
	// values.
	//
	// A zero-value here is more narrow than Go's definition. Here a zero-value
	// only refers to nil and empty strings.
	FieldValueZero *color.Color
	// ErrorKey sets the color attributes for the string key of the error added
	// via Event.WithError method for the logs.
	ErrorKey *color.Color
	// ErrorDelimiter sets the color attributes for the delimiter between the
	// string key and the formatted error string of the error added via
	// Event.WithError method for the logs.
	ErrorDelimiter *color.Color
	// ErrorValue sets the color attributes for the error string of the error
	// added via Event.WithError method for the logs.
	ErrorValue *color.Color
	// ErrorType sets the color attributes for the error type of the error
	// added via Event.WithError method for the logs.
	ErrorType *color.Color
}

// DefaultColorConfig is the config used in New to populate some values if left
// unset. Changing this global value also changes the fallback values used in
// New.
var DefaultColorConfig = ColorConfig{
	Date:                color.New(color.FgHiBlack),
	Scope:               color.New(color.FgCyan, color.Bold),
	CallerFile:          color.New(color.FgHiBlack),
	CallerDelimiter:     color.New(color.FgHiBlack),
	CallerLine:          color.New(color.FgHiBlack),
	PreMessageDelimiter: color.New(color.FgWhite),
	MessageDebug:        color.New(color.FgHiBlack, color.Italic),
	MessageInfo:         color.New(color.FgHiWhite),
	MessageWarn:         color.New(color.FgHiYellow),
	MessageError:        color.New(color.FgRed),
	MessagePanic:        color.New(color.FgHiRed, color.Bold),
	LevelDebug:          color.New(color.FgHiBlack, color.Italic),
	LevelInfo:           color.New(color.FgGreen),
	LevelWarn:           color.New(color.FgYellow),
	LevelError:          color.New(color.FgRed, color.Bold),
	LevelPanic:          color.New(color.FgHiWhite, color.BgRed, color.Bold),
	FieldKey:            color.New(color.FgHiBlack, color.Italic),
	FieldDelimiter:      color.New(color.FgHiBlack, color.Italic),
	FieldValue:          color.New(color.FgHiYellow),
	FieldValueZero:      color.New(color.FgYellow, color.Italic),
	ErrorKey:            color.New(color.FgRed, color.Italic, color.Bold),
	ErrorDelimiter:      color.New(color.FgRed, color.Italic),
	ErrorValue:          color.New(color.FgHiRed),
	ErrorType:           color.New(color.FgRed, color.Italic),
}

// Config lets you gradually configure the output of the logger by disabling
// certain features or changing the format of certain field types.
type Config struct {
	// Writer is the io.Writer target that the pretty-console logger will write
	// to. Defaults to using a github.com/mattn/go-colorable wrapper around
	// os.Stdout.
	Writer io.Writer

	// Coloring defines how certain parts of the logs are colored.
	Coloring *ColorConfig

	// DateFormat is the format to display the timestamp of when a logged
	// message was logged. This does not alter how Event.WithTime is rendered.
	DateFormat string

	// Prefix sets an optional string added to the beginning of the log message.
	//
	// When set to "" (empty string):
	// 	[INFO | 2006-01-02T15:04:05Z | example.go:20] Sample message.
	// When set to "foo:":
	// 	foo:[INFO | 2006-01-02T15:04:05Z | example.go:20] Sample message.
	Prefix string

	// DisableDate removes the date field from the log when set to true.
	//
	// When set to false:
	// 	[INFO | 2006-01-02T15:04:05Z | example.go:20] Sample message.
	// When set to true:
	// 	[INFO | example.go:20] Sample message.
	DisableDate bool

	// DisableCaller removee the caller file name and line fields from the log
	// when set to true.
	//
	// When set to false:
	// 	[INFO | 2006-01-02T15:04:05Z | example.go:20] Sample message.
	// With set to true:
	// 	[INFO | 2006-01-02T15:04:05Z] Sample message.
	DisableCaller bool

	// DisableCallerLine removee the just the caller line field from the log
	// when set to true, but leaves the caller file name as-is.
	//
	// When set to false:
	// 	[INFO | 2006-01-02T15:04:05Z | example.go:20] Sample message.
	// With set to true:
	// 	[INFO | 2006-01-02T15:04:05Z | example.go] Sample message.
	DisableCallerLine bool
}

// DefaultConfig is the config used in New to populate some values if left
// unset. Changing this global value also changes the fallback values used in
// New.
var DefaultConfig = Config{
	DateFormat: "Jan-02 15:04Z0700",
}

// Default is a logger Sink that outputs human-readable logs to the console
// using its default settings.
var Default = New(DefaultConfig)

// New creates a new pretty-console logging Sink and uses fallback values from
// DefaultConfig and DefaultColorConfig for certain configs. Namely:
//
// 	Config.Writer = DefaultConfig.Writer
// 	Config.DateFormat = DefaultConfig.DateFormat
//
// 	Config.Coloring = DefaultColorConfig
func New(conf Config) logger.Sink {
	if conf.Writer == nil {
		if DefaultConfig.Writer == nil {
			conf.Writer = colorable.NewColorableStdout()
		} else {
			conf.Writer = DefaultConfig.Writer
		}
	}
	if conf.Coloring == nil {
		conf.Coloring = &DefaultColorConfig
	}
	if conf.DateFormat == "" {
		conf.DateFormat = DefaultConfig.DateFormat
	}
	return sink{&conf}
}

type sink struct {
	config *Config
}

// NewContext creates a new pretty-console logging Context using the
// same configuration as the one given when creating the Sink.
func (s sink) NewContext() logger.Context {
	return context{
		Config: s.config,
	}
}

type context struct {
	*Config
	fields     []fieldPair
	scope      string
	callerFile string
	callerLine int
	err        error
}

type fieldPair struct {
	key   string
	value interface{}
}

func (c context) WriteOut(level logger.Level, message string) {
	var buf bytes.Buffer
	var coloring = c.Coloring
	if c.Prefix != "" {
		buf.WriteString(c.Prefix)
	}
	if !c.DisableDate {
		coloring.Date.Fprint(&buf, time.Now().Format(c.DateFormat))
		buf.WriteRune(' ')
	}
	coloring.PreMessageDelimiter.Fprint(&buf, "[")
	c.writeLevel(&buf, level)
	if c.scope != "" {
		coloring.PreMessageDelimiter.Fprint(&buf, " | ")
		coloring.Scope.Fprint(&buf, c.scope)
	}
	if c.callerFile != "" && !c.DisableCaller {
		coloring.PreMessageDelimiter.Fprint(&buf, " | ")
		coloring.CallerFile.Fprint(&buf, c.callerFile)
		if !c.DisableCallerLine {
			coloring.CallerDelimiter.Fprint(&buf, ":")
			coloring.CallerLine.Fprint(&buf, strconv.FormatInt(int64(c.callerLine), 10))
		}
	}
	coloring.PreMessageDelimiter.Fprint(&buf, "]")
	buf.WriteRune(' ')
	needsSeparator := false
	if message != "" {
		c.writeMessage(&buf, level, message)
		needsSeparator = true
	}
	for _, pair := range c.fields {
		if needsSeparator {
			buf.WriteString("  ")
		}
		coloring.FieldKey.Fprint(&buf, pair.key)
		coloring.FieldDelimiter.Fprint(&buf, "=")
		str, hasValue := getPrintableStringRepresentation(pair.value)
		if hasValue {
			coloring.FieldValue.Fprint(&buf, str)
		} else {
			coloring.FieldValueZero.Fprint(&buf, str)
		}
		needsSeparator = true
	}
	if c.err != nil {
		if needsSeparator {
			buf.WriteString("  ")
		}
		coloring.ErrorKey.Fprint(&buf, "error")
		coloring.ErrorDelimiter.Fprint(&buf, "=")
		coloring.ErrorValue.Fprint(&buf, c.err.Error())
		buf.WriteRune(' ')
		coloring.ErrorType.Fprintf(&buf, "(%T)", c.err)
	}
	buf.WriteRune('\n')
	io.Copy(c.Writer, &buf)
}

func getPrintableStringRepresentation(value interface{}) (str string, hasValue bool) {
	if value == nil {
		return "<nil>", false
	}
	switch v := value.(type) {
	case string:
		if v == "" {
			return "``", false
		}
		return escapeString(v), true
	default:
		return fmt.Sprint(value), true
	}
}

var escapeStringReplacer = strings.NewReplacer(
	"\a", `\a`,
	"\b", `\b`,
	"\f", `\f`,
	"\n", `\n`,
	"\r", `\r`,
	"\t", `\t`,
	"\v", `\v`,
	`\`, `\\`,
	"`", "\\`",
)

func escapeString(value string) string {
	if strings.ContainsAny(value, " \a\b\f\n\r\t\v\"\\`") {
		return fmt.Sprintf("`%s`", escapeStringReplacer.Replace(value))
	}
	return value
}

func (c context) SetScope(value string) logger.Context {
	c.scope = value
	return c
}

func (c context) SetCaller(file string, line int) logger.Context {
	c.callerFile = file
	c.callerLine = line
	return c
}

func (c context) SetError(value error) logger.Context {
	c.err = value
	return c
}

func (c context) AppendString(k string, v string) logger.Context          { return c.addField(k, v) }
func (c context) AppendRune(k string, v rune) logger.Context              { return c.addField(k, v) }
func (c context) AppendBool(k string, v bool) logger.Context              { return c.addField(k, v) }
func (c context) AppendInt(k string, v int) logger.Context                { return c.addField(k, v) }
func (c context) AppendInt32(k string, v int32) logger.Context            { return c.addField(k, v) }
func (c context) AppendInt64(k string, v int64) logger.Context            { return c.addField(k, v) }
func (c context) AppendUint(k string, v uint) logger.Context              { return c.addField(k, v) }
func (c context) AppendUint32(k string, v uint32) logger.Context          { return c.addField(k, v) }
func (c context) AppendUint64(k string, v uint64) logger.Context          { return c.addField(k, v) }
func (c context) AppendFloat32(k string, v float32) logger.Context        { return c.addField(k, v) }
func (c context) AppendFloat64(k string, v float64) logger.Context        { return c.addField(k, v) }
func (c context) AppendTime(k string, v time.Time) logger.Context         { return c.addField(k, v) }
func (c context) AppendDuration(k string, v time.Duration) logger.Context { return c.addField(k, v) }

func (c context) addField(key string, value interface{}) logger.Context {
	c.fields = append(c.fields, fieldPair{key, value})
	return c
}

func (c context) writeMessage(w io.Writer, level logger.Level, msg string) {
	var color *color.Color
	switch level {
	case logger.LevelDebug:
		color = c.Coloring.MessageDebug
	case logger.LevelInfo:
		color = c.Coloring.MessageInfo
	case logger.LevelWarn:
		color = c.Coloring.MessageWarn
	case logger.LevelError:
		color = c.Coloring.MessageError
	case logger.LevelPanic:
		color = c.Coloring.MessagePanic
	default:
		color = c.Coloring.MessageDebug
	}
	msg = strings.ReplaceAll(msg, "\n", "\n\t")
	color.Fprint(w, msg)
}

func (c context) writeLevel(w io.Writer, level logger.Level) {
	switch level {
	case logger.LevelDebug:
		c.Coloring.LevelDebug.Fprint(w, "DEBUG")
	case logger.LevelInfo:
		c.Coloring.LevelInfo.Fprint(w, "INFO ")
	case logger.LevelWarn:
		c.Coloring.LevelWarn.Fprint(w, "WARN ")
	case logger.LevelError:
		c.Coloring.LevelError.Fprint(w, "ERROR")
	case logger.LevelPanic:
		c.Coloring.LevelPanic.Fprint(w, "PANIC")
	default:
		c.Coloring.LevelDebug.Fprint(w, "???  ")
	}
}
