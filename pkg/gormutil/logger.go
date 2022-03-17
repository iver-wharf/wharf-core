package gormutil

import (
	"context"
	"errors"
	"time"

	"github.com/iver-wharf/wharf-core/v2/pkg/logger"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// LoggerConfig holds configuration for the GORM logging integration. Many
// configuration values of in the gorm.io/gorm/logger.Config will you also find
// in here.
type LoggerConfig struct {
	// Logger is the logger implementation used when logging. This defaults to
	// a new scoped logger with the scope "GORM".
	Logger logger.Logger
	// AlsoUseGORMLogLevel sets wether to honor GORM's own logging levels.
	//
	// If set to false (which is the default) then the logging level
	// configuration from the wharf-core logging library will be the only one
	// filtering logs.
	AlsoUseGORMLogLevel bool
	// IgnoreRecordNotFoundError will omit any "RecordNotFound" errors if set
	// to true.
	IgnoreRecordNotFoundError bool
	// SlowThreshold sets what duration is considered a slow SQL operation.
	// If an operation takes longer than this to complete then a warning log
	// message will be emitted.
	//
	// Set to 0 to disable.
	SlowThreshold time.Duration
}

type gormLog struct {
	LoggerConfig
	level gormlogger.LogLevel
}

// DefaultLogger is a GORM-compatible logger that uses wharf-core logging.
var DefaultLogger = NewLogger(LoggerConfig{
	IgnoreRecordNotFoundError: true,
	SlowThreshold:             200 * time.Millisecond,
})

// NewLogger creates a new gorm.io/gorm/logger.Interface compatible logger.
func NewLogger(config LoggerConfig) gormlogger.Interface {
	if config.Logger == nil {
		config.Logger = logger.NewScoped("GORM")
	}
	return gormLog{
		LoggerConfig: config,
		level:        gormlogger.Info,
	}
}

func (log gormLog) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	log.level = level
	return log
}

func (log gormLog) Info(_ context.Context, message string, args ...any) {
	if log.level >= gormlogger.Info || !log.AlsoUseGORMLogLevel {
		log.Logger.Info().Messagef(message, args...)
	}
}

func (log gormLog) Warn(_ context.Context, message string, args ...any) {
	if log.level >= gormlogger.Warn || !log.AlsoUseGORMLogLevel {
		log.Logger.Warn().Messagef(message, args...)
	}
}

func (log gormLog) Error(_ context.Context, message string, args ...any) {
	if log.level >= gormlogger.Error || !log.AlsoUseGORMLogLevel {
		log.Logger.Error().Messagef(message, args...)
	}
}

func (log gormLog) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if log.level <= gormlogger.Silent && log.AlsoUseGORMLogLevel {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case log.shouldLogError(err):
		sql, rowsAffected := fc()
		ev := log.Logger.Error()
		ev = withRowsAffected(ev, rowsAffected)
		ev.WithDuration("elapsed", elapsed).
			WithError(err).
			WithString("sql", sql).
			Message("Error in SQL.")
	case log.shouldLogWarnSlow(elapsed):
		sql, rowsAffected := fc()
		ev := log.Logger.Warn()
		ev = withRowsAffected(ev, rowsAffected)
		ev.WithDuration("elapsed", elapsed).
			WithDuration("threshold", log.SlowThreshold).
			WithString("sql", sql).
			Message("Slow SQL.")
	case log.shouldLogDebug():
		sql, rowsAffected := fc()
		ev := log.Logger.Debug()
		ev = withRowsAffected(ev, rowsAffected)
		ev.WithDuration("elapsed", elapsed).
			WithString("sql", sql).
			Message("")
	}
}

func withRowsAffected(ev logger.Event, rows int64) logger.Event {
	if rows == -1 {
		return ev.WithRune("rows", '-')
	}
	return ev.WithInt64("rows", rows)
}

func (log gormLog) shouldLogError(err error) bool {
	return err != nil &&
		(log.level >= gormlogger.Error || !log.AlsoUseGORMLogLevel) &&
		(!errors.Is(err, gorm.ErrRecordNotFound) ||
			!log.IgnoreRecordNotFoundError)
}

func (log gormLog) shouldLogWarnSlow(elapsed time.Duration) bool {
	return elapsed > log.SlowThreshold &&
		log.SlowThreshold != 0 &&
		(log.level >= gormlogger.Warn || !log.AlsoUseGORMLogLevel)
}

func (log gormLog) shouldLogDebug() bool {
	return log.level == gormlogger.Info || !log.AlsoUseGORMLogLevel
}
