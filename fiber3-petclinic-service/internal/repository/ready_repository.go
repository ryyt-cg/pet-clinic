package repository

import "gorm.io/gorm"

type ReadyRepositorier interface {
	// GetReadyStatus returns the readiness status of the service.
	GetReadyStatus() (bool, error)
}

// ReadyRepository searches owner from the database
type ReadyRepository struct {
	db *gorm.DB
}

// NewReadyRepository - OwnerRepository factory
func NewReadyRepository(db *gorm.DB) *ReadyRepository {
	return &ReadyRepository{
		db: db,
	}
}

// GetReadyStatus returns the readiness status of the service.
func (repository *ReadyRepository) GetReadyStatus() (bool, error) {
	var count int64
	err := repository.db.Model(&Owner{}).Count(&count).Error
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}
