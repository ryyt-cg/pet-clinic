package visit

import (
	"fiber-petclinic-service/internal/repository"
	"fiber-petclinic-service/internal/repository/model"
	"fiber-petclinic-service/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestResponse_FromVisit(t *testing.T) {
	testCases := []struct {
		fromVisit  *repository.Visit
		toResponse Response
	}{
		{fromVisit: &repository.Visit{
			Model:       gorm.Model{ID: 1},
			VisitDate:   test.ToDate("2023-07-12"),
			Description: "Flu Shot",
			PetID:       1,
		}, toResponse: Response{
			ID:          1,
			VisitDate:   "2023-07-12",
			Description: "Flu Shot",
			PetID:       1},
		},
		{fromVisit: &repository.Visit{
			Model:       gorm.Model{ID: 2},
			Description: "Flu Shot",
			PetID:       1,
		}, toResponse: Response{
			ID:          2,
			Description: "Flu Shot",
			PetID:       1},
		},
	}

	for _, tc := range testCases {
		result := &Response{}
		result.FromVisit(tc.fromVisit)
		assert.Equal(t, tc.toResponse, *result)
	}

}

// Test_FromVisits
// Return nil when fromResponses is nil or empty slice.
func Test_FromVisits(t *testing.T) {
	mockVisits := []repository.Visit{
		{Model: gorm.Model{ID: 10}, Description: "Flu Shot", VisitDate: test.ToDate("2023-07-12"), PetID: 1},
		{Model: gorm.Model{ID: 20}, Description: "Regular Checkup", VisitDate: test.ToDate("2023-07-15"), PetID: 2},
	}
	mockResponses := []Response{
		{ID: 10, Description: "Flu Shot", VisitDate: "2023-07-12", PetID: 1},
		{ID: 20, Description: "Regular Checkup", VisitDate: "2023-07-15", PetID: 2},
	}

	testCases := []struct {
		fromVisits  []repository.Visit
		toResponses []Response
	}{
		{
			fromVisits:  mockVisits,
			toResponses: mockResponses,
		},
		{
			fromVisits:  []repository.Visit{},
			toResponses: nil,
		},
		{
			fromVisits:  nil,
			toResponses: nil,
		},
	}

	for _, tc := range testCases {
		result := FromVisits(tc.fromVisits)
		assert.Equal(t, tc.toResponses, result)
	}
}

func Test_ToResponses(t *testing.T) {
	mockResponses := []Response{
		{ID: 10, Description: "Flu Shot", VisitDate: "2023-07-12", PetID: 1},
		{ID: 20, Description: "Regular Checkup", VisitDate: "2023-07-15", PetID: 2},
	}

	testCases := []struct {
		fromResponses []Response
		toResponses   *Responses
	}{
		{
			fromResponses: mockResponses,
			toResponses: &Responses{
				Context: model.Context{Count: 2},
				Visits:  mockResponses,
			},
		},
		{
			fromResponses: nil,
			toResponses: &Responses{
				Context: model.Context{Count: 0},
				Visits:  nil,
			},
		},
	}

	for _, tc := range testCases {
		result := ToResponses(tc.fromResponses)
		assert.Equal(t, tc.toResponses, result)
	}
}
