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
	// foo:[DEBUG|consolepretty/pretty_example_test.go] Sample message.
}

func ExampleConfig_ScopeMinLengthAuto() {
	defer logger.ClearOutputs()
	logger.AddOutput(logger.LevelDebug, consolepretty.New(consolepretty.Config{
		DisableDate:       true,
		DisableCallerLine: true,

		ScopeMinLengthAuto: true,
	}))

	log1 := logger.NewScoped("WHARF")
	log2 := logger.NewScoped("GORM-debug")
	log3 := logger.New()
	log1.Debug().Message("Sample message.")
	log2.Debug().Message("Sample message.")
	log3.Debug().Message("Sample message.")

	// Output:
	// [DEBUG|WHARF     |consolepretty/pretty_example_test.go] Sample message.
	// [DEBUG|GORM-debug|consolepretty/pretty_example_test.go] Sample message.
	// [DEBUG|          |consolepretty/pretty_example_test.go] Sample message.
}

func ExampleConfig_Ellipsis() {
	defer logger.ClearOutputs()
	logger.AddOutput(logger.LevelDebug, consolepretty.New(consolepretty.Config{
		DisableDate:       true,
		DisableCallerLine: true,

		CallerMaxLength: 15,
		//Ellipsis:        "…", // can be overridden
	}))

	logger.New().Debug().Message("Sample message.")

	// Output:
	// [DEBUG|…mple_test.go] Sample message.
}
