package repository

import (
	"time"

	"gorm.io/gorm"
)

// Pet species belongs to Species, SpeciesID is the foreign key
// Belongs To
// A belongs to association sets up a one-to-one connection with another model, Pet
//
// Notes:
//
//	can not add this: Owner     owner.Owner `gorm:"foreignKey:OwnerID"` because
//	Owner struct has Pets attribute relationship.  Therefore, it creates
//	circular relationship and cycle relationship not allowed
type Pet struct {
	gorm.Model
	Name      string     `gorm:"column:name"`
	Birthdate *time.Time //`gorm:"column:birth_date"`
	SpeciesID uint
	OwnerID   uint
	Species   Species `gorm:"foreignKey:SpeciesID"`
	Visits    []Visit
}

type Species struct {
	gorm.Model
	Name string
}
