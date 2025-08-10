package main

import (
	"context"
	"fiber-petclinic-service/api/info"
	"fiber-petclinic-service/api/owner"
	"fiber-petclinic-service/api/pet"
	"fiber-petclinic-service/api/vet"
	"fiber-petclinic-service/api/visit"
	"fiber-petclinic-service/config/app"
	"fiber-petclinic-service/pkg/dbase"
	"fiber-petclinic-service/pkg/repository"
	"os"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/circuitbreaker"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	fiberApp *fiber.App
	sqlite   *gorm.DB
)

// Instantiate zerolog
// Instantiate fiber router and middlewares
func loadConfig() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	log.Info().Msg("Go Fiber Pet Clinic starts")

	// load application configurations
	if err := app.LoadConfig("config"); err != nil {
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

	var err error
	sqlDB := dbase.Sqlite{}
	sqlite, err = sqlDB.Connect(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to connect the database.")
	}

	// fiberzerolog config
	fiberLog := fiberzerolog.Config{
		Logger: &logger,
	}

	fiberConfig := fiber.Config{
		AppName:           "fiber-petclinic-service",
		EnablePrintRoutes: app.Config.Server.EnablePrintRoutes,
	}
	// Create a new Fiber instance
	fiberApp = fiber.New(fiberConfig)

	prometheus := fiberprometheus.New("fiber-petclinic-service")
	prometheus.RegisterAt(fiberApp, app.Config.Server.BaseURL+"/metrics")
	prometheus.SetSkipPaths([]string{"/ping"}) // Optional: Remove some paths from metrics
	fiberApp.Use(prometheus.Middleware)

	// Monitor
	fiberApp.Get(app.Config.Server.BaseURL+"/monitor", monitor.New(monitor.Config{Title: "fiber-petclinic-service Monitor"}))
	fiberApp.Use(fiberzerolog.New(fiberLog))

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
	ipService := info.NewIPService()
	infoRouter := info.NewRouter(infoService, ipService)

	// Owner
	ownerRepository := repository.NewOwnerRepository(sqlite)
	ownerService := owner.NewService(ownerRepository)
	ownerRouter := owner.NewRouter(ownerService)

	// Pet
	petRepository := repository.NewPetRepository(sqlite)
	petService := pet.NewService(petRepository)
	petRouter := pet.NewRouter(petService)

	// Vet
	vetRepository := repository.NewVetRepository(sqlite)
	vetService := vet.NewService(vetRepository)
	vetRouter := vet.NewRouter(vetService)

	// Visit
	visitRepository := repository.NewVisitRepository(sqlite)
	visitService := visit.NewService(visitRepository)
	visitRouter := visit.NewRouter(visitService)

	// create a new group for the /api/gof endpoint
	home := fiberApp.Group(app.Config.Server.BaseURL)
	// Register the info router to the home group
	infoRouter.Register(home.Group("/info"))

	// create a new group for the /api/gof/v1 endpoint
	v1 := fiberApp.Group(app.Config.Server.BaseURL + "/v1")
	// Register the author router to the v1 group
	ownerRouter.Register(v1.Group("/owners"))
	petRouter.Register(v1.Group("/pets"))
	vetRouter.Register(v1.Group("/vets"))
	visitRouter.Register(v1.Group("/visits"))
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
