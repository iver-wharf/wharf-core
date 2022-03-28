package consolepretty

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/iver-wharf/wharf-core/v2/pkg/logger"
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
	FieldValue:          color.New(color.FgWhite),
	FieldValueZero:      color.New(color.FgHiBlack, color.Italic),
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
	// 	Jan 02 15:04Z [INFO |example.go:20] Sample message.
	// When set to "foo:":
	// 	foo:Jan 02 15:04Z [INFO |example.go:20] Sample message.
	Prefix string

	// DisableDate removes the date field from the log when set to true.
	//
	// When set to false:
	// 	Jan 02 15:04Z [INFO |example.go:20] Sample message.
	// When set to true:
	// 	[INFO |example.go:20] Sample message.
	DisableDate bool

	// DisableCaller removes the caller file name and line fields from the log
	// when set to true.
	//
	// When set to false:
	// 	Jan 02 15:04Z [INFO |example.go:20] Sample message.
	// With set to true:
	// 	Jan 02 15:04Z [INFO ] Sample message.
	DisableCaller bool

	// DisableCallerLine removes just the caller line field from the log
	// when set to true, but leaves the caller file name as-is.
	//
	// When set to false:
	// 	Jan 02 15:04Z [INFO |example.go:20] Sample message.
	// With set to true:
	// 	Jan 02 15:04Z [INFO |example.go] Sample message.
	DisableCallerLine bool

	// Ellipsis defines the string used when trimming the values, as an effect
	// of the caller or scope max length configs.
	//
	// Setting this to a value longer than the max length is considered
	// undefined behavior, and should be avoided.
	Ellipsis string

	// CallerMaxLength will trim the caller file and line down to this length
	// if set to a value of 1 or higher.
	//
	// When set to 0:
	// 	Jan 02 15:04Z [INFO |example.go:20] Sample message.
	// With set to 10:
	// 	Jan 02 15:04Z [INFO |…ple.go:20] Sample message.
	CallerMaxLength int

	// CallerMinLength will pad the caller file and line with spaces so that it
	// reaches the target character width.
	//
	// When set to 0:
	// 	Jan 02 15:04Z [INFO |example.go:20] Sample message.
	// 	Jan 02 15:04Z [INFO |test.go:20] Sample message.
	// With set to 13:
	// 	Jan 02 15:04Z [INFO |example.go:20] Sample message.
	// 	Jan 02 15:04Z [INFO |test.go:20   ] Sample message.
	CallerMinLength int

	// ScopeMaxLength will trim the scope down to this length if set to a value
	// of 1 or higher.
	//
	// When set to 0:
	// 	Jan 02 15:04Z [INFO |GORM-debug] Sample message.
	// With set to 5:
	// 	Jan 02 15:04Z [INFO |GORM…] Sample message.
	ScopeMaxLength int

	// ScopeMinLength will pad the scope with spaces so that it reaches the
	// target character width.
	//
	// When set to 0:
	// 	Jan 02 15:04Z [INFO |GORM] Sample message.
	// 	Jan 02 15:04Z [INFO |GORM-debug] Sample message.
	// With set to 12:
	// 	Jan 02 15:04Z [INFO |GORM        ] Sample message.
	// 	Jan 02 15:04Z [INFO |GORM-debug  ] Sample message.
	ScopeMinLength int

	// ScopeMinLengthAuto will automatically pad the scope with spaces to
	// accommodate for the longest scope created by logger.NewScoped.
	//
	// When set to false:
	// 	Jan 02 15:04Z [INFO |GORM] Sample message.
	// 	Jan 02 15:04Z [INFO |GORM-debug] Sample message.
	// With set to true:
	// 	Jan 02 15:04Z [INFO |GORM      ] Sample message.
	// 	Jan 02 15:04Z [INFO |GORM-debug] Sample message.
	ScopeMinLengthAuto bool
}

// DefaultConfig is the config used in New to populate some values if left
// unset. Changing this global value also changes the fallback values used in
// New.
var DefaultConfig = Config{
	Ellipsis:           "…",
	DateFormat:         "Jan-02 15:04Z0700",
	CallerMaxLength:    23,
	CallerMinLength:    23,
	ScopeMinLengthAuto: true,
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
	if conf.Ellipsis == "" {
		conf.Ellipsis = DefaultConfig.Ellipsis
	}
	return sink{
		config:      &conf,
		ellipsisLen: utf8.RuneCountInString(conf.Ellipsis),
	}
}

type sink struct {
	config      *Config
	ellipsisLen int
}

// NewContext creates a new pretty-console logging Context using the
// same configuration as the one given when creating the Sink.
func (s sink) NewContext(scope string) logger.Context {
	return context{
		Config:      s.config,
		scope:       scope,
		ellipsisLen: s.ellipsisLen,
	}
}

