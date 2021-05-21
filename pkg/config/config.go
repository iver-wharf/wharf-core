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

// Config type has methods for registering configuration sources, and then using
// those sources you can unmarshal into a struct to read the configuration.
//
// Later added config sources will merge on top of the previous on a
// per config field basis. Later added sources will override earlier added
// sources.
type Config interface {
	// AddConfigYAMLFile appends the path of a YAML file to the list of sources
	// for this configuration.
	//
	// Later added config sources will merge on top of the previous on a
	// per config field basis. Later added sources will override earlier added
	// sources.
	AddConfigYAMLFile(path string)

	// AddConfigYAML appends a byte reader for UTF-8 and YAML formatted content.
	// Useful for reading from embedded files, database stored configs, and from
	// HTTP response bodies.
	//
	// Later added config sources will merge on top of the previous on a
	// per config field basis. Later added sources will override earlier added
	// sources.
	AddConfigYAML(reader io.Reader)

	// AddEnvironmentVariables appends an environment variable source.
	//
	// Later added config sources will merge on top of the previous on a
	// per config field basis. Later added sources will override earlier added
	// sources.
	//
	// However, multiple environment variable sources cannot be added, due to
	// technical limitations with the implementation. Even if they use different
	// prefixes.
	//
	// Environment variables must be in all uppercase letters, and nested
	// structs use a single underscore "_" as delimiter. Example:
	//
	//  c := config.New(myConfigDefaults)
	//  c.AddEnvironmentVariables("FOO")
	//
	//  type MyConfig struct {
	//  	Bar        string // set via "FOO_BAR"
	//  	LoremIpsum string // set via "FOO_LOREMIPSUM"
	//  	Hello      struct {
	//  		World string // set via "FOO_HELLO_WORLD"
	//  	}
	//  }
	//
	// The prefix shall be without a trailing underscore "_" as this package
	// adds that in by itself. To not use a prefix, pass in an empty string as
	// prefix instead.
	AddEnvironmentVariables(prefix string)

	// Unmarshal applies the configuration, based on the numerous added sources,
	// on to an existing struct.
	//
	// For any field of type pointer and is set to nil, this function will
	// create a new instance and assign that before populating that branch.
	//
	// If none of the Config.Add...() functions has been called before this
	// function, then this function will effectively only apply the default
	// configuration onto this new object.
	//
	// The error that is returned is caused by any of the added config sources,
	// such as from invalid YAML syntax in an added YAML file.
	Unmarshal(config interface{}) error
}

// New creates a new Config based on a default configuration.
//
// Due to technical limitations, it's vital that this default configuration is
// of the same type that the config that you wish to unmarshal later, or at
// least that it contains fields with the same names.
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
