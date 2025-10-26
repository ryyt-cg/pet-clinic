package visit

import (
	"gin-petclinic-service/internal/repository"
	"gin-petclinic-service/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

//var request []Request = []Request{
//	{
//		ID:          1,
//		VisitDate:   "2024-12-05",
//		Description: "Nail Clip",
//		PetID:       21,
//	},
//	{
//		ID:          50,
//		VisitDate:   "2025-01-25",
//		Description: "Flu Shot",
//		PetID:       25,
//	},
//}

func Test_FromAddRequest(t *testing.T) {
	testCases := []struct {
		fromAddRequest *AddRequest
		toVisit        *repository.Visit
	}{
		{
			fromAddRequest: &AddRequest{
				VisitDate:   "2024-12-05",
				Description: "Nail Clip",
				PetID:       21,
			},
			toVisit: &repository.Visit{
				VisitDate:   test.ToDate("2024-12-05"),
				Description: "Nail Clip",
				PetID:       21,
			},
		},
		{
			fromAddRequest: &AddRequest{
				VisitDate:   "2025-03-32",
				Description: "Regular Checkup",
				PetID:       42,
			},
			toVisit: &repository.Visit{
				VisitDate:   test.ToDate("2025-03-32"),
				Description: "Regular Checkup",
				PetID:       42,
			},
		},
	}

	for _, tc := range testCases {
		result, err := FromAddRequest(tc.fromAddRequest)

		if err != nil {
			assert.Nil(t, result)
			assert.Equal(t, "parsing time \"2025-03-32\": day out of range", err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.toVisit, result)
		}
	}
}

func Test_FromUpdateRequest(t *testing.T) {
	testCases := []struct {
		fromUpdateRequest *UpdateRequest
		toVisit           *repository.Visit
	}{
		{
			fromUpdateRequest: &UpdateRequest{
				ID:          150,
				VisitDate:   "2024-12-05",
				Description: "Nail Clip",
				PetID:       21,
			},
			toVisit: &repository.Visit{
				Model:       gorm.Model{ID: 150},
				VisitDate:   test.ToDate("2024-12-05"),
				Description: "Nail Clip",
				PetID:       21,
			},
		},
		{
			fromUpdateRequest: &UpdateRequest{
				ID:          151,
				VisitDate:   "2025-03-32",
				Description: "Regular Checkup",
				PetID:       42,
			},
			toVisit: &repository.Visit{
				Model:       gorm.Model{ID: 151},
				VisitDate:   test.ToDate("2025-03-32"),
				Description: "Regular Checkup",
				PetID:       42,
			},
		},
	}

	for _, tc := range testCases {
		result, err := FromUpdateRequest(tc.fromUpdateRequest)

		if err != nil {
			assert.Nil(t, result)
			assert.Equal(t, "parsing time \"2025-03-32\": day out of range", err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.toVisit, result)
		}
	}
}
