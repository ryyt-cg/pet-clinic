package ready

import (
	"fiber-petclinic-service/config/app"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

// Config
// Configure your live and ready check here
func Config() healthcheck.Config {

	return healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: app.Config.Server.BaseURL + "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: app.Config.Server.BaseURL + "/ready",
	}
}
