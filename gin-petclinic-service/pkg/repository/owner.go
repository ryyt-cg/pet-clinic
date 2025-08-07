package repository

import (
	"gin-petclinic-service/pkg/repository/model"
	"gorm.io/gorm"
)

/*
Owner has many Pets, OwnerID is the foreign key
Has Many configuration
*/
type Owner struct {
	gorm.Model
	model.Person
	Address   string
	City      string
	Telephone string
	Pets      []Pet
}
