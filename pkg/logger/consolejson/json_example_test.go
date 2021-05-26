package consolejson_test

import (
	"time"

	"github.com/iver-wharf/wharf-core/pkg/logger"
	"github.com/iver-wharf/wharf-core/pkg/logger/consolejson"
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

func ExampleTimeFormat_unix() {
	defer logger.ClearOutputs()
	logger.AddOutput(logger.LevelDebug, consolejson.New(consolejson.Config{
		DisableDate:       true,
		DisableCallerLine: true,
		TimeFormat:        consolejson.TimeUnix,
	}))

	t := time.Date(2006, 1, 2, 3, 4, 5, 0, time.UTC)
	logger.New().Debug().WithTime("unix", t).Message("Sample message.")

	// Output:
	// {"level":"debug","caller":"consolejson/json_example_test.go","message":"Sample message.","unix":1136171045}
}

func ExampleTimeFormat_customFormat() {
	defer logger.ClearOutputs()
	logger.AddOutput(logger.LevelDebug, consolejson.New(consolejson.Config{
		DisableDate:       true,
		DisableCallerLine: true,
		TimeFormat:        time.Kitchen,
	}))

	t := time.Date(2006, 1, 2, 3, 4, 5, 0, time.UTC)
	logger.New().Debug().WithTime("kitchen", t).Message("Sample message.")

	// Output:
	// {"level":"debug","caller":"consolejson/json_example_test.go","message":"Sample message.","kitchen":"3:04AM"}
}
