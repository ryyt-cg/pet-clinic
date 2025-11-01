package owner

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

// NewRouter - Owner Router constructor
func NewRouter(service Servicer) *Router {
	return &Router{service}
}

// Register registers the router to the gin engine
func (r *Router) Register(router fiber.Router) {
	router.Get("all", r.allOwners)
	router.Get(":id", r.ownerById)
	router.Get(":id/pets", r.ownerByIdWithPets)
	router.Get("", r.ownersByLastName)
	router.Post("", r.addNewOwner)
	router.Put(":id", r.updateOwner)
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
func (r *Router) allOwners(c fiber.Ctx) error {
	log.Info().Msg("Getting all owners.")
	responses, err := r.service.getAllOwners()

	if err != nil {
		log.Error().Err(err).Msg("Fail to get all owners.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses)
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
func (r *Router) ownerById(c fiber.Ctx) error {
	strID := c.Params("id")
	log.Info().Str("id", strID).Msg("GET owner by ID")

	id, err := strconv.Atoi(strID)
	if err != nil {
		log.Error().Err(err).Str("id", strID).Msg("Invalid owner ID")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	response, err := r.service.getOwnerById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Int("id", id).Msg("Found no Owner.")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound(err.Error()))
		}

		log.Error().Err(err).Int("id", id).Msg("Fail to get owner by ID.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response)
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
func (r *Router) ownerByIdWithPets(c fiber.Ctx) error {
	strID := c.Params("id")
	log.Info().Str("id", strID).Msg("GET owner with pets by ID")

	id, err := strconv.Atoi(strID)
	if err != nil {
		log.Error().Err(err).Str("id", strID).Msg("Invalid owner ID")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	response, err := r.service.getOwnerByIdWithPets(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Err(err).Int("id", id).Msg("Find no owner.")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound(err.Error()))
		}

		log.Error().Err(err).Int("id", id).Msg("Fail to get owner by ID.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response)
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
func (r *Router) ownersByLastName(c fiber.Ctx) error {
	lastName := c.Query("last-name")
	log.Info().Str("lastName", lastName).Msg("GET owner by last name")

	response, err := r.service.getOwnerByLastName(lastName)
	if err != nil {
		log.Error().Err(err).Str("lastName", lastName).Msg("Fail to get owner by last name.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// addNewOwner - add new owner
// @Tags		owners
// @Summary		Insert new owner
//
// @Description	Insert new owner
// @Produce		json
// @Param		Request			body	addRequest	true	"Add owner"
// @Success		200	{object}	UpdateResponse
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/owners 		[post]
func (r *Router) addNewOwner(c fiber.Ctx) error {
	log.Info().Msg("Post a new owner.")
	ownerRequest := new(addRequest)
	err := c.Bind().Body(ownerRequest)
	if err != nil {
		log.Error().Err(err).Msg("Fail to Unmarshal owner JSON.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	newOwner, err := r.service.create(ownerRequest)
	if err != nil {
		log.Error().Err(err).Msg("Fail to create new owner.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(newOwner)
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
func (r *Router) updateOwner(c fiber.Ctx) error {
	strID := c.Params("id")
	log.Info().Str("id", strID).Msg("PUT update owner by ID")

	id, err := strconv.Atoi(strID)
	if err != nil {
		log.Error().Err(err).Str("id", strID).Msg("Invalid owner ID")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	var ownerRequest updateRequest
	err = c.Bind().Body(&ownerRequest)
	if err != nil {
		log.Error().Err(err).Msg("Fail to Unmarshal owner JSON.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	updatedOwner, err := r.service.update(uint(id), &ownerRequest)
	if err != nil {
		log.Error().Int("id", id).Msg("Fail to update owner.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))

	}

	return c.Status(fiber.StatusOK).JSON(updatedOwner)
}
