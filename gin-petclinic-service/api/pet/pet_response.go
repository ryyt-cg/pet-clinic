package pet

import (
	"gin-petclinic-service/api/visit"
	"gin-petclinic-service/pkg/infra/repository"
	"gin-petclinic-service/pkg/model"
)

type Response struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	Birthdate string           `json:"birthdate"`
	Type      string           `json:"type"`
	Visits    []visit.Response `json:"visits,omitempty"`
}

type Responses struct {
	Context model.Context `json:"context"`
	Pets    []Response    `json:"pets"`
}

type UpdateResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
	TypeID    uint   `json:"typeID"`
	OwnerID   uint   `json:"ownerID"`
}

type AddResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
	TypeID    uint   `json:"typeID"`
	OwnerID   uint   `json:"ownerID"`
}

func (pr *Response) FromPet(pet *repository.Pet) {
	pr.ID = pet.ID
	pr.Name = pet.Name
	pr.Birthdate = pet.Birthdate
	pr.Type = pet.Type.Name
	pr.Visits = visit.FromVisits(pet.Visits)
}

func FromPets(pets []repository.Pet) []Response {
	petResponses := make([]Response, len(pets))
	for i, v := range pets {
		petResponses[i].FromPet(&v)
	}
	return petResponses
}
