package env_test

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/iver-wharf/wharf-core/v2/pkg/env"
)

func ExampleBind() {
	os.Setenv("A", "1")
	os.Setenv("B", "2")

	type testType struct {
		A string
		B int
		C time.Duration
	}

	var x testType

	fmt.Printf("Before: A: %q\n", x.A)
	fmt.Printf("Before: B: %d\n", x.B)

	env.Bind(&x.A, "A")
	env.Bind(&x.B, "B")

	fmt.Printf("After: A: %q\n", x.A)
	fmt.Printf("After: B: %d\n", x.B)

	os.Setenv("C", "foo bar")
	err := env.Bind(&x.C, "C")
	fmt.Printf("Parse C: %s (is ErrParse? %t)\n", err, errors.Is(err, env.ErrParse))

	// Output:
	// Before: A: ""
	// Before: B: 0
	// After: A: "1"
	// After: B: 2
	// Parse C: env "C"="foo bar": time: invalid duration "foo bar" (is ErrParse? true)
}
