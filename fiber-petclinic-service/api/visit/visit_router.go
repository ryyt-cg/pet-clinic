package visit

import (
	resterr "fiber-petclinic-service/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"strconv"
)

type Router struct {
	service Servicer
}

func NewRouter(service Servicer) *Router {
	return &Router{service}
}

func (r *Router) Register(router fiber.Router) {
	router.Get("/all", r.allVisits)
	router.Get(":id", r.visitById)
	router.Post("", r.addNewVisit)
	router.Put(":id", r.updateVisit)
}

// visitById - get visit by ID
// @Tags		visits
// @Summary		Get visit by id
//
// @Description	Get visit by ID
// @Param		id	path	int	true	"Visit ID"
// @Produce		json
// @Success		200	{object}	Response
// @Failure		400	{object}	errors.ErrorResponse
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/visits/{id} 	[get]
func (r *Router) visitById(c *fiber.Ctx) error {
	strID := c.Params("id")
	log.Info().Str("id", strID).Msg("Visit by id.")

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error().Err(err).Str("id", strID).Msg("Invalid visit ID")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	response, err := r.service.getVisitById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// allVisits - get all visits
// @Tags	visits
// @Summary	List all visits
//
// @Description	Get all visits
// @Produce		json
// @Success		200	{object}	Responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/visits/all		[get]
func (r *Router) allVisits(c *fiber.Ctx) error {
	log.Info().Msg("All visits.")
	responses, err := r.service.getAllVisits()
	if err != nil {
		log.Error().Err(err).Msg("Fail to get all visits.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	//if responses.Context.Count == 0 {
	//	log.Warn().Msg("Find no visit.")
	//	return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound("Find no visit"))
	//}

	return c.Status(fiber.StatusOK).JSON(responses)
}

func (r *Router) addNewVisit(c *fiber.Ctx) error {
	var visit Request
	if err := c.BodyParser(&visit); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	response, err := r.service.create(&visit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// updateVisit - update visit
// @Tags		visits
// @Summary		Update visit
//
// @Description	Update visit
// @Produce		json
// @Param		id	path		int	true	"Visit ID"
// @Param		Request			body		UpdateRequest	true	"Update visit"
// @Success		200	{object}	Responses
// @Failure		404	{object}	errors.ErrorResponse
// @Failure		500	{object}	errors.ErrorResponse
// @Router		/visits/{id} 	[put]
func (r *Router) updateVisit(c *fiber.Ctx) error {
	var visit Request
	if err := c.BodyParser(&visit); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	visit.ID = id
	response, err := r.service.update(&visit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
