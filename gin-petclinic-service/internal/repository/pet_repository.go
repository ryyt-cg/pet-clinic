package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type PetRepositorier interface {
	FindAll() ([]Pet, error)
	FindById(id uint) (*Pet, error)
	FindByIdWithVisits(id uint) (*Pet, error)
	FindByName(name string) ([]Pet, error)
	Insert(pet *Pet) (*Pet, error)
	Update(pet *Pet) (*Pet, error)
}

// PetRepository searches pet from the database
type PetRepository struct {
	db *gorm.DB
}

func NewPetRepository(db *gorm.DB) *PetRepository {
	return &PetRepository{
		db: db,
	}
}

func (repository *PetRepository) FindAll() ([]Pet, error) {
	log.Info().Msg("Retrieve all pets")

	var pets []Pet
	result := repository.db.Find(&pets)
	return pets, result.Error
}

/*
FindById
Join query reduces into 1 query by joining 2 tables together through LEFT JOIN.

SELECT "pets"."id","pets"."created_at","pets"."updated_at","pets"."deleted_at","pets"."name","pets"."birth_date",

	"pets"."type_id","pets"."owner_id","Type"."id" AS "Type__id","Type"."created_at" AS "Type__created_at",
	"Type"."updated_at" AS "Type__updated_at","Type"."deleted_at" AS "Type__deleted_at","Type"."name" AS "Type__name"
	FROM "pets" LEFT JOIN "species" "Type" ON "pets"."type_id" = "Type"."id"
	WHERE pets.id = 2 AND "pets"."deleted_at" IS NULL ORDER BY "pets"."id" LIMIT 1
*/
func (repository *PetRepository) FindById(id uint) (*Pet, error) {
	log.Info().Uint("id", id).Msg("Search pet by id.")

	var pet Pet
	result := repository.db.Joins("Species").Where("pets.id = ?", id).First(&pet)
	return &pet, result.Error
}

func (repository *PetRepository) FindByIdWithVisits(id uint) (*Pet, error) {
	log.Info().Uint("id", id).Msg("Search pet by id.")

	var pet Pet
	result := repository.db.Joins("Species").Preload("Visits").Where("pets.id = ?", id).First(&pet)
	return &pet, result.Error
}

/*
FindByName
SELECT * FROM "species" WHERE "species"."id" = 1 AND "species"."deleted_at" IS NULL
SELECT * FROM "pets" WHERE name = 'Leo' AND "pets"."deleted_at" IS NULL
*/
func (repository *PetRepository) FindByName(name string) ([]Pet, error) {
	log.Info().Str("name", name).Msg("Search pet by name.")

	var pets []Pet
	result := repository.db.Preload("Type").Where("name = ?", name).Find(&pets)
	return pets, result.Error
}

// Insert - insert a new pet
func (repository *PetRepository) Insert(pet *Pet) (*Pet, error) {
	log.Info().Str("name", pet.Name).Msg("Insert a new pet.")

	err := repository.db.Create(&pet).Error
	if err != nil {
		log.Error().Err(err).Msg("Fail to insert new pet.")
		return nil, err
	}
	return pet, err
}

// Update - update a pet
func (repository *PetRepository) Update(pet *Pet) (*Pet, error) {
	log.Info().Uint("id", pet.ID).Msg("Update vet by id.")

	// Omit the column name from update...
	err := repository.db.Omit("created_at").Save(&pet).Error
	if err != nil {
		log.Error().Err(err).Msg("Fail to update pet.")
		return nil, err
	}

	return pet, err
}
