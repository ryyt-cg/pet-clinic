package visit

import (
	"fiber3-petclinic-service/internal/repository"
	"time"

	"gorm.io/gorm"
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

// FromAddRequest
// Map an AddRequest to repository.Visit
func FromAddRequest(request *AddRequest) (*repository.Visit, error) {
	visitEntify := &repository.Visit{
		Description: request.Description,
		PetID:       request.PetID,
	}

	visitDate, err := time.Parse(time.DateOnly, request.VisitDate)
	if err != nil {
		return nil, err
	}

	visitEntify.VisitDate = &visitDate
	return visitEntify, nil
}

// FromUpdateRequest
// Map a UpdateRequest to repository.Visit
func FromUpdateRequest(request *UpdateRequest) (*repository.Visit, error) {
	visitEntify := &repository.Visit{
		Model:       gorm.Model{ID: request.ID},
		Description: request.Description,
		PetID:       request.PetID,
	}

	visitDate, err := time.Parse(time.DateOnly, request.VisitDate)
	if err != nil {
		return nil, err
	}

	visitEntify.VisitDate = &visitDate
	return visitEntify, nil
}
