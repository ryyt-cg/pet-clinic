package vet

import (
	"github.com/rhtran/gin-petclinic-service/pkg/infra/repository"
	"gorm.io/gorm"
)

type specialtyResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type specialtiesResponse struct {
	Specialties []specialtyResponse `json:"specialties"`
}

func (s *specialtyResponse) FromSpecialty(specialty *repository.Specialty) {
	s.ID = specialty.ID
	s.Name = specialty.Name
}

func ToSpecialtyResponses(specialties []repository.Specialty) *[]specialtyResponse {
	specialtyResponses := make([]specialtyResponse, len(specialties))
	for i, s := range specialties {
		specialtyResponses[i].ID = s.ID
		specialtyResponses[i].Name = s.Name
	}
	return &specialtyResponses
}

func ToSpecialty(specialty *specialtyResponse) *repository.Specialty {
	return &repository.Specialty{
		Model: gorm.Model{
			ID: specialty.ID,
		},
		Name: specialty.Name,
	}
}

func (s *specialtiesResponse) ToSpecialties() *[]repository.Specialty {
	specialties := make([]repository.Specialty, len(s.Specialties))
	for i, s := range s.Specialties {
		specialties[i] = *ToSpecialty(&s)
	}
	return &specialties
}

func ToSpecialties(specialties []specialtyResponse) *[]repository.Specialty {
	specialtyEntities := make([]repository.Specialty, len(specialties))
	for i, s := range specialties {
		specialtyEntities[i] = *ToSpecialty(&s)
	}

	return &specialtyEntities
}
