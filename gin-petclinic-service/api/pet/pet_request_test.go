package pet

import (
	"gin-petclinic-service/internal/repository"
	"gin-petclinic-service/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_ToPet(t *testing.T) {
	petRequest := &request{
		Name:      "Five Miles",
		Birthdate: "2017-05-03",
		SpeciesID: 1,
		OwnerID:   20,
	}
	pet := &repository.Pet{
		Name:      "Five Miles",
		Birthdate: test.ToDate("2017-05-03"),
		SpeciesID: 1,
		OwnerID:   20,
	}

	result, err := toPet(petRequest)
	assert.NoError(t, err, "Expected no error when converting to Pet")
	assert.Equal(t, pet, result)
}

func Test_FromAddRequest(t *testing.T) {
	testCases := []struct {
		name        string
		petRequest  *addRequest
		expectedPet *repository.Pet
	}{
		{
			name: "Valid AddRequest",
			petRequest: &addRequest{
				Name:      "Five Miles",
				Birthdate: "2017-05-03",
				SpeciesID: 1,
				OwnerID:   20,
			},
			expectedPet: &repository.Pet{
				Name:      "Five Miles",
				Birthdate: test.ToDate("2017-05-03"),
				SpeciesID: 1,
				OwnerID:   20,
			},
		},
		{
			name: "Invalid Birthdate",
			petRequest: &addRequest{
				Name:      "Invalid Date",
				Birthdate: "invalid-date",
				SpeciesID: 1,
				OwnerID:   20,
			},
			expectedPet: nil, // Expecting an error, so no expected pet
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, _ := fromAddRequest(tc.petRequest)

			if tc.expectedPet == nil {
				assert.Nil(t, nil)
			} else {
				assert.Equal(t, tc.expectedPet, result)
			}
		})
	}
}

func Test_FromUpdateRequest(t *testing.T) {
	testCases := []struct {
		name        string
		petRequest  *updateRequest
		expectedPet *repository.Pet
	}{
		{
			name: "Valid Update Request",
			petRequest: &updateRequest{
				ID:        1,
				Name:      "Five Miles",
				Birthdate: "2017-05-03",
				SpeciesID: 1,
				OwnerID:   20,
			},
			expectedPet: &repository.Pet{
				Model: gorm.Model{
					ID: 1,
				},
				Name:      "Five Miles",
				Birthdate: test.ToDate("2017-05-03"),
				SpeciesID: 1,
				OwnerID:   20,
			},
		},
		{
			name: "Invalid Birthdate",
			petRequest: &updateRequest{
				Name:      "Invalid Date",
				Birthdate: "invalid-date",
				SpeciesID: 1,
				OwnerID:   20,
			},
			expectedPet: nil, // Expecting an error, so no expected pet
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := fromUpdateRequest(tc.petRequest)
			if tc.expectedPet == nil {
				assert.Error(t, err, "Expected an error for invalid request")
			} else {
				assert.NoError(t, err, "Expected no error for valid request")
				assert.Equal(t, tc.expectedPet, result)
			}
		})
	}
}
