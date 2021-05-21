package config

import (
	"io"
	"strings"

	"github.com/spf13/viper"
)

type Config interface {
	AddConfigYAMLFile(path string)
	AddConfigYAML(reader io.Reader)
	AddEnvironmentVariables(prefix string)
}

func New() Config {
	return &config{}
}

type config struct {
	sources []configSource
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

type configSource interface {
	apply(viper viper.Viper)
}

type yamlFileSource struct {
	path string
}

func (c yamlFileSource) apply(viper viper.Viper) {
	viper.AddConfigPath(c.path)
}

type yamlSource struct {
	reader io.Reader
}

func (c yamlSource) apply(viper viper.Viper) {
	viper.ReadConfig(c.reader)
}

type envVarsSource struct {
	prefix string
}

func (c envVarsSource) apply(viper viper.Viper) {
	viper.SetEnvPrefix(c.prefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
