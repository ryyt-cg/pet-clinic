package visit

import (
	"errors"
	"fiber-petclinic-service/pkg/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func Test_getById(t *testing.T) {
	testCases := []struct {
		name          string
		id            uint
		expectedVisit *repository.Visit
		expectedError error
	}{
		{
			name: "Test get visit by id success",
			id:   1,
			expectedVisit: &repository.Visit{
				Model: gorm.Model{
					ID: 1,
				},
				//VisitDate:   "2023-03-05",
				Description: "rabies shot",
			},
			expectedError: nil,
		},
		{
			name:          "Test get visit by id with error",
			id:            2,
			expectedVisit: nil,
			expectedError: errors.New("getById: unable to find visit by id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			visitMock := repository.NewMockVisitRepositorier(t)

			visitMock.EXPECT().FindById(tc.id).Return(tc.expectedVisit, tc.expectedError)
			visitService := NewService(visitMock)
			result, err := visitService.getVisitById(1)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedVisit.ID, result.ID)
				assert.Equal(t, tc.expectedVisit.Description, result.Description)
				//assert.Equal(t, tc.expectedVisit.VisitDate, result.VisitDate)
				assert.Nil(t, err)
			}

			visitMock.AssertExpectations(t)
			visitMock.AssertNumberOfCalls(t, "FindById", 1)
		})
	}
}

// Test_getAllVisits tests the getAllVisits function
func Test_getAllVisits(t *testing.T) {
	testCases := []struct {
		name           string
		expectedVisits []repository.Visit
		expectedError  error
	}{
		{
			name: "Test get all visits success",
			expectedVisits: []repository.Visit{
				{
					Model: gorm.Model{
						ID: 1,
					},
					//VisitDate:   "2023-03-05",
					Description: "rabies shot",
				},
				{
					Model: gorm.Model{
						ID: 2,
					},
					//VisitDate:   "2023-03-05",
					Description: "rabies shot",
				},
			},
			expectedError: nil,
		},
		{
			name:           "Test get all visits with error",
			expectedVisits: nil,
			expectedError:  errors.New("getAllVisits: unable to find all visits"),
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			visitMock := repository.NewMockVisitRepositorier(t)

			visitMock.EXPECT().FindAll().Return(tc.expectedVisits, tc.expectedError)
			visitService := NewService(visitMock)
			result, err := visitService.getAllVisits()

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				//for i, visit := range tc.expectedVisits {
				//	assert.Equal(t, visit.ID, result[i].ID)
				//	assert.Equal(t, visit.Description, result[i].Description)
				//	assert.Equal(t, visit.VisitDate, result[i].VisitDate)
				//	assert.Nil(t, err)
				//}
			}

			visitMock.AssertExpectations(t)
			visitMock.AssertNumberOfCalls(t, "FindAll", 1)
		})
	}
}
