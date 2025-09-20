package pet

import (
	"fiber3-petclinic-service/api/visit"
	"fiber3-petclinic-service/internal/repository"
	"fiber3-petclinic-service/internal/repository/model"
	"time"
)

type Response struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	Birthdate string           `json:"birthdate"`
	Species   string           `json:"species"`
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

	if pet.Birthdate == nil {
		pr.Birthdate = ""
	} else {
		// Format the birthdate as a string in the format "YYYY-MM-DD"
		pr.Birthdate = pet.Birthdate.Format(time.DateOnly)
	}
	pr.Species = pet.Species.Name
	pr.Visits = visit.FromVisits(pet.Visits)
}

func FromPets(pets []repository.Pet) []Response {
	petResponses := make([]Response, len(pets))
	for i, v := range pets {
		petResponses[i].FromPet(&v)
	}
	return petResponses
}

func ToResponses(pets []repository.Pet) *Responses {
	petResponses := FromPets(pets)
	contextJson := model.Context{Count: len(petResponses)}
	return &Responses{Pets: petResponses, Context: contextJson}
}
