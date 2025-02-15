package vet

import (
	"github.com/rhtran/gin-petclinic-service/pkg/infra/repository"
	"github.com/rhtran/gin-petclinic-service/pkg/model"
)

type Response struct {
	ID          uint                `json:"id"`
	FirstName   string              `json:"firstName"`
	LastName    string              `json:"lastName"`
	Specialties []specialtyResponse `json:"specialties,omitempty"`
}

type Responses struct {
	Context model.Context `json:"context"`
	Vets    []Response    `json:"vets"`
}

func (vr *Response) FromVet(vet *repository.Vet) {
	vr.ID = vet.ID
	vr.FirstName = vet.FirstName
	vr.LastName = vet.LastName
	vr.Specialties = *ToSpecialtyResponses(vet.Specialties)
}

func FromVets(vets []repository.Vet) []Response {
	vetResponses := make([]Response, len(vets))
	for i, v := range vets {
		vetResponses[i].FromVet(&v)
	}
	return vetResponses
}
