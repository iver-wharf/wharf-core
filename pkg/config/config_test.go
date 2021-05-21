package config_test

import (
	"os"
	"strings"
	"testing"

	"github.com/iver-wharf/wharf-core/pkg/config"
)

type InnerConfig struct {
	Foo string
}

type RootConfig struct {
	Moo InnerConfig
	Daz string
}

func DefaultConfig() RootConfig {
	return RootConfig{
		Daz: "default daz",
	}
}

var sampleConfYaml = `
#daz: conf daz
moo:
  foo: conf foo
`

func TestEnvConfig(t *testing.T) {
	var cfg RootConfig

	c := config.New(DefaultConfig())

	c.AddConfigYAML(strings.NewReader(sampleConfYaml))
	c.AddConfigYAMLFile("wharf-config.yml")
	c.AddEnvironmentVariables("")

	os.Setenv("MOO_FOO", "moo foo")

	if err := c.Unmarshal(&cfg); err != nil {
		t.Fatal("unmarshal config:", err)
	}

	if cfg.Moo.Foo != "moo foo" {
		t.Errorf("wanted MOO_FOO to be 'moo foo', got: %+v", cfg)
	}

	t.Errorf("%+v", cfg)
}
