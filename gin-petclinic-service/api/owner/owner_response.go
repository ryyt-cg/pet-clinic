package owner

import (
	"gin-petclinic-service/api/pet"
	"gin-petclinic-service/middleware/errors"
	"gin-petclinic-service/pkg/infra/repository"
	"gin-petclinic-service/pkg/model"
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

// FromUpdateEntity - convert owner entity to response
func (ur *UpdateResponse) FromUpdateEntity(owner *repository.Owner) {
	ur.ID = owner.ID
	ur.FirstName = owner.FirstName
	ur.LastName = owner.LastName
	ur.Address = owner.Address
	ur.City = owner.City
	ur.Telephone = owner.Telephone
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

type exceptionResponse struct {
	Code          int                  `json:"code"`
	ErrorResponse errors.ErrorResponse `json:"error"`
}

func (or *Response) FromOwner(owner *repository.Owner) {
	or.ID = owner.ID
	or.FirstName = owner.FirstName
	or.LastName = owner.LastName
	or.Address = owner.Address
	or.City = owner.City
	or.Telephone = owner.Telephone
	or.Pets = pet.FromPets(owner.Pets)
}

func FromOwners(owners []repository.Owner) []Response {
	ownerResponses := make([]Response, len(owners))
	for i, v := range owners {
		ownerResponses[i].FromOwner(&v)
	}
	return ownerResponses
}
