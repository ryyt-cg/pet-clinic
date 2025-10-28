package owner

import (
	"gin-petclinic-service/api/pet"
	"gin-petclinic-service/internal/repository"
	"gin-petclinic-service/internal/repository/model"
	"gin-petclinic-service/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_ResponseFromOwner(t *testing.T) {
	owner := &repository.Owner{
		Model: gorm.Model{
			ID: 1,
		},
		Person: model.Person{
			FirstName: "John",
			LastName:  "Doe",
		},
		Address:   "123 Main St",
		City:      "Anytown",
		Telephone: "1234567890",
		Pets:      nil,
	}

	expectedResponse := Response{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Address:   "123 Main St",
		City:      "Anytown",
		Telephone: "1234567890",
		Pets:      nil,
	}

	response := ToResponse(owner)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.FirstName, response.FirstName)
	assert.Equal(t, expectedResponse.LastName, response.LastName)
	assert.Equal(t, expectedResponse.Address, response.Address)
	assert.Equal(t, expectedResponse.City, response.City)
	assert.Equal(t, expectedResponse.Telephone, response.Telephone)
	assert.Equal(t, expectedResponse.Pets, response.Pets)
}

func Test_FromOwners(t *testing.T) {
	owners := []repository.Owner{
		{
			Model: gorm.Model{ID: 1},
			Person: model.Person{
				FirstName: "John",
				LastName:  "Doe",
			},
			Address:   "123 Main St",
			City:      "Anytown",
			Telephone: "1234567890",
		},
		{
			Model: gorm.Model{ID: 2},
			Person: model.Person{FirstName: "Jane",
				LastName: "Doe"},
			Address:   "456 Main St",
			City:      "Anytown",
			Telephone: "0987654321",
		},
		{
			Model: gorm.Model{ID: 3},
			Person: model.Person{FirstName: "Jane",
				LastName: "Doe"},
			Address:   "456 Main St",
			City:      "Anytown",
			Telephone: "0987654321",
			Pets: []repository.Pet{
				{Model: gorm.Model{ID: 1}, Name: "Pet1", Birthdate: test.ToDate("2019-04-07")},
			},
		},
	}

	ownerResponses := []Response{
		{ID: 1, FirstName: "John", LastName: "Doe", Address: "123 Main St", City: "Anytown", Telephone: "1234567890",
			Pets: nil,
		},
		{ID: 2, FirstName: "Jane", LastName: "Doe", Address: "456 Main St", City: "Anytown", Telephone: "0987654321",
			Pets: nil,
		},
		{ID: 3, FirstName: "Jane", LastName: "Doe", Address: "456 Main St", City: "Anytown", Telephone: "0987654321",
			Pets: []pet.Response{
				{ID: 1, Name: "Pet1", Birthdate: "2019-04-07", Visits: nil},
			},
		},
	}

	responses := FromOwners(owners)

	assert.Equal(t, len(ownerResponses), len(responses))
	for i, owner := range ownerResponses {
		assert.Equal(t, owner, responses[i])
	}
}

func Test_FromUpdateEntity(t *testing.T) {
	tests := []struct {
		updateEntity   *repository.Owner
		updateResponse *UpdateResponse
	}{
		{
			updateEntity: &repository.Owner{
				Model: gorm.Model{ID: 202},
				Person: model.Person{
					FirstName: "Ronald",
					LastName:  "Petersen",
				},
				Address: "123 Main St",
				City:    "Anytown",
			},
			updateResponse: &UpdateResponse{
				ID:        202,
				FirstName: "Ronald",
				LastName:  "Petersen",
				Address:   "123 Main St",
				City:      "Anytown",
			},
		},
	}

	for _, tc := range tests {
		result := ToUpdateResponse(tc.updateEntity)
		assert.Equal(t, result, tc.updateResponse)
	}
}
