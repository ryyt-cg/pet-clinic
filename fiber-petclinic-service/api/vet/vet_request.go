package vet

import (
	"fiber-petclinic-service/internal/repository"
	"fiber-petclinic-service/internal/repository/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

type addRequest struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Specialties []uint `json:"specialties" binding:"required"`
}

type updateRequest struct {
	ID          uint   `json:"id" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Specialties []uint `json:"specialties" binding:"required"`
}

func (vr addRequest) Validate() error {
	return validation.ValidateStruct(&vr,
		validation.Field(&vr.FirstName, validation.Required),
		validation.Field(&vr.LastName, validation.Required),
		//validation.Field(&vr.Specialties, validation.Required),
	)
}

// fromAddRequest
// Map an addRequest to repository.Vet
func fromAddRequest(vetRequest *addRequest) *repository.Vet {
	return &repository.Vet{
		Person: model.Person{
			FirstName: vetRequest.FirstName,
			LastName:  vetRequest.LastName,
		},
		//Specialties: *toSpecialties(vetRequest.Specialties),
	}
}

// fromUpdateRequest
// Map an updateRequest to repository.Vet
func fromUpdateRequest(vetRequest *updateRequest) *repository.Vet {
	return &repository.Vet{
		Model: gorm.Model{
			ID: vetRequest.ID,
		},
		Person: model.Person{
			FirstName: vetRequest.FirstName,
			LastName:  vetRequest.LastName,
		},
		//Specialties: *toSpecialties(vetRequest.Specialties),
	}
}
