package pet

import (
	"fiber-petclinic-service/pkg/repository"
)

type Request struct {
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	TypeID    int    `json:"typeID" binding:"required"`
	OwnerID   int    `json:"ownerID" binding:"required"`
}

type AddRequest struct {
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	TypeID    int    `json:"typeID" binding:"required"`
	OwnerID   int    `json:"ownerID" binding:"required"`
}

type UpdateRequest struct {
	ID        uint   `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
	TypeID    int    `json:"typeID" binding:"required"`
	OwnerID   int    `json:"ownerID" binding:"required"`
}

func ToPet(petRequest *Request) *repository.Pet {
	//birthday, _ := time.Parse(time.DateOnly, petRequest.Birthdate)
	return &repository.Pet{
		Name:      petRequest.Name,
		Birthdate: petRequest.Birthdate,
		TypeID:    petRequest.TypeID,
		OwnerID:   petRequest.OwnerID,
	}
}
