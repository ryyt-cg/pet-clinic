package visit

import (
	"fiber-petclinic-service/pkg/repository"
	"time"
)

type Request struct {
	ID          int    `json:"id"`
	VisitDate   string `json:"visitDate" binding:"required"`
	Description string `json:"description" binding:"required"`
	PetID       int    `json:"petId" binding:"required"`
}

type AddRequest struct {
	VisitDate   string `json:"visitDate" binding:"required"`
	Description string `json:"description" binding:"required"`
	PetID       int    `json:"petId" binding:"required"`
}

type UpdateRequest struct {
	ID          int    `json:"id"`
	VisitDate   string `json:"visitDate" binding:"required"`
	Description string `json:"description" binding:"required"`
	PetID       int    `json:"petId" binding:"required"`
}

func ToVisit(request *Request) *repository.Visit {
	visitDate, _ := time.Parse(time.DateOnly, request.VisitDate)
	return &repository.Visit{
		VisitDate:   visitDate,
		Description: request.Description,
		PetID:       request.PetID,
	}
}
