package vet

import (
	"fiber3-petclinic-service/internal/repository"
	"fiber3-petclinic-service/internal/repository/model"
)

type response struct {
	ID          uint                 `json:"id"`
	FirstName   string               `json:"firstName"`
	LastName    string               `json:"lastName"`
	Specialties *[]specialtyResponse `json:"specialties,omitempty"`
}

type responses struct {
	Context model.Context `json:"context"`
	Vets    []response    `json:"vets"`
}

// fromVet
// Map a response to repository.Vet
func (vr *response) fromVet(vet *repository.Vet) {
	vr.ID = vet.ID
	vr.FirstName = vet.FirstName
	vr.LastName = vet.LastName
	vr.Specialties = toSpecialtyResponses(vet.Specialties)
}

// fromVets
// Map a list of repository.Vet to a list of response
func fromVets(vets []repository.Vet) []response {
	if len(vets) == 0 {
		return nil
	}

	vetResponses := make([]response, len(vets))
	for i, v := range vets {
		vetResponses[i].fromVet(&v)
	}
	return vetResponses
}

// toResponses
// Map a list of repository.Vet to a responses
func toResponses(vets []repository.Vet) *responses {
	vetsJson := fromVets(vets)
	contextJson := model.Context{Count: len(vetsJson)}
	return &responses{Vets: vetsJson, Context: contextJson}
}
