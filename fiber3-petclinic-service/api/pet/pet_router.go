package pet

import (
	"errors"
	resterr "fiber3-petclinic-service/internal/errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
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
// @Success		200	{object}	responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/pets/all 		[get]
func (router *Router) getAll(c fiber.Ctx) error {
	log.Info().Msg("GET all pets")
	responses, err := router.service.getAllPets()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("Found no pet.")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound(err.Error()))
		}

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
// @Success		200	{object}	responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/pets/{id} 		[get]
func (router *Router) getById(c fiber.Ctx) error {
	strID := c.Params("id")
	log.Info().Str("id", strID).Msg("GET pet by ID")

	id, err := strconv.Atoi(strID)
	if err != nil {
		log.Error().Err(err).Str("id", c.Params("id")).Msg("Invalid pet ID")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	response, err := router.service.getPetById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("find no pet by this id.")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound("Pet not found"))
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
// @Success		200	{object}	responses
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/pets/{id}/visits 		[get]
func (router *Router) getWithVisitsById(c fiber.Ctx) error {
	strID := c.Params("id")
	log.Info().Str("id", strID).Msg("GET pet with visits by ID")

	id, err := strconv.Atoi(strID)
	if err != nil {
		log.Error().Err(err).Str("id", c.Params("id")).Msg("Invalid pet ID")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	response, err := router.service.getPetWithVisitsById(uint(id))
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
// @Success		200	{object}	responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/pets 			[get]
func (router *Router) getByQueryParam(c fiber.Ctx) error {
	log.Info().Msg("GET pet by query param")
	val := c.Query("name")
	if val == "" {
		log.Error().Msg("param name is empty.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest("pet name is empty"))
	}

	return router.getByName(c, val)
}

// petByName - get pet by name
func (router *Router) getByName(c fiber.Ctx, param string) error {
	log.Info().Str("name", param).Msg("GET pet by name")

	responses, err := router.service.getPetsByName(param)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("name", param).Msg("find no pet by this name")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound("No pets found with name: " + param))
		}
		log.Error().Err(err).Msg("Unable to get pets by name.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}

// addNewPet - add new pet
// @Tags		pets
// @Summary		Add a new pet
//
// @Description	Add a new pet
// @Param		Request	body	addRequest	true	"Add pet"
// @Produce		json
// @Success		200	{object}	responses
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/pets	 		[post]
func (router *Router) create(c fiber.Ctx) error {
	log.Info().Msg("Add a new pet.")
	var request *addRequest
	err := c.Bind().Body(&request)

	if err != nil {
		log.Error().Err(err).Msg("Unable to Unmarshal JSON.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	petEntity, err := fromAddRequest(request)
	if err != nil {
		log.Error().Err(err).Msg("unable to convert request to pet entity.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	petResponse, err := router.service.create(petEntity)
	if err != nil {
		log.Error().Err(err).Msg("unable to create pet.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(petResponse)
}

// updatePet - update pet
// @Tags		pets
// @Summary		update a pet
//
// @Description	update pet
// @Param		id	path	int	true	"Pet ID"
// @Param		Request	body	updateRequest	true	"Update pet"
// @Produce		json
// @Success		200	{object}	responses
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/pets/{id}	 	[put]
func (router *Router) update(c fiber.Ctx) error {
	log.Info().Msg("Update a new pet.")
	var request request

	strID := c.Params("id")
	log.Info().Str("id", strID).Msg("GET pet by ID")

	id, err := strconv.Atoi(strID)
	if err != nil {
		log.Error().Err(err).Str("id", c.Params("id")).Msg("Invalid pet ID")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	err = c.Bind().Body(&request)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Unmarshal JSON.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	log.Info().Str("name", request.Name).Msg("Update a pet.")
	petEntity, err := toPet(&request)
	if err != nil {
		log.Error().Err(err).Msg("Invalid request data.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}
	petEntity.ID = uint(id)
	petResponse, err := router.service.update(petEntity)

	if err != nil {
		log.Error().Err(err).Msg("Unable to update pet.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(petResponse)
}
