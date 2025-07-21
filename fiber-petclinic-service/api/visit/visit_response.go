package visit

import (
	"fiber-petclinic-service/pkg/repository"
	"fiber-petclinic-service/pkg/repository/model"
)

type Response struct {
	ID          uint   `json:"id"`
	VisitDate   string `json:"visitDate"`
	Description string `json:"description"`
}

type UpdateResponse struct {
	ID          uint   `json:"id"`
	VisitDate   string `json:"visitDate"`
	Description string `json:"description"`
	PetID       uint   `json:"petId"`
}

type Responses struct {
	Context model.Context `json:"context"`
	Visits  []Response    `json:"visits"`
}

func (vr *Response) FromVisit(visit *repository.Visit) {
	vr.ID = visit.ID
	// Format(time.DateOnly)
	vr.VisitDate = visit.VisitDate
	vr.Description = visit.Description
}

func FromVisits(visits []repository.Visit) []Response {
	visitResponses := make([]Response, len(visits))
	for i, v := range visits {
		visitResponses[i].FromVisit(&v)
	}
	return visitResponses
}
