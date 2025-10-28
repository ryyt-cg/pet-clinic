package vet

import (
	"gin-petclinic-service/internal/repository"
	"gin-petclinic-service/internal/repository/model"
)

type Response struct {
	ID          uint                 `json:"id"`
	FirstName   string               `json:"firstName"`
	LastName    string               `json:"lastName"`
	Specialties *[]specialtyResponse `json:"specialties,omitempty"`
}

type Responses struct {
	Context model.Context `json:"context"`
	Vets    []Response    `json:"vets"`
}

// FromVet
// Map a Response to repository.Vet
func (vr *Response) FromVet(vet *repository.Vet) {
	vr.ID = vet.ID
	vr.FirstName = vet.FirstName
	vr.LastName = vet.LastName
	vr.Specialties = ToSpecialtyResponses(vet.Specialties)
}

// FromVets
// Map a list of repository.Vet to a list of Response
func FromVets(vets []repository.Vet) []Response {
	if len(vets) == 0 {
		return nil
	}

	vetResponses := make([]Response, len(vets))
	for i, v := range vets {
		vetResponses[i].FromVet(&v)
	}
	return vetResponses
}

// ToResponses
// Map a list of repository.Vet to a Responses
func ToResponses(vets []repository.Vet) *Responses {
	vetsJson := FromVets(vets)
	contextJson := model.Context{Count: len(vetsJson)}
	return &Responses{Vets: vetsJson, Context: contextJson}
}
