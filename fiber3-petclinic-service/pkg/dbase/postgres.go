package dbase

import (
	"context"
	"gorm.io/gorm"
)

type Database interface {
	Connect(context.Context) (*gorm.DB, error)
}

type Postgres struct {
}

// Connect
// Create connection pooling using GORM postgres driver.
//func Connect(ctx context.Context) (*gorm.DB, error) {
//	log, _ := zap.NewProduction()
//	db, err := gorm.Open(postgres.Open(app.Config.Databases["secondary"].Dsn), &gorm.Config{
//		Logger: logger.Default.LogMode(logger.Info),
//		//NamingStrategy: schema.NamingStrategy{SingularTable: true},
//	})
//	if err != nil {
//		log.Fatal("Error connecting Postgres database.", zap.String("error", err.Error()))
//		panic(err)
//	}
//
//	// connection pooling configuration
//	sqlDB, _ := db.DB()
//	sqlDB.SetMaxIdleConns(app.Config.Database.MaxIdleConns)
//	sqlDB.SetMaxOpenConns(app.Config.Database.MaxOpenConns)
//	sqlDB.SetConnMaxIdleTime(time.Duration(app.Config.Database.MaxIdleTime) * time.Second)
//
//	return db
//}
