package visit

import (
	"github.com/rhtran/gin-petclinic-service/pkg/infra/repository"
	"go.uber.org/zap"
)

type Servicer interface {
	getVisitById(id int) (*Response, error)
	getAllVisits() ([]Response, error)
	create(visit *Request) (*Response, error)
	update(visit *Request) (*Response, error)
}

type Service struct {
	logger     *zap.Logger
	repository repository.VisitRepositorier
}

func NewService(logger *zap.Logger, repository repository.VisitRepositorier) *Service {
	return &Service{logger: logger, repository: repository}
}

func (service *Service) getVisitById(id int) (*Response, error) {
	visit, err := service.repository.FindById(id)
	if err != nil {
		service.logger.Error("Fail to retrieve visit by id.",
			zap.Int("id", id), zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromVisit(visit)
	return response, nil
}

func (service *Service) getAllVisits() ([]Response, error) {
	visits, err := service.repository.FindAll()

	if err != nil {
		service.logger.Error("Fail to retrieve all visits.", zap.String("error", err.Error()))
		return nil, err
	}
	service.logger.Info("counts of all visits", zap.Int("count", len(visits)))
	return FromVisits(visits), nil
}

func (service *Service) create(request *Request) (*Response, error) {
	visitEntity := ToVisit(request)
	newVisit, err := service.repository.Insert(visitEntity)
	if err != nil {
		service.logger.Error("Fail to create visit.", zap.String("error", err.Error()))
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
		service.logger.Error("Fail to update visit.", zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromVisit(updatedVisit)
	return response, nil
}
