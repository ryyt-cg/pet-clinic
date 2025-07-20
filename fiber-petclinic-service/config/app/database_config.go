package app

type DatabaseConfig struct {
	Driver      string               `yaml:"driver" validate:"required"`
	Host        string               `yaml:"host" validate:"required"`
	Port        int                  `yaml:"port" validate:"required"`
	Username    string               `yaml:"username" validate:"required"`
	Password    string               `yaml:"password" validate:"required"`
	ConnectPool ConnectionPoolConfig `yaml:"connectPool"`
}

type ConnectionPoolConfig struct {
	MaxIdleConns int `yaml:"maxIdleConns" validate:"required"`
	MaxOpenConns int `yaml:"maxOpenConns" validate:"required"`
	MaxIdleTime  int `yaml:"maxIdleTime" validate:"required"` // in seconds
}
