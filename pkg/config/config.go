package config

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const configTypeYAML = "yaml"

type Config interface {
	AddConfigYAMLFile(path string)
	AddConfigYAML(reader io.Reader)
	AddEnvironmentVariables(prefix string)
	Unmarshal(config interface{}) error
}

func New(defaultConfig interface{}) Config {
	return &config{
		defaultConfig: defaultConfig,
	}
}

type config struct {
	defaultConfig interface{}
	sources       []configSource
}

type configSource interface {
	name() string
	apply(v *viper.Viper) error
}

func (c *config) AddConfigYAMLFile(path string) {
	c.sources = append(c.sources, yamlFileSource{path})
}

func (c *config) AddConfigYAML(reader io.Reader) {
	c.sources = append(c.sources, yamlSource{reader})
}

func (c *config) AddEnvironmentVariables(prefix string) {
	c.sources = append(c.sources, envVarsSource{prefix})
}

func (c *config) Unmarshal(config interface{}) error {
	v := viper.New()
	initDefaults(v, c.defaultConfig)
	for _, s := range c.sources {
		if err := s.apply(v); err != nil {
			return fmt.Errorf("applying config source: %s: %w", s.name(), err)
		}
	}
	return v.Unmarshal(config)
}

func initDefaults(v *viper.Viper, defaultConfig interface{}) error {
	// Uses a workaround to force viper to read environment variables
	// by making it aware of all fields that exists so it can later map
	// environment variables correctly.
	// https://github.com/spf13/viper/issues/188#issuecomment-413368673
	b, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("setting config defaults: %w", err)
	}
	defaultCfg := bytes.NewReader(b)
	v.SetConfigType(configTypeYAML)
	if err := v.MergeConfig(defaultCfg); err != nil {
		return fmt.Errorf("setting config defaults: %w", err)
	}
	return nil
}

type yamlFileSource struct {
	path string
}

func (s yamlFileSource) name() string {
	return s.path
}

func (s yamlFileSource) apply(v *viper.Viper) error {
	if s.path == "" {
		// viper does not set config file if its empty, so viper.MergeInConfig()
		// would then use the previously set config path value
		return nil
	}
	v.SetConfigType(configTypeYAML)
	v.SetConfigFile(s.path)
	return v.MergeInConfig()
}

type yamlSource struct {
	reader io.Reader
}

func (s yamlSource) name() string {
	return "YAML io.Reader"
}

func (s yamlSource) apply(v *viper.Viper) error {
	v.SetConfigType(configTypeYAML)
	return v.MergeConfig(s.reader)
}

type envVarsSource struct {
	prefix string
}

func (s envVarsSource) name() string {
	return "environment variables"
}

func (s envVarsSource) apply(v *viper.Viper) error {
	v.SetEnvPrefix(s.prefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return nil
}
