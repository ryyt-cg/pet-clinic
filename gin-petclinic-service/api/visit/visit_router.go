package visit

import (
	"gin-petclinic-service/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Router struct {
	service Servicer
}

func NewRouter(service Servicer) *Router {
	return &Router{service}
}

func (r *Router) Register(router *gin.RouterGroup) {
	router.GET("/all", r.allVisits)
	router.GET(":id", r.visitById)
	router.POST("", r.addNewVisit)
	router.PUT(":id", r.updateVisit)
}

// visitById - get visit by ID
// @Tags		visits
// @Summary		Get visit by id
//
// @Description	Get visit by ID
// @Param		id	path	int	true	"Visit ID"
// @Produce		json
// @Success		200	{object}	Response
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/visits/{id} 	[get]
func (r *Router) visitById(c *gin.Context) {
	pathID := c.Param("id")
	log.Info().Str("pathID", pathID).Msg("Get pet by id")

	id, err := strconv.Atoi(pathID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest(err.Error()))
		return
	}

	response, err := r.service.getVisitById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

// allVisits - get all visits
// @Tags	visits
// @Summary	List all visits
//
// @Description	Get all visits
// @Produce		json
// @Success		200	{object}	Responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/visits/all		[get]
func (r *Router) allVisits(c *gin.Context) {
	log.Info().Msg("get all visits")
	response, err := r.service.getAllVisits()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *Router) addNewVisit(c *gin.Context) {
	log.Info().Msg("add new visit")
	var visit Request
	if err := c.ShouldBindJSON(&visit); err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest(err.Error()))
		return
	}

	response, err := r.service.create(&visit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// updateVisit - update visit
// @Tags		visits
// @Summary		Update visit
//
// @Description	Update visit
// @Produce		json
// @Param		id	path		int	true	"Visit ID"
// @Param		Request			body		UpdateRequest	true	"Update visit"
// @Success		200	{object}	Responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/visits/{id} 	[put]
func (r *Router) updateVisit(c *gin.Context) {
	log.Info().Msg("update visit")
	var visit Request
	if err := c.ShouldBindJSON(&visit); err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest(err.Error()))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest(err.Error()))
		return
	}

	visit.ID = id
	response, err := r.service.update(&visit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}
