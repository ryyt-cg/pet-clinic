package visit

import (
	"fiber-petclinic-service/pkg/repository"
	"gorm.io/gorm"
	"time"
)

type Request struct {
	ID          uint   `json:"id"`
	VisitDate   string `json:"visitDate" binding:"required"`
	Description string `json:"description" binding:"required"`
	PetID       uint   `json:"petID" binding:"required"`
}

type AddRequest struct {
	VisitDate   string `json:"visitDate" binding:"required"`
	Description string `json:"description" binding:"required"`
	PetID       uint   `json:"petID" binding:"required"`
}

type UpdateRequest struct {
	ID          uint   `json:"id"`
	VisitDate   string `json:"visitDate" binding:"required"`
	Description string `json:"description" binding:"required"`
	PetID       uint   `json:"petID" binding:"required"`
}

func FromAddRequest(request *AddRequest) (*repository.Visit, error) {
	visitDate, err := time.Parse(time.DateOnly, request.VisitDate)
	if err != nil {
		return nil, err
	}

	return &repository.Visit{
		VisitDate:   &visitDate,
		Description: request.Description,
		PetID:       request.PetID,
	}, nil
}

func FromUpdateRequest(request *UpdateRequest) (*repository.Visit, error) {
	visitDate, err := time.Parse(time.DateOnly, request.VisitDate)
	if err != nil {
		return nil, err
	}

	return &repository.Visit{
		Model:       gorm.Model{ID: request.ID},
		VisitDate:   &visitDate,
		Description: request.Description,
		PetID:       request.PetID,
	}, nil
}
