package config_test

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/iver-wharf/wharf-core/pkg/config"
)

type Logging struct {
	LogLevel string
}

type DBConfig struct {
	Host string
	Port int
}

type Config struct {
	// Both mapstructure:",squash" and yaml:",inline" needs to be set in
	// embedded structs, even if you're only reading environment variables
	Logging  `mapstructure:",squash" yaml:",inline"`
	Username string
	Password string
	DB       DBConfig
}

var defaultConfig = Config{
	Logging: Logging{
		LogLevel: "Warning",
	},
	Username: "postgres",
	DB: DBConfig{
		Host: "localhost",
		Port: 5432,
	},
}

//go:embed testdata/embedded-config.yml
var embeddedConfig []byte

func ExampleConfig() {
	c := config.New(defaultConfig)
	c.AddConfigYAML(bytes.NewReader(embeddedConfig))
	c.AddConfigYAMLFile("/etc/my-app/config.yml")
	c.AddConfigYAMLFile("$HOME/.config/my-app/config.yml")
	c.AddConfigYAMLFile("my-app-config.yml") // from working directory
	c.AddEnvironmentVariables("MYAPP")

	var config Config
	if err := c.Unmarshal(&config); err != nil {
		fmt.Println("Failed to read config:", err)
		return
	}

	fmt.Println("Log level:", config.LogLevel)
	fmt.Println("Username: ", config.Username)
	fmt.Println("Password: ", config.Password)
	fmt.Println("DB host:  ", config.DB.Host)
	fmt.Println("DB port:  ", config.DB.Port)

	// Output:
	// Log level: Info
	// Username:  postgres
	// Password:  Sommar2020
	// DB host:   localhost
	// DB port:   8080
}

func ExampleConfig_AddEnvironmentVariables() {
	c := config.New(defaultConfig)
	c.AddEnvironmentVariables("")

	// Environment variables can, but should not, be set like this.
	// Recommended to set them externally instead.
	os.Setenv("DB_PORT", "8080")
	os.Setenv("PASSWORD", "Sommar2020")
	os.Setenv("LOGLEVEL", "Info")
	// Environment variables must be all uppercase
	os.Setenv("Username", "not used")

	var config Config
	if err := c.Unmarshal(&config); err != nil {
		fmt.Println("Failed to read config:", err)
		return
	}

	fmt.Println("Log level:", config.LogLevel)
	fmt.Println("Username: ", config.Username)
	fmt.Println("Password: ", config.Password)
	fmt.Println("DB host:  ", config.DB.Host)
	fmt.Println("DB port:  ", config.DB.Port)

	// Output:
	// Log level: Info
	// Username:  postgres
	// Password:  Sommar2020
	// DB host:   localhost
	// DB port:   8080
}

func ExampleConfig_AddConfigYAML() {
	// This content could come from go:embed or a HTTP response body
	yamlBytes := `
logLevel: Info
# YAML key names are case-insensitive
pAssWOrD: Sommar2020
db:
  port: 8080
`
	c := config.New(defaultConfig)
	c.AddConfigYAML(strings.NewReader(yamlBytes))

	var config Config
	if err := c.Unmarshal(&config); err != nil {
		fmt.Println("Failed to read config:", err)
		return
	}

	fmt.Println("Log level:", config.LogLevel)
	fmt.Println("Username: ", config.Username)
	fmt.Println("Password: ", config.Password)
	fmt.Println("DB host:  ", config.DB.Host)
	fmt.Println("DB port:  ", config.DB.Port)

	// Output:
	// Log level: Info
	// Username:  postgres
	// Password:  Sommar2020
	// DB host:   localhost
	// DB port:   8080
}

func ExampleConfig_AddConfigYAMLFile() {
	c := config.New(defaultConfig)
	// The file referenced here contains the following:
	//
	// 	logLevel: Info
	// 	# YAML key names are case-insensitive
	// 	pAssWOrD: Sommar2020
	// 	db:
	// 	  port: 8080
	c.AddConfigYAMLFile("testdata/add-config-yaml-file.yml")

	var config Config
	if err := c.Unmarshal(&config); err != nil {
		fmt.Println("Failed to read config:", err)
		return
	}

	fmt.Println("Log level:", config.LogLevel)
	fmt.Println("Username: ", config.Username)
	fmt.Println("Password: ", config.Password)
	fmt.Println("DB host:  ", config.DB.Host)
	fmt.Println("DB port:  ", config.DB.Port)

	// Output:
	// Log level: Info
	// Username:  postgres
	// Password:  Sommar2020
	// DB host:   localhost
	// DB port:   8080
}
