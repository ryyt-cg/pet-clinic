package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type VisitRepositorier interface {
	FindById(id int) (*Visit, error)
	FindAll() ([]Visit, error)
	Insert(visit *Visit) (*Visit, error)
	Update(visit *Visit) (*Visit, error)
}

// VisitRepository searches vet from the database
type VisitRepository struct {
	pg *gorm.DB
}

func NewVisitRepository(pg *gorm.DB) *VisitRepository {
	return &VisitRepository{
		pg: pg,
	}
}

/*
FindById - search visit by id
SELECT * FROM "pets" WHERE "pets"."id" = 8 AND "pets"."deleted_at" IS NULL
SELECT * FROM "types" WHERE "types"."id" = 1 AND "types"."deleted_at" IS NULL
SELECT * FROM "visits" WHERE "visits"."id" = 2 AND "visits"."deleted_at" IS NULL ORDER BY "visits"."id" LIMIT 1
*/
func (repository *VisitRepository) FindById(id int) (*Visit, error) {
	log.Info("Search visit by id.", zap.Int("id", id))

	var visit Visit
	err := repository.pg.Preload("Pet.Type").First(&visit, id).Error
	if err != nil {
		log.Error("Fail to find visit by id.",
			zap.Int("id", id), zap.String("error", err.Error()))
		return nil, err
	}

	return &visit, err
}

/*
FindAll - search all visits
SELECT * FROM "pets" WHERE "pets"."id" IN (7,8) AND "pets"."deleted_at" IS NULL
SELECT * FROM "types" WHERE "types"."id" = 1 AND "types"."deleted_at" IS NULL
SELECT * FROM "visits" WHERE "visits"."deleted_at" IS NULL
*/
func (repository *VisitRepository) FindAll() ([]Visit, error) {
	log.Info("get list of visits")

	var visits []Visit
	// not needed to preload for all visits - will be a performance if get large result.
	err := repository.pg.Preload("Pet.Type").Find(&visits).Error
	if err != nil {
		log.Error("Fail to find all visits.", zap.String("error", err.Error()))
		return nil, err
	}

	return visits, err
}

// Insert - insert visit
// INSERT INTO "visits" ("created_at","updated_at","deleted_at","pet_id","visit_date","description")
//
//	VALUES ('2021-08-29 15:00:00','2021-08-29 15:00:00',NULL,8,'2021-08-29 15:00:00','test')
func (repository *VisitRepository) Insert(visit *Visit) (*Visit, error) {
	log.Info("Fail visit.", zap.Any("visit", visit))

	err := repository.pg.Create(visit).Error
	if err != nil {
		log.Error("Fail to insert visit.", zap.String("error", err.Error()))
		return nil, err
	}

	return visit, err
}

// Update - update visit
// UPDATE "visits" SET "visit_date" = '2021-08-01 00:00:00', "description" = 'test', "pet_id" = 8 WHERE "id" = 2
func (repository *VisitRepository) Update(visit *Visit) (*Visit, error) {
	log.Info("Update visits.", zap.Any("visit", visit))

	err := repository.pg.Save(visit).Error
	if err != nil {
		log.Error("Fail to update visit,.", zap.String("error", err.Error()))
		return nil, err
	}

	return visit, err
}
