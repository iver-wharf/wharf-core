package consolejson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInefficientlyEscapeJSON(t *testing.T) {
	var testCases = []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "nothing to escape",
			input: "foo bar",
			want:  "foo bar",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "quotes at start",
			input: `"foo bar`,
			want:  `\"foo bar`,
		},
		{
			name:  "quotes at end",
			input: `foo bar"`,
			want:  `foo bar\"`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := inefficientlyEscapeJSON(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNew_defaults(t *testing.T) {
	jsonSink := New(Config{}).(sink)
	assert.Equal(t, false, jsonSink.config.DisableDate)
	assert.Equal(t, "date", jsonSink.config.DateField)
	assert.Equal(t, "message", jsonSink.config.MessageField)
	assert.Equal(t, "level", jsonSink.config.LevelField)
}

func TestNew_escaping(t *testing.T) {
	conf := Config{
		DisableDate:  true,
		DateField:    `my "time" here`,
		MessageField: `"simon says"`,
		LevelField:   `lävel`,
	}
	jsonSink := New(conf).(sink)
	assert.Equal(t, true, jsonSink.config.DisableDate)
	assert.Equal(t, `my \"time\" here`, jsonSink.config.DateField)
	assert.Equal(t, `\"simon says\"`, jsonSink.config.MessageField)
	assert.Equal(t, `lävel`, jsonSink.config.LevelField)
}
