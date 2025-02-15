package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PetRepositorier interface {
	FindAll() ([]Pet, error)
	FindById(id int) (*Pet, error)
	FindByIdWithVisits(id int) (*Pet, error)
	FindByName(name string) ([]Pet, error)
	Insert(pet *Pet) (*Pet, error)
	Update(pet *Pet) (*Pet, error)
}

// PetRepository searches pet from the database
type PetRepository struct {
	logger *zap.Logger
	pg     *gorm.DB
}

func NewPetRepository(logger *zap.Logger, pg *gorm.DB) *PetRepository {
	return &PetRepository{
		logger: logger,
		pg:     pg,
	}
}

func (repository *PetRepository) FindAll() ([]Pet, error) {
	repository.logger.Info("Retrieve all pets")

	var pets []Pet
	result := repository.pg.Find(&pets)
	return pets, result.Error
}

/*
FindById
Join query reduces into 1 query by joining 2 tables together through LEFT JOIN.

SELECT "pets"."id","pets"."created_at","pets"."updated_at","pets"."deleted_at","pets"."name","pets"."birth_date",

	"pets"."type_id","pets"."owner_id","Type"."id" AS "Type__id","Type"."created_at" AS "Type__created_at",
	"Type"."updated_at" AS "Type__updated_at","Type"."deleted_at" AS "Type__deleted_at","Type"."name" AS "Type__name"
	FROM "pets" LEFT JOIN "types" "Type" ON "pets"."type_id" = "Type"."id"
	WHERE pets.id = 2 AND "pets"."deleted_at" IS NULL ORDER BY "pets"."id" LIMIT 1
*/
func (repository *PetRepository) FindById(id int) (*Pet, error) {
	repository.logger.Info("Search pet by id.", zap.Int("id", id))

	var pet Pet
	result := repository.pg.Joins("Type").Where("pets.id = ?", id).First(&pet)
	return &pet, result.Error
}

func (repository *PetRepository) FindByIdWithVisits(id int) (*Pet, error) {
	repository.logger.Info("Search pet by id.", zap.Int("id", id))

	var pet Pet
	result := repository.pg.Joins("Type").Preload("Visits").Where("pets.id = ?", id).First(&pet)
	return &pet, result.Error
}

/*
FindByName
SELECT * FROM "types" WHERE "types"."id" = 1 AND "types"."deleted_at" IS NULL
SELECT * FROM "pets" WHERE name = 'Leo' AND "pets"."deleted_at" IS NULL
*/
func (repository *PetRepository) FindByName(name string) ([]Pet, error) {
	repository.logger.Info("Search pet by name.", zap.String("name", name))

	var pets []Pet
	result := repository.pg.Preload("Type").Where("name = ?", name).Find(&pets)
	return pets, result.Error
}

// Insert - insert a new pet
func (repository *PetRepository) Insert(pet *Pet) (*Pet, error) {
	repository.logger.Info("insert a new pet.", zap.String("name", pet.Name))

	err := repository.pg.Create(&pet).Error
	if err != nil {
		repository.logger.Error("fail to insert new pet.", zap.String("error", err.Error()))
		return nil, err
	}
	return pet, err
}

// Update - update a pet
func (repository *PetRepository) Update(pet *Pet) (*Pet, error) {
	repository.logger.Info("update vet id.", zap.Uint("id", pet.ID))

	// Omit the column name from update...
	err := repository.pg.Omit("created_at").Save(&pet).Error
	if err != nil {
		repository.logger.Error("fails to update pet.", zap.String("error", err.Error()))
		return nil, err
	}

	return pet, err
}
