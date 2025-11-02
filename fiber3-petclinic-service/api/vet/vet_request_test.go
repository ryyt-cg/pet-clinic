package vet

import (
	"fiber3-petclinic-service/internal/repository"
	"fiber3-petclinic-service/internal/repository/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_AddRequestValidation(t *testing.T) {
	vetRequest := &addRequest{
		FirstName: "John",
		LastName:  "Doe",
		//Specialties: []specialtyResponse{
		//	{ID: 1, Name: "Cardiology"},
		//	{ID: 2, Name: "Dentistry"},
		//},
	}

	err := vetRequest.Validate()

	assert.Nil(t, err)
}

func Test_AddRequestValidationWithEmptyFields(t *testing.T) {
	vetRequest := &addRequest{
		FirstName: "",
		LastName:  "",
		//Specialties: []Specialty{},
	}

	err := vetRequest.Validate()

	assert.NotNil(t, err)
}

func Test_FromAddRequest(t *testing.T) {
	vetRequest := &addRequest{
		FirstName: "John",
		LastName:  "Johnson",
		//Specialties: []Specialty{
		//	{ID: 1, Name: "Cardiology"},
		//	{ID: 2, Name: "Dentistry"},
		//},
	}
	expectedVet := &repository.Vet{
		Person: model.Person{
			FirstName: "John",
			LastName:  "Johnson",
		},
	}

	vet := fromAddRequest(vetRequest)
	assert.Equal(t, expectedVet, vet)
}

func Test_FromUpdateRequest(t *testing.T) {
	vetRequest := &updateRequest{
		ID:        1,
		FirstName: "John",
		LastName:  "Johnson",
		//Specialties: []Specialty{
		//	{ID: 1, Name: "Cardiology"},
		//	{ID: 2, Name: "Dentistry"},
		//},
	}
	expectedVet := &repository.Vet{
		Model: gorm.Model{
			ID: 1,
		},
		Person: model.Person{
			FirstName: "John",
			LastName:  "Johnson",
		},
	}

	vet := fromUpdateRequest(vetRequest)
	assert.Equal(t, expectedVet, vet)
}
