package info

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Router struct {
	logger      *zap.Logger
	infoService Servicer
	ipService   IPServicer
}

// NewRouter creates a new Router
func NewRouter(logger *zap.Logger, infoService Servicer, ipService IPServicer) *Router {
	return &Router{logger,
		infoService, ipService}
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

	addrs, err := infoRouter.ipService.lookupIP("localhost")
	if err != nil {
		result.Ip = "Unknown host"
	} else {
		for _, ia := range addrs {
			result.Ip += ia.String() + " "
		}
		result.Ip = strings.TrimSpace(result.Ip)
	}

	c.JSON(http.StatusOK, result)
}
