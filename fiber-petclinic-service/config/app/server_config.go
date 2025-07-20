package app

type ServerConfig struct {
	BaseURL  string `yaml:"baseURL" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	HttpPort string `yaml:"httpPort" validate:"required"`
	LogLevel string `yaml:"logLevel" validate:"required"`
}
