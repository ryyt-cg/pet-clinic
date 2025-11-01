package owner

import (
	"fiber-petclinic-service/api/pet"
	"fiber-petclinic-service/internal/repository"
	"fiber-petclinic-service/internal/repository/model"
)

// Owner Responses - A collection of responses (output contracts) for the owner API.
// 1. Composes the validation functions to enforce the input contracts.
// 2. Composes the transformation functions to transform the input contracts into the domain model.
// 3. Composes the error handling functions to handle the errors.

// UpdateResponse - owner add/update response
type UpdateResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Telephone string `json:"telephone"`
}

// ToUpdateResponse
// Map a repository.Owner to UpdateResponse
func ToUpdateResponse(owner *repository.Owner) *UpdateResponse {
	if owner == nil {
		return nil
	}

	return &UpdateResponse{
		ID:        owner.ID,
		FirstName: owner.FirstName,
		LastName:  owner.LastName,
		Address:   owner.Address,
		City:      owner.City,
		Telephone: owner.Telephone,
	}
}

// Response - owner response
type Response struct {
	ID        uint           `json:"id"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Address   string         `json:"address"`
	City      string         `json:"city"`
	Telephone string         `json:"telephone"`
	Pets      []pet.Response `json:"pets,omitempty"`
}

// Responses - list of owners
type Responses struct {
	Context model.Context `json:"context"`
	Owners  []Response    `json:"owners"`
}

// ToResponse
// Map a repository.Owner to Response
func ToResponse(owner *repository.Owner) *Response {
	if owner == nil {
		return nil
	}

	return &Response{
		ID:        owner.ID,
		FirstName: owner.FirstName,
		LastName:  owner.LastName,
		Address:   owner.Address,
		City:      owner.City,
		Telephone: owner.Telephone,
		Pets:      pet.FromPets(owner.Pets),
	}
}

func FromOwners(owners []repository.Owner) []Response {
	ownerResponses := make([]Response, len(owners))
	for i, v := range owners {
		ownerResponses[i] = *ToResponse(&v)
	}
	return ownerResponses
}
