package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CheckRouter struct {
	healthCheckService Servicer
}

func NewRouter(healthCheckService Servicer) *CheckRouter {
	return &CheckRouter{healthCheckService}
}

// Register registers the router to the gin engine
func (healthCheckRouter *CheckRouter) Register(router *gin.RouterGroup) {
	router.GET("", healthCheckRouter.healthCheck)
}

// healthCheck	Show health check
func (healthCheckRouter *CheckRouter) healthCheck(c *gin.Context) {
	result, err := healthCheckRouter.healthCheckService.check()

	if err != nil {
		c.JSON(http.StatusBadRequest, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}
