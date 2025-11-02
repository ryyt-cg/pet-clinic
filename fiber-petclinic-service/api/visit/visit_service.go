package visit

import (
	"errors"
	"fiber-petclinic-service/internal/repository"
	"fiber-petclinic-service/internal/repository/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Servicer interface {
	getVisitById(id uint) (*response, error)
	getAllVisits() (*responses, error)
	create(visit *repository.Visit) (*response, error)
	update(visit *repository.Visit) (*response, error)
}

type Service struct {
	repository repository.VisitRepositorier
}

func NewService(repository repository.VisitRepositorier) *Service {
	return &Service{repository: repository}
}

func (service *Service) getVisitById(id uint) (*response, error) {
	visit, err := service.repository.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Uint("id", id).Msg("Visit by id not found.")
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Uint("id", id).Msg("Fail to retrieve visit by id.")
		return nil, err
	}

	response := &response{}
	response.fromVisit(visit)
	return response, nil
}

func (service *Service) getAllVisits() (*responses, error) {
	visits, err := service.repository.FindAll()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &responses{
				Context: model.Context{},
				Visits:  []response{},
			}, nil
		}

		log.Error().Err(err).Msg("Fail to retrieve all visits.")
		return nil, err
	}
	log.Debug().Int("count", len(visits)).Msg("counts of all visits")

	// convert to []response & model.Context
	responses := fromVisitsToResponses(visits)
	return responses, nil
}

func (service *Service) create(visit *repository.Visit) (*response, error) {
	newVisit, err := service.repository.Insert(visit)
	if err != nil {
		log.Error().Err(err).Msg("Fail to create visit.")
		return nil, err
	}

	response := &response{}
	response.fromVisit(newVisit)
	return response, nil
}

func (service *Service) update(visit *repository.Visit) (*response, error) {
	updatedVisit, err := service.repository.Update(visit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Uint("id", visit.ID).Msg("No visit found for update.")
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Msg("Fail to update visit.")
		return nil, err
	}

	response := &response{}
	response.fromVisit(updatedVisit)
	return response, nil
}
