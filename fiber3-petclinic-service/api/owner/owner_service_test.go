package owner

import (
	"errors"
	"fiber3-petclinic-service/api/pet"
	"fiber3-petclinic-service/api/visit"
	"fiber3-petclinic-service/internal/repository"
	"fiber3-petclinic-service/internal/repository/model"
	"fiber3-petclinic-service/internal/test"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Test_getById tests the getOwnerById function
func Test_getById(t *testing.T) {
	owner := &repository.Owner{
		Model: gorm.Model{
			ID: 1,
		},
		Person: model.Person{
			FirstName: "Nat",
			LastName:  "Cole",
		},
		Address: "123 Main St.",
		City:    "New York",
	}

	tests := []struct {
		id             uint
		mockResult     *repository.Owner
		mockError      error
		expectedResult *Response
		expectedError  error
	}{
		{
			id:         1,
			mockResult: owner,
			mockError:  nil,
			expectedResult: &Response{
				ID:        1,
				FirstName: "Nat",
				LastName:  "Cole",
				Address:   "123 Main St.",
				City:      "New York",
			},
			expectedError: nil,
		},
		{
			id:             1,
			mockResult:     nil,
			mockError:      errors.New("unable to find owner by id: 1"),
			expectedResult: nil,
			expectedError:  errors.New("unable to find owner by id: 1"),
		},
	}

	for _, testCase := range tests {
		ownerMock := repository.NewMockOwnerRepositorier(t)
		ownerMock.EXPECT().FindById(testCase.id).Return(testCase.mockResult, testCase.mockError)

		ownerService := NewService(ownerMock)
		result, err := ownerService.getOwnerById(testCase.id)
		ownerMock.AssertExpectations(t)
		ownerMock.AssertNumberOfCalls(t, "FindById", 1)

		if err != nil {
			assert.Equal(t, testCase.expectedError.Error(), err.Error())
			assert.Nil(t, result)
		} else {
			assert.Equal(t, testCase.id, result.ID)
			assert.Equal(t, testCase.expectedResult.FirstName, result.FirstName)
			assert.Equal(t, testCase.expectedResult.LastName, result.LastName)
			assert.Equal(t, testCase.expectedResult.Address, result.Address)
			assert.Equal(t, testCase.expectedResult.City, result.City)
			assert.Nil(t, err)
		}
		ownerMock.AssertNumberOfCalls(t, "FindById", 1)
	}

}

// Test_getByLastName tests the getOwnerByLastName function
func Test_getByLastName(t *testing.T) {
	tests := []struct {
		tcName         string
		lastName       string
		mockResult     []repository.Owner
		mockError      error
		expectedResult *Responses
		expectedError  error
	}{
		{
			tcName:   "Get Owner By Last Name",
			lastName: "DiCaprio",
			mockResult: []repository.Owner{
				{
					Model: gorm.Model{
						ID: 1,
					},
					Person: model.Person{
						FirstName: "Leo",
						LastName:  "DiCaprio",
					},
					City: "Boston",
				},
			},
			mockError: nil,
			expectedResult: &Responses{
				Context: model.Context{
					Count: 1,
				},
				Owners: []Response{
					{ID: 1, FirstName: "Leo", LastName: "DiCaprio", City: "Boston"},
				},
			},
			expectedError: nil,
		},
		{
			tcName:         "Fail to get Owner By Last Name",
			lastName:       "DiCaprio",
			mockResult:     nil,
			mockError:      errors.New("unable to find owner by last name: DiCaprio"),
			expectedResult: nil,
			expectedError:  errors.New("unable to find owner by last name: DiCaprio"),
		}, // Add more test cases here
	}

	for _, testCase := range tests {
		t.Run(testCase.tcName, func(t *testing.T) {
			ownerMock := repository.NewMockOwnerRepositorier(t)
			ownerMock.EXPECT().FindByLastName(testCase.lastName).Return(testCase.mockResult, testCase.mockError)
			ownerMock.On("FindByLastName", testCase.lastName).Return(testCase.expectedResult, testCase.expectedError)

			ownerService := NewService(ownerMock)
			result, err := ownerService.getOwnerByLastName(testCase.lastName)

			if err != nil {
				assert.Equal(t, testCase.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, testCase.expectedResult.Context.Count, result.Context.Count)
				for i, owner := range testCase.expectedResult.Owners {
					assert.Equal(t, owner.ID, result.Owners[i].ID)
					assert.Equal(t, owner.FirstName, result.Owners[i].FirstName)
					assert.Equal(t, owner.LastName, result.Owners[i].LastName)
					assert.Equal(t, owner.City, result.Owners[i].City)
					assert.Nil(t, err)
				}
			}

			ownerMock.AssertExpectations(t)
			// Verify that the FindByLastName method is called once
			// for each test case. Therefore, the number of calls = i+1
			ownerMock.AssertNumberOfCalls(t, "FindByLastName", 1)
		})
	}
}

func Test_getAllOwners(t *testing.T) {
	tests := []struct {
		tcName         string
		mockResult     []repository.Owner
		mockError      error
		expectedResult *Responses
		expectedError  error
	}{
		{
			tcName: "Test get all owners success",
			mockResult: []repository.Owner{
				{
					Model: gorm.Model{
						ID: 1,
					},
					Person: model.Person{
						FirstName: "Leo",
						LastName:  "DiCaprio",
					},
				},
				{
					Model: gorm.Model{
						ID: 2,
					},
					Person: model.Person{
						FirstName: "Tom",
						LastName:  "Hanks",
					},
				},
			},
			mockError: nil,
			expectedResult: &Responses{
				Context: model.Context{
					Count: 2,
				},
				Owners: []Response{
					{ID: 1, FirstName: "Leo", LastName: "DiCaprio"},
					{ID: 2, FirstName: "Tom", LastName: "Hanks"},
				},
			},
			expectedError: nil,
		},
		{
			tcName:         "Unable to get all owners",
			mockResult:     nil,
			mockError:      errors.New("unable to get all owners"),
			expectedResult: nil,
			expectedError:  errors.New("unable to get all owners"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.tcName, func(t *testing.T) {
			ownerMock := repository.NewMockOwnerRepositorier(t)
			ownerMock.EXPECT().FindAll().Return(tc.mockResult, tc.mockError)

			ownerService := NewService(ownerMock)
			result, err := ownerService.getAllOwners()

			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedResult.Context.Count, result.Context.Count)
				for i, owner := range tc.expectedResult.Owners {
					assert.Equal(t, owner.ID, result.Owners[i].ID)
					assert.Equal(t, owner.FirstName, result.Owners[i].FirstName)
					assert.Equal(t, owner.LastName, result.Owners[i].LastName)
				}
				assert.Nil(t, err)

				ownerMock.AssertExpectations(t)
				ownerMock.AssertNumberOfCalls(t, "FindAll", 1)
			}
		})
	}
}

