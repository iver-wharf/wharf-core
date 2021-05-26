package config_test

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"

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
	cfgBuilder := config.NewBuilder(defaultConfig)
	cfgBuilder.AddConfigYAML(bytes.NewReader(embeddedConfig))
	cfgBuilder.AddConfigYAMLFile("/etc/my-app/config.yml")
	cfgBuilder.AddConfigYAMLFile("$HOME/.config/my-app/config.yml")
	cfgBuilder.AddConfigYAMLFile("my-app-config.yml") // from working directory
	cfgBuilder.AddEnvironmentVariables("MYAPP")

	os.Setenv("MYAPP_PASSWORD", "Sommar2020")

	var cfg Config
	if err := cfgBuilder.Unmarshal(&cfg); err != nil {
		fmt.Println("Failed to read config:", err)
		return
	}

	fmt.Println("Log level:", cfg.LogLevel) // set from embeddedConfig
	fmt.Println("Username: ", cfg.Username) // uses defaultConfig.Username
	fmt.Println("Password: ", cfg.Password) // set from environment variable
	fmt.Println("DB host:  ", cfg.DB.Host)  // uses defaultConfig.DB.Host
	fmt.Println("DB port:  ", cfg.DB.Port)  // set from embeddedConfig

	// Output:
	// Log level: Info
	// Username:  postgres
	// Password:  Sommar2020
	// DB host:   localhost
	// DB port:   8080
}
