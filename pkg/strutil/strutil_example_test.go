package strutil_test

import (
	"fmt"

	"github.com/iver-wharf/wharf-core/pkg/strutil"
)

func ExampleFirstRuneUpper() {
	fmt.Println(strutil.FirstRuneUpper("hello world"))
	// Output:
	// Hello world
}

func ExampleFirstRuneLower() {
	fmt.Println(strutil.FirstRuneLower("HELLO WORLD"))
	// Output:
	// hELLO WORLD
}
