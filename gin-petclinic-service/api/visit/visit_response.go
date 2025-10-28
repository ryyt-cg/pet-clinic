package visit

import (
	"gin-petclinic-service/internal/repository"
	"gin-petclinic-service/internal/repository/model"
	"time"
)

type Response struct {
	ID          uint   `json:"id"`
	VisitDate   string `json:"visitDate"`
	Description string `json:"description"`
	PetID       uint   `json:"petID"`
}

type UpdateResponse struct {
	ID          uint   `json:"id"`
	VisitDate   string `json:"visitDate"`
	Description string `json:"description"`
	PetID       uint   `json:"petID"`
}

type Responses struct {
	Context model.Context `json:"context"`
	Visits  []Response    `json:"visits"`
}

// FromVisit
// Map from repository.Visit to json Response
func (vr *Response) FromVisit(visit *repository.Visit) {
	vr.ID = visit.ID
	if visit.VisitDate == nil {
		vr.VisitDate = ""
	} else {
		vr.VisitDate = visit.VisitDate.Format(time.DateOnly)
	}

	vr.Description = visit.Description
	vr.PetID = visit.PetID
}

// FromVisits
// Map from repository.Visit list to Response list
func FromVisits(visits []repository.Visit) []Response {
	visitResponses := make([]Response, len(visits))
	for i, v := range visits {
		visitResponses[i].FromVisit(&v)
	}

	return visitResponses
}

// FromVisitsToResponses
// Map from repository.Visit list to Responses
func FromVisitsToResponses(visits []repository.Visit) *Responses {
	if len(visits) == 0 {
		return &Responses{
			Context: model.Context{
				Count: len(visits),
			},
			Visits: nil,
		}
	}

	visitResponses := make([]Response, len(visits))
	for i, v := range visits {
		visitResponses[i].FromVisit(&v)
	}

	responses := &Responses{
		Context: model.Context{
			Count: len(visits),
		},
		Visits: visitResponses,
	}
	return responses
}

// ToResponses
// Map list of Response to Responses
func ToResponses(responses []Response) *Responses {
	if len(responses) == 0 {
		return &Responses{
			Context: model.Context{
				Count: len(responses),
			},
			Visits: nil,
		}
	}

	result := make([]Response, len(responses))
	for i, v := range responses {
		result[i].ID = v.ID
		result[i].VisitDate = v.VisitDate
		result[i].Description = v.Description
		result[i].PetID = v.PetID
	}

	return &Responses{
		Context: model.Context{
			Count: len(responses),
		},
		Visits: result,
	}
}
