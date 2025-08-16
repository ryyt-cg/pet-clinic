package info

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Router struct {
	infoService Servicer
	ipService   IPServicer
}

// NewRouter creates a new Router
func NewRouter(infoService Servicer, ipService IPServicer) *Router {
	return &Router{infoService, ipService}
}

// Register registers the router to the gin engine
func (infoRouter *Router) Register(router *gin.RouterGroup) {
	router.GET("", infoRouter.appInfo)
}

// appInfo	Show app info
func (infoRouter *Router) appInfo(c *gin.Context) {
	result, err := infoRouter.infoService.getAppInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ips, err := infoRouter.ipService.lookupIP("localhost")
	if err != nil {
		result.Ip = "Unknown host"
	} else {
		for _, ip := range ips {
			result.Ip += ip.String() + "; "
		}
		result.Ip = strings.TrimSpace(result.Ip)
	}

	c.JSON(http.StatusOK, result)
}
