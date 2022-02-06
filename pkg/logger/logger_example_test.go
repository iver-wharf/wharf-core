package logger_test

import (
	"github.com/iver-wharf/wharf-core/v2/pkg/logger"
	"github.com/iver-wharf/wharf-core/v2/pkg/logger/consolejson"
	"github.com/iver-wharf/wharf-core/v2/pkg/logger/consolepretty"
)

var prettyConf = consolepretty.Config{
	DisableDate:   true,
	DisableCaller: true,
}
var jsonConf = consolejson.Config{
	DisableDate:   true,
	DisableCaller: true,
}

func ExampleAddOutput() {
	defer logger.ClearOutputs()

	logger.AddOutput(logger.LevelDebug, consolepretty.New(prettyConf))
	logger.AddOutput(logger.LevelInfo, consolejson.New(jsonConf))

	// will not be used due to too high logger.Level
	logger.AddOutput(logger.LevelError, consolepretty.New(prettyConf))

	log := logger.New()

	log.Info().Message(`first "log".`)
	log.Info().WithInt("id", 5).Message("second log.")
	log.Info().WithString("hello", "world").Message("third log.")

	// Output:
	// [INFO ] first "log".
	// {"level":"info","message":"first \"log\"."}
	// [INFO ] second log.  id=5
	// {"level":"info","message":"second log.","id":5}
	// [INFO ] third log.  hello=world
	// {"level":"info","message":"third log.","hello":"world"}
}

func ExampleNewScoped() {
	var log = logger.NewScoped("example")

	defer logger.ClearOutputs()
	logger.AddOutput(logger.LevelInfo, consolepretty.New(prettyConf))

	log.Info().Message("first log.")

	// Output:
	// [INFO |example] first log.
}
