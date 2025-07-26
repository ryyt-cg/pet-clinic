package health

type Servicer interface {
	check() (*Check, error)
}

type CheckService struct {
}

func NewService() *CheckService {
	return &CheckService{}
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
