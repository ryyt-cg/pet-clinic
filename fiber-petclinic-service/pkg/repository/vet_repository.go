package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type VetRepositorier interface {
	FindAllSpecialties() ([]Specialty, error)
	FindById(id int) (*Vet, error)
	FindByIdWithSpecialties(id int) (*Vet, error)
	FindByLastName(lastName string) ([]Vet, error)
	FindAll() ([]Vet, error)
	FindAllPreload() ([]Vet, error)
	Insert(vet *Vet) (*Vet, error)
	Update(vet *Vet) (*Vet, error)
}

// VetRepository searches vet from the database
type VetRepository struct {
	pg *gorm.DB
}

func NewVetRepository(pg *gorm.DB) *VetRepository {
	return &VetRepository{
		pg: pg,
	}
}

// FindAllSpecialties - retrieve all veterinarian's specialties
// SELECT * FROM "specialties" WHERE "specialties"."deleted_at" IS NULL
func (repository *VetRepository) FindAllSpecialties() ([]Specialty, error) {
	log.Info("get all specialties")

	var specialties []Specialty
	err := repository.pg.Find(&specialties).Error
	if err != nil {
		log.Error("Fail to find all specialties.", zap.String("error", err.Error()))
		return nil, err
	}
	return specialties, err
}

// FindById - find a veterinarian by id
/*
Many-to-Many relation vets <-- vet_specialties --> specialties
Preload - eagerly loading the specialties relation and hence,
3 queries are executed.
SELECT * FROM "vet_specialties" WHERE "vet_specialties"."vet_id" = 2
SELECT * FROM "specialties" WHERE "specialties"."id" = 1 AND "specialties"."deleted_at" IS NULL
SELECT * FROM "vets" WHERE "vets"."id" = 2 AND "vets"."deleted_at" IS NULL ORDER BY "vets"."id" LIMIT 1
*/
func (repository *VetRepository) FindById(id int) (*Vet, error) {
	log.Info("Search vet by id.", zap.Int("id", id))

	var vet Vet
	err := repository.pg.First(&vet, id).Error
	if err != nil {
		log.Error("Fail to find vet by id.",
			zap.Int("id", id), zap.String("error", err.Error()))
		return nil, err
	}

	return &vet, err
}

func (repository *VetRepository) FindByIdWithSpecialties(id int) (*Vet, error) {
	log.Info("Search vet by id.", zap.Int("id", id))

	var vet Vet
	err := repository.pg.Preload("Specialties").First(&vet, id).Error
	if err != nil {
		log.Error("Fail to find vet by id.",
			zap.Int("id", id), zap.String("error", err.Error()))
		return nil, err
	}

	return &vet, err
}

// FindByLastName
// SELECT * FROM "vet_specialties" WHERE "vet_specialties"."vet_id" = 5
// SELECT * FROM "specialties" WHERE "specialties"."id" = 1 AND "specialties"."deleted_at" IS NULL
// SELECT * FROM "vets" WHERE last_name = 'Stevens' AND "vets"."deleted_at" IS NULL
func (repository *VetRepository) FindByLastName(lastName string) ([]Vet, error) {
	log.Info("Search vet by last name.", zap.String("lastName", lastName))

	var vets []Vet
	err := repository.pg.Preload("Specialties").Where("last_name = ?", lastName).Find(&vets).Error
	if err != nil {
		log.Error("Fail to find vet by last name.",
			zap.String("lastName", lastName), zap.String("error", err.Error()))
		return nil, err
	}

	return vets, nil
}

/*
FindAll
SELECT * FROM "vets" WHERE "vets"."deleted_at" IS NULL
*/
func (repository *VetRepository) FindAll() ([]Vet, error) {
	log.Info("get all vets")

	var vets []Vet
	err := repository.pg.Find(&vets).Error
	if err != nil {
		log.Error("Fail to find all vets.", zap.String("error", err.Error()))
		return nil, err
	}
	return vets, err
}

/*
FindAllPreload
Expect there would be n + 1 queries executed just like in Java Hibernate.  Surprise,
only 3 queries actually executed.

	SELECT * FROM "vet_specialties" WHERE "vet_specialties"."vet_id" IN (1,2,3,4,5,6)
	SELECT * FROM "specialties" WHERE "specialties"."id" IN (1,2,3) AND "specialties"."deleted_at" IS NULL
	SELECT * FROM "vets" WHERE "vets"."deleted_at" IS NULL
*/
func (repository *VetRepository) FindAllPreload() ([]Vet, error) {
	log.Info("get all vets with relations preloaded.")

	var vets []Vet
	err := repository.pg.Preload("Specialties").Find(&vets).Error
	if err != nil {
		log.Error("Fail to find all vets.", zap.String("error", err.Error()))
		return nil, err
	}
	return vets, err
}

// Insert - insert a new vet
// INSERT INTO "vets" ("first_name","last_name","created_at","updated_at")
//
//	VALUES ('James','Carter','2021-07-25 15:00:00','2021-07-25 15:00:00') RETURN
func (repository *VetRepository) Insert(vet *Vet) (*Vet, error) {
	log.Info("Fail a new vet.", zap.Any("vet", vet))

	err := repository.pg.Create(&vet).Error
	if err != nil {
		log.Error("Fail to insert new vet.", zap.String("error", err.Error()))
		return nil, err
	}
	return vet, err
}

// Update - update vet
// UPDATE "vets" SET "first_name" = 'James', "last_name" = 'Carter' WHERE "id" = 1
func (repository *VetRepository) Update(vet *Vet) (*Vet, error) {
	log.Info("Fail vet id.", zap.Uint("id", vet.ID))

	// Omit the column name from update...
	err := repository.pg.Omit("created_at").Save(&vet).Error
	if err != nil {
		log.Error("Fail to update vet.", zap.String("error", err.Error()))
		return nil, err
	}

	return vet, err
}
