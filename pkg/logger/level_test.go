package logger

import (
	"testing"
)

func TestParseLevel(t *testing.T) {
	var levels = []Level{
		LevelDebug,
		LevelInfo,
		LevelWarn,
		LevelError,
		LevelPanic,
	}
	for _, lvl := range levels {
		t.Run(lvl.String(), func(t *testing.T) {
			parsed, err := ParseLevel(lvl.String())
			if err != nil {
				t.Errorf("wanted %s, got error: %s", lvl, err)
			} else if parsed != lvl {
				t.Errorf("wanted %s, got: %s", lvl, parsed)
			}
		})
	}
}
