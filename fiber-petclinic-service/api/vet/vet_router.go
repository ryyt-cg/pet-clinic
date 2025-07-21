package vet

import (
	"errors"
	resterr "fiber-petclinic-service/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
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

func (vetRouter *Router) allSpecialties(c *fiber.Ctx) error {
	response, err := vetRouter.service.getAllSpecialties()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, response)
}

// vetById - retrieve vet by id
func (vetRouter *Router) vetById(c *fiber.Ctx) {
	pathID := c.Param("id")

	id, err := strconv.Atoi(pathID)
	if err != nil {
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	response, err := vetRouter.service.getVetById(id)
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

func (vetRouter *Router) vetByParam(c *fiber.Ctx) {
	lastName := c.Query("last-name")

	vetRouter.vetByLastName(c, lastName)
	return
}

// vetByLastName - retrieve vet by last name
func (vetRouter *Router) vetByLastName(c *fiber.Ctx, lastName string) {
	response, err := vetRouter.service.getVetByLastName(lastName)
	if err != nil {
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
func (vetRouter *Router) allVets(c *fiber.Ctx) {
	response, err := vetRouter.service.getAllVets()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, response)
}

// AllVetsWithSpecialties - retrieve all vets with specialties
func (vetRouter *Router) getVetByIdWithSpecialties(c *fiber.Ctx) {
	pathID := c.Param("id")

	id, err := strconv.Atoi(pathID)
	if err != nil {
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	response, err := vetRouter.service.getVetByIdWithSpecialties(id)
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
func (vetRouter *Router) create(c *fiber.Ctx) {
	var vetRequest Request
	err := c.ShouldBindJSON(&vetRequest)
	if err != nil {
		log.Error("Unable to Unmarshal JSON.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	newVet, err := vetRouter.service.create(ToVet(&vetRequest))
	c.JSON(http.StatusCreated, newVet)
}

func (vetRouter *Router) update(c *fiber.Ctx) {
	var vetRequest Request
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err := c.ShouldBindJSON(&vetRequest)

	if err != nil {
		log.Error("Unable to Unmarshal JSON.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	vetEntity := ToVet(&vetRequest)
	vetEntity.ID = uint(id)
	newVet, err := vetRouter.service.update(vetEntity)
	c.JSON(http.StatusCreated, newVet)
}
