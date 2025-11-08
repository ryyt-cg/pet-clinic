package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPetFields(t *testing.T) {
	pet := &Pet{
		Name:    "Fido",
		Species: Species{Name: "Dog"},
	}

	assert.Equal(t, "Fido", pet.Name)
	assert.Equal(t, "Dog", pet.Species.Name)
}

func TestPetFieldsWithEmptyName(t *testing.T) {
	pet := &Pet{
		Name:    "",
		Species: Species{Name: "Dog"},
	}

	assert.Equal(t, "", pet.Name)
	assert.Equal(t, "Dog", pet.Species.Name)
}

func TestPetFieldsWithEmptyType(t *testing.T) {
	pet := &Pet{
		Name:    "Fido",
		Species: Species{Name: ""},
	}

	assert.Equal(t, "Fido", pet.Name)
	assert.Equal(t, "", pet.Species.Name)
}