func Test_getOwnerByIdWithPets(t *testing.T) {
	mockOwner := &repository.Owner{
		Model: gorm.Model{
			ID: 1,
		},
		Person: model.Person{
			FirstName: "Leo",
			LastName:  "DiCaprio",
		},
		Pets: []repository.Pet{
			{
				Name:      "Buddy",
				Birthdate: test.ToDate("2015-02-05"),
				Species: repository.Species{
					Name: "Dog",
				},
			},
		},
	}

	expectedResult := &Response{
		ID:        1,
		FirstName: "Leo",
		LastName:  "DiCaprio",
		Pets: []pet.Response{
			{
				Name:      "Buddy",
				Birthdate: "2015-02-05",
				Species:   "Dog",
				Visits:    []visit.Response{},
			},
		},
	}

	testCases := []struct {
		name           string
		id             uint
		mockOwner      *repository.Owner
		mockError      error
		expectedResult *Response
		expectedError  error
	}{
		{
			name:           "get owner by id with pets",
			id:             1,
			mockOwner:      mockOwner,
			mockError:      nil,
			expectedResult: expectedResult,
			expectedError:  nil,
		},
		{
			name:           "fail to get owner by id with pets",
			id:             1,
			mockOwner:      nil,
			mockError:      errors.New("unable to get owner by id"),
			expectedResult: nil,
			expectedError:  errors.New("unable to get owner by id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := repository.NewMockOwnerRepositorier(t)
			ownerMock.EXPECT().FindByIdWithPets(tc.id).Return(tc.mockOwner, tc.mockError)

			ownerService := NewService(ownerMock)
			result, err := ownerService.getOwnerByIdWithPets(tc.id)

			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedResult, result)
				assert.Nil(t, err)
			}

			ownerMock.AssertExpectations(t)
			ownerMock.AssertNumberOfCalls(t, "FindByIdWithPets", 1)
		})
	}
}

