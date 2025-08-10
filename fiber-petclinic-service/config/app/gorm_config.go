package app

type GormConfig struct {
	LogLevel      int  `yaml:"logLevel" default:"1" validate:"required"`
	SingularTable bool `yaml:"singularTable" default:"false" validate:"required"`
}
