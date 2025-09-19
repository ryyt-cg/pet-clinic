package ready

import (
	"fiber3-petclinic-service/config/app"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

// Config
// Configure your live and ready check here
func Config() healthcheck.Config {

	return healthcheck.Config{
		Probe: func(c fiber.Ctx) bool {
			err := app.PingRepo.Ping()
			if err != nil {
				return false
			}
			return true
		},
	}
}
