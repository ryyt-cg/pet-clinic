package app

type DatabaseConfig struct {
	Driver         string               `yaml:"driver" validate:"required"`
	Host           string               `yaml:"host" validate:"required"`
	Port           int                  `yaml:"port"`
	Name           string               `yaml:"name" validate:"required"`
	Username       string               `yaml:"username" envconfig:"DB_USERNAME" validate:"required"`
	Password       string               `yaml:"password" envconfig:"DB_PASSWORD"`
	SslMode        string               `yaml:"sslMode" envconfig:"DB_SSL_MODE"`
	ConnectionPool ConnectionPoolConfig `yaml:"connectionPool"`
}

type ConnectionPoolConfig struct {
	MaxIdleConnections int `yaml:"maxIdleConnections" validate:"required" envconfig:"MAX_IDLE_CONNECTION"`
	MaxOpenConnections int `yaml:"maxOpenConnections" validate:"required" envconfig:"MAX_OPEN_CONNECTION"`
	MaxIdleTime        int `yaml:"maxIdleTime" validate:"required" envconfig:"MAX_IDLE_TIME"` // in seconds
}
