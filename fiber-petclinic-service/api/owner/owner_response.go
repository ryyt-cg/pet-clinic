package owner

import (
	"fiber-petclinic-service/internal/repository"
	"fiber-petclinic-service/internal/repository/model"
	"time"
)

// Owner Responses - A collection of responses (output contracts) for the owner API.
// 1. Composes the validation functions to enforce the input contracts.
// 2. Composes the transformation functions to transform the input contracts into the domain model.
// 3. Composes the error handling functions to handle the errors.

// UpdateResponse - owner add/update response
type updateResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Telephone string `json:"telephone"`
}

// toUpdateResponse
// Map a repository.Owner to UpdateResponse
func toUpdateResponse(owner *repository.Owner) *updateResponse {
	if owner == nil {
		return nil
	}

	return &updateResponse{
		ID:        owner.ID,
		FirstName: owner.FirstName,
		LastName:  owner.LastName,
		Address:   owner.Address,
		City:      owner.City,
		Telephone: owner.Telephone,
	}
}

// Response - owner response
type response struct {
	ID        uint          `json:"id"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Address   string        `json:"address"`
	City      string        `json:"city"`
	Telephone string        `json:"telephone"`
	Pets      []petResponse `json:"pets,omitempty"`
}

type petResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
	Species   string `json:"species"`
}

// Responses - list of owners
type responses struct {
	Context model.Context `json:"context"`
	Owners  []response    `json:"owners"`
}

// toResponse
// Map a repository.Owner to Response
func toResponse(owner *repository.Owner) *response {
	if owner == nil {
		return nil
	}

	return &response{
		ID:        owner.ID,
		FirstName: owner.FirstName,
		LastName:  owner.LastName,
		Address:   owner.Address,
		City:      owner.City,
		Telephone: owner.Telephone,
		Pets:      toPetResponses(owner.Pets),
	}
}

// toResponse
// Map a repository.Owner to Response
func toResponses(owner *repository.Owner) *response {
	if owner == nil {
		return nil
	}

	return &response{
		ID:        owner.ID,
		FirstName: owner.FirstName,
		LastName:  owner.LastName,
		Address:   owner.Address,
		City:      owner.City,
		Telephone: owner.Telephone,
		Pets:      toPetResponses(owner.Pets),
	}
}

func fromOwners(owners []repository.Owner) []response {
	ownerResponses := make([]response, len(owners))
	for i, v := range owners {
		ownerResponses[i] = *toResponse(&v)
	}

	return ownerResponses
}

// petResponse

// ToPetResponse
// Map a repository.Pet to Response
func toPetResponse(pet *repository.Pet) *petResponse {
	if pet == nil {
		return nil
	}

	var responseBirthday string
	if pet.Birthdate != nil {
		responseBirthday = pet.Birthdate.Format(time.DateOnly)
	}

	return &petResponse{
		ID:        pet.ID,
		Name:      pet.Name,
		Birthdate: responseBirthday,
		Species:   pet.Species.Name,
	}
}

func toPetResponses(pets []repository.Pet) []petResponse {
	if len(pets) == 0 {
		return nil
	}

	petResponses := make([]petResponse, len(pets))
	for i, pet := range pets {
		petResponses[i] = *toPetResponse(&pet)
	}
	return petResponses
}
