package repository

import (
	"gorm.io/gorm"
	"time"
)

/*
Belongs To
A belongs to association sets up a one-to-one connection with another model,
Pet

Notes:
	can not add this: Owner     owner.Owner `gorm:"foreignKey:OwnerID"` because
	Owner struct has Pets attribute relationship.  Therefore, it creates
	circular relationship and cycle relationship not allowed
*/
// Pet type belongs to Type, TypeID is the foreign key
type Pet struct {
	gorm.Model
	Name      string    `gorm:"column:name"`
	Birthdate time.Time //`gorm:"column:birth_date"`
	TypeID    int
	OwnerID   int
	Type      Type `gorm:"foreignKey:TypeID"`
	Visits    []Visit
}

type Type struct {
	gorm.Model
	Name string
}
