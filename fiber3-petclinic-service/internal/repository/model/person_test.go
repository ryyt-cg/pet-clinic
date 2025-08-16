package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPersonCreation(t *testing.T) {
	person := &Person{
		FirstName: "John",
		LastName:  "Doe",
	}

	assert.Equal(t, "John", person.FirstName)
	assert.Equal(t, "Doe", person.LastName)
}

func TestPersonCreationWithEmptyFirstName(t *testing.T) {
	person := &Person{
		FirstName: "",
		LastName:  "Doe",
	}

	assert.Equal(t, "", person.FirstName)
	assert.Equal(t, "Doe", person.LastName)
}

func TestPersonCreationWithEmptyLastName(t *testing.T) {
	person := &Person{
		FirstName: "John",
		LastName:  "",
	}

	assert.Equal(t, "John", person.FirstName)
	assert.Equal(t, "", person.LastName)
}

func TestPersonCreationWithEmptyFirstNameAndLastName(t *testing.T) {
	person := &Person{
		FirstName: "",
		LastName:  "",
	}

	assert.Equal(t, "", person.FirstName)
	assert.Equal(t, "", person.LastName)
}
