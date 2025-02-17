package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// OwnerRepositorier - owner repository interface
// It composes the interfaces to interact with the database
type OwnerRepositorier interface {
	FindById(id uint) (*Owner, error)
	FindByLastName(lastName string) ([]Owner, error)
	FindAll() ([]Owner, error)
	FindByIdWithPets(id uint) (*Owner, error)
	Insert(owner *Owner) (*Owner, error)
	Update(owner *Owner) (*Owner, error)
}

// OwnerRepository searches owner from the database
type OwnerRepository struct {
	logger *zap.Logger
	pg     *gorm.DB
}

// NewOwnerRepository - OwnerRepository factory
func NewOwnerRepository(logger *zap.Logger, pg *gorm.DB) *OwnerRepository {
	return &OwnerRepository{
		logger: logger,
		pg:     pg,
	}
}

// FindById - find the owner by id
//
// SELECT * FROM "pets" WHERE "pets"."owner_id" = 3 AND "pets"."deleted_at" IS NULL
// SELECT * FROM `visits` WHERE `visits`.`pet_id` = 6 AND `visits`.`deleted_at` IS NULL
// SELECT * FROM "types" WHERE "types"."id" = 2 AND "types"."deleted_at" IS NULL
// SELECT * FROM "owners" WHERE "owners"."id" = 1 AND "owners"."deleted_at" IS NULL ORDER BY "owners"."id" LIMIT 1
func (repository *OwnerRepository) FindById(id uint) (*Owner, error) {
	repository.logger.Info("search owner by id.", zap.Uint("id", id))

	var owner Owner
	// Pets.Type - Nested Preloading (Eager Loading)
	err := repository.pg.First(&owner, id).Error
	if err != nil {
		repository.logger.Error("Fail to find owner by id.",
			zap.Uint("id", id), zap.String("error", err.Error()))
		return nil, err
	}
	return &owner, nil
}

// FindByLastName - find owner by last name
//
// SELECT * FROM "owners" WHERE last_name = 'Rodriquez' AND "owners"."deleted_at" IS NULL
func (repository *OwnerRepository) FindByLastName(lastName string) ([]Owner, error) {
	repository.logger.Info("Search owner by last name",
		zap.String("lastName", lastName))

	// Pets.Type - Nested Preloading (Eager Loading)
	var owners []Owner
	err := repository.pg.Where("last_name = ?", lastName).Find(&owners).Error
	if err != nil {
		repository.logger.Error("Fail to find owners by last name.",
			zap.String("lastName", lastName), zap.String("error", err.Error()))
		return nil, err
	}

	return owners, nil
}

// FindAll
// SELECT * FROM "owners" WHERE "owners"."deleted_at" IS NULL
func (repository *OwnerRepository) FindAll() ([]Owner, error) {
	repository.logger.Info("get list of owners")

	var owners []Owner
	// Get all owners
	// SELECT * FROM "owners" WHERE "owners"."deleted_at" IS NULL
	err := repository.pg.Find(&owners).Error
	if err != nil {
		repository.logger.Error("Fail to find all owners", zap.String("error", err.Error()))
		return nil, err
	}

	return owners, nil
}

// FindByIdWithPets Find all owners with pets & visits and its nested associations.
//
// SELECT * FROM "pets" WHERE "pets"."owner_id" IN (1,2,3,4,5,6,7,8,9,10) AND "pets"."deleted_at" IS NULL
// SELECT * FROM "types" WHERE "types"."id" IN (1,6,2,3,4,5) AND "types"."deleted_at" IS NULL
// SELECT * FROM "owners" WHERE "owners"."deleted_at" IS NULL
func (repository *OwnerRepository) FindByIdWithPets(id uint) (*Owner, error) {
	repository.logger.Info("get list of owners with pets")

	var owner Owner
	/*
	 clause.Associations won’t preload nested associations, but can use it with Nested Preloading together
	*/
	err := repository.pg.Preload("Pets.Type").Preload("Pets.Visits").Preload(clause.Associations).First(&owner, id).Error
	if err != nil {
		repository.logger.Error("Fail to find owners with pets by id.", zap.String("error", err.Error()))
		return nil, err
	}

	return &owner, nil
}

// Insert - insert a new owner
// INSERT INTO "owners" ("first_name","last_name","address","city","telephone","created_at","updated_at")
//
//	VALUES ('George','Franklin','110 W. Liberty St.','Madison','6085551023','2021-09-26 15:00:00','2021-09-26 15:00:00')
func (repository *OwnerRepository) Insert(owner *Owner) (*Owner, error) {
	repository.logger.Info("Fail a new owner.", zap.Any("owner", owner))
	now := time.Now()
	owner.UpdatedAt = now
	owner.CreatedAt = now

	err := repository.pg.Create(&owner).Error
	if err != nil {
		repository.logger.Error("Fail to insert new owner", zap.String("error", err.Error()))
		return nil, err
	}
	return owner, nil
}

// Update - update the owner
// UPDATE "owners" SET "first_name" = 'George', "last_name" = 'Franklin', "address" = '110 W. Liberty St.',
func (repository *OwnerRepository) Update(owner *Owner) (*Owner, error) {
	repository.logger.Info("Fail owner", zap.Uint("id", owner.ID))
	owner.UpdatedAt = time.Now()

	// Omit the column name from update...
	err := repository.pg.Omit("created_at").Save(&owner).Error
	if err != nil {
		repository.logger.Error("Fail to update owner.",
			zap.Uint("id", owner.ID), zap.String("error", err.Error()))
		return nil, err
	}

	return owner, nil
}
