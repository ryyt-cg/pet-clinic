package vet

import (
	"fiber3-petclinic-service/internal/repository"
	"fiber3-petclinic-service/internal/repository/model"

	"gorm.io/gorm"
)

type specialtyResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type specialtiesResponse struct {
	Context     model.Context       `json:"context"`
	Specialties []specialtyResponse `json:"specialties"`
}

func (s *specialtyResponse) FromSpecialty(specialty *repository.Specialty) {
	s.ID = specialty.ID
	s.Name = specialty.Name
}

func toSpecialtyResponses(specialties []repository.Specialty) *[]specialtyResponse {
	if len(specialties) == 0 {
		return nil
	}

	specialtyResponses := make([]specialtyResponse, len(specialties))
	for i, s := range specialties {
		specialtyResponses[i].ID = s.ID
		specialtyResponses[i].Name = s.Name
	}
	return &specialtyResponses
}

func toSpecialtiesResponses(specialties []repository.Specialty) *specialtiesResponse {
	specialtyResponses := make([]specialtyResponse, len(specialties))
	for i, s := range specialties {
		specialtyResponses[i].ID = s.ID
		specialtyResponses[i].Name = s.Name
	}
	specialtiesResponse := &specialtiesResponse{
		Context:     model.Context{Count: len(specialtyResponses)},
		Specialties: specialtyResponses,
	}
	return specialtiesResponse
}

func toSpecialty(specialty *specialtyResponse) *repository.Specialty {
	return &repository.Specialty{
		Model: gorm.Model{
			ID: specialty.ID,
		},
		Name: specialty.Name,
	}
}

func toSpecialties(specialties []specialtyResponse) *[]repository.Specialty {
	specialtyEntities := make([]repository.Specialty, len(specialties))
	for i, s := range specialties {
		specialtyEntities[i] = *toSpecialty(&s)
	}

	return &specialtyEntities
}
