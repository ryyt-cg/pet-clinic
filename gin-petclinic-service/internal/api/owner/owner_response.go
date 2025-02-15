package owner

import (
	"github.com/rhtran/gin-petclinic-service/internal/api/pet"
	"github.com/rhtran/gin-petclinic-service/pkg/infra/repository"
	"github.com/rhtran/gin-petclinic-service/pkg/model"
)

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
