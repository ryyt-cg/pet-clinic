package repository

import (
	"gorm.io/gorm"
	"time"
)

// Visit
// represents a visit made by a pet to the veterinarian.
// entity:model visit
type Visit struct {
	gorm.Model
	VisitDate   *time.Time // return nil if no visit date presents
	Description string
	PetID       uint
	Pet         Pet
}
