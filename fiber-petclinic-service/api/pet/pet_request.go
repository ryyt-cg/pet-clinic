package pet

import (
	"fiber-petclinic-service/internal/repository"
	"time"

	"gorm.io/gorm"
)

type Request struct {
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	SpeciesID uint   `json:"speciesID" binding:"required"`
	OwnerID   uint   `json:"ownerID" binding:"required"`
}

type AddRequest struct {
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	SpeciesID uint   `json:"speciesID" binding:"required"`
	OwnerID   uint   `json:"ownerID" binding:"required"`
}

type UpdateRequest struct {
	ID        uint   `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	SpeciesID uint   `json:"speciesID" binding:"required"`
	OwnerID   uint   `json:"ownerID" binding:"required"`
}

func ToPet(petRequest *Request) (*repository.Pet, error) {
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

func FromAddRequest(petRequest *AddRequest) *repository.Pet {
	birthday, err := time.Parse(time.DateOnly, petRequest.Birthdate)
	if err != nil {
		return nil
	}

	return &repository.Pet{
		Name:      petRequest.Name,
		Birthdate: &birthday,
		SpeciesID: petRequest.SpeciesID,
		OwnerID:   petRequest.OwnerID,
	}
}

func FromUpdateRequest(petRequest *UpdateRequest) (*repository.Pet, error) {
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
