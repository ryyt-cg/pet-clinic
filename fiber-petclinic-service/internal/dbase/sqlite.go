package dbase

import (
	"context"
	"fiber-petclinic-service/config/app"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Sqlite struct {
}

func (s *Sqlite) Connect(ctx context.Context) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger:         logger.Default.LogMode(logger.LogLevel(app.Config.Gorm.LogLevel)),
		NamingStrategy: schema.NamingStrategy{SingularTable: app.Config.Gorm.SingularTable},
	}

	db, err := gorm.Open(sqlite.Open(app.Config.Databases["primary"].Name), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}
