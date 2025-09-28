package dbase

import (
	"context"
	"fmt"
	"gin-petclinic-service/config/app"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Database interface {
	Connect(context.Context) (*gorm.DB, error)
}

type Postgres struct {
}

// Connect
// Create connection pooling using GORM postgres driver.
func (p Postgres) Connect(ctx context.Context) (*gorm.DB, error) {
	// data service name (DSN) example:
	// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=UTC"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		app.Config.Databases["postgres"].Host,
		app.Config.Databases["postgres"].Username,
		app.Config.Databases["postgres"].Password,
		app.Config.Databases["postgres"].Name,
		app.Config.Databases["postgres"].Port,
		app.Config.Databases["postgres"].SslMode,
		"UTC") // TimeZone is hardcoded to UTC, can be modified as needed
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.LogLevel(app.Config.Gorm.LogLevel)),
		NamingStrategy: schema.NamingStrategy{SingularTable: app.Config.Gorm.SingularTable},
	})
	if err != nil {
		return nil, err
	}

	// connection pooling configuration
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(app.Config.Databases["postgres"].ConnectionPool.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(app.Config.Databases["postgres"].ConnectionPool.MaxOpenConnections)
	sqlDB.SetConnMaxIdleTime(time.Duration(app.Config.Databases["postgres"].ConnectionPool.MaxIdleTime) * time.Second)

	return db, nil
}
