package owner

import (
	"gin-petclinic-service/internal/repository"
	"gin-petclinic-service/internal/repository/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_ToOwner(t *testing.T) {
	ownerRequest := &addRequest{
		FirstName: "John",
		LastName:  "Doe",
		Address:   "123 Main St",
		City:      "Anytown",
		Telephone: "1234567890",
	}

	owner := toOwnerEntity(ownerRequest)

	assert.Equal(t, ownerRequest.FirstName, owner.Person.FirstName)
	assert.Equal(t, ownerRequest.LastName, owner.Person.LastName)
	assert.Equal(t, ownerRequest.Address, owner.Address)
	assert.Equal(t, ownerRequest.City, owner.City)
	assert.Equal(t, ownerRequest.Telephone, owner.Telephone)
}

func Test_ToOwnerEntityWithNil(t *testing.T) {
	owner := toOwnerEntity(nil)
	assert.Nil(t, owner)
}

func Test_ToOwnerWithEmptyFields(t *testing.T) {
	ownerRequest := &addRequest{
		FirstName: "",
		LastName:  "",
		Address:   "",
		City:      "",
		Telephone: "",
	}

	owner := toOwnerEntity(ownerRequest)

	assert.Equal(t, ownerRequest.FirstName, owner.Person.FirstName)
	assert.Equal(t, ownerRequest.LastName, owner.Person.LastName)
	assert.Equal(t, ownerRequest.Address, owner.Address)
	assert.Equal(t, ownerRequest.City, owner.City)
	assert.Equal(t, ownerRequest.Telephone, owner.Telephone)
}

func Test_ToOwnerWithPartialFields(t *testing.T) {
	ownerRequest := &addRequest{
		FirstName: "John",
		LastName:  "Doe",
		Address:   "",
		City:      "",
		Telephone: "",
	}

	owner := toOwnerEntity(ownerRequest)

	assert.Equal(t, ownerRequest.FirstName, owner.Person.FirstName)
	assert.Equal(t, ownerRequest.LastName, owner.Person.LastName)
	assert.Equal(t, ownerRequest.Address, owner.Address)
	assert.Equal(t, ownerRequest.City, owner.City)
	assert.Equal(t, ownerRequest.Telephone, owner.Telephone)
}

func Test_ToOwnerEntityFromUpdateRequest(t *testing.T) {
	tests := []struct {
		inputUpdateRequest *updateRequest
		ownerEntity        *repository.Owner
	}{
		{inputUpdateRequest: &updateRequest{
			ID:        101,
			FirstName: "John",
			LastName:  "Five",
			Address:   "123 Main St",
			City:      "Anytown",
			Telephone: "1234567890",
		}, ownerEntity: &repository.Owner{
			Model: gorm.Model{
				ID: 101,
			},
			Person: model.Person{
				FirstName: "John",
				LastName:  "Five",
			},
			Address:   "123 Main St",
			City:      "Anytown",
			Telephone: "1234567890",
		}},
		{inputUpdateRequest: &updateRequest{}, ownerEntity: &repository.Owner{}},
		{inputUpdateRequest: nil, ownerEntity: nil},
	}

	for _, tc := range tests {
		result := toOwnerEntityFromUpdateRequest(tc.inputUpdateRequest)
		assert.Equal(t, tc.ownerEntity, result)
	}

}
