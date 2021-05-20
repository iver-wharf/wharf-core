package app

import (
	"time"

	"gopkg.in/yaml.v2"
)

// Version holds common version fields used in the different Wharf components
// to distinguish it from other versions. This metadata can be commonly viewed
// through the build application via an endpoint or commandline flag.
type Version struct {
	// Version is the version of this API build. A SemVer2.0.0 formatted version
	// prefixed with a single "v" is expected, but not enforced.
	//
	// For local development versions a value of "local dev", "local docker",
	// or something alike is recommended.
	Version string `json:"version" yaml:"version" example:"v1.0.0"`

	// BuildGitCommit is the Git commit that this version of the API was
	// built from.
	BuildGitCommit string `json:"buildGitCommit" yaml:"buildGitCommit" example:"10aaf36a71ffe4f021b3d85341f684931f333040"`

	// BuildDate is the date on which this version of the API was built.
	BuildDate time.Time `json:"buildDate" yaml:"buildDate" format:"date-time"`

	// BuildRef is the Wharf build ID/reference from which this version of
	// the API was build in.
	BuildRef uint `json:"buildRef" yaml:"buildRef"`
}

// UnmarshalVersionYAML reads a YAML formatted file body and returns the
// parsed Version.
func UnmarshalVersionYAML(in []byte) (Version, error) {
	var version Version
	err := yaml.Unmarshal(in, &version)
	return version, err
}
