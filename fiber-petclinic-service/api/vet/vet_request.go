package vet

import (
	"fiber-petclinic-service/pkg/repository"
	"fiber-petclinic-service/pkg/repository/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Request struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	//Specialties []Specialty `json:"specialties" binding:"required"`
}

func (vr Request) Validate() error {
	return validation.ValidateStruct(&vr,
		validation.Field(&vr.FirstName, validation.Required),
		validation.Field(&vr.LastName, validation.Required),
		//validation.Field(&vr.Specialties, validation.Required),
	)
}

func ToVet(vetRequest *Request) *repository.Vet {
	return &repository.Vet{
		Person: model.Person{
			FirstName: vetRequest.FirstName,
			LastName:  vetRequest.LastName,
		},
		//Specialties: *ToSpecialties(vetRequest.Specialties),
	}
}
