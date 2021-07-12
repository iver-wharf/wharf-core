package logger

import (
	"fmt"
	"strings"
)

// Level is the enum type of different logging levels used throughout this
// package to categorize and filter different log events.
type Level byte

const (
	// LevelDebug is the "debugging" logging level, and also the lowest logging
	// level available.
	LevelDebug Level = iota
	// LevelInfo is the "information" logging level
	LevelInfo
	// LevelWarn is the "warning" logging level
	LevelWarn
	// LevelError is the "error" logging level
	LevelError
	// LevelPanic is the "panic" logging level, and also the highest logging
	// level available.
	LevelPanic
)

// String returns a readable representation of the logging level.
func (lvl Level) String() string {
	switch lvl {
	case LevelDebug:
		return "Debugging"
	case LevelInfo:
		return "Information"
	case LevelWarn:
		return "Warning"
	case LevelError:
		return "Error"
	case LevelPanic:
		return "Panic"
	default:
		return fmt.Sprintf("Level(%d)", byte(lvl))
	}
}

// ParseLevel tries to convert a string to a logging level value. It supports
// all the outputs from the logging level String() method, and some more.
func ParseLevel(lvl string) (Level, error) {
	switch strings.TrimSpace(strings.ToLower(lvl)) {
	case "d", "debug", "debugging":
		return LevelDebug, nil
	case "i", "info", "information":
		return LevelInfo, nil
	case "w", "warn", "warning":
		return LevelWarn, nil
	case "e", "error":
		return LevelError, nil
	case "p", "panic":
		return LevelPanic, nil
	default:
		return LevelDebug, fmt.Errorf("invalid logging level string: %q", lvl)
	}
}
