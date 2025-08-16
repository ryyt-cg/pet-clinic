package vet

import (
	"errors"
	resterr "fiber3-petclinic-service/internal/errors"

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

// Register
// register vet endpoints
func (vetRouter *Router) Register(route fiber.Router) {
	route.Get("specialties", vetRouter.allSpecialties)
	route.Get("all", vetRouter.allVets)
	route.Get(":id", vetRouter.vetById)
	route.Get(":id/specialties", vetRouter.getVetByIdWithSpecialties)
	route.Post("", vetRouter.create)
	route.Put(":id", vetRouter.update)
	//route.Get("", vetRouter.vetByParam)
}

// allSpecialties - retrieve all specialties

func (vetRouter *Router) allSpecialties(c fiber.Ctx) error {
	log.Info().Msg("Retrieving all specialties")
	responses, err := vetRouter.service.getAllSpecialties()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Msg("No specialties found.")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound("No specialties found"))
		}
		log.Error().Err(err).Msg("Unable to get all specialties.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(responses)
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
func (vetRouter *Router) allVets(c fiber.Ctx) error {
	log.Info().Msg("Retrieving all vets.")
	responses, err := vetRouter.service.getAllVets()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Msg("No vets found.")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound("No vets found"))
		}
		log.Error().Err(err).Msg("Fail to get all vets.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))

	}

	return c.Status(fiber.StatusOK).JSON(responses)
}

// vetById - retrieve vet by id
func (vetRouter *Router) vetById(c fiber.Ctx) error {
	pathID := c.Params("id")
	log.Info().Str("id", pathID).Msg("Retrieving vet by id")

	id := fiber.Params[int](c, "id")
	//id, err := c.ParamsInt("id")
	//if err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	//}

	response, err := vetRouter.service.getVetById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound(err.Error()))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// AllVetsWithSpecialties - retrieve all vets with specialties
func (vetRouter *Router) getVetByIdWithSpecialties(c fiber.Ctx) error {
	strID := c.Params("id")
	log.Info().Str("id", strID).Msg("GET vet with specialties by ID")

	id := fiber.Params[int](c, "id")
	//id, err := c.ParamsInt("id")
	//if err != nil {
	//	log.Error().Err(err).Str("id", strID).Msg("Invalid vet ID")
	//	return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	//}

	response, err := vetRouter.service.getVetByIdWithSpecialties(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Int("id", id).Msg("Found no vet.")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound(err.Error()))
		}

		log.Error().Err(err).Int("id", id).Msg("Fail to get vet with specialties by ID.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// addNewVet - add new vet
func (vetRouter *Router) create(c fiber.Ctx) error {
	var vetRequest AddRequest
	err := c.Bind().Body(&vetRequest)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Unmarshal JSON.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	newVet, err := vetRouter.service.create(FromAddRequest(&vetRequest))
	if err != nil {
		log.Error().Err(err).Msg("Unable to create new vet.")
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}
	return c.Status(fiber.StatusCreated).JSON(newVet)
}

func (vetRouter *Router) update(c fiber.Ctx) error {
	id := fiber.Params[int](c, "id")
	//id, err := c.ParamsInt("id")
	//if err != nil {
	//	log.Error().Err(err).Msg("Unable to convert ID to integer.")
	//	return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	//}

	var vetRequest UpdateRequest
	err := c.Bind().Body(&vetRequest)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Unmarshal JSON.")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	vetEntity := FromUpdateRequest(&vetRequest)
	vetEntity.ID = uint(id)
	newVet, err := vetRouter.service.update(vetEntity)
	if err != nil {
		log.Error().Err(err).Msg("Unable to update vet.")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound("Vet not found with ID: " + c.Params("id")))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(newVet)
}

//func (vetRouter *Router) vetByParam(c fiber.Ctx) error {
//	var nameParam api.NameParam
//	err := c.BodyParser(&nameParam)
//	if err != nil {
//		log.Error().Err(err).Msg("Unable to bind query param.")
//		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
//	}
//
//	return vetRouter.byLastName(c, nameParam)
//}

// vetByLastName - retrieve vet by last name
//func (vetRouter *Router) byLastName(c fiber.Ctx, param api.NameParam) error {
//	response, err := vetRouter.service.getVetByLastName(param.Name)
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(resterr.InternalServerError(err.Error()))
//	}
//
//	c.JSON(http.StatusOK, response)
//}
