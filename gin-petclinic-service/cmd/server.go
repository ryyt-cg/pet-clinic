package main

import (
	"gin-petclinic-service/api/health"
	"gin-petclinic-service/api/info"
	"gin-petclinic-service/api/owner"
	"gin-petclinic-service/api/pet"
	"gin-petclinic-service/api/vet"
	"gin-petclinic-service/api/visit"
	"gin-petclinic-service/config/app"
	_ "gin-petclinic-service/docs"
	"gin-petclinic-service/middleware"
	"gin-petclinic-service/pkg/dbase"
	"gin-petclinic-service/pkg/ds"
	"gin-petclinic-service/pkg/infra/repository"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"os"
)

var (
	g      errgroup.Group
	sqlite *gorm.DB
	logger *zap.Logger
	r      *gin.Engine
)

func loadConfig() {
	logger, _ = zap.NewProduction()
	logger.Info("pet-clinic starts")

	// load application configurations
	if err := app.LoadConfig("./config"); err != nil {
		logger.Fatal("Fail to load application configuration.",
			zap.String("error", err.Error()))
	}

	//pg = dbase.PgConnect()
	var err error
	sqlite, err = dbase.SqliteConnect()
	if err != nil {
		logger.Fatal("Fail to connect the database.",
			zap.String("error", err.Error()))
		//os.Exit(-1)
	}

	// Creates a router without any middleware by default
	r = gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default, gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.Use(location.Default())
	r.Use(middleware.SetRequestUUID())
}

// Component initialization
func loadComponents() {
	logger.Info("Component initialization starts")
	healthCheckService := health.NewService(logger)
	healthCheckRouter := health.NewRouter(healthCheckService, logger)

	infoService := info.NewService(logger)
	ipService := info.NewIPService(logger)
	infoRouter := info.NewRouter(logger, infoService, ipService)

	// Owner
	ownerRepository := repository.NewOwnerRepository(logger, sqlite)
	ownerService := owner.NewService(logger, ownerRepository)
	ownerRouter := owner.NewRouter(logger, ownerService)

	// Pet
	petRepository := repository.NewPetRepository(logger, sqlite)
	petService := pet.NewService(logger, petRepository)
	petRouter := pet.NewRouter(logger, petService)

	// Vet
	vetRepository := repository.NewVetRepository(logger, sqlite)
	vetService := vet.NewService(logger, vetRepository)
	vetRouter := vet.NewRouter(logger, vetService)

	// Visit
	visitRepository := repository.NewVisitRepository(logger, sqlite)
	visitService := visit.NewService(logger, visitRepository)
	visitRouter := visit.NewRouter(logger, visitService)

	//authenService := service.NewAuthenService(logger)

	// config prometheus & endpoint group
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
	home := r.Group("/")
	v1 := r.Group("/v1")
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

//	@Schemes	https
//  @Host		localhost:8443
//	@BasePath	/v1

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	loadConfig()
	loadComponents()
	// add swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//r.Run(":8080")
	httpServer := ds.NewHttpServer(r)
	httpRouter := httpServer.HttpRouter()

	g.Go(func() error {
		return httpRouter.ListenAndServeTLS(app.Config.Server.CertFile, app.Config.Server.KeyFile)
	})

	if err := g.Wait(); err != nil {
		logger.Fatal("Fail to run http server.", zap.String("error", err.Error()))
		os.Exit(-1)
	}
}
