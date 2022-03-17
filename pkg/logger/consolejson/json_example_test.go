package consolejson_test

import (
	"time"

	"github.com/iver-wharf/wharf-core/v2/pkg/logger"
	"github.com/iver-wharf/wharf-core/v2/pkg/logger/consolejson"
)

func ExampleNew() {
	defer logger.ClearOutputs()
	logger.AddOutput(logger.LevelDebug, consolejson.New(consolejson.Config{
		DisableDate:       true,
		DisableCallerLine: true,
	}))

	logger.New().Debug().Message("Sample message.")

	// Output:
	// {"level":"debug","caller":"consolejson/json_example_test.go","message":"Sample message."}
}

func ExampleTimeFormat() {
	defer logger.ClearOutputs()
	logger.AddOutput(logger.LevelDebug, consolejson.New(consolejson.Config{
		DisableDate:   true,
		DisableCaller: true,
		TimeFormat:    consolejson.TimeUnix,
	}))
	logger.AddOutput(logger.LevelDebug, consolejson.New(consolejson.Config{
		DisableDate:   true,
		DisableCaller: true,
		TimeFormat:    time.Kitchen, // any format string is supported
	}))

	t := time.Date(2006, 1, 2, 3, 4, 5, 0, time.UTC)
	logger.New().Debug().WithTime("sample", t).Message("Sample message.")

	// Output:
	// {"level":"debug","message":"Sample message.","sample":1136171045}
	// {"level":"debug","message":"Sample message.","sample":"3:04AM"}
}
