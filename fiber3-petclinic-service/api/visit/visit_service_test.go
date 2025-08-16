package visit

import (
	"errors"
	"fiber3-petclinic-service/internal/repository"
	"fiber3-petclinic-service/internal/repository/model"
	"fiber3-petclinic-service/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_getById(t *testing.T) {
	mockVisit := &repository.Visit{
		Model: gorm.Model{
			ID: 1,
		},
		VisitDate:   test.ToDate("2023-03-05"),
		Description: "rabies shot",
		PetID:       101,
	}
	expectedResponse := &Response{
		ID:          1,
		VisitDate:   "2023-03-05",
		Description: "rabies shot",
		PetID:       101,
	}

	testCases := []struct {
		name           string
		id             uint
		mockVisit      *repository.Visit
		mockError      error
		expectedResult interface{}
	}{
		{
			name:           "get visit by id",
			id:             1,
			mockVisit:      mockVisit,
			mockError:      nil,
			expectedResult: expectedResponse,
		},
		{
			name:           "get not visit by id",
			id:             2,
			mockVisit:      nil,
			mockError:      gorm.ErrRecordNotFound,
			expectedResult: gorm.ErrRecordNotFound,
		},
		{
			name:           "fail to get visit by id",
			id:             2,
			mockVisit:      nil,
			mockError:      errors.New("getById: unable to find visit by id"),
			expectedResult: errors.New("getById: unable to find visit by id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			visitMock := repository.NewMockVisitRepositorier(t)

			visitMock.EXPECT().FindById(tc.id).Return(tc.mockVisit, tc.mockError)
			visitService := NewService(visitMock)
			result, err := visitService.getVisitById(tc.id)

			if err != nil {
				assert.Equal(t, tc.mockError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedResult.(*Response), result)
				assert.Nil(t, err)
			}

			visitMock.AssertExpectations(t)
			visitMock.AssertNumberOfCalls(t, "FindById", 1)
		})
	}
}

// Test_getAllVisits tests the getAllVisits function
func Test_getAllVisits(t *testing.T) {
	mockVisits := []repository.Visit{
		{
			Model: gorm.Model{
				ID: 1,
			},
			VisitDate:   test.ToDate("2023-03-05"),
			Description: "rabies shot",
			PetID:       201,
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			//VisitDate:   "2023-03-05",
			Description: "rabies shot",
			PetID:       202,
		},
	}
	expectedVisits := &Responses{
		Context: model.Context{
			Count: 2,
		},
		Visits: FromVisits(mockVisits),
	}

	testCases := []struct {
		name           string
		mockVisits     []repository.Visit
		mockError      error
		expectedVisits interface{}
	}{
		{
			name:           "get all visits",
			mockVisits:     mockVisits,
			mockError:      nil,
			expectedVisits: expectedVisits,
		},
		{
			name:           "get no visits",
			mockVisits:     nil,
			mockError:      gorm.ErrRecordNotFound,
			expectedVisits: gorm.ErrRecordNotFound,
		},
		{
			name:           "fail get all visits",
			mockVisits:     nil,
			mockError:      errors.New("getAllVisits: unable to find all visits"),
			expectedVisits: errors.New("getAllVisits: unable to find all visits"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			visitMock := repository.NewMockVisitRepositorier(t)

			visitMock.EXPECT().FindAll().Return(tc.mockVisits, tc.mockError)
			visitService := NewService(visitMock)
			result, err := visitService.getAllVisits()

			if err != nil {
				assert.Equal(t, tc.expectedVisits, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedVisits.(*Responses), result)
				assert.Nil(t, err)
			}
			visitMock.AssertExpectations(t)
			visitMock.AssertNumberOfCalls(t, "FindAll", 1)
		})
	}
}

func Test_create(t *testing.T) {
	mockVisit := &repository.Visit{
		Model: gorm.Model{
			ID: 1,
		},
		VisitDate:   test.ToDate("2023-03-05"),
		Description: "rabies shot",
		PetID:       101,
	}
	expectedResponse := &Response{
		ID:          1,
		VisitDate:   "2023-03-05",
		Description: "rabies shot",
		PetID:       101,
	}

	testCases := []struct {
		name           string
		mockVisit      *repository.Visit
		mockError      error
		expectedResult interface{}
	}{
		{
			name:           "create visit",
			mockVisit:      mockVisit,
			mockError:      nil,
			expectedResult: expectedResponse,
		},
		{
			name:           "fail to create visit",
			mockVisit:      mockVisit,
			mockError:      errors.New("create: unable to create visit"),
			expectedResult: errors.New("create: unable to create visit"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			visitMock := repository.NewMockVisitRepositorier(t)
			if tc.mockError != nil {
				visitMock.EXPECT().Insert(tc.mockVisit).Return(nil, tc.mockError)
			} else {
				visitMock.EXPECT().Insert(tc.mockVisit).Return(tc.mockVisit, tc.mockError)
			}
			visitMock.EXPECT().Insert(tc.mockVisit).Return(tc.mockVisit, tc.mockError)

			visitService := NewService(visitMock)
			result, err := visitService.create(tc.mockVisit)

			if err != nil {
				assert.Equal(t, tc.expectedResult, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedResult.(*Response), result)
				assert.Nil(t, err)
			}
			visitMock.AssertExpectations(t)
			visitMock.AssertNumberOfCalls(t, "Insert", 1)
		})
	}
}

func Test_update(t *testing.T) {
	mockVisit := &repository.Visit{
		Model: gorm.Model{
			ID: 1,
		},
		VisitDate:   test.ToDate("2023-03-05"),
		Description: "rabies shot",
		PetID:       101,
	}
	expectedResponse := &Response{
		ID:          1,
		VisitDate:   "2023-03-05",
		Description: "rabies shot",
		PetID:       101,
	}

	testCases := []struct {
		name           string
		mockVisit      *repository.Visit
		mockError      error
		expectedResult interface{}
	}{
		{
			name:           "update visit",
			mockVisit:      mockVisit,
			mockError:      nil,
			expectedResult: expectedResponse,
		},
		{
			name:           "get no visits",
			mockVisit:      mockVisit,
			mockError:      gorm.ErrRecordNotFound,
			expectedResult: gorm.ErrRecordNotFound,
		},
		{
			name:           "fail to update visit",
			mockVisit:      mockVisit,
			mockError:      errors.New("update: unable to update visit"),
			expectedResult: errors.New("update: unable to update visit"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			visitMock := repository.NewMockVisitRepositorier(t)
			if tc.mockError != nil {
				visitMock.EXPECT().Update(tc.mockVisit).Return(nil, tc.mockError)
			} else {
				visitMock.EXPECT().Update(tc.mockVisit).Return(tc.mockVisit, tc.mockError)
			}
			visitService := NewService(visitMock)
			result, err := visitService.update(tc.mockVisit)

			if err != nil {
				assert.Equal(t, tc.expectedResult, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedResult.(*Response), result)
				assert.Nil(t, err)
			}
			visitMock.AssertExpectations(t)
			visitMock.AssertNumberOfCalls(t, "Update", 1)
		})
	}
}
