package pet

import (
	"errors"
	"fiber-petclinic-service/api/visit"
	"fiber-petclinic-service/pkg/repository"
	"fiber-petclinic-service/pkg/repository/model"
	"fiber-petclinic-service/pkg/test"
	"gorm.io/gorm"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// for learning purpose, how manual mocking can be done.
// petRepositoryMock
//
//	requires mocking all functions; otherwise, compilation errors.
type petRepositoryMock struct {
	mock.Mock
}

func (petM *petRepositoryMock) FindAll() ([]repository.Pet, error) {
	args := petM.Called()
	intf := args.Get(0)
	val := intf.([]repository.Pet)

	return val, args.Error(1)
}

// FindById is a method on the petRepositoryMock struct. It simulates the behavior of
// fetching a Pet object from a data repository based on its unique identifier.
//
// This method takes an integer id as an argument, which is expected to be the unique
// identifier of a Pet object in a data repository.
//
// The method starts by calling the Called method with the id argument. This is a method
// provided by the mock package, which records the fact that a method has been called with
// specific arguments. It returns a mock.Call object.
//
// The mock.Call object has a Get method that retrieves the argument at the specified index.
// In this case, it retrieves the first argument (index 0), which is the Pet object that the
// FindById method is expected to return. This Pet object is then type-asserted to the Pet type.
//
// The method then returns a pointer to the Pet object and an error. The error is expected to be
// nil if the method operates as expected, or an instance of error if something goes wrong.
func (petM *petRepositoryMock) FindById(id uint) (*repository.Pet, error) {
	args := petM.Called(id)
	intf := args.Get(0)
	val := intf.(*repository.Pet)
	ptr := val

	return ptr, args.Error(1)
}

func (petM *petRepositoryMock) FindByIdWithVisits(id uint) (*repository.Pet, error) {
	args := petM.Called(id)
	intf := args.Get(0)
	val := intf.(*repository.Pet)
	ptr := val

	return ptr, args.Error(1)
}

func (petM *petRepositoryMock) FindByName(name string) ([]repository.Pet, error) {
	args := petM.Called(name)
	intf := args.Get(0)
	val := intf.([]repository.Pet)
	ptr := val

	return ptr, args.Error(1)
}

func (petM *petRepositoryMock) Insert(pet *repository.Pet) (*repository.Pet, error) {
	args := petM.Called(pet)
	intf := args.Get(0)
	val := intf.(*repository.Pet)
	ptr := val

	return ptr, args.Error(1)
}

func (petM *petRepositoryMock) Update(pet *repository.Pet) (*repository.Pet, error) {
	args := petM.Called(pet)
	intf := args.Get(0)
	val := intf.(*repository.Pet)
	ptr := val

	return ptr, args.Error(1)
}

// mocking ends

func Test_retrieveAllPets(t *testing.T) {
	mockPet := []repository.Pet{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name:      "Nash",
			Birthdate: *test.ToDate("2014-10-07"),
			Type: repository.Type{
				Model: gorm.Model{
					ID: 1,
				},
				Name: "Dog",
			},
		},
	}

	expectedPets := &Responses{
		Context: model.Context{
			Count: len(mockPet),
		},
		Pets: []Response{
			{
				ID:        1,
				Name:      "Nash",
				Birthdate: "2014-10-07",
				Type:      "Dog",
				Visits:    []visit.Response{},
			},
		},
	}

	testCases := []struct {
		name          string
		mockPets      []repository.Pet
		mockError     error
		expectedPets  *Responses
		expectedError error
	}{
		{
			name:          "get all pets",
			mockPets:      mockPet,
			mockError:     nil,
			expectedPets:  expectedPets,
			expectedError: nil,
		},
		{
			name:          "get no pet",
			mockPets:      nil,
			mockError:     gorm.ErrRecordNotFound,
			expectedPets:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "fail to get all pets",
			mockPets:      nil,
			mockError:     errors.New("fail to retrieve all pets"),
			expectedPets:  nil,
			expectedError: errors.New("fail to retrieve all pets"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			petMock := petRepositoryMock{}
			petMock.On("FindAll").Return(tc.mockPets, tc.mockError)

			petService := NewService(&petMock)
			result, err := petService.getAllPets()

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedPets, result)
				assert.Nil(t, err)
			}

			petMock.AssertExpectations(t)
			petMock.AssertNumberOfCalls(t, "FindAll", 1)
		})
	}
}

