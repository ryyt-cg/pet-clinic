package main

import (
	"context"
	"gin-petclinic-service/api/health"
	"gin-petclinic-service/api/info"
	"gin-petclinic-service/api/owner"
	"gin-petclinic-service/api/pet"
	"gin-petclinic-service/api/vet"
	"gin-petclinic-service/api/visit"

	//"gin-petclinic-service/api/pet"
	//"gin-petclinic-service/api/vet"
	//"gin-petclinic-service/api/visit"
	"gin-petclinic-service/config/app"
	_ "gin-petclinic-service/docs"
	"gin-petclinic-service/internal/dbase"
	"gin-petclinic-service/internal/ds"
	"gin-petclinic-service/internal/repository"
	"gin-petclinic-service/middleware"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

var (
	g  errgroup.Group
	db *gorm.DB
	r  *gin.Engine
)

func loadConfig() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("gin-pet-clinic starts")

	// load application configurations
	if err := app.LoadConfig("./config"); err != nil {
		log.Fatal().Err(err).Msg("Fail to load application configuration.")
	}

	sqlite := dbase.Sqlite{}
	var err error
	db, err = sqlite.Connect(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to connect the database.")
	}

	// Creates a router without any middleware by default
	r = gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default, gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())
	r.Use(middleware.JSONContentType())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.Use(location.Default())
	r.Use(middleware.SetRequestUUID())
}

// Component initialization
func loadComponents() {
	log.Info().Msg("Component initialization starts")
	healthCheckService := health.NewService()
	healthCheckRouter := health.NewRouter(healthCheckService)

	infoService := info.NewService()
	ipService := info.NewIPService()
	infoRouter := info.NewRouter(infoService, ipService)

	// Owner
	ownerRepository := repository.NewOwnerRepository(db)
	ownerService := owner.NewService(ownerRepository)
	ownerRouter := owner.NewRouter(ownerService)

	// Pet
	petRepository := repository.NewPetRepository(db)
	petService := pet.NewService(petRepository)
	petRouter := pet.NewRouter(petService)

	// Vet
	vetRepository := repository.NewVetRepository(db)
	vetService := vet.NewService(vetRepository)
	vetRouter := vet.NewRouter(vetService)

	// Visit
	visitRepository := repository.NewVisitRepository(db)
	visitService := visit.NewService(visitRepository)
	visitRouter := visit.NewRouter(visitService)

	//authenService := service.NewAuthenService(logger)

	// config prometheus & endpoint group
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
	home := r.Group(app.Config.Server.BaseURL)
	v1 := r.Group(app.Config.Server.BaseURL + "/v1")
	//v1.Use(middleware.Authenticate(authenService))

	healthCheckRouter.Register(home.Group("/health"))
	infoRouter.Register(home.Group("/info"))
	ownerRouter.Register(v1.Group("/owners"))
	petRouter.Register(v1.Group("/pets"))
	vetRouter.Register(v1.Group("/vets"))
	visitRouter.Register(v1.Group("/visits"))
}

//	@title			Pet Clinic API
//	@version		1.0
//	@description	This is a pet clinic API server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@Schemas	https
//  @Host		localhost:8092
//	@BasePath	/api/pet-clinic/v1

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	loadConfig()
	loadComponents()
	// add swagger endpoint
	r.GET(app.Config.Server.BaseURL+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//r.Run(":8080")
	httpServer := ds.NewHttpServer(r)
	httpRouter := httpServer.HttpRouter()

	g.Go(func() error {
		return httpRouter.ListenAndServeTLS(app.Config.Server.CertFile, app.Config.Server.KeyFile)
		//return httpRouter.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal().Err(err).Msg("Fail to run http server.")
		os.Exit(-1)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpRouter.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server Shutdown:")
	}
	log.Info().Msg("Server exiting")
}
