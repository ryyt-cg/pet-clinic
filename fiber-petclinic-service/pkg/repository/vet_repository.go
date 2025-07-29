package repository

import (
	"github.com/rs/zerolog/log"
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
	db *gorm.DB
}

func NewVetRepository(db *gorm.DB) *VetRepository {
	return &VetRepository{
		db: db,
	}
}

// FindAllSpecialties - retrieve all veterinarian's specialties
// SELECT * FROM "specialties" WHERE "specialties"."deleted_at" IS NULL
func (repository *VetRepository) FindAllSpecialties() ([]Specialty, error) {
	log.Debug().Msg("get all specialties")

	var specialties []Specialty
	err := repository.db.Find(&specialties).Error
	if err != nil {
		log.Error().Err(err).Msg("Fail to find all specialties.")
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
	log.Debug().Int("id", id).Msg("Search vet by id.")

	var vet Vet
	err := repository.db.First(&vet, id).Error
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Fail to find vet by id.")
		return nil, err
	}

	return &vet, err
}

func (repository *VetRepository) FindByIdWithSpecialties(id int) (*Vet, error) {
	log.Debug().Int("id", id).Msg("Search vet by id.")

	var vet Vet
	err := repository.db.Preload("Specialties").First(&vet, id).Error
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Fail to find vet by id.")
		return nil, err
	}

	return &vet, err
}

// FindByLastName
// SELECT * FROM "vet_specialties" WHERE "vet_specialties"."vet_id" = 5
// SELECT * FROM "specialties" WHERE "specialties"."id" = 1 AND "specialties"."deleted_at" IS NULL
// SELECT * FROM "vets" WHERE last_name = 'Stevens' AND "vets"."deleted_at" IS NULL
func (repository *VetRepository) FindByLastName(lastName string) ([]Vet, error) {
	log.Debug().Str("lastName", lastName).Msg("Search vet by last name.")

	var vets []Vet
	err := repository.db.Preload("Specialties").Where("last_name = ?", lastName).Find(&vets).Error
	if err != nil {
		log.Error().Err(err).Str("lastName", lastName).Msg("Fail to find vet by last name.")
		return nil, err
	}

	return vets, nil
}

/*
FindAll
SELECT * FROM "vets" WHERE "vets"."deleted_at" IS NULL
*/
func (repository *VetRepository) FindAll() ([]Vet, error) {
	log.Debug().Msg("get all vets")

	var vets []Vet
	err := repository.db.Find(&vets).Error
	if err != nil {
		log.Error().Err(err).Msg("Fail to find all vets.")
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
	log.Debug().Msg("get all vets with relations preloaded.")

	var vets []Vet
	err := repository.db.Preload("Specialties").Find(&vets).Error
	if err != nil {
		log.Error().Err(err).Msg("Fail to find all vets.")
		return nil, err
	}
	return vets, err
}

// Insert - insert a new vet
// INSERT INTO "vets" ("first_name","last_name","created_at","updated_at")
//
//	VALUES ('James','Carter','2021-07-25 15:00:00','2021-07-25 15:00:00') RETURN
func (repository *VetRepository) Insert(vet *Vet) (*Vet, error) {
	log.Debug().Any("vet", vet).Msg("Fail a new vet.")

	err := repository.db.Create(&vet).Error
	if err != nil {
		log.Error().Err(err).Msg("Fail to insert new vet.")
		return nil, err
	}
	return vet, err
}

// Update - update vet
// UPDATE "vets" SET "first_name" = 'James', "last_name" = 'Carter' WHERE "id" = 1
func (repository *VetRepository) Update(vet *Vet) (*Vet, error) {
	//log.Debug().Int("id", vet.ID).Msg("Fail vet id.")

	// Omit the column name from update...
	err := repository.db.Omit("created_at").Save(&vet).Error
	if err != nil {
		log.Error().Err(err).Msg("Fail to update vet.")
		return nil, err
	}

	return vet, err
}
