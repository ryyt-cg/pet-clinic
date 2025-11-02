package pet

import (
	"fiber-petclinic-service/internal/repository"
	"fiber-petclinic-service/internal/repository/model"
	"time"
)

type ownerResponse struct {
}

type visitResponse struct {
	ID          uint   `json:"id"`
	VisitDate   string `json:"visitDate"`
	Description string `json:"description"`
	PetID       uint   `json:"petID"`
}

type response struct {
	ID        uint            `json:"id"`
	Name      string          `json:"name"`
	Birthdate string          `json:"birthdate"`
	Species   string          `json:"species"`
	Visits    []visitResponse `json:"visits,omitempty"`
}

type responses struct {
	Context model.Context `json:"context"`
	Pets    []response    `json:"pets"`
}

type updateResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
	SpeciesID uint   `json:"speciesID"`
	OwnerID   uint   `json:"ownerID"`
}

type addResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
	SpeciesID uint   `json:"speciesID"`
	OwnerID   uint   `json:"ownerID"`
}

// toResponse
// Map a repository.Pet to Response
func toResponse(pet *repository.Pet) *response {
	if pet == nil {
		return nil
	}

	var responseBirthday string
	if pet.Birthdate != nil {
		responseBirthday = pet.Birthdate.Format(time.DateOnly)
	}

	return &response{
		ID:        pet.ID,
		Name:      pet.Name,
		Birthdate: responseBirthday,
		Species:   pet.Species.Name,
		Visits:    toVisitResponses(pet.Visits),
	}
}

// toADDResponse
// Map a repository.Pet to addResponse
func toAddResponse(pet *repository.Pet) *addResponse {
	if pet == nil {
		return nil
	}

	var responseBirthday string
	if pet.Birthdate != nil {
		responseBirthday = pet.Birthdate.Format(time.DateOnly)
	}

	return &addResponse{
		ID:        pet.ID,
		Name:      pet.Name,
		Birthdate: responseBirthday,
		SpeciesID: pet.Species.ID,
		OwnerID:   pet.OwnerID,
	}
}

// toUpdateResponse
// Map a repository.Pet to updateResponse
func toUpdateResponse(pet *repository.Pet) *updateResponse {
	if pet == nil {
		return nil
	}

	var responseBirthday string
	if pet.Birthdate != nil {
		responseBirthday = pet.Birthdate.Format(time.DateOnly)
	}

	return &updateResponse{
		ID:        pet.ID,
		Name:      pet.Name,
		Birthdate: responseBirthday,
		SpeciesID: pet.Species.ID,
		OwnerID:   pet.OwnerID,
	}
}

func fromPets(pets []repository.Pet) []response {
	if len(pets) == 0 {
		return nil
	}

	petResponses := make([]response, len(pets))
	for i, pet := range pets {
		petResponses[i] = *toResponse(&pet)
	}
	return petResponses
}

// toResponses
// Map list of repository.Pet to responses
func toResponses(pets []repository.Pet) *responses {
	petResponses := fromPets(pets)
	contextJson := model.Context{Count: len(petResponses)}
	return &responses{Pets: petResponses, Context: contextJson}
}

func toVisitResponse(visit repository.Visit) *visitResponse {
	return &visitResponse{
		ID:          visit.ID,
		VisitDate:   visit.VisitDate.Format(time.DateOnly),
		Description: visit.Description,
		PetID:       visit.PetID,
	}
}

func toVisitResponses(visits []repository.Visit) []visitResponse {
	if len(visits) == 0 {
		return nil
	}

	visitResponses := make([]visitResponse, len(visits))
	for i, v := range visits {
		visitResponses[i] = *toVisitResponse(v)
	}

	return visitResponses
}
