package visit

import (
	"errors"
	"fiber3-petclinic-service/internal/repository"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Servicer interface {
	getVisitById(id uint) (*Response, error)
	getAllVisits() (*Responses, error)
	create(visit *repository.Visit) (*Response, error)
	update(visit *repository.Visit) (*Response, error)
}

type Service struct {
	repository repository.VisitRepositorier
}

func NewService(repository repository.VisitRepositorier) *Service {
	return &Service{repository: repository}
}

func (service *Service) getVisitById(id uint) (*Response, error) {
	visit, err := service.repository.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Uint("id", id).Msg("Visit by id not found.")
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Uint("id", id).Msg("Fail to retrieve visit by id.")
		return nil, err
	}

	response := &Response{}
	response.FromVisit(visit)
	return response, nil
}

func (service *Service) getAllVisits() (*Responses, error) {
	visits, err := service.repository.FindAll()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Msg("No visits found.")
			return nil, gorm.ErrRecordNotFound
		}

		log.Error().Err(err).Msg("Fail to retrieve all visits.")
		return nil, err
	}
	log.Debug().Int("count", len(visits)).Msg("counts of all visits")

	// convert to []Response & model.Context
	responses := FromVisitsToResponses(visits)
	return responses, nil
}

func (service *Service) create(visit *repository.Visit) (*Response, error) {
	newVisit, err := service.repository.Insert(visit)
	if err != nil {
		log.Error().Err(err).Msg("Fail to create visit.")
		return nil, err
	}

	response := &Response{}
	response.FromVisit(newVisit)
	return response, nil
}

func (service *Service) update(visit *repository.Visit) (*Response, error) {
	updatedVisit, err := service.repository.Update(visit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Uint("id", visit.ID).Msg("No visit found for update.")
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Msg("Fail to update visit.")
		return nil, err
	}

	response := &Response{}
	response.FromVisit(updatedVisit)
	return response, nil
}
