package gormutil

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/iver-wharf/wharf-core/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	logInfoFunc = func(ctx context.Context, logger gormlogger.Interface, message string) {
		logger.Info(ctx, message)
	}
	logWarnFunc = func(ctx context.Context, logger gormlogger.Interface, message string) {
		logger.Warn(ctx, message)
	}
	logErrorFunc = func(ctx context.Context, logger gormlogger.Interface, message string) {
		logger.Error(ctx, message)
	}
)

func TestLoggerSilencedByWrongGORMLogLevel(t *testing.T) {
	testCases := []struct {
		name      string
		logLevels []gormlogger.LogLevel
		logFunc   func(context.Context, gormlogger.Interface, string)
	}{
		{
			name:      "info",
			logLevels: []gormlogger.LogLevel{gormlogger.Silent, gormlogger.Warn, gormlogger.Error},
			logFunc:   logInfoFunc,
		},
		{
			name:      "warn",
			logLevels: []gormlogger.LogLevel{gormlogger.Silent, gormlogger.Error},
			logFunc:   logWarnFunc,
		},
		{
			name:      "error",
			logLevels: []gormlogger.LogLevel{gormlogger.Silent},
			logFunc:   logErrorFunc,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logMock := logger.NewMock()
			for _, logLevel := range tc.logLevels {
				t.Run(logLevelStr(logLevel), func(t *testing.T) {
					log := NewLogger(LoggerConfig{
						Logger:              logMock,
						AlsoUseGORMLogLevel: true,
					}).LogMode(logLevel)

					tc.logFunc(context.Background(), log, "some message")

					assert.Empty(t, logMock.Logs)
				})
			}
		})
	}
}

func TestLoggerUsingGORMLogLevel(t *testing.T) {
	testCases := []struct {
		name         string
		logLevels    []gormlogger.LogLevel
		wantLogLevel logger.Level
		logFunc      func(context.Context, gormlogger.Interface, string)
	}{
		{
			name:         "info",
			logLevels:    []gormlogger.LogLevel{gormlogger.Info},
			wantLogLevel: logger.LevelInfo,
			logFunc:      logInfoFunc,
		},
		{
			name:         "warn",
			logLevels:    []gormlogger.LogLevel{gormlogger.Warn, gormlogger.Info},
			wantLogLevel: logger.LevelWarn,
			logFunc:      logWarnFunc,
		},
		{
			name:         "error",
			logLevels:    []gormlogger.LogLevel{gormlogger.Error, gormlogger.Warn, gormlogger.Info},
			wantLogLevel: logger.LevelError,
			logFunc:      logErrorFunc,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logMock := logger.NewMock()
			for _, logLevel := range tc.logLevels {
				t.Run(logLevelStr(logLevel), func(t *testing.T) {
					log := NewLogger(LoggerConfig{
						Logger:              logMock,
						AlsoUseGORMLogLevel: true,
					}).LogMode(logLevel)

					tc.logFunc(context.Background(), log, "some message")

					require.NotEmpty(t, logMock.Logs)

					assert.Equal(t, tc.wantLogLevel, logMock.Logs[0].Level, "logged level")
					assert.Equal(t, "some message", logMock.Logs[0].Message)
				})
			}
		})
	}
}

func TestLoggerTraceLogsFields(t *testing.T) {
	var (
		baseFields = []string{"caller", "line", "rows", "sql", "elapsed"}

		testCases = []struct {
			name           string
			logLevel       gormlogger.LogLevel
			begin          time.Time
			err            error
			wantFieldNames []string
			wantLogLevel   logger.Level
		}{
			{
				name:           "non-'record not found' error",
				logLevel:       gormlogger.Error,
				begin:          time.Now(),
				err:            errors.New("this is not a 'record not found' error"),
				wantFieldNames: append(baseFields, "error"),
				wantLogLevel:   logger.LevelError,
			},
			{
				name:           "slow SQL warn",
				logLevel:       gormlogger.Warn,
				begin:          time.Now().Add(-time.Minute),
				err:            nil,
				wantFieldNames: append(baseFields, "threshold"),
				wantLogLevel:   logger.LevelWarn,
			},
			{
				name:           "SQL debug",
				logLevel:       gormlogger.Info,
				begin:          time.Now(),
				err:            nil,
				wantFieldNames: baseFields,
				wantLogLevel:   logger.LevelDebug,
			},
		}

		fakeSQL            = "SELECT * FROM null"
		affectedRows int64 = 42

		fc = func() (string, int64) {
			return fakeSQL, affectedRows
		}
	)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logMock := logger.NewMock()
			log := NewLogger(LoggerConfig{
				Logger:              logMock,
				AlsoUseGORMLogLevel: true,
				SlowThreshold:       time.Millisecond * 200,
			}).LogMode(tc.logLevel)

			log.Trace(context.Background(), tc.begin, fc, tc.err)

			require.NotEmpty(t, logMock.Logs)
			assert.Equal(t, tc.wantLogLevel, logMock.Logs[0].Level, "logged level")
			assert.ElementsMatch(t, tc.wantFieldNames, logMock.Logs[0].FieldNames, "logged field names")
			assert.Equal(t, fakeSQL, logMock.Logs[0].Fields["sql"], "logged 'sql' field")
			assert.Equal(t, affectedRows, logMock.Logs[0].Fields["rows"], "logged 'rows' field")
		})
	}
}

func TestLoggerTraceNoLogsWhenSilenced(t *testing.T) {
	var (
		logLevel = gormlogger.Silent
		err      = errors.New("should not be logged")
		fc       = func() (string, int64) {
			return "SELECT * FROM null", 42
		}
		logMock = logger.NewMock()
		log     = NewLogger(LoggerConfig{
			Logger:              logMock,
			AlsoUseGORMLogLevel: true,
		}).LogMode(logLevel)
	)
	log.Trace(context.Background(), time.Now(), fc, err)
	assert.Empty(t, logMock.Logs)
}

func TestLoggerOutput(t *testing.T) {
	var (
		wantFieldNames = []string{"caller", "line", "rows", "elapsed", "sql"}
		log            = logger.NewMock()
		db, err        = gorm.Open(postgres.Open("host=localhost"), &gorm.Config{
			DryRun:               true,
			DisableAutomaticPing: true,
			Logger: NewLogger(LoggerConfig{
				Logger: log,
			}),
		})
	)

	require.Nil(t, err)

	type User struct {
		gorm.Model
		Name string `gorm:"size:256"`
	}

	db.Find(&User{}, 1)

	require.NotEmpty(t, log.Logs, "logged messages")
	assert.Equal(t, 1, len(log.Logs), "logged message count")
	assert.ElementsMatch(t, wantFieldNames, log.Logs[0].FieldNames)
}

func logLevelStr(lvl gormlogger.LogLevel) string {
	switch lvl {
	case gormlogger.Silent:
		return "Silent"
	case gormlogger.Info:
		return "Info"
	case gormlogger.Warn:
		return "Warn"
	case gormlogger.Error:
		return "Error"
	default:
		return fmt.Sprintf("LogLevel(%d)", int(lvl))
	}
}
