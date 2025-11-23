package owner

import (
	"errors"
	resterr "fiber-petclinic-service/internal/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

var tracer = otel.Tracer("owners-router")

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
// @Summary	List all owners
//
// @Description	Get all owners
// @Produce		json
// @Success		200	{object}	Responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/owners/all		[get]
func (r *Router) allOwners(c *fiber.Ctx) error {
	_, span := tracer.Start(c.UserContext(), "ownerById", oteltrace.WithAttributes())
	defer span.End()

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
func (r *Router) ownerById(c *fiber.Ctx) error {
	strID := c.Params("id")
	_, span := tracer.Start(c.UserContext(), "ownerById", oteltrace.WithAttributes(attribute.String("id", strID)))
	defer span.End()

	log.Info().Str("id", strID).Msg("GET owner by ID")
	id, err := c.ParamsInt("id")

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
func (r *Router) ownerByIdWithPets(c *fiber.Ctx) error {
	strID := c.Params("id")

	_, span := tracer.Start(c.UserContext(), "ownerByIdWithPets", oteltrace.WithAttributes(attribute.String("id", strID)))
	defer span.End()

	log.Info().Str("id", strID).Msg("GET owner with pets by ID")
	id, err := c.ParamsInt("id")

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
func (r *Router) ownersByLastName(c *fiber.Ctx) error {
	lastName := c.Query("last-name")

	_, span := tracer.Start(c.UserContext(), "ownersByLastName", oteltrace.WithAttributes(attribute.String("lastName", lastName)))
	defer span.End()

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
// @Param		Request			body	AddRequest	true	"Add owner"
// @Success		200	{object}	UpdateResponse
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/owners 		[post]
func (r *Router) addNewOwner(c *fiber.Ctx) error {
	log.Info().Msg("Post a new owner.")

	_, span := tracer.Start(c.UserContext(), "addNewOwner", oteltrace.WithAttributes())
	defer span.End()

	ownerRequest := new(addRequest)
	err := c.BodyParser(ownerRequest)
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
func (r *Router) updateOwner(c *fiber.Ctx) error {
	strID := c.Params("id")

	_, span := tracer.Start(c.UserContext(), "updateOwner", oteltrace.WithAttributes(attribute.String("id", strID)))
	defer span.End()

	log.Info().Str("id", strID).Msg("PUT update owner by ID")
	id, err := c.ParamsInt("id")

	if err != nil {
		log.Error().Err(err).Str("id", strID).Msg("Invalid owner ID")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	var ownerRequest updateRequest
	err = c.BodyParser(&ownerRequest)
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
