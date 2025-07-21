package pet

import (
	"errors"
	"fiber-petclinic-service/api"
	resterr "fiber-petclinic-service/pkg/errors"
	"fiber-petclinic-service/pkg/repository/model"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Router struct {
	service Servicer
}

// NewRouter - create new router
func NewRouter(service Servicer) *Router {
	return &Router{service}
}

// Register - register routes
func (router *Router) Register(route fiber.Router) {
	route.Get("all", router.getAll)
	route.Get(":id", router.getById)
	route.Get(":id/visits", router.getWithVisitsById)
	route.Get("", router.getByQueryParam)
	route.Post("", router.create)
	route.Put(":id", router.update)
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
func (router *Router) getAll(c *fiber.Ctx) error {
	log.Info().Msg("GET all pets")
	responses, err := router.service.getAllPets()
	if err != nil {
		log.Error().Err(err).Msg("Unable to get all pets.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))

	}

	return c.Status(fiber.StatusOK).JSON(responses)
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
func (router *Router) getById(c *fiber.Ctx) error {
	strID := c.Params("id")
	log.Info().Str("id", strID).Msg("GET pet by ID")

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error().Err(err).Str("id", c.Params("id")).Msg("Invalid pet ID")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	response, err := router.service.getPetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("find no pet by this id.")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound(err.Error()))
		}
		log.Error().Err(err).Msg("Unable to get pet by id.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response)
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
func (router *Router) getWithVisitsById(c *fiber.Ctx) error {
	strID := c.Params("id")
	log.Info().Str("id", strID).Msg("GET pet with visits by ID")

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error().Err(err).Str("id", c.Params("id")).Msg("Invalid pet ID")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	response, err := router.service.getPetWithVisitsById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("find no pet by this id.")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound(err.Error()))
		}
		log.Error().Err(err).Msg("Unable to get pet by id.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response)
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
func (router *Router) getByQueryParam(c *fiber.Ctx) error {
	var nameParam api.NameParam
	err := c.BodyParser(&nameParam)
	if err != nil {
		log.Error().Err(err).Msg("Unable to bind query param.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	return router.getByName(c, nameParam)
}

// petByName - get pet by name
func (router *Router) getByName(c *fiber.Ctx, param api.NameParam) error {
	pets, err := router.service.getPetsByName(param.Name)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	if len(pets) == 0 {
		log.Warn().Str("name", param.Name).Msg("find no pet by this name")
		return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound("No pets found with name: " + param.Name))
	}

	responses := Responses{
		Context: model.Context{
			Count: len(pets),
		},
		Pets: pets,
	}
	return c.Status(fiber.StatusOK).JSON(responses)
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
func (router *Router) create(c *fiber.Ctx) error {
	log.Info().Msg("Add a new pet.")
	var request Request
	err := c.BodyParser(&request)

	if err != nil {
		log.Error().Err(err).Msg("Unable to Unmarshal JSON.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	petResponse, err := router.service.create(ToPet(&request))
	return c.Status(fiber.StatusCreated).JSON(petResponse)
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
func (router *Router) update(c *fiber.Ctx) error {
	log.Info().Msg("Update a new pet.")
	var request Request

	strID := c.Params("id")
	log.Info().Str("id", strID).Msg("GET pet by ID")

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error().Err(err).Str("id", c.Params("id")).Msg("Invalid pet ID")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	err = c.BodyParser(&request)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Unmarshal JSON.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	log.Info().Str("name", request.Name).Msg("Update a pet.")
	petEntity := ToPet(&request)
	petEntity.ID = uint(id)
	petResponse, err := router.service.update(petEntity)

	if err != nil {
		log.Error().Err(err).Msg("Unable to update pet.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(petResponse)
}
