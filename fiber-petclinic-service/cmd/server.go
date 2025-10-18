package main

import (
	"context"
	"fiber-petclinic-service/api/info"
	"fiber-petclinic-service/api/owner"
	"fiber-petclinic-service/api/pet"
	"fiber-petclinic-service/api/ready"
	"fiber-petclinic-service/api/vet"
	"fiber-petclinic-service/api/visit"
	"fiber-petclinic-service/config/app"
	"fiber-petclinic-service/internal/dbase"
	"fiber-petclinic-service/internal/repository"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "fiber-petclinic-service/docs"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/circuitbreaker"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	fiberApp *fiber.App
	gdb      *gorm.DB
)

// Instantiate zerolog
// Instantiate fiber router and middlewares
func loadConfig() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Info().Msg("Fiber Pet Clinic starts")

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
	//sqlDB := dbase.Postgres{}
	sqlDB := dbase.Sqlite{}
	gdb, err = sqlDB.Connect(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to connect the database.")
	}

	// fiberzerolog config
	fiberLog := fiberzerolog.Config{
		Logger: &log.Logger,
	}

	fiberConfig := fiber.Config{
		Prefork:           true,
		AppName:           "fiber-petclinic-service",
		EnablePrintRoutes: app.Config.Server.EnablePrintRoutes,
	}
	// Create a new Fiber instance
	fiberApp = fiber.New(fiberConfig)

	// Use the CORS middleware
	fiberApp.Use(cors.New())
	// Initialize default config
	fiberApp.Use(favicon.New())

	// Enable Compression
	if app.Config.Server.EnableCompression {
		fiberApp.Use(compress.New(compress.Config{
			Level: compress.Level(app.Config.Server.CompressionLevel),
		}))
	}

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
	app.PingRepo = repository.NewPingRepository(gdb)
	fiberApp.Use(healthcheck.New(ready.Config()))
	fiberApp.Use(recover.New())   // Recover from panics and continue
	fiberApp.Use(requestid.New()) // Generate a unique request ID for each request
	fiberApp.Get(app.Config.Server.BaseURL+"/swagger/*", swagger.HandlerDefault)
}

func loadComponents() {
	infoService := info.NewService()
	ipService := info.NewIPService()
	infoRouter := info.NewRouter(infoService, ipService)

	// Owner
	ownerRepository := repository.NewOwnerRepository(gdb)
	ownerService := owner.NewService(ownerRepository)
	ownerRouter := owner.NewRouter(ownerService)

	// Pet
	petRepository := repository.NewPetRepository(gdb)
	petService := pet.NewService(petRepository)
	petRouter := pet.NewRouter(petService)

	// Vet
	vetRepository := repository.NewVetRepository(gdb)
	vetService := vet.NewService(vetRepository)
	vetRouter := vet.NewRouter(vetService)

	// Visit
	visitRepository := repository.NewVisitRepository(gdb)
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

// @title           Pet Clinic API
// @version         1.0
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8091
// @BasePath  /api/pet-clinic/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	loadConfig()
	loadComponents()

	// Start the server on port define in yaml in a goroutine
	go func() {
		err := fiberApp.ListenTLS(app.Config.Server.HttpPort, app.Config.Server.CertFile, app.Config.Server.KeyFile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to start the server")
			return
		}
	}()

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-quit

	log.Info().Msg("Shutting down server gracefully...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the Fiber app
	if err := fiberApp.ShutdownWithContext(ctx); err != nil {
		log.Error().Err(err).Msg("error during server shutdown")
	}

	log.Info().Msg("Fiber Pet Clinic stops.")
}
