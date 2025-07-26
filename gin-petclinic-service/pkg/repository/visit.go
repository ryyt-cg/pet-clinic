package repository

import (
	"gorm.io/gorm"
)

// Visit
// represents a visit made by a pet to the veterinarian.
// entity:model visit
type Visit struct {
	gorm.Model
	VisitDate   string
	Description string
	PetID       int
	Pet         Pet
}
