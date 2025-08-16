package vet

import (
	"errors"
	resterr "gin-petclinic-service/internal/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Router struct {
	service Servicer
}

func NewRouter(service Servicer) *Router {
	return &Router{service}
}

// Register
// register vet endpoints
func (vetRouter *Router) Register(router *gin.RouterGroup) {
	router.GET("specialties", vetRouter.allSpecialties)
	router.GET("all", vetRouter.allVets)
	router.GET(":id", vetRouter.vetById)
	router.GET(":id/specialties", vetRouter.getVetByIdWithSpecialties)
	router.GET("", vetRouter.vetByParam)
	router.POST("", vetRouter.create)
	router.PUT(":id", vetRouter.update)
}

// allSpecialties - retrieve all specialties

func (vetRouter *Router) allSpecialties(c *gin.Context) {
	log.Info().Msg("get all specialties")
	response, err := vetRouter.service.getAllSpecialties()
	if err != nil {
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

// vetById - retrieve vet by id
func (vetRouter *Router) vetById(c *gin.Context) {
	pathID := c.Param("id")

	id, err := strconv.Atoi(pathID)
	if err != nil {
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	response, err := vetRouter.service.getVetById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, resterr.NotFound(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (vetRouter *Router) vetByParam(c *gin.Context) {
	lastName := c.Query("last-name")

	vetRouter.vetByLastName(c, lastName)
	return
}

// vetByLastName - retrieve vet by last name
func (vetRouter *Router) vetByLastName(c *gin.Context, lastName string) {
	response, err := vetRouter.service.getVetByLastName(lastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

// AllVets - retrieve all vets
// @Tags	vets
//
//	@Summary	List all vets
//
// @Description	Get all vets
// @Produce		json
// @Success		200	{object}	Responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/vets/all 		[get]
func (vetRouter *Router) allVets(c *gin.Context) {
	response, err := vetRouter.service.getAllVets()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, response)
}

// AllVetsWithSpecialties - retrieve all vets with specialties
func (vetRouter *Router) getVetByIdWithSpecialties(c *gin.Context) {
	pathID := c.Param("id")

	id, err := strconv.Atoi(pathID)
	if err != nil {
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	response, err := vetRouter.service.getVetByIdWithSpecialties(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, resterr.NotFound(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

// addNewVet - add new vet
func (vetRouter *Router) create(c *gin.Context) {
	var vetRequest AddRequest
	err := c.ShouldBindJSON(&vetRequest)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Unmarshal JSON.")
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	vetEntity := FromAddRequest(&vetRequest)
	newVet, err := vetRouter.service.create(vetEntity)
	if err != nil {
		log.Error().Err(err).Msg("Fail to create vet.")
		if err != nil {
			c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
			return
		}
	}
	c.JSON(http.StatusCreated, newVet)
}

func (vetRouter *Router) update(c *gin.Context) {
	var vetRequest UpdateRequest
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	err = c.ShouldBindJSON(&vetRequest)
	// ShouldBindJSON will return an error if the request body is not valid JSON
	if err != nil {
		log.Error().Err(err).Msg("Unable to Unmarshal JSON.")
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	vetEntity := FromUpdateRequest(&vetRequest)
	vetEntity.ID = uint(id)
	newVet, err := vetRouter.service.update(vetEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, resterr.NotFound(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, newVet)
}
