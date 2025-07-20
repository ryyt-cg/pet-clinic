package author

import (
	"errors"
	"fiber-petclinic-service/exception"
	resterr "fiber-petclinic-service/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"net/url"
)

type Router struct {
	authorService Servicer
}

// NewRouter creates a new Router
func NewRouter(authorService Servicer) *Router {
	return &Router{
		authorService,
	}
}

// Register registers the router to the gin engine
func (authorRouter *Router) Register(router fiber.Router) {
	router.Get(":id", authorRouter.authorByID)
	router.Get("/names/:name", authorRouter.authorByName)
	// Handle query parameter for author name
	router.Get("", authorRouter.authorByQuery)
}

// authorByID fetches an author by ID
func (authorRouter *Router) authorByID(c *fiber.Ctx) error {
	log.Info().Str("id", c.Params("id")).Msg("GET author by ID")
	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error().Err(err).Str("id", c.Params("id")).Msg("Invalid author ID")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	log.Debug().Int("ID", id).Msg("Fetching author by ID")
	result, err := authorRouter.authorService.getAuthorByID(id)
	if err != nil {
		if errors.Is(err, exception.ErrorNotFound) {
			log.Warn().Int("ID", id).Msg("Author not found by ID")
			return c.Status(fiber.StatusNotFound).JSON(resterr.NotFound(err.Error()))
		}

		log.Error().Err(err).Msg("Failed to fetch author by ID")
		return c.Status(fiber.StatusInternalServerError).JSON(
			resterr.InternalServerError("Failed to fetch author by ID"))
	}

	return c.JSON(result)
}

// authorByName fetches an author by name
func (authorRouter *Router) authorByName(c *fiber.Ctx) error {
	name := c.Params("name")
	log.Info().Str("name", name).Msg("GET author by name")
	return funcName(c, name, authorRouter)
}

// authorByQuery fetches an author by name
func (authorRouter *Router) authorByQuery(c *fiber.Ctx) error {
	// Get the name from query parameters
	name := c.Query("name")
	log.Info().Str("name", name).Msg("GET author by query string 'name'")
	return funcName(c, name, authorRouter)
}

func funcName(c *fiber.Ctx, name string, authorRouter *Router) error {
	// Decode the name parameter to handle URL encoding
	decodedName, err := url.QueryUnescape(name)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to decode author name")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest(err.Error()))
	}

	if name == "" {
		log.Error().Str("name", decodedName).Msg("Invalid author name")
		return c.Status(fiber.StatusBadRequest).JSON(resterr.BadRequest("Invalid author name"))
	}

	log.Debug().Str("name", decodedName).Msg("Fetching author by name")
	result, err := authorRouter.authorService.getAuthorByName(decodedName)
	if err != nil {
		if errors.Is(err, exception.ErrorNotFound) {
			log.Warn().Str("name", decodedName).Msg("Author not found by name")
			return c.Status(fiber.StatusNotFound).JSON(
				resterr.NotFound("Author not found by name: " + decodedName))
		}

		log.Error().Err(err).Msg("Failed to fetch author by name")
		return c.Status(fiber.StatusInternalServerError).JSON(
			resterr.InternalServerError("Failed to fetch author by name"))
	}

	return c.JSON(result)
}
