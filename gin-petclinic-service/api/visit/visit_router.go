package visit

import (
	"errors"
	resterr "gin-petclinic-service/internal/errors"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

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
// @Success		200	{object}	response
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/visits/{id} 	[get]
func (r *Router) visitById(c *gin.Context) {
	pathID := c.Param("id")
	log.Info().Str("pathID", pathID).Msg("Get pet by id")

	id, err := strconv.Atoi(pathID)
	if err != nil {
		log.Error().Err(err).Str("id", pathID).Msg("Invalid visit ID")
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	response, err := r.service.getVisitById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Str("id", pathID).Msg("Visit by id not found")
			c.JSON(http.StatusNotFound, resterr.NotFound(err.Error()))
			return
		}

		log.Error().Err(err).Str("id", pathID).Msg("fail to get a visit by id")
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
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
// @Success		200	{object}	responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/visits/all		[get]
func (r *Router) allVisits(c *gin.Context) {
	log.Info().Msg("get all visits")
	response, err := r.service.getAllVisits()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("get no visits found")
			c.JSON(http.StatusNotFound, resterr.NotFound(err.Error()))
			return
		}

		log.Error().Err(err).Msg("get all visits failed")
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *Router) addNewVisit(c *gin.Context) {
	log.Info().Msg("add new visit")
	var visit addRequest
	if err := c.ShouldBindJSON(&visit); err != nil {
		log.Error().Err(err).Msg("fail to bind json request")
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	visitEntity, err := fromAddRequest(&visit)
	if err != nil {
		log.Error().Err(err).Msg("fail to bind json request")
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}
	response, err := r.service.create(visitEntity)
	if err != nil {
		log.Error().Err(err).Msg("fail to create visit")
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
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
// @Param		request			body		updateRequest	true	"Update visit"
// @Success		200	{object}	responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/visits/{id} 	[put]
func (r *Router) updateVisit(c *gin.Context) {
	log.Info().Msg("update visit")
	var visit updateRequest
	if err := c.ShouldBindJSON(&visit); err != nil {
		log.Error().Err(err).Msg("fail to bind json request")
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("fail to convert id to int")
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	visitEntity, err := fromUpdateRequest(&visit)
	if err != nil {
		log.Error().Err(err).Msg("fail to convert to visit entity")
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}
	visitEntity.ID = uint(id)
	response, err := r.service.update(visitEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("update visit by id not found")
			c.JSON(http.StatusNotFound, resterr.NotFound(err.Error()))
		}
		log.Error().Err(err).Msg("fail to update visit")
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}
