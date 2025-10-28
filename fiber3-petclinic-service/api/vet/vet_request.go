package vet

import (
	"fiber3-petclinic-service/internal/repository"
	"fiber3-petclinic-service/internal/repository/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

type AddRequest struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Specialties []uint `json:"specialties" binding:"required"`
}

type UpdateRequest struct {
	ID          uint   `json:"id" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Specialties []uint `json:"specialties" binding:"required"`
}

func (vr AddRequest) Validate() error {
	return validation.ValidateStruct(&vr,
		validation.Field(&vr.FirstName, validation.Required),
		validation.Field(&vr.LastName, validation.Required),
		//validation.Field(&vr.Specialties, validation.Required),
	)
}

// FromAddRequest
// Map an AddRequest to repository.Vet
func FromAddRequest(vetRequest *AddRequest) *repository.Vet {
	return &repository.Vet{
		Person: model.Person{
			FirstName: vetRequest.FirstName,
			LastName:  vetRequest.LastName,
		},
		//Specialties: *ToSpecialties(vetRequest.Specialties),
	}
}

// FromUpdateRequest
// Map an UpdateRequest to repository.Vet
func FromUpdateRequest(vetRequest *UpdateRequest) *repository.Vet {
	return &repository.Vet{
		Model: gorm.Model{
			ID: vetRequest.ID,
		},
		Person: model.Person{
			FirstName: vetRequest.FirstName,
			LastName:  vetRequest.LastName,
		},
		//Specialties: *ToSpecialties(vetRequest.Specialties),
	}
}
