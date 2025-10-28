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
	SpeciesID uint   `json:"speciesID"`
	OwnerID   uint   `json:"ownerID"`
}

// ToResponse
// Map a repository.Pet to Response
func ToResponse(pet *repository.Pet) *Response {
	if pet == nil {
		return nil
	}

	responseBirthday := ""

	if pet.Birthdate != nil {
		responseBirthday = pet.Birthdate.Format(time.DateOnly)
	}

	return &Response{
		ID:        pet.ID,
		Name:      pet.Name,
		Birthdate: responseBirthday,
		Species:   pet.Species.Name,
		Visits:    visit.FromVisits(pet.Visits),
	}
}

func FromPets(pets []repository.Pet) []Response {
	if len(pets) == 0 {
		return nil
	}
	
	petResponses := make([]Response, len(pets))
	for i, pet := range pets {
		petResponses[i] = *ToResponse(&pet)
	}
	return petResponses
}

func ToResponses(pets []repository.Pet) *Responses {
	petResponses := FromPets(pets)
	contextJson := model.Context{Count: len(petResponses)}
	return &Responses{Pets: petResponses, Context: contextJson}
}
