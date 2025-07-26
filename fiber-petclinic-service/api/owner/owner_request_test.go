package owner

import (
	"fiber-petclinic-service/pkg/repository"
	"fiber-petclinic-service/pkg/repository/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func Test_ToOwner(t *testing.T) {
	ownerRequest := &AddRequest{
		FirstName: "John",
		LastName:  "Doe",
		Address:   "123 Main St",
		City:      "Anytown",
		Telephone: "1234567890",
	}

	owner := ToOwnerEntity(ownerRequest)

	assert.Equal(t, ownerRequest.FirstName, owner.Person.FirstName)
	assert.Equal(t, ownerRequest.LastName, owner.Person.LastName)
	assert.Equal(t, ownerRequest.Address, owner.Address)
	assert.Equal(t, ownerRequest.City, owner.City)
	assert.Equal(t, ownerRequest.Telephone, owner.Telephone)
}

func Test_ToOwnerWithEmptyFields(t *testing.T) {
	ownerRequest := &AddRequest{
		FirstName: "",
		LastName:  "",
		Address:   "",
		City:      "",
		Telephone: "",
	}

	owner := ToOwnerEntity(ownerRequest)

	assert.Equal(t, ownerRequest.FirstName, owner.Person.FirstName)
	assert.Equal(t, ownerRequest.LastName, owner.Person.LastName)
	assert.Equal(t, ownerRequest.Address, owner.Address)
	assert.Equal(t, ownerRequest.City, owner.City)
	assert.Equal(t, ownerRequest.Telephone, owner.Telephone)
}

func Test_ToOwnerWithPartialFields(t *testing.T) {
	ownerRequest := &AddRequest{
		FirstName: "John",
		LastName:  "Doe",
		Address:   "",
		City:      "",
		Telephone: "",
	}

	owner := ToOwnerEntity(ownerRequest)

	assert.Equal(t, ownerRequest.FirstName, owner.Person.FirstName)
	assert.Equal(t, ownerRequest.LastName, owner.Person.LastName)
	assert.Equal(t, ownerRequest.Address, owner.Address)
	assert.Equal(t, ownerRequest.City, owner.City)
	assert.Equal(t, ownerRequest.Telephone, owner.Telephone)
}

func Test_ToOwnerEntityFromUpdateRequest(t *testing.T) {
	tests := []struct {
		updateRequest *UpdateRequest
		ownerEntity   *repository.Owner
	}{
		{updateRequest: &UpdateRequest{
			Id:        101,
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
		{updateRequest: &UpdateRequest{}, ownerEntity: &repository.Owner{}},
	}

	for _, tc := range tests {
		result := ToOwnerEntityFromUpdateRequest(tc.updateRequest)
		assert.Equal(t, tc.ownerEntity, result)
	}

}
