package owner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToOwner(t *testing.T) {
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

func TestToOwnerWithEmptyFields(t *testing.T) {
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

func TestToOwnerWithPartialFields(t *testing.T) {
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
