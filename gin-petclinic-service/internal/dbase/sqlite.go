package dbase

import (

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Sqlite struct{}

func (s *Sqlite) Connect() (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger:         logger.Default.LogMode(logger.LogLevel(app.Config.Gorm.LogLevel)),
		NamingStrategy: schema.NamingStrategy{SingularTable: app.Config.Gorm.SingularTable},
	}

	db, err := gorm.Open(sqlite.Open(app.Config.Databases["sqlite"].Name), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}
