package dbase

import (
	"context"
	"gin-petclinic-service/config/app"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	Connect(context.Context) (*gorm.DB, error)
}

type Postgres struct {
}

// Connect
// Create connection pooling using GORM postgres driver.
func (p *Postgres) Connect(ctx context.Context) *gorm.DB {
	log, _ := zap.NewProduction()
	db, err := gorm.Open(postgres.Open(app.Config.Databases["postgres"].Username), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		//NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		log.Fatal("Error connecting Postgres database.", zap.String("error", err.Error()))
		panic(err)
	}

	// connection pooling configuration
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(app.Config.Databases["postgres"].ConnectionPool.MaxIdleConnection)
	sqlDB.SetMaxOpenConns(app.Config.Databases["postgres"].ConnectionPool.MaxOpenConnection)
	sqlDB.SetConnMaxIdleTime(time.Duration(app.Config.Databases["postgres"].ConnectionPool.MaxIdleTime) * time.Second)

	return db
}
