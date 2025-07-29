package pet

import (
	"fiber-petclinic-service/pkg/repository"
	"fiber-petclinic-service/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ToPet(t *testing.T) {
	petRequest := &Request{
		Name:      "Five Miles",
		Birthdate: "2017-05-03",
		TypeID:    1,
		OwnerID:   20,
	}
	pet := &repository.Pet{
		Name:      "Five Miles",
		Birthdate: *test.ToDate("2017-05-03"),
		TypeID:    1,
		OwnerID:   20,
	}

	result := ToPet(petRequest)
	assert.Equal(t, pet.Name, result.Name)
	assert.Equal(t, pet.Birthdate, result.Birthdate)
	assert.Equal(t, pet, result)
}
