package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/go-playground/validator/v10"
	"github.com/vrischmann/envconfig"
)

// Config stores the application-wide configurations
var (
	Config   AppConfig
	validate *validator.Validate
)

type AppConfig struct {
	AppInfo   AppInfoConfig             `yaml:"appInfo" validate:"required" envconfig:"-"`
	Databases map[string]DatabaseConfig `yaml:"databases" validate:"required" envconfig:"-"`
	Gorm      GormConfig                `yaml:"gorm" envconfig:"-"`
	Server    ServerConfig              `yaml:"server" validate:"required" envconfig:"-"`
}

// Validate all config required values are populated.
func (config AppConfig) Validate() error {
	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(config.AppInfo); err != nil {
		panic(err.Error())
	}
	return nil
}

// LoadConfig loads configuration from the given list of paths and populates it into the Config variable.
// The configuration file(s) should be named as app.yaml.
func LoadConfig(configPath string) error {
	env := strings.ToUpper(os.Getenv("ENV"))
	configFile := configPath + "/" + getConfigFile(env) + ".yaml"

	loader := aconfig.LoaderFor(&Config, aconfig.Config{
		SkipFlags: true,
		Files:     []string{configFile},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(), // Register the YAML decoder
		},
	})
	if err := loader.Load(); err != nil {
		return fmt.Errorf("failed to load configuration file %s: %w", configFile, err)
	}

	// Override with environment variables
	err := envconfig.Init(&Config)
	if err != nil {
		return fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return Config.Validate()
}

func getConfigFile(env string) string {
	switch env {
	case "PRD":
		return "app-prd"
	case "DEV":
		return "app-dev"
	default:
		return "app"
	}
}
