package logger_test

import (
	"fmt"

	"github.com/iver-wharf/wharf-core/v2/pkg/logger"
)

func ExampleLevel_String() {
	for _, level := range []logger.Level{
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelPanic,
	} {
		fmt.Println(level.String())
	}

	// Output:
	// Debugging
	// Information
	// Warning
	// Error
	// Panic
}
