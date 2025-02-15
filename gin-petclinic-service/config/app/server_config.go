package app

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

type ServerConfig struct {
	HttpPort string `yaml:"httpPort"`
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

func (sc ServerConfig) Validate() error {
	return validation.ValidateStruct(&sc,
		validation.Field(&sc.HttpPort, validation.Required),
		validation.Field(&sc.HttpPort, validation.Match(regexp.MustCompile("^:\\d+$")).Error("must be in the format :dddd")),
	)
}
