package vet

import (
	"fiber-petclinic-service/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_FromSpecialty(t *testing.T) {
	specialty := &repository.Specialty{Model: gorm.Model{ID: 1}, Name: "Cardiology"}
	s := &specialtyResponse{}
	s.FromSpecialty(specialty)

	assert.Equal(t, specialty.ID, s.ID)
	assert.Equal(t, specialty.Name, s.Name)
}

func Test_ToSpecialtyResponses(t *testing.T) {
	testCases := []struct {
		name        string
		specialties []repository.Specialty
		expected    *[]specialtyResponse
	}{
		{
			name: "non-empty specialties",
			specialties: []repository.Specialty{
				{Model: gorm.Model{ID: 1}, Name: "Cardiology"},
				{Model: gorm.Model{ID: 2}, Name: "Dentistry"},
				{Model: gorm.Model{ID: 3}, Name: "Empty Specialty"},
			},
			expected: &[]specialtyResponse{
				{ID: 1, Name: "Cardiology"},
				{ID: 2, Name: "Dentistry"},
				{ID: 3, Name: "Empty Specialty"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToSpecialtyResponses(tc.specialties)
			assert.Equal(t, len(*result), len(*tc.expected))
			for i, specialty := range *result {
				assert.Equal(t, (*tc.expected)[i].ID, specialty.ID)
				assert.Equal(t, (*tc.expected)[i].Name, specialty.Name)
			}
		})
	}
}

func Test_ToSpecialtyResponsesWithEmptySlice(t *testing.T) {
	testCases := []struct {
		name        string
		specialties []repository.Specialty
		expected    *[]specialtyResponse
	}{
		{
			name:        "empty specialties",
			specialties: []repository.Specialty{},
			expected:    nil,
		},
		{
			name: "has some specialties",
			specialties: []repository.Specialty{
				{Model: gorm.Model{ID: 1}, Name: "Cardiology"},
				{Model: gorm.Model{ID: 2}, Name: "Dentistry"},
				{Model: gorm.Model{ID: 3}, Name: "Empty Specialty"},
			},
			expected: &[]specialtyResponse{
				{ID: 1, Name: "Cardiology"},
				{ID: 2, Name: "Dentistry"},
				{ID: 3, Name: "Empty Specialty"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToSpecialtyResponses(tc.specialties)
			if tc.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, len(*result), len(*tc.expected))
				for i, specialty := range *result {
					assert.Equal(t, (*tc.expected)[i].ID, specialty.ID)
					assert.Equal(t, (*tc.expected)[i].Name, specialty.Name)
				}
			}
		})
	}
}

func Test_ToSpecialties(t *testing.T) {
	specialties := []specialtyResponse{
		{ID: 1, Name: "Cardiology"},
		{ID: 2, Name: "Dentistry"},
		{ID: 3, Name: "Empty Specialty"},
	}

	result := ToSpecialties(specialties)

	assert.Equal(t, len(specialties), len(*result))
	for i, specialty := range *result {
		assert.Equal(t, specialties[i].ID, specialty.ID)
		assert.Equal(t, specialties[i].Name, specialty.Name)
	}
}

func Test_ToSpecialtiesWithEmptySlice(t *testing.T) {
	specialties := []specialtyResponse{}

	result := ToSpecialties(specialties)
	assert.Equal(t, 0, len(*result))
}
