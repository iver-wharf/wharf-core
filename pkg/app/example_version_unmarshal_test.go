package app_test

import (
	"fmt"

	"github.com/iver-wharf/wharf-core/v2/pkg/app"
)

func ExampleUnmarshalVersionYAML() {
	var body = []byte(`
version: v1.0.0
buildGitCommit: 10aaf36a71ffe4f021b3d85341f684931f333040
buildDate: 2021-05-20T14:27:11+01:00
buildRef: 123
`)

	var version app.Version
	if err := app.UnmarshalVersionYAML(body, &version); err != nil {
		fmt.Println("Unexpected error:", err)
	}

	fmt.Println("Version:         ", version.Version)
	fmt.Println("Build Git commit:", version.BuildGitCommit)
	fmt.Println("Build date:      ", version.BuildDate)
	fmt.Println("Build reference: ", version.BuildRef)

	// Output:
	// Version:          v1.0.0
	// Build Git commit: 10aaf36a71ffe4f021b3d85341f684931f333040
	// Build date:       2021-05-20 14:27:11 +0100 +0100
	// Build reference:  123
}
