package pet

import (
	"fiber3-petclinic-service/pkg/repository"
	"gorm.io/gorm"
	"time"
)

type Request struct {
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	TypeID    uint   `json:"typeID" binding:"required"`
	OwnerID   uint   `json:"ownerID" binding:"required"`
}

type AddRequest struct {
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	TypeID    uint   `json:"typeID" binding:"required"`
	OwnerID   uint   `json:"ownerID" binding:"required"`
}

type UpdateRequest struct {
	ID        uint   `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	TypeID    uint   `json:"typeID" binding:"required"`
	OwnerID   uint   `json:"ownerID" binding:"required"`
}

func ToPet(petRequest *Request) (*repository.Pet, error) {
	birthday, err := time.Parse(time.DateOnly, petRequest.Birthdate)
	if err != nil {
		return nil, err
	}

	petEntity := &repository.Pet{
		Name:      petRequest.Name,
		Birthdate: &birthday,
		TypeID:    petRequest.TypeID,
		OwnerID:   petRequest.OwnerID,
	}

	return petEntity, nil
}

func FromAddRequest(petRequest *AddRequest) (*repository.Pet, error) {
	birthday, err := time.Parse(time.DateOnly, petRequest.Birthdate)
	if err != nil {
		return nil, err
	}

	return &repository.Pet{
		Name:      petRequest.Name,
		Birthdate: &birthday,
		TypeID:    petRequest.TypeID,
		OwnerID:   petRequest.OwnerID,
	}, nil
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
		TypeID:    petRequest.TypeID,
		OwnerID:   petRequest.OwnerID,
	}, nil
}
