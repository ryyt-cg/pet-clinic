package pet

import (
	"errors"
	"gin-petclinic-service/api"
	resterr "gin-petclinic-service/pkg/errors"
	"gin-petclinic-service/pkg/model"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Router struct {
	service Servicer
}

// NewRouter - create new router
func NewRouter(service Servicer) *Router {
	return &Router{service}
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
	log.Info().Msg("Retrieve all pets")
	responses, err := router.service.getAllPets()
	if err != nil {
		log.Error().Err(err).Msg("Unable to get all pets.")
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
	log.Info().Msg("Retrieve pet by id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error().Err(err).Msg("Unable to convert to number.")
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	response, err := router.service.getPetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Int("id", id).Msg("find no pet by this id.")
			c.JSON(http.StatusNotFound, resterr.NotFound(err.Error()))
			return
		}
		log.Error().Err(err).Int("id", id).Msg("Unable to get pet by id.")
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
		log.Error().Err(err).Msg("Unable to convert to number.")
		c.JSON(http.StatusBadRequest, resterr.BadRequest(err.Error()))
		return
	}

	response, err := router.service.getPetWithVisitsById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("find no pet by this id.")
			c.JSON(http.StatusNotFound, resterr.NotFound(err.Error()))
			return
		}
		log.Error().Err(err).Msg("Unable to get pet by id.")
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
		log.Error().Err(err).Msg("Unable to bind query param.")
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
		log.Error().Err(err).Str("name", param.Name).Msg("find no pet by this name")
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
		log.Error().Err(err).Msg("Unable to Unmarshal JSON.")
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	log.Info().Str("name", request.Name).Msg("Add a new pet.")
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
		log.Error().Err(err).Msg("Unable to convert to number.")
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	err = c.ShouldBind(&request)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Unmarshal JSON.")
		c.JSON(http.StatusBadRequest, resterr.BadRequestWithDetails(err))
		return
	}

	log.Info().Str("name", request.Name).Msg("Update a pet.")
	petEntity := ToPet(&request)
	petEntity.ID = uint(id)
	petResponse, err := router.service.update(petEntity)

	if err != nil {
		log.Error().Err(err).Msg("Unable to update pet.")
		c.JSON(http.StatusInternalServerError, resterr.InternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, petResponse)
}
