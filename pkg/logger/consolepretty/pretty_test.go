package consolepretty

import (
	"fmt"
	"math"
	"testing"

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
