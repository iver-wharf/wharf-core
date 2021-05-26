package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestLogging struct {
	LogLevel string
}

type TestDBConfig struct {
	Host string
	Port int
}

type TestConfig struct {
	// Both mapstructure:",squash" and yaml:",inline" needs to be set in
	// embedded structs, even if you're only reading environment variables
	TestLogging `mapstructure:",squash" yaml:",inline"`
	Username    string
	Password    string
	DB          TestDBConfig
}

const (
	defaultLogLevel = "default log level"
	defaultUsername = "default username"
	defaultPassword = "default password"
	defaultDBHost   = "default db host"
	defaultDBPort   = 12345
	updatedLogLevel = "updated log level"
	updatedPassword = "updated password"
	updatedPort     = 8080
)

var defaultConfig = TestConfig{
	TestLogging: TestLogging{
		LogLevel: defaultLogLevel,
	},
	Username: defaultUsername,
	Password: defaultPassword,
	DB: TestDBConfig{
		Host: defaultDBHost,
		Port: defaultDBPort,
	},
}

func assertUnmarshaledConfig(t *testing.T, c Builder) {
	var cfg TestConfig
	require.Nil(t, c.Unmarshal(&cfg), "failed to read config")
	assert.Equal(t, updatedLogLevel, cfg.LogLevel)
	assert.Equal(t, defaultUsername, cfg.Username)
	assert.Equal(t, updatedPassword, cfg.Password)
	assert.Equal(t, defaultDBHost, cfg.DB.Host)
	assert.Equal(t, updatedPort, cfg.DB.Port)
}

func TestConfig_AddEnvironmentVariables(t *testing.T) {
	cb := NewBuilder(defaultConfig)
	cb.AddEnvironmentVariables("")

	os.Clearenv()
	os.Setenv("DB_PORT", strconv.FormatInt(updatedPort, 10))
	os.Setenv("PASSWORD", updatedPassword)
	os.Setenv("LOGLEVEL", updatedLogLevel)
	// Environment variables must be all uppercase
	os.Setenv("Username", "not used")

	assertUnmarshaledConfig(t, cb)
}

func TestConfig_AddConfigYAML(t *testing.T) {
	yamlContent := fmt.Sprintf(`
logLevel: %s
# YAML key names are case-insensitive
pAssWOrD: %s
db:
  port: %d
`, updatedLogLevel, updatedPassword, updatedPort)
	cb := NewBuilder(defaultConfig)
	cb.AddConfigYAML(strings.NewReader(yamlContent))
	assertUnmarshaledConfig(t, cb)
}

func TestConfig_AddConfigYAMLFile(t *testing.T) {
	cb := NewBuilder(defaultConfig)
	cb.AddConfigYAMLFile("testdata/add-config-yaml-file.yml")
	assertUnmarshaledConfig(t, cb)
}