func Test_getById(t *testing.T) {
	testCases := []struct {
		name          string
		id            uint
		mockPet       *repository.Pet
		mockError     error
		route         string
		expectedPet   *Response
		expectedError error
	}{
		{
			name: "get pet by id",
			mockPet: &repository.Pet{
				Model: gorm.Model{
					ID: 1,
				},
				Name:      "Nash",
				Birthdate: *test.ToDate("2014-10-07"),
				Type: repository.Type{
					Model: gorm.Model{
						ID: 1,
					},
					Name: "Dog",
				},
			},
			mockError: nil,
			route:     "/v1/pets/1",
			expectedPet: &Response{
				ID:        1,
				Name:      "Nash",
				Birthdate: "2014-10-07",
				Type:      "Dog",
				Visits:    []visit.Response{},
			},
			expectedError: nil,
		},
		{
			name:          "found no pet by id",
			id:            1,
			mockPet:       nil,
			mockError:     gorm.ErrRecordNotFound,
			route:         "/v1/pets/1",
			expectedPet:   nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "fail to get pet by id",
			id:            1,
			mockPet:       nil,
			mockError:     errors.New("fail to retrieve pet by id"),
			route:         "/v1/pets/1",
			expectedPet:   nil,
			expectedError: errors.New("fail to retrieve pet by id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			petMock := petRepositoryMock{}
			petMock.On("FindById", tc.id).Return(tc.mockPet, tc.mockError)

			petService := NewService(&petMock)
			result, err := petService.getPetById(tc.id)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedPet, result)
				assert.Nil(t, err)
			}

			petMock.AssertExpectations(t)
			petMock.AssertNumberOfCalls(t, "FindById", 1)
		})
	}
}

func Test_getByIdWithVisits(t *testing.T) {
	testCases := []struct {
		name          string
		id            uint
		mockPet       *repository.Pet
		mockError     error
		route         string
		expectedPet   *Response
		expectedError error
	}{
		{
			name: "get pet by id with visits",
			mockPet: &repository.Pet{
				Model: gorm.Model{
					ID: 1,
				},
				Name:      "Nash",
				Birthdate: *test.ToDate("2014-10-07"),
				Type: repository.Type{
					Model: gorm.Model{
						ID: 1,
					},
					Name: "Dog",
				},
				Visits: []repository.Visit{
					{
						Model: gorm.Model{
							ID: 1,
						},
						VisitDate:   test.ToDate("2014-10-07"),
						Description: "First visit",
					},
				},
			},
			mockError: nil,
			route:     "/v1/pets/1/visits",
			expectedPet: &Response{
				ID:        1,
				Name:      "Nash",
				Birthdate: "2014-10-07",
				Type:      "Dog",
				Visits: []visit.Response{
					{
						ID:          1,
						VisitDate:   "2014-10-07",
						Description: "First visit",
					},
				},
			},
			expectedError: nil,
		},
		{
			name:          "found no pet by id",
			id:            1,
			mockPet:       nil,
			mockError:     gorm.ErrRecordNotFound,
			route:         "/v1/pets/1/visits",
			expectedPet:   nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "fail to get pet by id",
			id:            1,
			mockPet:       nil,
			mockError:     errors.New("fail to retrieve pet by id"),
			route:         "/v1/pets/1/visits",
			expectedPet:   nil,
			expectedError: errors.New("fail to retrieve pet by id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			petMock := petRepositoryMock{}
			petMock.On("FindByIdWithVisits", tc.id).Return(tc.mockPet, tc.mockError)

			petService := NewService(&petMock)
			result, err := petService.getPetWithVisitsById(tc.id)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedPet, result)
				assert.Nil(t, err)
			}

			petMock.AssertExpectations(t)
			petMock.AssertNumberOfCalls(t, "FindByIdWithVisits", 1)
		})
	}
}