// Test_update tests the update function
func Test_update(t *testing.T) {
	mockOwner := &repository.Owner{
		Model: gorm.Model{
			ID: 1,
		},
		Person: model.Person{
			FirstName: "Nat",
			LastName:  "Cole",
		},
	}

	updatedOwner := &repository.Owner{
		Model: gorm.Model{
			ID:        1,
			UpdatedAt: time.Now(),
		},
		Person: model.Person{
			FirstName: "Nat",
			LastName:  "Cole",
		},
	}

	updateRequest := &UpdateRequest{
		FirstName: "Nat",
		LastName:  "Cole",
	}

	expectedResult := &Response{
		ID:        1,
		FirstName: "Nat",
		LastName:  "Cole",
	}
	testCases := []struct {
		name           string
		id             uint
		mockOwner      *repository.Owner
		mockError      error
		updatedOwner   *repository.Owner
		updateRequest  *UpdateRequest
		expectedResult *Response
		expectedError  error
	}{
		{
			name:           "Test update owner success",
			id:             1,
			mockOwner:      mockOwner,
			mockError:      nil,
			updatedOwner:   updatedOwner,
			updateRequest:  updateRequest,
			expectedResult: expectedResult,
			expectedError:  nil,
		},
		{
			name:           "fail to update owner",
			id:             1,
			mockOwner:      mockOwner,
			mockError:      errors.New("unable to update owner"),
			updatedOwner:   nil,
			updateRequest:  updateRequest,
			expectedResult: nil,
			expectedError:  errors.New("unable to update owner"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := repository.NewMockOwnerRepositorier(t)
			ownerMock.EXPECT().Update(tc.mockOwner).Return(tc.updatedOwner, tc.expectedError)

			ownerService := NewService(ownerMock)
			result, err := ownerService.update(tc.id, tc.updateRequest)

			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedResult.ID, result.ID)
				assert.Equal(t, tc.expectedResult.FirstName, result.FirstName)
				assert.Equal(t, tc.expectedResult.LastName, result.LastName)
				assert.Nil(t, err)
			}

			ownerMock.AssertExpectations(t)
			ownerMock.AssertNumberOfCalls(t, "Update", 1)
		})
	}
}

func Test_create(t *testing.T) {
	mockOwner := &repository.Owner{
		//Model: gorm.Model{
		//	ID: 1,
		//},
		Person: model.Person{
			FirstName: "Nat",
			LastName:  "Cole",
		},
	}

	newOwner := &repository.Owner{
		Model: gorm.Model{
			//ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Person: model.Person{
			FirstName: "Nat",
			LastName:  "Cole",
		},
	}

	addRequest := &AddRequest{
		FirstName: "Nat",
		LastName:  "Cole",
	}

	expectedResult := &Response{
		//ID:        int(1),
		FirstName: "Nat",
		LastName:  "Cole",
	}
	testCases := []struct {
		name           string
		mockOwner      *repository.Owner
		mockError      error
		newOwner       *repository.Owner
		addRequest     *AddRequest
		expectedResult *Response
		expectedError  error
	}{
		{
			name:           "add a new owner",
			mockOwner:      mockOwner,
			mockError:      nil,
			newOwner:       newOwner,
			addRequest:     addRequest,
			expectedResult: expectedResult,
			expectedError:  nil,
		},
		{
			name:           "fail to add a new owner",
			mockOwner:      mockOwner,
			mockError:      errors.New("unable to add new owner"),
			newOwner:       nil,
			addRequest:     addRequest,
			expectedResult: nil,
			expectedError:  errors.New("unable to add new owner"),
		},
	}

	for _, tc := range testCases {
		ownerMock := repository.NewMockOwnerRepositorier(t)
		ownerMock.EXPECT().Insert(tc.mockOwner).Return(tc.newOwner, tc.mockError)

		ownerService := NewService(ownerMock)
		result, err := ownerService.create(tc.addRequest)

		if err != nil {
			assert.Equal(t, tc.expectedError, err)
			assert.Nil(t, result)
		} else {
			assert.Equal(t, tc.expectedResult.ID, result.ID)
			assert.Equal(t, tc.expectedResult.FirstName, result.FirstName)
			assert.Equal(t, tc.expectedResult.LastName, result.LastName)
			assert.Nil(t, err)
		}
		ownerMock.AssertExpectations(t)
		ownerMock.AssertNumberOfCalls(t, "Insert", 1)
	}
}
