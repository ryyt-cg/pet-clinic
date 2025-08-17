package repository

import "gorm.io/gorm"

type PingRepositorier interface {
	Ping() error
}
type PingRepository struct {
	db *gorm.DB
}

// NewPingRepository creates a new instance of PingRepository
func NewPingRepository(db *gorm.DB) *PingRepository {
	return &PingRepository{
		db: db,
	}
}

// Ping checks the database connection by executing a simple query
func (repository *PingRepository) Ping() error {
	if err := repository.db.Exec("SELECT 1").Error; err != nil {
		return err
	}
	return nil
}
