package vet

import (
	"errors"
	"fiber-petclinic-service/pkg/repository"
	"fiber-petclinic-service/pkg/repository/model"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getAllSpecialties(t *testing.T) {
	testCases := []struct {
		name          string
		mockSpecs     []repository.Specialty
		mockError     error
		expectedSpecs *specialtiesResponse
		expectedError error
	}{
		{
			name: "get all specialties",
			mockSpecs: []repository.Specialty{
				{Model: gorm.Model{
					ID: uint(1),
				}, Name: "Surgery"},
				{Model: gorm.Model{
					ID: uint(2),
				}, Name: "Dentistry"},
			},
			mockError: nil,
			expectedSpecs: &specialtiesResponse{
				Context: model.Context{
					Count: 2,
				},
				Specialties: []specialtyResponse{
					{ID: uint(1), Name: "Surgery"},
					{ID: uint(2), Name: "Dentistry"},
				},
			},
			expectedError: nil,
		},
		{
			name:          "get no specialties",
			mockSpecs:     nil,
			mockError:     gorm.ErrRecordNotFound,
			expectedSpecs: nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "fail to get all specialties",
			mockSpecs:     nil,
			mockError:     errors.New("getAllSpecialties: unable to retrieve specialties"),
			expectedSpecs: nil,
			expectedError: errors.New("getAllSpecialties: unable to retrieve specialties"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vetMock := repository.NewMockVetRepositorier(t)

			vetMock.EXPECT().FindAllSpecialties().Return(tc.mockSpecs, tc.expectedError)
			vetService := NewService(vetMock)
			result, err := vetService.getAllSpecialties()

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedSpecs, result)
				assert.Nil(t, err)
			}

			vetMock.AssertExpectations(t)
			vetMock.AssertNumberOfCalls(t, "FindAllSpecialties", 1)
		})
	}
}

// Test_getById tests the getVetById function
func Test_getById(t *testing.T) {
	testCases := []struct {
		name          string
		id            uint
		mockVet       *repository.Vet
		mockError     error
		expectedVet   *Response
		expectedError error
	}{
		{
			name: "get vet by id",
			id:   uint(1),
			mockVet: &repository.Vet{
				Model: gorm.Model{
					ID: 1,
				},
				Person: model.Person{
					FirstName: "Nat",
					LastName:  "Cole",
				},
			},
			mockError: nil,
			expectedVet: &Response{
				ID:        1,
				FirstName: "Nat",
				LastName:  "Cole",
			},
			expectedError: nil,
		},
		{
			name:          "get no vet by id",
			id:            2,
			mockVet:       nil,
			mockError:     gorm.ErrRecordNotFound,
			expectedVet:   nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "get vet by id with error",
			id:            2,
			mockVet:       nil,
			mockError:     errors.New("getById: unable to find vet by id"),
			expectedVet:   nil,
			expectedError: errors.New("getById: unable to find vet by id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vetMock := repository.NewMockVetRepositorier(t)

			vetMock.EXPECT().FindById(tc.id).Return(tc.mockVet, tc.mockError)
			vetService := NewService(vetMock)
			result, err := vetService.getVetById(tc.id)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedVet.ID, result.ID)
				assert.Equal(t, tc.expectedVet.FirstName, result.FirstName)
				assert.Equal(t, tc.expectedVet.LastName, result.LastName)
				assert.Nil(t, err)
			}

			vetMock.AssertExpectations(t)
			vetMock.AssertNumberOfCalls(t, "FindById", 1)
		})
	}
}

func Test_getByIdWithSpecialties(t *testing.T) {
	testCases := []struct {
		name          string
		id            uint
		mockVet       *repository.Vet
		mockError     error
		expectedVet   *Response
		expectedError error
	}{
		{
			name: "get vet by id",
			id:   1,
			mockVet: &repository.Vet{
				Model: gorm.Model{
					ID: 1,
				},
				Person: model.Person{
					FirstName: "Nat",
					LastName:  "Cole",
				},
				Specialties: []repository.Specialty{
					{
						Model: gorm.Model{
							ID: 1,
						},
						Name: "Surgery",
					},
					{
						Model: gorm.Model{
							ID: 2,
						},
						Name: "Dentistry",
					},
				},
			},
			mockError: nil,
			expectedVet: &Response{
				ID:        1,
				FirstName: "Nat",
				LastName:  "Cole",
				Specialties: &[]specialtyResponse{
					{ID: 1, Name: "Surgery"},
					{ID: 2, Name: "Dentistry"},
				},
			},
			expectedError: nil,
		},
		{
			name:          "get no vet by id with specialties",
			id:            2,
			mockVet:       nil,
			mockError:     gorm.ErrRecordNotFound,
			expectedVet:   nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "fail to get vet by id with specialties",
			id:            2,
			mockVet:       nil,
			mockError:     errors.New("getById: unable to find vet by id"),
			expectedVet:   nil,
			expectedError: errors.New("getById: unable to find vet by id"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vetMock := repository.MockVetRepositorier{}

			vetMock.EXPECT().FindByIdWithSpecialties(tc.id).Return(tc.mockVet, tc.mockError)
			vetService := NewService(&vetMock)
			result, err := vetService.getVetByIdWithSpecialties(tc.id)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedVet, result)
				assert.Nil(t, err)
			}

			vetMock.AssertExpectations(t)
			vetMock.AssertNumberOfCalls(t, "FindByIdWithSpecialties", 1)
		})
	}
}

