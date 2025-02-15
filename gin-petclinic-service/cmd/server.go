package main

import (
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	health2 "github.com/rhtran/gin-petclinic-service/api/health"
	info2 "github.com/rhtran/gin-petclinic-service/api/info"
	owner2 "github.com/rhtran/gin-petclinic-service/api/owner"
	pet2 "github.com/rhtran/gin-petclinic-service/api/pet"
	vet2 "github.com/rhtran/gin-petclinic-service/api/vet"
	visit2 "github.com/rhtran/gin-petclinic-service/api/visit"
	"github.com/rhtran/gin-petclinic-service/config/app"
	_ "github.com/rhtran/gin-petclinic-service/docs"
	"github.com/rhtran/gin-petclinic-service/middleware"
	"github.com/rhtran/gin-petclinic-service/pkg/dbase"
	"github.com/rhtran/gin-petclinic-service/pkg/ds"
	"github.com/rhtran/gin-petclinic-service/pkg/infra/repository"
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
		logger.Fatal("failed to load application configuration.",
			zap.String("error", err.Error()))
	}

	//pg = dbase.PgConnect()
	var err error
	sqlite, err = dbase.SqliteConnect()
	if err != nil {
		logger.Fatal("failed to connect the database.",
			zap.String("error", err.Error()))
		//os.Exit(-1)
	}

	//ctx, client, err := okta.NewClient(context.TODO(), okta.WithOrgUrl("https://dev-293522.okta.com/"), okta.WithToken("{apiToken}"))

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
	healthCheckService := health2.NewService(logger)
	healthCheckRouter := health2.NewRouter(healthCheckService, logger)

	infoService := info2.NewService(logger)
	ipService := info2.NewIPService(logger)
	infoRouter := info2.NewRouter(logger, infoService, ipService)

	// Owner
	ownerRepository := repository.NewOwnerRepository(logger, sqlite)
	ownerService := owner2.NewService(logger, ownerRepository)
	ownerRouter := owner2.NewRouter(logger, ownerService)

	// Pet
	petRepository := repository.NewPetRepository(logger, sqlite)
	petService := pet2.NewService(logger, petRepository)
	petRouter := pet2.NewRouter(logger, petService)

	// Vet
	vetRepository := repository.NewVetRepository(logger, sqlite)
	vetService := vet2.NewService(logger, vetRepository)
	vetRouter := vet2.NewRouter(logger, vetService)

	// Visit
	visitRepository := repository.NewVisitRepository(logger, sqlite)
	visitService := visit2.NewService(logger, visitRepository)
	visitRouter := visit2.NewRouter(logger, visitService)

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
		logger.Fatal("fail to run http server.", zap.String("error", err.Error()))
		os.Exit(-1)
	}
}
