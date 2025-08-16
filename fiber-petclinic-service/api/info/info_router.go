package info

import (
	"fiber-petclinic-service/config/app"
	resterr "fiber-petclinic-service/internal/errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Router struct {
	infoService Servicer
	ipService   IPServicer
}

// NewRouter creates a new Router
func NewRouter(infoService Servicer, ipService IPServicer) *Router {
	return &Router{
		infoService: infoService,
		ipService:   ipService,
	}
}

// Register registers the router to the gin engine
func (infoRouter *Router) Register(router fiber.Router) {
	router.Get("", infoRouter.appInfo)
}

// appInfo	Show app info
func (infoRouter *Router) appInfo(c *fiber.Ctx) error {
	log.Info().Msg("Fetching app info")
	result, err := infoRouter.infoService.getAppInfo()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching app info")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	ips, err := infoRouter.ipService.lookupIP(app.Config.Server.Host)
	if err != nil {
		log.Error().Err(err).Msg("Failed to lookup IP address")
		result.Ip = "Unknown host"
	} else {
		// Convert IP addresses to string and join them with commas
		strIPs := make([]string, len(ips))
		for i, ip := range ips {
			strIPs[i] = ip.String()
		}
		log.Debug().Str("host", app.Config.Server.Host).Strs("ips", strIPs).Msg("IP addresses for host")
		result.Ip = strings.Join(strIPs, ",")
	}

	// Content-Type will be application/json by c.JSON
	return c.Status(fiber.StatusOK).JSON(result)
}
