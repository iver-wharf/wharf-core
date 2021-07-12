package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func reset() {
	minGlobalLevel = LevelDebug
	minScopedLevels = make(map[string]Level)
	ClearOutputs()
}

func TestSetLevel(t *testing.T) {
	t.Cleanup(reset)

	mock := NewMock()
	SetLevel(LevelWarn)
	AddOutput(LevelDebug, mock)

	log := New()
	log.Debug().Message("Suppressed")
	log.Warn().Message("Logged")

	assert.ElementsMatch(t, mock.LogMessages, []string{"Logged"})
}

func TestSetLevelScoped(t *testing.T) {
	t.Cleanup(reset)

	mock := NewMock()
	SetLevelScoped(LevelWarn, "MY-SCOPE")
	AddOutput(LevelDebug, mock)

	NewScoped("MY-SCOPE").Debug().Message("Suppressed")
	NewScoped("MY-SCOPE").Warn().Message("Logged1")

	New().Debug().Message("Logged2")
	NewScoped("OTHER-SCOPE").Debug().Message("Logged3")

	assert.ElementsMatch(t, mock.LogMessages, []string{"Logged1", "Logged2", "Logged3"})
}

func TestLevelSilence(t *testing.T) {
	t.Cleanup(reset)

	mock := NewMock()
	AddOutput(LevelSilence, mock)

	New().Error().Message("Suppressed")

	assert.Len(t, mock.LogMessages, 0)
}

func TestLevelSilenceScoped(t *testing.T) {
	t.Cleanup(reset)

	mock := NewMock()
	AddOutput(LevelDebug, mock)
	SetLevelScoped(LevelSilence, "MY-SCOPE")

	NewScoped("MY-SCOPE").Error().Message("Suppressed")
	New().Error().Message("Logged")

	assert.ElementsMatch(t, mock.LogMessages, []string{"Logged"})
}
