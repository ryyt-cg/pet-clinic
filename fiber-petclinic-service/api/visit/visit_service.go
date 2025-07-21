package visit

import (
	"fiber-petclinic-service/pkg/repository"
	"fiber-petclinic-service/pkg/repository/model"
	"github.com/rs/zerolog/log"
)

type Servicer interface {
	getVisitById(id int) (*Response, error)
	getAllVisits() (*Responses, error)
	create(visit *Request) (*Response, error)
	update(visit *Request) (*Response, error)
}

type Service struct {
	repository repository.VisitRepositorier
}

func NewService(repository repository.VisitRepositorier) *Service {
	return &Service{repository: repository}
}

func (service *Service) getVisitById(id int) (*Response, error) {
	visit, err := service.repository.FindById(id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Fail to retrieve visit by id.")
		return nil, err
	}

	response := &Response{}
	response.FromVisit(visit)
	return response, nil
}

func (service *Service) getAllVisits() (*Responses, error) {
	visits, err := service.repository.FindAll()

	if err != nil {
		log.Error().Err(err).Msg("Fail to retrieve all visits.")
		return nil, err
	}
	log.Debug().Int("count", len(visits)).Msg("counts of all visits")

	// convert to []Response & model.Context
	visitsJson := FromVisits(visits)
	contextJson := model.Context{Count: len(visitsJson)}
	return &Responses{Visits: visitsJson, Context: contextJson}, nil
}

func (service *Service) create(request *Request) (*Response, error) {
	visitEntity := ToVisit(request)
	newVisit, err := service.repository.Insert(visitEntity)
	if err != nil {
		log.Error().Err(err).Msg("Fail to create visit.")
		return nil, err
	}

	response := &Response{}
	response.FromVisit(newVisit)
	return response, nil
}

func (service *Service) update(request *Request) (*Response, error) {
	visitEntity := ToVisit(request)
	visitEntity.ID = uint(request.ID)
	updatedVisit, err := service.repository.Update(visitEntity)
	if err != nil {
		log.Error().Err(err).Msg("Fail to update visit.")
		return nil, err
	}

	response := &Response{}
	response.FromVisit(updatedVisit)
	return response, nil
}
