package owner

import (
	"errors"
	"fiber-petclinic-service/pkg/repository"
	"fiber-petclinic-service/pkg/repository/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
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

//// Test_getById_withError tests the getOwnerById function with an error
//func Test_getById_withError(t *testing.T) {
//	ownerMock := repository.MockOwnerRepositorier{}
//	err := errors.New("getById: unable to find owner by id")
//	ownerMock.On("FindById", 1).Return(nil, err)
//
//	ownerService := NewService(&ownerMock)
//	result, err := ownerService.getOwnerById(1)
//	ownerMock.AssertExpectations(t)
//	ownerMock.AssertNumberOfCalls(t, "FindById", 1)
//
//	assert.Equal(t, "getById: unable to find owner by id", err.Error())
//	assert.Nil(t, result)
//}

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
			ownerMock.On("FindAll").Return(tc.mockResult, tc.mockError)

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

//
//func Test_getAllOwnersWithPets(t *testing.T) {
//	testCases := []struct {
//		name           string
//		expectedOwners []repository.Owner
//		expectedError  error
//	}{
//		{
//			name: "Test get all owners with pets success",
//			expectedOwners: []repository.Owner{
//				{
//					Model: gorm.Model{
//						ID: 1,
//					},
//					Person: model.Person{
//						FirstName: "Leo",
//						LastName:  "DiCaprio",
//					},
//					Pets: []repository.Pet{
//						{
//							Name:      "Buddy",
//							BirthDate: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
//							Type: repository.Type{
//								Name: "Dog",
//							},
//						},
//					},
//				},
//				{
//					Model: gorm.Model{
//						ID: 2,
//					},
//					Person: model.Person{
//						FirstName: "Tom",
//						LastName:  "Hanks",
//					},
//					Pets: []repository.Pet{
//						{
//							Name:      "Kitty",
//							BirthDate: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
//							Type: repository.Type{
//								Name: "Cat",
//							},
//						},
//					},
//				},
//			},
//			expectedError: nil,
//		},
//		{
//			name:           "Test get all owners with pets with error",
//			expectedOwners: nil,
//			expectedError:  errors.New("getAllOwnersWithPets: unable to retrieve owners"),
//		},
//		// Add more test cases here
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			logger := log.New().With(nil, "function", "Test_getAllOwnersWithPets")
//			ownerMock := repository.MockOwnerRepositorier{}
//			ownerMock.On("FindAllWithPets").Return(tc.expectedOwners, tc.expectedError)
//
//			ownerService := NewService(logger, &ownerMock)
//			result, err := ownerService.getAllOwnersWithPets()
//
//			if tc.expectedError != nil {
//				assert.Equal(t, tc.expectedError, err)
//				assert.Nil(t, result)
//			} else {
//				assert.Equal(t, len(tc.expectedOwners), len(result))
//				for i, owner := range tc.expectedOwners {
//					assert.Equal(t, owner.ID, result[i].ID)
//					assert.Equal(t, owner.FirstName, result[i].FirstName)
//					assert.Equal(t, owner.LastName, result[i].LastName)
//					assert.Equal(t, len(owner.Pets), len(result[i].Pets))
//					assert.Nil(t, err)
//				}
//			}
//
//			ownerMock.AssertExpectations(t)
//			ownerMock.AssertNumberOfCalls(t, "FindAllWithPets", 1)
//		})
//	}
//}

//// Test_update tests the update function
////func Test_update(t *testing.T) {
////	testCases := []struct {
////		name          string
////		expectedOwner *repository.Owner
////		expectedError error
////	}{
////		{
////			name: "Test update owner success",
////			expectedOwner: &repository.Owner{
////				Model: gorm.Model{
////					ID: 1,
////				},
////				Person: model.Person{
////					FirstName: "Nat",
////					LastName:  "Cole",
////				},
////			},
////			expectedError: nil,
////		},
////		{
////			name:          "Test update owner with error",
////			expectedOwner: nil,
////			expectedError: errors.New("update: unable to update owner"),
////		},
////		// Add more test cases here
////	}
////
////	for _, tc := range testCases {
////		t.Run(tc.name, func(t *testing.T) {
////			logger := zap.NewNop()
////			//logger := log.New().With(nil, "function", "Test_update")
////			ownerMock := repository.MockOwnerRepositorier{}
////			ownerMock.On("Update", tc.expectedOwner).Return(tc.expectedOwner, tc.expectedError)
////
////			ownerService := NewService(logger, &ownerMock)
////			result, err := ownerService.update(tc.expectedOwner)
////
////			if tc.expectedError != nil {
////				assert.Equal(t, tc.expectedError, err)
////				assert.Nil(t, result)
////			} else {
////				assert.Equal(t, tc.expectedOwner.ID, result.ID)
////				assert.Equal(t, tc.expectedOwner.FirstName, result.FirstName)
////				assert.Equal(t, tc.expectedOwner.LastName, result.LastName)
////				assert.Nil(t, err)
////			}
////
////			ownerMock.AssertExpectations(t)
////			ownerMock.AssertNumberOfCalls(t, "Update", 1)
////		})
////	}
////}
//
///*
//func Test_create(t *testing.T) {
//	logger := log.New().With(nil, "function", "Test_create")
//	ownerMock := repository.MockOwnerRepositorier{}
//	owner := &repository.Owner{
//		Model: gorm.Model{
//			ID: 1,
//		},
//		Person: model.Person{
//			FirstName: "Nat",
//			LastName:  "Cole",
//		},
//	}
//	ownerMock.On("Create", owner).Return(owner, nil)
//
//	ownerService := NewService(logger, &ownerMock)
//	result, err := ownerService.create(owner)
//	ownerMock.AssertExpectations(t)
//	ownerMock.AssertNumberOfCalls(t, "Create", 1)
//
//	assert.Equal(t, owner.ID, result.ID)
//	assert.Equal(t, owner.FirstName, result.FirstName)
//	assert.Equal(t, owner.LastName, result.LastName)
//	assert.Nil(t, err)
//}
//
