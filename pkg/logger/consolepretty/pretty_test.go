package consolepretty

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/fatih/color"
	"github.com/iver-wharf/wharf-core/v2/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func printedIntLenSlow(number int) int {
	if number == 0 {
		return 1
	}
	if number < 0 {
		return printedIntLenSlow(-number) + 1 // +1 for the sign symbol
	}
	return int(math.Floor(math.Log10(float64(number)) + 1))
}

func TestPrintedIntLen(t *testing.T) {
	var testCases = []struct {
		input int
		want  int
	}{
		{0, 1},
		{5, 1},
		{-5, 2},
		{1251, 4},
		{-3356, 5},
		{math.MaxInt32, 10},
		{math.MinInt32, 11},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.input), func(t *testing.T) {
			assert.Equal(t, tc.want, printedIntLenSlow(tc.input), "printedIntLenSlow")
			assert.Equal(t, tc.want, printedIntLenFast(tc.input), "printedIntLenFast")
		})
	}
}

func TestContextWriteScope(t *testing.T) {
	tests := []struct {
		name    string
		scope   string
		config  Config
		longest int
		want    string
	}{
		{
			name:    "no scope",
			scope:   "",
			longest: 0,
			want:    "",
		},
		{
			name:    "with scope",
			scope:   "abc",
			longest: 0,
			want:    "|abc",
		},
		{
			name:    "padded",
			scope:   "abc",
			config:  Config{ScopeMinLength: 6},
			longest: 0,
			want:    "|abc   ",
		},
		{
			name:    "autopadded",
			scope:   "abc",
			config:  Config{ScopeMinLengthAuto: true},
			longest: 6,
			want:    "|abc   ",
		},
		{
			name:    "maxxed",
			scope:   "abcdef",
			config:  Config{ScopeMaxLength: 3},
			longest: 0,
			want:    "|abc",
		},
	}
	color.NoColor = true
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger.LongestScopeNameLength = tc.longest
			tc.config.Coloring = &DefaultColorConfig
			ctx := context{
				scope:  tc.scope,
				Config: &tc.config,
			}
			var buf bytes.Buffer
			ctx.writeScope(&buf)
			assert.Equal(t, tc.want, buf.String())
		})
	}
	logger.LongestScopeNameLength = 0
}

var varThatDisablesCompilerOptimizations int

func BenchmarkPrintedIntLenSlow(b *testing.B) {
	var r = varThatDisablesCompilerOptimizations
	for n := 0; n < b.N; n++ {
		r = printedIntLenSlow(125)
	}
	varThatDisablesCompilerOptimizations = r
}

func BenchmarkPrintedIntLenFast(b *testing.B) {
	var r = varThatDisablesCompilerOptimizations
	for n := 0; n < b.N; n++ {
		r = printedIntLenFast(125)
	}
	varThatDisablesCompilerOptimizations = r
}
