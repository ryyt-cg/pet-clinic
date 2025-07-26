package visit

import (
	"gin-petclinic-service/pkg/repository"
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
	return &repository.Visit{
		VisitDate:   request.VisitDate,
		Description: request.Description,
		PetID:       request.PetID,
	}
}
