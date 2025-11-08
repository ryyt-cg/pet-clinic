package repository

import (
	"errors"
	"fiber3-petclinic-service/internal/test"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestPetRepository_FindById
// Three outcomes when finding pet by id
// 1. Found a pet with provided id
// 2. Found no pet with provided id
// 3. System error caused by number things
func TestPetRepository_FindById(t *testing.T) {
	// setting up the mock
	gdb, mock := test.NewMockPostgresDB()

	testCases := []struct {
		name          string
		id            uint
		expected      *Pet
		errorExpected error
	}{
		{
			name: "Found a pet by id",
			id:   1,
			expected: &Pet{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: *test.ToDate("2025-11-02"),
					UpdatedAt: *test.ToDate("2025-11-02"),
				},
				Name:      "Leo",
				Birthdate: test.ToDate("2020-09-07"),
				SpeciesID: 1,
				OwnerID:   2,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			row := sqlmock.NewRows([]string{"id", "name", "birthdate", "species_id", "owner_id", "created_at", "updated_at", "deleted_at"})
			row.AddRow(tc.expected.ID, tc.expected.Name, tc.expected.Birthdate, tc.expected.SpeciesID, tc.expected.OwnerID,
				tc.expected.CreatedAt, tc.expected.UpdatedAt, tc.expected.DeletedAt)
			mock.ExpectQuery(`SELECT (.+) FROM "pets" LEFT JOIN "species" "Species" ON "pets"."species_id" = "Species"."id" AND "Species"."deleted_at" IS NULL (.+)`).
				WithArgs(tc.id, 1).WillReturnRows(row)

			petRepo := NewPetRepository(gdb)
			result, err := petRepo.FindById(tc.id)
			require.Equal(t, err, tc.errorExpected)
			require.Equal(t, result, tc.expected)
		})
	}
}

func TestPetRepository_FindAll(t *testing.T) {
	gdb, mock := test.NewMockPostgresDB()

	var pets []Pet

	testCases := []struct {
		name          string
		expected      []Pet
		errorExpected error
	}{
		{
			name: "All pets",
			expected: []Pet{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: *test.ToDate("2025-11-02"),
						UpdatedAt: *test.ToDate("2025-11-02"),
					},
					Name:      "Leo",
					Birthdate: test.ToDate("2020-09-07"),
					SpeciesID: 1,
					OwnerID:   2,
				},
				{
					Model: gorm.Model{
						ID:        2,
						CreatedAt: *test.ToDate("2025-11-04"),
						UpdatedAt: *test.ToDate("2025-11-04"),
					},
					Name:      "Downy",
					Birthdate: test.ToDate("2022-02-02"),
					SpeciesID: 2,
					OwnerID:   4,
				},
			},
		},
		{
			name:          "No pets",
			expected:      pets,
			errorExpected: gorm.ErrRecordNotFound,
		},
		{
			name:          "Unknown error",
			expected:      pets,
			errorExpected: errors.New("unknown error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"id", "name", "birthdate", "species_id", "owner_id", "created_at", "updated_at", "deleted_at"})
			if tc.errorExpected != nil {
				mock.ExpectQuery(`SELECT (.+) FROM "pets" (.+)`).WithArgs().WillReturnError(tc.errorExpected)
			} else {
				rows.AddRow(tc.expected[0].ID, tc.expected[0].Name, tc.expected[0].Birthdate, tc.expected[0].SpeciesID, tc.expected[0].OwnerID, tc.expected[0].CreatedAt, tc.expected[0].UpdatedAt, tc.expected[0].DeletedAt)
				rows.AddRow(tc.expected[1].ID, tc.expected[1].Name, tc.expected[1].Birthdate, tc.expected[1].SpeciesID, tc.expected[1].OwnerID, tc.expected[1].CreatedAt, tc.expected[1].UpdatedAt, tc.expected[1].DeletedAt)
				mock.ExpectQuery(`SELECT (.+) FROM "pets" (.+)`).WithArgs().WillReturnRows(rows)
			}

			petRepo := NewPetRepository(gdb)
			result, err := petRepo.FindAll()
			require.Equal(t, err, tc.errorExpected)
			require.Equal(t, result, tc.expected)
		})
	}
}

func TestPetRepository_FindByIdWithVisits(t *testing.T) {
	expected := &Pet{
		Model: gorm.Model{
			ID:        5,
			CreatedAt: *test.ToDate("2025-11-01"),
			UpdatedAt: *test.ToDate("2025-11-01"),
		},
		Name:      "Johnny",
		Birthdate: test.ToDate("2022-09-07"),
		SpeciesID: 1,
		OwnerID:   2,
	}
	gdb, mock := test.NewMockPostgresDB()

	row := sqlmock.NewRows([]string{"id", "name", "birthdate", "species_id", "owner_id", "created_at", "updated_at", "deleted_at"})
	row.AddRow(expected.ID, expected.Name, expected.Birthdate, expected.SpeciesID, expected.OwnerID, expected.CreatedAt, expected.UpdatedAt, expected.DeletedAt)
	mock.ExpectQuery(`SELECT (.+) FROM "pets" LEFT JOIN "species" ON "pets"."species_id" = "species"."id" (.+) WHERE (.+)`).WithArgs(5, 1).WillReturnRows(row)

	petRepo := NewPetRepository(gdb)
	result, err := petRepo.FindByIdWithVisits(5)
	require.Nil(t, err)
	require.Equal(t, result, expected)
}

func TestPetRepository_FindByName(t *testing.T) {

}

func TestPetRepository_Insert(t *testing.T) {

}

func TestPetRepository_Update(t *testing.T) {}
