package owner

import (
	"errors"
	resterr "github.com/rhtran/gin-petclinic-service/middleware/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Router struct {
	logger  *zap.Logger
	service Servicer
}

func NewRouter(logger *zap.Logger, service Servicer) *Router {
	return &Router{logger, service}
}

// Register registers the router to the gin engine
func (r *Router) Register(routerGroup *gin.RouterGroup) {
	routerGroup.GET("all", r.allOwners)
	routerGroup.GET(":id", r.ownerById)
	routerGroup.GET(":id/pets", r.ownerByIdWithPets)
	routerGroup.GET("", r.ownersByLastName)
	routerGroup.POST("", r.addNewOwner)
	routerGroup.PUT(":id", r.updateOwner)
}

// ownerById  get owner by ID
// @Tags		owners
// @Summary		Search owner by ID
//
// @Description	Get owner by ID
// @Param		id	path	int	true	"Owner ID"
// @Produce		json
// @Success		200	{object}	Response
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/owners/{id} [get]
func (r *Router) ownerById(c *gin.Context) {
	pathID := c.Param("id")

	id, err := strconv.Atoi(pathID)
	if err != nil {
		r.logger.Error("Fail to convert ID to int.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	response, err := r.service.getOwnerById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error("Found no Owner.", zap.Int("id", id), zap.String("error", err.Error()))
			c.JSON(http.StatusNotFound, resterr.NotFound(err.Error()))
			return
		}

		r.logger.Error("Fail to get owner by ID.", zap.Int("id", id), zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

// ownerByLastName - get owners by last name
// @Tags		owners
// @Summary		Search owners by last name
//
// @Description	Get owners by last name
// @Param		last-name	query	string	true	"Owner last name"
// @Produce		json
// @Success		200	{object}	Responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/owners [get]
func (r *Router) ownersByLastName(c *gin.Context) {
	lastName := c.Query("last-name")
	response, err := r.service.getOwnerByLastName(lastName)
	if err != nil {
		r.logger.Error("Fail to get owner by last name.",
			zap.String("lastName", lastName), zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	if response.Context.Count == 0 {
		r.logger.Error("Find no owners with this last name.", zap.String("lastName", lastName))
		c.JSON(http.StatusNotFound, resterr.NotFound("Find no owner with last name: "+lastName))
		return
	}

	c.JSON(http.StatusOK, response)
}

// allOwners - get all owners
// @Tags		owners
//
//	@Summary	List all owners
//
// @Description	Get all owners
// @Produce		json
// @Success		200	{object}	Responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/owners/all		[get]
func (r *Router) allOwners(c *gin.Context) {
	responses, err := r.service.getAllOwners()
	if err != nil {
		r.logger.Error("Fail to get all owners.", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	if responses.Context.Count == 0 {
		r.logger.Warn("Find no owner.")
		c.JSON(http.StatusNotFound, resterr.NotFound("Find no owner"))
		return
	}

	c.JSON(http.StatusOK, responses)
}

// ownerByIdWithPets - get owner by ID with pets
// @Tags		owners
// @Summary		Search owner that has pets by ID
//
// @Description	Get owner that has pets by ID
// @Param		id	path	int	true	"Owner ID"
// @Produce		json
// @Success		200	{object}	Response
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/owners/{id}/pets [get]
func (r *Router) ownerByIdWithPets(c *gin.Context) {
	pathID := c.Param("id")

	id, err := strconv.Atoi(pathID)
	if err != nil {
		r.logger.Error("Fail to convert ID to int.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	response, err := r.service.getOwnerByIdWithPets(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("Find no owner.",
				zap.Int("id", id), zap.String("error", err.Error()))
			c.JSON(http.StatusNotFound, resterr.NotFound(err.Error()))
			return
		}

		r.logger.Error("Fail to get owner by ID.",
			zap.Int("id", id), zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

// addNewOwner - add new owner
// @Tags		owners
// @Summary		Insert new owner
//
// @Description	Insert new owner
// @Produce		json
// @Param		Request			body	AddRequest	true	"Add owner"
// @Success		200	{object}	UpdateResponse
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/owners 		[post]
func (r *Router) addNewOwner(c *gin.Context) {
	var ownerRequest AddRequest
	err := c.ShouldBindJSON(&ownerRequest)
	if err != nil {
		r.logger.Error("Fail to Unmarshal JSON.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	newOwner, err := r.service.create(&ownerRequest)
	if err != nil {
		r.logger.Error("Fail to create new owner.", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, newOwner)
}

// updateOwner - update owner
// @Tags		owners
// @Summary		Update owner
//
// @Description	Update new owner
// @Produce		json
// @Param		id	path	int	true	"Owner ID"
// @Param		Request			body	UpdateRequest	true	"Update owner"
// @Success		200	{object}	UpdateResponse
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/owners/{id} 	[put]
func (r *Router) updateOwner(c *gin.Context) {
	var ownerRequest UpdateRequest
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.logger.Warn("Fail to convert ID to int.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	err = c.ShouldBindJSON(&ownerRequest)
	if err != nil {
		r.logger.Error("Fail to Unmarshal JSON.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	updatedOwner, err := r.service.update(uint(id), &ownerRequest)
	if err != nil {
		r.logger.Error("Fail to update owner.",
			zap.Int("id", id), zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, updatedOwner)
}