func Test_getByName(t *testing.T) {
	testCases := []struct {
		name          string
		expectedPets  []repository.Pet
		expectedError error
	}{
		{
			name: "Test get pet by name success",
			expectedPets: []repository.Pet{
				{
					Model: gorm.Model{
						ID: 1,
					},
					Name:      "Leo",
					Birthdate: time.Date(2017, 07, 02, 0, 0, 0, 0, time.UTC),
					Type: repository.Type{
						Model: gorm.Model{
							ID: 2,
						},
						Name: "Cat",
					},
				},
			},
			expectedError: nil,
		},
		{
			name:          "Test get pet by name failed",
			expectedPets:  nil,
			expectedError: errors.New("Fail to retrieve pet by name"),
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//logger := log.New().With(nil, "function", "Test_getByName")
			petMock := petRepositoryMock{}
			petMock.On("FindByName", "Leo").Return(tc.expectedPets, tc.expectedError)

			petService := NewService(&petMock)
			result, err := petService.getPetsByName("Leo")

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, len(tc.expectedPets), len(result))
				for i, pet := range tc.expectedPets {
					assert.Equal(t, pet.ID, result[i].ID, "Pet ID should be the same")
					assert.Equal(t, pet.Name, result[i].Name, "Pet Name should be the same")
					assert.Nil(t, err)
				}
			}

			petMock.AssertExpectations(t)
			petMock.AssertNumberOfCalls(t, "FindByName", 1)
		})
	}
}

// Test_update is a test function that tests the update method of the PetService struct.
func Test_update(t *testing.T) {
	testCases := []struct {
		name          string
		expectedPet   *repository.Pet
		expectedError error
	}{
		{
			name: "Test update pet success",
			expectedPet: &repository.Pet{
				Model: gorm.Model{
					ID: 1,
				},
				Name:      "Leo",
				Birthdate: time.Date(2017, 7, 02, 0, 0, 0, 0, time.UTC),
				Type: repository.Type{
					Model: gorm.Model{
						ID: 2,
					},
					Name: "Cat",
				},
			},
			expectedError: nil,
		},
		{
			name:          "Test update pet failed",
			expectedPet:   nil,
			expectedError: errors.New("Fail to update pet"),
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//logger := log.New().With(nil, "function", "Test_update")
			petMock := petRepositoryMock{}
			petMock.On("Update", tc.expectedPet).Return(tc.expectedPet, tc.expectedError)

			petService := NewService(&petMock)
			result, err := petService.update(tc.expectedPet)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedPet.ID, result.ID, "Pet ID should be the same")
				assert.Equal(t, tc.expectedPet.Name, result.Name, "Pet Name should be the same")
				assert.Nil(t, err)
			}

			petMock.AssertExpectations(t)
			petMock.AssertNumberOfCalls(t, "Update", 1)
		})
	}
}

//func Test_create(t *testing.T) {
//	testCases := []struct {
//		name          string
//		expectedPet   *repository.Pet
//		expectedError error
//	}{
//		{
//			name: "Test create pet success",
//			expectedPet: &repository.Pet{
//				Model: gorm.Model{
//					ID: 1,
//				},
//				Name:      "Leo",
//				BirthDate: time.Date(2017, 7, 2, 0, 0, 0, 0, time.UTC),
//				Type: repository.Type{
//					Model: gorm.Model{
//						ID: 2,
//					},
//					Name: "Cat",
//				},
//			},
//			expectedError: nil,
//		},
//		{
//			name:          "Test create pet failed",
//			expectedPet:   nil,
//			expectedError: errors.New("Fail to insert pet"),
//		},
//		// Add more test cases here
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			logger := log.New().With(nil, "function", "Test_insert")
//			petMock := petRepositoryMock{}
//			petMock.On("Fail", tc.expectedPet).Return(tc.expectedPet, tc.expectedError)
//
//			petService := NewService(logger, &petMock)
//			result, err := petService.create(tc.expectedPet)
//
//			if tc.expectedError != nil {
//				assert.Equal(t, tc.expectedError, err)
//				assert.Nil(t, result)
//			} else {
//				assert.Equal(t, tc.expectedPet.ID, result.ID, "Pet ID should be the same")
//				assert.Equal(t, tc.expectedPet.Name, result.Name, "Pet Name should be the same")
//				assert.Nil(t, err)
//			}
//
//			petMock.AssertExpectations(t)
//			petMock.AssertNumberOfCalls(t, "Fail", 1)
//		})
//	}
//}

/*
	mockPets := []repository.Pet{
		{Model: gorm.Model{ID: 1}, Name: "Tom",
			Birthdate: time.Date(2015, 11, 19, 0, 0, 0, 00, time.UTC),
			TypeID:    19, OwnerID: 7},
		{Model: gorm.Model{ID: 2}, Name: "Mike",
			Birthdate: time.Date(2018, 4, 17, 0, 0, 0, 0, time.UTC),
			TypeID:    20, OwnerID: 7},
	}
*/
