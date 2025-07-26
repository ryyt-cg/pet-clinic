package visit

import (
	"errors"
	repository2 "gin-petclinic-service/pkg/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"testing"
)

func Test_getById(t *testing.T) {
	testCases := []struct {
		name          string
		expectedVisit *repository2.Visit
		expectedError error
	}{
		{
			name: "Test get visit by id success",
			expectedVisit: &repository2.Visit{
				Model: gorm.Model{
					ID: 1,
				},
				VisitDate:   "2023-03-05",
				Description: "rabies shot",
			},
			expectedError: nil,
		},
		{
			name:          "Test get visit by id with error",
			expectedVisit: nil,
			expectedError: errors.New("getById: unable to find visit by id"),
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := zap.NewNop()
			//logger := log.New().With(nil, "function", "Test_getById")
			visitMock := repository2.MockVisitRepositorier{}

			visitMock.On("FindById", 1).Return(tc.expectedVisit, tc.expectedError)
			visitService := NewService(logger, &visitMock)
			result, err := visitService.getVisitById(1)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedVisit.ID, result.ID)
				assert.Equal(t, tc.expectedVisit.Description, result.Description)
				assert.Equal(t, tc.expectedVisit.VisitDate, result.VisitDate)
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
		expectedVisits []repository2.Visit
		expectedError  error
	}{
		{
			name: "Test get all visits success",
			expectedVisits: []repository2.Visit{
				{
					Model: gorm.Model{
						ID: 1,
					},
					VisitDate:   "2023-03-05",
					Description: "rabies shot",
				},
				{
					Model: gorm.Model{
						ID: 2,
					},
					VisitDate:   "2023-03-05",
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
			//logger := log.New().With(nil, "function", "Test_getAllVisits")
			logger := zap.NewNop()
			visitMock := repository2.MockVisitRepositorier{}

			visitMock.On("FindAll").Return(tc.expectedVisits, tc.expectedError)
			visitService := NewService(logger, &visitMock)
			result, err := visitService.getAllVisits()

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				for i, visit := range tc.expectedVisits {
					assert.Equal(t, visit.ID, result[i].ID)
					assert.Equal(t, visit.Description, result[i].Description)
					assert.Equal(t, visit.VisitDate, result[i].VisitDate)
					assert.Nil(t, err)
				}
			}

			visitMock.AssertExpectations(t)
			visitMock.AssertNumberOfCalls(t, "FindAll", 1)
		})
	}
}