type context struct {
	*Config
	fields      []fieldPair
	scope       string
	callerFile  string
	callerLine  int
	err         error
	ellipsisLen int
}

type fieldPair struct {
	key   string
	value any
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
	if c.scope != "" || c.Config.ScopeMinLength > 0 || c.Config.ScopeMinLengthAuto {
		coloring.PreMessageDelimiter.Fprint(&buf, "|")
		scopeWrittenWidth := len(c.scope)
		if c.Config.ScopeMaxLength > 0 {
			scopeWrittenWidth = c.writeTrimmedRight(&buf, coloring.Scope, c.scope, c.Config.ScopeMaxLength)
		} else {
			coloring.Scope.Fprint(&buf, c.scope)
		}
		scopeMinWidth := c.Config.ScopeMinLength
		if c.Config.ScopeMinLengthAuto {
			scopeMinWidth = logger.LongestScopeNameLength
		}
		for i := scopeWrittenWidth; i < scopeMinWidth; i++ {
			buf.WriteRune(' ')
		}
	}
	if c.callerFile != "" && !c.DisableCaller {
		coloring.PreMessageDelimiter.Fprint(&buf, "|")
		writtenWidth := 0
		maxFileWidth := c.Config.CallerMaxLength
		if maxFileWidth > 0 {
			if !c.DisableCallerLine {
				maxFileWidth-- // for the delimiter
				maxFileWidth -= printedIntLenFast(c.callerLine)
			}
			writtenWidth = c.writeTrimmedLeft(&buf, coloring.CallerFile, c.callerFile, maxFileWidth)
		} else {
			coloring.CallerFile.Fprint(&buf, c.callerFile)
			writtenWidth = len(c.callerFile)
		}
		if !c.DisableCallerLine {
			coloring.CallerDelimiter.Fprint(&buf, ":")
			lineStr := strconv.FormatInt(int64(c.callerLine), 10)
			coloring.CallerLine.Fprint(&buf, lineStr)
			writtenWidth += len(lineStr) + 1
		}
		for i := writtenWidth; i < c.Config.CallerMinLength; i++ {
			buf.WriteRune(' ')
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
		str, _ := getPrintableStringRepresentation(strings.TrimSpace(c.err.Error()))
		coloring.ErrorValue.Fprint(&buf, str)
		buf.WriteRune(' ')
		coloring.ErrorType.Fprintf(&buf, "(%T)", c.err)
	}
	buf.WriteRune('\n')
	io.Copy(c.Writer, &buf)
}

func getPrintableStringRepresentation(value any) (str string, hasValue bool) {
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

func (c context) addField(key string, value any) logger.Context {
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

func (c context) writeTrimmedRight(w io.Writer, col *color.Color, value string, maxLen int) int {
	if written, ok := c.writeUntrimmedString(w, col, value, maxLen); ok {
		return written
	}
	sliceLen := maxLen - c.ellipsisLen
	col.Fprint(w, value[:sliceLen], c.Ellipsis)
	return maxLen
}

func (c context) writeTrimmedLeft(w io.Writer, col *color.Color, value string, maxLen int) int {
	if written, ok := c.writeUntrimmedString(w, col, value, maxLen); ok {
		return written
	}
	sliceStartIndex := len(value) - maxLen + c.ellipsisLen
	col.Fprint(w, c.Ellipsis, value[sliceStartIndex:])
	return maxLen
}

func (c context) writeUntrimmedString(w io.Writer, col *color.Color, value string, maxLen int) (int, bool) {
	valueLen := len(value)
	if valueLen > maxLen {
		return 0, false
	}
	switch {
	case valueLen == 0 || maxLen <= 0:
		// do nothing
		return 0, true
	case maxLen <= c.ellipsisLen && valueLen > c.ellipsisLen:
		col.Fprint(w, c.Ellipsis)
		return c.ellipsisLen, true
	default:
		col.Fprint(w, value)
		return valueLen, true
	}
}

func printedIntLenFast(number int) int {
	// could do log10(number), but as the benchmark shows, that's approx 8-10
	// times slower
	switch {
	case number < 0:
		return printedIntLenFast(-number) + 1 // +1 for the sign symbol
	case number < 10:
		return 1
	case number < 100:
		return 2
	case number < 1000:
		return 3
	case number < 10000:
		return 4
	case number < 100000:
		return 5
	case number < 1000000:
		return 6
	case number < 10000000:
		return 7
	case number < 100000000:
		return 8
	case number < 1000000000:
		return 9
	// for our purposes here, handling >int32 max value is not needed
	default:
		return 10
	}
}
