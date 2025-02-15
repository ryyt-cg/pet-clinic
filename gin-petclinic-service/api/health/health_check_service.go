package health

import (
	"go.uber.org/zap"
)

type Servicer interface {
	check() (*Check, error)
}

type CheckService struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) *CheckService {
	return &CheckService{logger}
}

// check
// Returns the health status of the service ("UP" or "DOWN")
func (service *CheckService) check() (*Check, error) {
	healthCheck := &Check{
		Status: "UP",
	}

	// TODO: Add more health checks here
	// External services, database, etc.

	return healthCheck, nil
}
