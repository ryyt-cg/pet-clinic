package pet

import (
	"fiber3-petclinic-service/internal/repository"
	"time"

	"gorm.io/gorm"
)

type request struct {
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	SpeciesID uint   `json:"speciesID" binding:"required"`
	OwnerID   uint   `json:"ownerID" binding:"required"`
}

type addRequest struct {
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	SpeciesID uint   `json:"speciesID" binding:"required"`
	OwnerID   uint   `json:"ownerID" binding:"required"`
}

type updateRequest struct {
	ID        uint   `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	SpeciesID uint   `json:"speciesID" binding:"required"`
	OwnerID   uint   `json:"ownerID" binding:"required"`
}

// requestSchema
// Validate Pet request payload
//var requestSchema = z.Struct(z.Shape{
//	"name":      z.String().Required(),
//	"birthdate": z.String(),
//	"speciesID": z.Uint(),
//	"ownerID":   z.Uint(),
//})

// toPet
// Map Request to repository.Pet
func toPet(petRequest *request) (*repository.Pet, error) {
	// TODO need to add data validation
	birthday, err := time.Parse(time.DateOnly, petRequest.Birthdate)
	if err != nil {
		return nil, err
	}

	petEntity := &repository.Pet{
		Name:      petRequest.Name,
		Birthdate: &birthday,
		SpeciesID: petRequest.SpeciesID,
		OwnerID:   petRequest.OwnerID,
	}

	return petEntity, nil
}

// fromAddRequest
// Map AddRequest to repository.Pet
func fromAddRequest(petRequest *addRequest) (*repository.Pet, error) {
	if petRequest == nil {
		return nil, nil
	}

	birthday, err := time.Parse(time.DateOnly, petRequest.Birthdate)
	if err != nil {
		return nil, err
	}

	return &repository.Pet{
		Name:      petRequest.Name,
		Birthdate: &birthday,
		SpeciesID: petRequest.SpeciesID,
		OwnerID:   petRequest.OwnerID,
	}, nil
}

// fromUpdateRequest
// Map UpdateRequest to repository.Pet
func fromUpdateRequest(petRequest *updateRequest) (*repository.Pet, error) {
	birthday, err := time.Parse(time.DateOnly, petRequest.Birthdate)
	if err != nil {
		return nil, err
	}

	return &repository.Pet{
		Model: gorm.Model{
			ID: petRequest.ID,
		},
		Name:      petRequest.Name,
		Birthdate: &birthday,
		SpeciesID: petRequest.SpeciesID,
		OwnerID:   petRequest.OwnerID,
	}, nil
}
