package consolepretty_test

import (
	"github.com/iver-wharf/wharf-core/pkg/logger"
	"github.com/iver-wharf/wharf-core/pkg/logger/consolepretty"
)

func ExampleNew() {
	defer logger.ClearOutputs()
	logger.AddOutput(logger.LevelDebug, consolepretty.New(consolepretty.Config{
		Prefix:            "foo:",
		DisableDate:       true,
		DisableCallerLine: true,
	}))

	logger.New().Debug().Message("Sample message.")

	// Output:
	// foo:[DEBUG | consolepretty/pretty_example_test.go] Sample message.
}
