package app

type DatabaseConfig struct {
	Driver         string               `yaml:"driver" validate:"required"`
	Host           string               `yaml:"host" validate:"required"`
	Port           int                  `yaml:"port" validate:"required"`
	Name           string               `yaml:"name" validate:"required"`
	Username       string               `yaml:"username" validate:"required"`
	Password       string               `yaml:"password" validate:"required"`
	ConnectionPool ConnectionPoolConfig `yaml:"connectionPool"`
}

type ConnectionPoolConfig struct {
	MaxIdleConnection int `yaml:"maxIdleConnection" validate:"required"`
	MaxOpenConnection int `yaml:"maxOpenConnection" validate:"required"`
	MaxIdleTime       int `yaml:"maxIdleTime" validate:"required"` // in seconds
}
