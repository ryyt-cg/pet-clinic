package info

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Router struct {
	infoService Servicer
}

// NewRouter creates a new Router
func NewRouter(infoService Servicer) *Router {
	return &Router{
		infoService,
	}
}

// Register registers the router to the gin engine
func (infoRouter *Router) Register(router fiber.Router) {
	router.Get("", infoRouter.appInfo)
}

// appInfo	Show app info
func (infoRouter *Router) appInfo(c *fiber.Ctx) error {
	log.Debug().Str("appInfo", "appInfo").Msg("Fetching app info")
	result, err := infoRouter.infoService.getAppInfo()
	if err != nil {
		err := c.JSON(err)
		if err != nil {
			return err
		}
	}

	// Content-Type will be application/json by c.JSON
	return c.JSON(result)
}
