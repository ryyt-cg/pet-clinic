package visit

import (
	"errors"
	resterr "fiber3-petclinic-service/pkg/errors"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
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
func (r *Router) visitById(c fiber.Ctx) error {
	strID := c.Params("id")
	log.Info().Str("id", strID).Msg("Visit by id.")

	id := fiber.Params[int](c, "id") // Ensure the parameter is parsed as an integer
	//id, err := c.ParamsInt("id")
	//if id == 0 {
	//	log.Error().Err(err).Str("id", strID).Msg("Invalid visit ID")
	//	return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	//}

	response, err := r.service.getVisitById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Str("id", strID).Msg("Visit by id not found")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound("get no visit by id"))
		}
		log.Error().Err(err).Str("id", strID).Msg("fail to get a visit by id")
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
func (r *Router) allVisits(c fiber.Ctx) error {
	log.Info().Msg("All visits.")
	responses, err := r.service.getAllVisits()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Msg("No visits found.")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound("No visits found"))
		}
		log.Error().Err(err).Msg("Fail to get all visits.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}

func (r *Router) addNewVisit(c fiber.Ctx) error {
	log.Info().Msg("Post a new visit.")
	var addRequest *AddRequest
	if err := c.Bind().Body(&addRequest); err != nil {
		log.Error().Err(err).Msg("Fail to Unmarshal visit JSON.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	addVisit, err := FromAddRequest(addRequest)
	if err != nil {
		log.Error().Err(err).Msg("Invalid Add Request.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	response, err := r.service.create(addVisit)
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
func (r *Router) updateVisit(c fiber.Ctx) error {
	log.Info().Msg("Update a visit.")
	strID := c.Params("id")

	id := fiber.Params[int](c, "id")
	//id, err := c.ParamsInt("id")
	//if err != nil {
	//	log.Error().Err(err).Str("id", strID).Msg("Invalid visit ID")
	//	return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	//}

	var updateRequest *UpdateRequest
	if err := c.Bind().Body(&updateRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	visit, err := FromUpdateRequest(updateRequest)
	if err != nil {
		log.Error().Err(err).Msg("Invalid Update Request.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}
	visit.ID = uint(id)
	response, err := r.service.update(visit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Str("id", strID).Msg("Visit by id not found")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound("visit not found"))
		}
		log.Error().Err(err).Msg("Fail to update visit.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
