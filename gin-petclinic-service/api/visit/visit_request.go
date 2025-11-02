package visit

import (
	"gin-petclinic-service/internal/repository"
	"time"

	"gorm.io/gorm"
)

type request struct {
	ID          uint   `json:"id"`
	VisitDate   string `json:"visitDate" binding:"required"`
	Description string `json:"description" binding:"required"`
	PetID       uint   `json:"petID" binding:"required"`
}

type addRequest struct {
	VisitDate   string `json:"visitDate" binding:"required"`
	Description string `json:"description" binding:"required"`
	PetID       uint   `json:"petID" binding:"required"`
}

type updateRequest struct {
	ID          uint   `json:"id"`
	VisitDate   string `json:"visitDate" binding:"required"`
	Description string `json:"description" binding:"required"`
	PetID       uint   `json:"petID" binding:"required"`
}

// fromAddRequest
// Map an addRequest to repository.Visit
func fromAddRequest(request *addRequest) (*repository.Visit, error) {
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

// fromUpdateRequest
// Map a updateRequest to repository.Visit
func fromUpdateRequest(request *updateRequest) (*repository.Visit, error) {
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
