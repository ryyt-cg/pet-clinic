package owner

import (
	"gin-petclinic-service/internal/repository"
	"gin-petclinic-service/internal/repository/model"

	"gorm.io/gorm"
)

// Owner Requests - A collection of requests (input contracts) for the owner API
// 1. Composes the validation functions to enforce the input contracts.
// 2. Composes the transformation functions to transform the input contracts into the domain model.
// 3. Composes the error handling functions to handle the errors.

type AddRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Telephone string `json:"telephone"`
}

type UpdateRequest struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Telephone string `json:"telephone"`
}

// ToOwnerEntity
// Map an AddRequest to repository.Owner
func ToOwnerEntity(ownerRequest *AddRequest) *repository.Owner {
	if ownerRequest == nil {
		return nil
	}

	return &repository.Owner{
		Person: model.Person{
			FirstName: ownerRequest.FirstName,
			LastName:  ownerRequest.LastName,
		},

		Address:   ownerRequest.Address,
		City:      ownerRequest.City,
		Telephone: ownerRequest.Telephone,
	}
}

// ToOwnerEntityFromUpdateRequest
// Map an UpdateRequest to repository.Owner
func ToOwnerEntityFromUpdateRequest(ownerRequest *UpdateRequest) *repository.Owner {
	if ownerRequest == nil {
		return nil
	}

	return &repository.Owner{
		Model: gorm.Model{
			ID: ownerRequest.ID,
		},
		Person: model.Person{
			FirstName: ownerRequest.FirstName,
			LastName:  ownerRequest.LastName,
		},

		Address:   ownerRequest.Address,
		City:      ownerRequest.City,
		Telephone: ownerRequest.Telephone,
	}
}
