package repository

import (
	"fiber-petclinic-service/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVisitFields(t *testing.T) {
	visit := &Visit{
		VisitDate:   test.ToDate("2021-01-01"),
		Description: "Regular checkup",
		PetID:       1,
	}

	assert.Equal(t, "Regular checkup", visit.Description)
	assert.Equal(t, uint(1), visit.PetID)
}