// Test_getByLastName tests the getVetByLastName function
func Test_getByLastName(t *testing.T) {
	testCases := []struct {
		name          string
		lastName      string
		mockVets      []repository.Vet
		mockError     error
		expectedVets  *Responses
		expectedError error
	}{
		{
			name:     "get vet by last name",
			lastName: "DiCaprio",
			mockVets: []repository.Vet{
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
						FirstName: "Tiger",
						LastName:  "DiCaprio",
					},
				},
			},
			mockError: nil,
			expectedVets: &Responses{
				Context: model.Context{
					Count: 2,
				},
				Vets: []Response{
					{ID: 1, FirstName: "Leo", LastName: "DiCaprio"},
					{ID: 2, FirstName: "Tiger", LastName: "DiCaprio"}, // Specialties: {}},
				},
			},
			expectedError: nil,
		},
		{
			name:          "get no vet by last namer",
			lastName:      "DiCaprio",
			mockVets:      nil,
			mockError:     gorm.ErrRecordNotFound,
			expectedVets:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "fail to get vet by last namer",
			lastName:      "DiCaprio",
			mockVets:      nil,
			mockError:     errors.New("getByLastName: unable to find vet by last name"),
			expectedVets:  nil,
			expectedError: errors.New("getByLastName: unable to find vet by last name"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vetMock := repository.NewMockVetRepositorier(t)
			vetMock.EXPECT().FindByLastName(tc.lastName).Return(tc.mockVets, tc.mockError)

			vetService := NewService(vetMock)
			result, err := vetService.getVetByLastName(tc.lastName)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedVets, result)
				assert.Nil(t, err)
			}
			vetMock.AssertExpectations(t)
			vetMock.AssertNumberOfCalls(t, "FindByLastName", 1)
		})
	}
}

func Test_getAllVets(t *testing.T) {
	testCases := []struct {
		name          string
		mockVets      []repository.Vet
		mockError     error
		expectedVets  *Responses
		expectedError error
	}{
		{
			name: "get all vets",
			mockVets: []repository.Vet{
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
						FirstName: "Tiger",
						LastName:  "DiCaprio",
					},
				},
			},
			mockError: nil,
			expectedVets: &Responses{
				Context: model.Context{
					Count: 2,
				},
				Vets: []Response{
					{ID: 1, FirstName: "Leo", LastName: "DiCaprio"},
					{ID: 2, FirstName: "Tiger", LastName: "DiCaprio"}, // Specialties: {}},
				},
			},
			expectedError: nil,
		},
		{
			name:          "get no vets",
			mockVets:      nil,
			mockError:     gorm.ErrRecordNotFound,
			expectedVets:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "fail to get all vets",
			mockVets:      nil,
			mockError:     errors.New("getAllVets: unable to retrieve all vets"),
			expectedVets:  nil,
			expectedError: errors.New("getAllVets: unable to retrieve all vets"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vetMock := repository.NewMockVetRepositorier(t)
			vetMock.EXPECT().FindAll().Return(tc.mockVets, tc.mockError)

			vetService := NewService(vetMock)
			result, err := vetService.getAllVets()

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedVets, result)
				assert.Nil(t, err)
			}
			vetMock.AssertExpectations(t)
			vetMock.AssertNumberOfCalls(t, "FindAll", 1)
		})
	}
}

