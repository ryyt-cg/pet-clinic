package visit

import (
	"fiber-petclinic-service/internal/repository"
	"fiber-petclinic-service/internal/repository/model"
	"time"
)

type response struct {
	ID          uint   `json:"id"`
	VisitDate   string `json:"visitDate"`
	Description string `json:"description"`
	PetID       uint   `json:"petID"`
}

type updateResponse struct {
	ID          uint   `json:"id"`
	VisitDate   string `json:"visitDate"`
	Description string `json:"description"`
	PetID       uint   `json:"petID"`
}

type responses struct {
	Context model.Context `json:"context"`
	Visits  []response    `json:"visits"`
}

// fromVisit
// Map from repository.Visit to json response
func (vr *response) fromVisit(visit *repository.Visit) {
	vr.ID = visit.ID
	if visit.VisitDate == nil {
		vr.VisitDate = ""
	} else {
		vr.VisitDate = visit.VisitDate.Format(time.DateOnly)
	}

	vr.Description = visit.Description
	vr.PetID = visit.PetID
}

// fromVisits
// Map from repository.Visit list to response list
func fromVisits(visits []repository.Visit) []response {
	if len(visits) == 0 {
		return nil
	}

	visitResponses := make([]response, len(visits))
	for i, v := range visits {
		visitResponses[i].fromVisit(&v)
	}

	return visitResponses
}

// fromVisitsToResponses
// Map from repository.Visit list to responses
func fromVisitsToResponses(visits []repository.Visit) *responses {
	if len(visits) == 0 {
		return &responses{
			Context: model.Context{
				Count: len(visits),
			},
			Visits: nil,
		}
	}

	visitResponses := make([]response, len(visits))
	for i, v := range visits {
		visitResponses[i].fromVisit(&v)
	}

	responses := &responses{
		Context: model.Context{
			Count: len(visits),
		},
		Visits: visitResponses,
	}
	return responses
}

// toResponses
// Map list of response to responses
func toResponses(visitResponses []response) *responses {
	if len(visitResponses) == 0 {
		return &responses{
			Context: model.Context{
				Count: len(visitResponses),
			},
			Visits: nil,
		}
	}

	result := make([]response, len(visitResponses))
	for i, v := range visitResponses {
		result[i].ID = v.ID
		result[i].VisitDate = v.VisitDate
		result[i].Description = v.Description
		result[i].PetID = v.PetID
	}

	return &responses{
		Context: model.Context{
			Count: len(visitResponses),
		},
		Visits: result,
	}
}
