package logger

import "fmt"

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
