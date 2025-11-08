package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewMockPostgresDB
// Mock Postgres Database
func NewMockPostgresDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal().Err(err).Msg("An error was not expected when opening a stub database connection")
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal().Err(err).Msg("An error was not expected when opening gorm database")
	}

	return gormDB, mock
}

// NewMockMySQLDB
// Mock MySQL Database
func NewMockMySQLDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal().Err(err).Msg("An error was not expected when opening a stub database connection")
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal().Err(err).Msg("An error was not expected when opening gorm database")
	}

	return gormDB, mock
}
