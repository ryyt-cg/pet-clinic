package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisitFields(t *testing.T) {
	visit := &Visit{
		VisitDate:   "2021-01-01",
		Description: "Regular checkup",
		PetID:       1,
	}

	assert.Equal(t, "Regular checkup", visit.Description)
	assert.Equal(t, 1, visit.PetID)
}

func TestVisitFieldsWithEmptyDescription(t *testing.T) {
	visit := &Visit{
		VisitDate:   "2022-04-01",
		Description: "",
		PetID:       1,
	}

	assert.Equal(t, "", visit.Description)
	assert.Equal(t, 1, visit.PetID)
}

func TestVisitFieldsWithNoPetID(t *testing.T) {
	visit := &Visit{
		VisitDate:   "2022-04-23",
		Description: "Regular checkup",
		PetID:       0,
	}

	assert.Equal(t, "Regular checkup", visit.Description)
	assert.Equal(t, 0, visit.PetID)
}

func TestVisitFieldsWithEmptyDescriptionAndNoPetID(t *testing.T) {
	visit := &Visit{
		VisitDate:   "2022-07-23",
		Description: "",
		PetID:       0,
	}

	assert.Equal(t, "", visit.Description)
	assert.Equal(t, 0, visit.PetID)
}
