package owner

import (
	"fiber-petclinic-service/api/pet"
	"fiber-petclinic-service/pkg/repository"
	"fiber-petclinic-service/pkg/repository/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
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
		Pets:      []repository.Pet{},
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

	response := &Response{}
	response.FromOwner(owner)

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
			Pets:      nil,
		},
		{
			Model: gorm.Model{ID: 2},
			Person: model.Person{FirstName: "Jane",
				LastName: "Doe"},
			Address:   "456 Main St",
			City:      "Anytown",
			Telephone: "0987654321",
			Pets:      []repository.Pet{},
		},
		{
			Model: gorm.Model{ID: 3},
			Person: model.Person{FirstName: "Jane",
				LastName: "Doe"},
			Address:   "456 Main St",
			City:      "Anytown",
			Telephone: "0987654321",
			Pets: []repository.Pet{
				{Model: gorm.Model{ID: 1}, Name: "Pet1"},
			},
		},
	}

	exceptionResponses := []Response{
		{ID: 1, FirstName: "John", LastName: "Doe", Address: "123 Main St", City: "Anytown", Telephone: "1234567890"},
		{ID: 2, FirstName: "Jane", LastName: "Doe", Address: "456 Main St", City: "Anytown", Telephone: "0987654321"},
		{ID: 3, FirstName: "Jane", LastName: "Doe", Address: "456 Main St", City: "Anytown", Telephone: "0987654321",
			Pets: &[]pet.Response{
				{ID: 1, Name: "Pet1"},
			}},
	}

	responses := FromOwners(owners)

	assert.Equal(t, len(exceptionResponses), len(responses))
	for i, owner := range exceptionResponses {
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
		result := &UpdateResponse{}
		result.FromUpdateEntity(tc.updateEntity)
		assert.Equal(t, result, tc.updateResponse)
	}
}
