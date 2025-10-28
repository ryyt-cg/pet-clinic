package repository

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type VisitRepositorier interface {
	FindById(id uint) (*Visit, error)
	FindAll() ([]Visit, error)
	Insert(visit *Visit) (*Visit, error)
	Update(visit *Visit) (*Visit, error)
}

// VisitRepository searches vet from the database
type VisitRepository struct {
	db *gorm.DB
}

func NewVisitRepository(db *gorm.DB) *VisitRepository {
	return &VisitRepository{
		db: db,
	}
}

/*
FindById - search visit by id
SELECT * FROM "pets" WHERE "pets"."id" = 8 AND "pets"."deleted_at" IS NULL
SELECT * FROM "species" WHERE "species"."id" = 1 AND "species"."deleted_at" IS NULL
SELECT * FROM "visits" WHERE "visits"."id" = 2 AND "visits"."deleted_at" IS NULL ORDER BY "visits"."id" LIMIT 1
*/
func (repository *VisitRepository) FindById(id uint) (*Visit, error) {
	log.Info().Uint("id", id).Msg("Search visit by id.")

	var visit Visit
	err := repository.db.Preload("Pet.Species").First(&visit, id).Error
	if err != nil {
		log.Error().Err(err).Uint("id", id).Msg("Fail to find visit by id.")
		return nil, err
	}

	return &visit, err
}

/*
FindAll - search all visits
SELECT * FROM "pets" WHERE "pets"."id" IN (7,8) AND "pets"."deleted_at" IS NULL
SELECT * FROM "species" WHERE "species"."id" = 1 AND "species"."deleted_at" IS NULL
SELECT * FROM "visits" WHERE "visits"."deleted_at" IS NULL
*/
func (repository *VisitRepository) FindAll() ([]Visit, error) {
	log.Info().Msg("get list of visits")

	var visits []Visit
	// not needed to preload for all visits - will be a performance if get large result.
	err := repository.db.Preload("Pet.Species").Find(&visits).Error
	if err != nil {
		log.Error().Err(err).Msg("Fail to find all visits.")
		return nil, err
	}

	return visits, err
}

// Insert - insert visit
// INSERT INTO "visits" ("created_at","updated_at","deleted_at","pet_id","visit_date","description")
//
//	VALUES ('2021-08-29 15:00:00','2021-08-29 15:00:00',NULL,8,'2021-08-29 15:00:00','test')
func (repository *VisitRepository) Insert(visit *Visit) (*Visit, error) {
	log.Info().Any("visit", visit).Msg("insert a new visit.")

	err := repository.db.Create(visit).Error
	if err != nil {
		log.Error().Err(err).Msg("Fail to insert visit.")
		return nil, err
	}

	return visit, err
}

// Update - update visit
// UPDATE "visits" SET "visit_date" = '2021-08-01 00:00:00', "description" = 'test', "pet_id" = 8 WHERE "id" = 2
func (repository *VisitRepository) Update(visit *Visit) (*Visit, error) {
	log.Info().Any("visit", visit).Msg("Update visits.")

	err := repository.db.Save(visit).Error
	if err != nil {
		log.Error().Err(err).Msg("Fail to update visit,.")
		return nil, err
	}

	return visit, err
}
