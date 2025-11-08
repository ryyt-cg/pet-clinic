package repository

import (
	"fiber3-petclinic-service/internal/repository/model"
	"fiber3-petclinic-service/internal/test"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestOwnerRepository_FindById(t *testing.T) {
	gdb, mock := test.NewMockPostgresDB()
	expected := &Owner{
		Model: gorm.Model{
			ID: 2,
		},
		Person: model.Person{
			FirstName: "Betty",
			LastName:  "Davis",
		},
		Address:   "123 Main Street",
		City:      "New York",
		Telephone: "0123456789",
	}
	row := sqlmock.NewRows([]string{"id", "first_name", "last_name", "address", "city", "telephone", "created_at", "updated_at", "deleted_at"}).
		AddRow(expected.ID, expected.FirstName, expected.LastName, expected.Address, expected.City, expected.Telephone, expected.CreatedAt, expected.UpdatedAt, expected.DeletedAt)
	mock.ExpectQuery(`SELECT (.+) FROM "owners" WHERE "owners"."id" = (.+)`).WithArgs(expected.ID, 1).WillReturnRows(row)

	ownerRepo := NewOwnerRepository(gdb)
	result, err := ownerRepo.FindById(expected.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, result, expected)

}

func TestOwnerRepository_FindByIdWithPets(t *testing.T) {
	expected := &Owner{
		Model: gorm.Model{
			ID: 2,
		},
		Person: model.Person{
			FirstName: "Peter",
			LastName:  "McTavish",
		},
		Address:   "123 Main Street",
		City:      "New York",
		Telephone: "0123456789",
		Pets: []Pet{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Name:      "Peter",
				Birthdate: test.ToDate("2021-05-19"),
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
				Name:      "John",
				Birthdate: test.ToDate("2021-06-10"),
			},
		},
	}

	gdb, mock := test.NewMockPostgresDB()
	row := sqlmock.NewRows([]string{"id", "first_name", "last_name", "address", "city", "telephone", "created_at", "updated_at", "deleted_at"}).
		AddRow(expected.ID, expected.FirstName, expected.LastName, expected.Address, expected.City, expected.Telephone, expected.CreatedAt, expected.UpdatedAt, expected.DeletedAt)
	mock.ExpectQuery(`SELECT (.+) FROM "owners" WHERE "owners"."id" = (.+)`).WithArgs(expected.ID, 1).WillReturnRows(row)

	ownerRepo := NewOwnerRepository(gdb)
	result, err := ownerRepo.FindByIdWithPets(expected.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, result, expected)
}

func TestOwnerRepository_FindByLastName(t *testing.T) {
	expected := &Owner{
		Model: gorm.Model{
			ID: 20,
		},
		Person: model.Person{
			FirstName: "Jean",
			LastName:  "Coleman",
		},
	}

	gdb, mock := test.NewMockPostgresDB()
	row := sqlmock.NewRows([]string{"id", "first_name", "last_name", "address", "city", "telephone", "created_at", "updated_at", "deleted_at"}).
		AddRow(expected.ID, expected.FirstName, expected.LastName, expected.Address, expected.City, expected.Telephone, expected.CreatedAt, expected.UpdatedAt, expected.DeletedAt)
	mock.ExpectQuery(`SELECT (.+) FROM "owners" WHERE "owners"."last_name" = (.+)`).WithArgs(expected.LastName, 1).WillReturnRows(row)

	ownerRepo := NewOwnerRepository(gdb)
	result, err := ownerRepo.FindByLastName(expected.LastName)
	assert.Equal(t, err, nil)
	assert.Equal(t, result, expected)
}

func TestOwnerRepository_Insert(t *testing.T) {
	expected := &Owner{
		Person: model.Person{
			FirstName: "David",
			LastName:  "Schroeder",
		},
	}

	gdb, mock := test.NewMockPostgresDB()
	row := sqlmock.NewRows([]string{"id", "first_name", "last_name", "address", "city", "telephone", "created_at", "updated_at", "deleted_at"}).
		AddRow(expected.ID, expected.FirstName, expected.LastName, expected.Address, expected.City, expected.Telephone, expected.CreatedAt, expected.UpdatedAt, expected.DeletedAt)
	mock.ExpectQuery(`SELECT (.+) FROM "owners" WHERE "owners"."id" = (.+)`).WithArgs(expected.ID, 1).WillReturnRows(row)

	ownerRepo := NewOwnerRepository(gdb)
	result, err := ownerRepo.FindById(expected.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, result, expected)
}

func TestOwnerRepository_Update(t *testing.T) {

}