func Test_getAllVetsWithSpecialties(t *testing.T) {
	testCases := []struct {
		name          string
		mockVets      []repository.Vet
		mockError     error
		expectedVets  *Responses
		expectedError error
	}{
		{
			name: "get all vets with specialties",
			mockVets: []repository.Vet{
				{
					Model: gorm.Model{
						ID: 1,
					},
					Person: model.Person{
						FirstName: "Leo",
						LastName:  "DiCaprio",
					},
					Specialties: []repository.Specialty{
						{
							Model: gorm.Model{
								ID: 1,
							},
							Name: "Surgery",
						},
					},
				},
				{
					Model: gorm.Model{
						ID: 2,
					},
					Person: model.Person{
						FirstName: "Tiger",
						LastName:  "DiCaprio",
					},
				},
			},
			mockError: nil,
			expectedVets: &Responses{
				Context: model.Context{
					Count: 2,
				},
				Vets: []Response{
					{ID: 1, FirstName: "Leo", LastName: "DiCaprio",
						Specialties: &[]specialtyResponse{
							{ID: 1, Name: "Surgery"},
						}},
					{ID: 2, FirstName: "Tiger", LastName: "DiCaprio", Specialties: nil}, // Specialties: {}},
				},
			},
			expectedError: nil,
		},
		{
			name:          "get no vets with specialties",
			mockVets:      nil,
			mockError:     gorm.ErrRecordNotFound,
			expectedVets:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "fail to get all vets",
			mockVets:      nil,
			mockError:     errors.New("getAllVets: unable to retrieve all vets"),
			expectedVets:  nil,
			expectedError: errors.New("getAllVets: unable to retrieve all vets"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vetMock := repository.NewMockVetRepositorier(t)
			vetMock.EXPECT().FindAllPreload().Return(tc.mockVets, tc.mockError)

			vetService := NewService(vetMock)
			result, err := vetService.getAllVetsWithSpecialties()

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedVets, result)
				assert.Nil(t, err)
			}
			vetMock.AssertExpectations(t)
			vetMock.AssertNumberOfCalls(t, "FindAllPreload", 1)
		})
	}
}

func Test_create(t *testing.T) {
	mockVet := &repository.Vet{
		Model: gorm.Model{
			ID: 1,
		},
		Person: model.Person{
			FirstName: "Leo",
			LastName:  "DiCaprio",
		},
	}

	testCases := []struct {
		name          string
		mockVet       *repository.Vet
		mockResult    *repository.Vet
		mockError     error
		expectedVet   *Response
		expectedError error
	}{
		{
			name:       "create vet",
			mockVet:    mockVet,
			mockResult: mockVet,
			mockError:  nil,
			expectedVet: &Response{
				ID:        1,
				FirstName: "Leo",
				LastName:  "DiCaprio",
			},
			expectedError: nil,
		},
		{
			name:          "fail to create vet",
			mockVet:       mockVet,
			mockResult:    nil,
			mockError:     errors.New("createVet: unable to create vet"),
			expectedVet:   nil,
			expectedError: errors.New("createVet: unable to create vet"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vetMock := repository.NewMockVetRepositorier(t)
			vetMock.EXPECT().Insert(tc.mockVet).Return(tc.mockResult, tc.mockError)

			vetService := NewService(vetMock)
			result, err := vetService.create(tc.mockVet)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedVet, result)
				assert.Nil(t, err)
			}
			vetMock.AssertExpectations(t)
			vetMock.AssertNumberOfCalls(t, "Insert", 1)
		})
	}
}

func Test_update(t *testing.T) {
	testCases := []struct {
		name          string
		mockVet       *repository.Vet
		mockError     error
		expectedVet   *Response
		expectedError error
	}{
		{
			name: "update vet",
			mockVet: &repository.Vet{
				Model: gorm.Model{
					ID: 1,
				},
				Person: model.Person{
					FirstName: "Leo",
					LastName:  "DiCaprio",
				},
			},
			mockError: nil,
			expectedVet: &Response{
				ID:        1,
				FirstName: "Leo",
				LastName:  "DiCaprio",
			},
			expectedError: nil,
		},
		{
			name:          "found no vet to update",
			mockVet:       nil,
			mockError:     gorm.ErrRecordNotFound,
			expectedVet:   nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "fail to update vet",
			mockVet:       nil,
			mockError:     errors.New("updateVet: unable to update vet"),
			expectedVet:   nil,
			expectedError: errors.New("updateVet: unable to update vet"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vetMock := repository.NewMockVetRepositorier(t)
			vetMock.EXPECT().Update(tc.mockVet).Return(tc.mockVet, tc.mockError)

			vetService := NewService(vetMock)
			result, err := vetService.update(tc.mockVet)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tc.expectedVet, result)
				assert.Nil(t, err)
			}
			vetMock.AssertExpectations(t)
			vetMock.AssertNumberOfCalls(t, "Update", 1)
		})
	}
}
