package config_test

import (
	"bytes"
	"fmt"
	"os"

	"github.com/iver-wharf/wharf-core/pkg/config"
)

func ExampleConfig_AddEnvironmentVariables() {
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
	defaultConfig := Config{
		Logging: Logging{
			LogLevel: "Warning",
		},
		Username: "postgres",
		DB: DBConfig{
			Host: "localhost",
			Port: 5432,
		},
	}
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
	c.Unmarshal(&config)

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
	defaultConfig := Config{
		Logging: Logging{
			LogLevel: "Warning",
		},
		Username: "postgres",
		DB: DBConfig{
			Host: "localhost",
			Port: 5432,
		},
	}
	// This content could come from go:embed or a HTTP response body
	yamlBytes := []byte(`
logLevel: Info
# YAML key names are case-insensitive
pAssWOrD: Sommar2020
db:
  port: 8080
`)
	c := config.New(defaultConfig)
	c.AddConfigYAML(bytes.NewReader(yamlBytes))

	var config Config
	c.Unmarshal(&config)

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
	defaultConfig := Config{
		Logging: Logging{
			LogLevel: "Warning",
		},
		Username: "postgres",
		DB: DBConfig{
			Host: "localhost",
			Port: 5432,
		},
	}
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
	c.Unmarshal(&config)

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
