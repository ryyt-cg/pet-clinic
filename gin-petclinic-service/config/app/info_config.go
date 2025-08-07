package app

import validation "github.com/go-ozzo/ozzo-validation/v4"

type AppInfoConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

func (aic AppInfoConfig) Validate() error {
	return validation.ValidateStruct(&aic,
		validation.Field(&aic.Name, validation.Required),
		validation.Field(&aic.Description, validation.Required),
		validation.Field(&aic.Version, validation.Required),
	)
}
