package app

type CircuitBreakerConfig struct {
	FailureThreshold int `yaml:"failureThreshold" validate:"required"`
	Timeout          int `yaml:"timeout" validate:"required"` // in seconds
	SuccessThreshold int `yaml:"successThreshold" validate:"required"`
}
