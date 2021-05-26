package ginutil

import (
	"errors"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iver-wharf/wharf-core/pkg/logger"
)

// LoggerConfig holds configuration for the Gin logging integration.
type LoggerConfig struct {
	// Level is the logging level that each log message uses. Defaults to the
	// zero value of logger.Level, which is logger.LevelDebug.
	Level logger.Level
	// Logger is the logger implementation used when logging.
	Logger logger.Logger
	// OmitClientIP leaves out the client IP address that issued the web request
	// from the logs when set to true.
	OmitClientIP bool
	// OmitLatency leaves out the server cost in time for processing a request
	// from the logs when set to true.
	OmitLatency bool
	// OmitMethod leaves out the HTTP method (GET, POST, DELETE, HEAD, etc.)
	// that was used in the web request from the logs when set to true.
	OmitMethod bool
	// OmitPath leaves out the web request path (the URL without the protocol,
	// hostname, query parameters, and such) from the logs when set to true.
	OmitPath bool
	// OmitStatus leaves out the HTTP status (200 OK, 404 Not Found, etc.) of
	// the web response from the logs when set to true.
	OmitStatus bool
	// OmitError leaves out any Go errors that were thrown when processing the
	// web request from the logs when set to true.
	OmitError bool
	// SkipPaths is a url path array which logs are not written. Useful for
	// disabling logs issued by health checks.
	SkipPaths []string
}

// DefaultLoggerHandler is a Gin-compatible logger that uses wharf-core logging.
var DefaultLoggerHandler = LoggerWithConfig(LoggerConfig{
	Level: logger.LevelDebug,
})

// LoggerWithConfig creates a Gin middleware handler function that logs all
// requests using the logging level from the config.
//
// If the logger implementation inside the config is unset, then it defaults to
// a new scoped logger with the scope "GIN".
func LoggerWithConfig(config LoggerConfig) gin.HandlerFunc {
	if config.Logger == nil {
		config.Logger = logger.NewScoped("GIN")
	}
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: config.SkipPaths,
		Formatter: func(param gin.LogFormatterParams) string {
			ev := logger.NewEventFromLogger(config.Logger, config.Level)
			if !config.OmitClientIP {
				ev = ev.WithString("clientIp", param.ClientIP)
			}
			if !config.OmitMethod {
				ev = ev.WithString("method", param.Method)
			}
			if !config.OmitPath {
				ev = ev.WithString("path", param.Path)
			}
			if !config.OmitStatus {
				ev = ev.WithInt("status", param.StatusCode)
			}
			if !config.OmitLatency {
				ev = ev.WithDuration("latency", param.Latency)
			}
			if param.ErrorMessage != "" && !config.OmitError {
				ev = ev.WithError(errors.New(param.ErrorMessage))
			}
			ev.Message("")
			return ""
		},
		// if writer is not set then it defaults to os.Stdout
		Output: nopWriter{},
	})
}

type nopWriter struct{}

func (w nopWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

// DefaultLoggerWriter is an io.Writer that logs all written messages using
// appropriate logging levels on a logger with the scope "GIN-debug".
//
// Any [GIN-debug] messages are trimmed away.
//
// Any messages starting with [WARNING] or [ERROR] are logged with the
// appropriate logging levels, and any other logs will use debug logging.
var DefaultLoggerWriter = NewLoggerWriter(logger.NewScoped("GIN-debug"), logger.LevelDebug)

type loggerWriter struct {
	logger       logger.Logger
	defaultLevel logger.Level
}

// NewLoggerWriter creates a logger that channels everything written to it via a
// wharf-core logger.
//
// Any [GIN-debug] messages are trimmed away.
//
// Any messages starting with [WARNING] or [ERROR] are logged with the
// appropriate logging levels, and any other logs will use the default logging
// level provided to this function.
func NewLoggerWriter(log logger.Logger, defaultLevel logger.Level) io.Writer {
	return loggerWriter{log, defaultLevel}
}

func (w loggerWriter) Write(p []byte) (n int, err error) {
	const (
		prefixGinDebug = "[GIN-debug] "
		prefixWarning  = "[WARNING] "
		prefixError    = "[ERROR] "
	)

	var message = strings.TrimPrefix(strings.TrimRight(string(p), "\n"), prefixGinDebug)
	var level = w.defaultLevel

	if strings.HasPrefix(message, prefixWarning) {
		message = strings.TrimPrefix(message, prefixWarning)
		level = logger.LevelWarn
	} else if strings.HasPrefix(message, prefixError) {
		message = strings.TrimPrefix(message, prefixError)
		level = logger.LevelError
	}

	logger.NewEventFromLogger(w.logger, level).Message(message)
	return len(p), nil
}
