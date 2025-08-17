package app

type DatabaseConfig struct {
	Driver         string               `yaml:"driver" validate:"required"`
	Host           string               `yaml:"host" validate:"required"`
	Port           int                  `yaml:"port"`
	Name           string               `yaml:"name" validate:"required"`
	Username       string               `yaml:"username"`
	Password       string               `yaml:"password"`
	SslMode        string               `yaml:"sslMode"`
	ConnectionPool ConnectionPoolConfig `yaml:"connectionPool"`
}

type ConnectionPoolConfig struct {
	MaxIdleConnections int `yaml:"maxIdleConnections" validate:"required"`
	MaxOpenConnections int `yaml:"maxOpenConnections" validate:"required"`
	MaxIdleTime        int `yaml:"maxIdleTime" validate:"required"` // in seconds
}
