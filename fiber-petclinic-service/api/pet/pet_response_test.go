package pet

import (
	"fiber-petclinic-service/api/visit"
	"fiber-petclinic-service/internal/repository"
	"fiber-petclinic-service/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_FromPet(t *testing.T) {
	testCases := []struct {
		name     string
		pet      *repository.Pet
		expected Response
	}{
		{
			name: "Valid Pet",
			pet: &repository.Pet{
				Model: gorm.Model{
					ID: 1,
				},
				Name:      "Buddy",
				Birthdate: nil,
				Visits:    []repository.Visit{}, // No visits
			},
			expected: Response{
				ID:        1,
				Name:      "Buddy",
				Birthdate: "",                 // Empty birthdate
				Visits:    []visit.Response{}, // No visits
			},
		},
		{
			name: "Pet with Birthdate",
			pet: &repository.Pet{
				Model: gorm.Model{
					ID: 2,
				},
				Name:      "Mittens",
				Birthdate: test.ToDate("2023-01-01"),
				Visits:    []repository.Visit{}, // No visits
			},
			expected: Response{
				ID:        2,
				Name:      "Mittens",
				Birthdate: "2023-01-01",       // Formatted birthdate
				Visits:    []visit.Response{}, // No visits
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var response Response
			response.FromPet(tc.pet)
			assert.Equal(t, tc.expected, response)

		})
	}
}

func Test_FromPets(t *testing.T) {}
