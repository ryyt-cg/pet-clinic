package main

import (
	"context"
	"fiber3-petclinic-service/api/info"
	"fiber3-petclinic-service/api/owner"
	"fiber3-petclinic-service/api/pet"
	"fiber3-petclinic-service/api/ready"
	"fiber3-petclinic-service/api/vet"
	"fiber3-petclinic-service/api/visit"
	"fiber3-petclinic-service/config/app"
	_ "fiber3-petclinic-service/docs"
	"fiber3-petclinic-service/internal/dbase"
	"fiber3-petclinic-service/internal/repository"
	"fiber3-petclinic-service/internal/telemetry"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	otelfiber "github.com/gofiber/contrib/v3/otel"

	"github.com/gofiber/contrib/monitor"
	"github.com/gofiber/contrib/v3/circuitbreaker"
	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	fiberApp *fiber.App
	gdb      *gorm.DB
)

func logConfig() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}

// Instantiate zerolog
// Instantiate fiber router and middlewares
func loadConfig() {
	logConfig()
	log.Info().Msg("Fiber 3 Pet Clinic starts")

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

	var err error
	sqlDB := dbase.Sqlite{}
	gdb, err = sqlDB.Connect(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to connect the database.")
	}

	//sqlDB := dbase.Postgres{}
	//gdb, err = sqlDB.Connect(context.Background())
	//if err != nil {
	//	log.Fatal().Err(err).Msg("Fail to connect the database.")
	//}

	fiberConfig := fiber.Config{
		AppName: "fiber3-petclinic-service",
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

	// Monitor
	fiberApp.Get(app.Config.Server.BaseURL+"/monitor", monitor.New(monitor.Config{Title: "fiber3-petclinic-service Monitor"}))
	fiberApp.Use(otelfiber.Middleware())

	// Create a new Circuit Breaker with custom configuration
	cb := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: app.Config.CircuitBreaker.FailureThreshold,                     // Max failures before opening the circuit
		Timeout:          time.Duration(app.Config.CircuitBreaker.Timeout) * time.Second, // Wait time before retrying
		SuccessThreshold: app.Config.CircuitBreaker.SuccessThreshold,                     // Required successes to move back to closed state
	})
	// Apply Circuit Breaker to ALL routes
	fiberApp.Use(circuitbreaker.Middleware(cb))

	// Middleware for Enforcing Accept only application/json requests
	fiberApp.Use(func(c fiber.Ctx) error {
		if offer := c.Accepts(fiber.MIMEApplicationJSON); offer == "" {
			return c.Status(fiber.StatusNotAcceptable).SendString("Only application/json is accepted.")
		}
		return c.Next()
	})

	// Apply global middlewares
	app.PingRepo = repository.NewPingRepository(gdb)
	fiberApp.Use(recover.New())   // Recover from panics and continue
	fiberApp.Use(requestid.New()) // Generate a unique request ID for each request

	fiberApp.Get(app.Config.Server.BaseURL+"/ready", healthcheck.New(ready.Config()))
	// Mount the UI with the default configuration under /swagger
	fiberApp.Get(app.Config.Server.BaseURL+"/swagger/*", swaggo.HandlerDefault)
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

// @host      	localhost:8092
// @BasePath	/api/pet-clinic/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	loadConfig()
	loadComponents()

	tp := telemetry.InitTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("Error shutting down tracer provider")
		}
	}()

	// Start the server on port define in yaml in a goroutine
	go func() {
		err := fiberApp.Listen(app.Config.Server.HttpPort, fiber.ListenConfig{
			EnablePrefork: app.Config.Server.EnablePrefork,
			CertFile:      app.Config.Server.CertFile,
			CertKeyFile:   app.Config.Server.KeyFile,
		})
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

	log.Info().Msg("Fiber 3 Pet Clinic stops.")
}
