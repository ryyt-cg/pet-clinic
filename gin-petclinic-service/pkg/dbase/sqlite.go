package dbase

import (
	"github.com/rhtran/gin-petclinic-service/config/app"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SqliteConnect() (*gorm.DB, error) {

	db, err := gorm.Open(sqlite.Open(app.Config.Database.Sqlite.Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// NamingStrategy: schema.NamingStrategy{SingularTable: true},

	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
