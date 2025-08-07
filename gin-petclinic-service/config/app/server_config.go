package app

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

type ServerConfig struct {
	BaseURL           string `yaml:"baseURL" validate:"required"`
	Host              string `yaml:"host" validate:"required"`
	HttpPort          string `yaml:"httpPort" validate:"required"`
	CertFile          string `yaml:"certFile"`
	KeyFile           string `yaml:"keyFile"`
	LogLevel          string `yaml:"logLevel" validate:"required"`
	EnablePrintRoutes bool   `yaml:"enablePrintRoutes"`
}

func (sc ServerConfig) Validate() error {
	return validation.ValidateStruct(&sc,
		validation.Field(&sc.HttpPort, validation.Required),
		validation.Field(&sc.HttpPort, validation.Match(regexp.MustCompile("^:\\d+$")).Error("must be in the format :dddd")),
	)
}
