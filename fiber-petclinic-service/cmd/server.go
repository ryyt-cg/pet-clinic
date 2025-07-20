package main

import (
	"fiber-petclinic-service/api/author"
	"fiber-petclinic-service/api/info"
	"fiber-petclinic-service/config/app"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/circuitbreaker"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

var (
	fiberApp *fiber.App
)

// Instantiate zerolog
// Instantiate fiber router and middlewares
func loadConfig() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Go Fiber 01 starts")

	// load application configurations
	if err := app.LoadConfig("./config"); err != nil {
		log.Fatal().Err(err).
			Msg("Fail to load application configuration.")
	}

	// Set the log level based on the configuration
	switch app.Config.Server.LogLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	// Create a new Fiber instance
	fiberApp = fiber.New()

	prometheus := fiberprometheus.New("fiber-petclinic-service")
	prometheus.RegisterAt(fiberApp, app.Config.Server.BaseURL+"/metrics")
	prometheus.SetSkipPaths([]string{"/ping"}) // Optional: Remove some paths from metrics
	fiberApp.Use(prometheus.Middleware)

	// Monitor
	fiberApp.Get(app.Config.Server.BaseURL+"/monitor", monitor.New(monitor.Config{Title: "fiber-petclinic-service Monitor"}))

	cb := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: app.Config.CircuitBreaker.FailureThreshold,
		Timeout:          time.Duration(app.Config.CircuitBreaker.Timeout) * time.Second,
		SuccessThreshold: app.Config.CircuitBreaker.SuccessThreshold,
	})

	fiberApp.Use(circuitbreaker.Middleware(cb))
	// Middleware for Enforcing Accept only application/json requests
	fiberApp.Use(func(c *fiber.Ctx) error {
		if offer := c.Accepts(fiber.MIMEApplicationJSON); offer == "" {
			return c.Status(fiber.StatusNotAcceptable).SendString("Only application/json is accepted.")
		}
		return c.Next()
	})
	// Apply global middlewares
	fiberApp.Use(healthcheck.New())
	fiberApp.Use(recover.New())   // Recover from panics and continue
	fiberApp.Use(requestid.New()) // Generate a unique request ID for each request
}

func loadComponents() {
	infoService := info.NewService()
	infoRouter := info.NewRouter(infoService)

	authorService := author.NewService()
	authorRouter := author.NewRouter(authorService)

	// create a new group for the /api/gof endpoint
	home := fiberApp.Group(app.Config.Server.BaseURL)
	// Register the info router to the home group
	infoRouter.Register(home.Group("/info"))

	// create a new group for the /api/gof/v1 endpoint
	v1 := fiberApp.Group(app.Config.Server.BaseURL + "/v1")
	// Register the author router to the v1 group
	authorRouter.Register(v1.Group("/authors"))
}

func main() {
	loadConfig()
	loadComponents()

	// Start the server on port 3000
	err := fiberApp.Listen(app.Config.Server.HttpPort)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server")
		return
	}
}
