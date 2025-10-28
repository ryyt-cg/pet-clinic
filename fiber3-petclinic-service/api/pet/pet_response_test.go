package pet

import (
	"fiber3-petclinic-service/api/visit"
	"fiber3-petclinic-service/internal/repository"
	"fiber3-petclinic-service/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_Response(t *testing.T) {
	testCases := []struct {
		name     string
		pet      *repository.Pet
		expected *Response
	}{
		{
			name: "Valid Pet",
			pet: &repository.Pet{
				Model: gorm.Model{
					ID: 1,
				},
				Name:      "Buddy",
				Birthdate: nil,
			},
			expected: &Response{
				ID:        1,
				Name:      "Buddy",
				Birthdate: "", // Empty birthdate
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
			},
			expected: &Response{
				ID:        2,
				Name:      "Mittens",
				Birthdate: "2023-01-01", // Formatted birthdate
			},
		},
		{
			name:     "Nil Pet",
			pet:      nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := ToResponse(tc.pet)
			assert.Equal(t, tc.expected, response)
		})
	}
}

func Test_FromPets(t *testing.T) {

}
