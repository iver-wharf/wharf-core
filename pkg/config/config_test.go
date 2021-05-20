package config_test

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
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
		Moo: InnerConfig{
			Foo: "default foo",
		},
	}
}

func TestEnvConfig(t *testing.T) {
	var cfg RootConfig

	os.Setenv("MOO_FOO", "moo bla")
	os.Setenv("MOO", "moo")
	os.Setenv("FOO", "foo")
	os.Setenv("DAZ", "daz")

	// May be ugly, but it works
	// https://github.com/spf13/viper/issues/188#issuecomment-413368673
	b, err := yaml.Marshal(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	defaultCfg := bytes.NewReader(b)
	viper.SetConfigType("yaml")
	if err := viper.MergeConfig(defaultCfg); err != nil {
		t.Fatal(err)
	}

	viper.SetConfigFile("/etc/my-app")
	if err := viper.MergeInConfig(); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("%T: %v", err, err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Unmarshal(&cfg); err != nil {
		t.Fatal("unmarshal config:", err)
	}

	viper.Debug()

	if cfg.Moo.Foo != "moo foo" {
		t.Errorf("wanted MOO_FOO to be 'moo foo', got: %+v", cfg)
	}
	t.Errorf("%+v", cfg)
}
