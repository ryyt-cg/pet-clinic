package pet

import (
	"errors"
	"github.com/rhtran/gin-petclinic-service/internal/api"
	resterr "github.com/rhtran/gin-petclinic-service/middleware/errors"
	"github.com/rhtran/gin-petclinic-service/pkg/model"
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

// NewRouter - create new router
func NewRouter(logger *zap.Logger, service Servicer) *Router {
	return &Router{logger, service}
}

// Register - register routes
func (router *Router) Register(routerGroup *gin.RouterGroup) {
	routerGroup.GET("all", router.getAll)
	routerGroup.GET(":id", router.getById)
	routerGroup.GET(":id/visits", router.getWithVisitsById)
	routerGroup.GET("", router.getByQueryParam)
	routerGroup.POST("", router.create)
	routerGroup.PUT(":id", router.update)
}

// allPets - retrieve all pets
// @Tags		pets
// @Summary	List all pets
//
// @Description	Get all pets
// @Produce		json
// @Success		200	{object}	Responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/pets/all 		[get]
func (router *Router) getAll(c *gin.Context) {
	responses, err := router.service.getAllPets()
	if err != nil {
		router.logger.Error("Unable to get all pets.", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses)
}

// petById - retrieve pet by id
// @Tags		pets
// @Summary		Get pet by ID
//
// @Description	Get pet by ID
// @Param		id	path	int	true	"Pet ID"
// @Produce		json
// @Success		200	{object}	Responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/pets/{id} 		[get]
func (router *Router) getById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.logger.Error("Unable to convert to number.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	response, err := router.service.getPetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			router.logger.Error("find no pet by this id.", zap.String("error", err.Error()))
			c.JSON(http.StatusNotFound, resterr.NotFound(err.Error()))
			return
		}
		router.logger.Error("Unable to get pet by id.", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

// petWithVisitsById - get pet with visits by id
// @Tags		pets
// @Summary		Get pet with visits by ID
//
// @Description	Get pet with visits by ID
// @Param		id	path	int	true	"Pet ID"
// @Produce		json
// @Success		200	{object}	Responses
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/pets/{id}/visits 		[get]
func (router *Router) getWithVisitsById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.logger.Error("Unable to convert to number.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	response, err := router.service.getPetWithVisitsById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			router.logger.Error("find no pet by this id.", zap.String("error", err.Error()))
			c.JSON(http.StatusNotFound, resterr.NotFound(err.Error()))
			return
		}
		router.logger.Error("Unable to get pet by id.", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

// petByQueryParam -
// @Tags		pets
// @Summary		Get pet by name
//
// @Description	Get pet by name
// @Param		name	query	string	true	"Pet Name"
// @Produce		json
// @Success		200	{object}	Responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/pets 			[get]
func (router *Router) getByQueryParam(c *gin.Context) {
	var nameParam api.NameParam
	err := c.BindQuery(&nameParam)
	if err != nil {
		router.logger.Error("Unable to bind query param.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	router.getByName(c, nameParam)
}

// petByName - get pet by name
func (router *Router) getByName(c *gin.Context, param api.NameParam) {
	pets, err := router.service.getPetsByName(param.Name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(""))
		return
	}

	if len(pets) == 0 {
		router.logger.Error("find no pet by this name", zap.String("name", param.Name))
		c.JSON(http.StatusNotFound, resterr.NotFound("No pets found with name: "+param.Name))
		return
	}

	responses := Responses{
		Context: model.Context{
			Count: len(pets),
		},
		Pets: pets,
	}
	c.JSON(http.StatusOK, responses)
}

// addNewPet - add new pet
// @Tags		pets
// @Summary		Add a new pet
//
// @Description	Add a new pet
// @Param		Request	body	AddRequest	true	"Add pet"
// @Produce		json
// @Success		200	{object}	Responses
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/pets	 		[post]
func (router *Router) create(c *gin.Context) {
	var request Request
	err := c.ShouldBind(&request)
	if err != nil {
		router.logger.Error("Unable to Unmarshal JSON.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	router.logger.Info("Add a new pet.", zap.String("name", request.Name))
	petResponse, err := router.service.create(ToPet(&request))
	c.JSON(http.StatusCreated, petResponse)
}

// updatePet - update pet
// @Tags		pets
// @Summary		update a pet
//
// @Description	update pet
// @Param		id	path	int	true	"Pet ID"
// @Param		Request	body	UpdateRequest	true	"Update pet"
// @Produce		json
// @Success		200	{object}	Responses
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/pets/{id}	 	[put]
func (router *Router) update(c *gin.Context) {
	var request Request
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		router.logger.Error("Unable to convert to number.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	err = c.ShouldBind(&request)
	if err != nil {
		router.logger.Error("Unable to Unmarshal JSON.", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	router.logger.Info("Update a pet.", zap.String("name", request.Name))
	petEntity := ToPet(&request)
	petEntity.ID = uint(id)
	petResponse, err := router.service.update(petEntity)

	if err != nil {
		router.logger.Error("Unable to update pet.", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, petResponse)
}
