package app

type ServerConfig struct {
	BaseURL           string `yaml:"baseURL" validate:"required"`
	Host              string `yaml:"host" validate:"required"`
	HttpPort          string `yaml:"httpPort" validate:"required"`
	CertFile          string `yaml:"certFile"`
	KeyFile           string `yaml:"keyFile"`
	EnableCompression bool   `yaml:"enableCompression" default:"false"`
	CompressionLevel  int    `yaml:"compressionLevel" default:"0"`
	LogLevel          string `yaml:"logLevel" validate:"required"`
}
