package vet

import (
	"github.com/rhtran/gin-petclinic-service/pkg/infra/repository"
	"go.uber.org/zap"
)

type Servicer interface {
	getAllSpecialties() (*specialtiesResponse, error)
	getVetById(id int) (*Response, error)
	getVetByIdWithSpecialties(id int) (*Response, error)
	getVetByLastName(lastName string) ([]Response, error)
	getAllVets() ([]Response, error)
	getAllVetsWithSpecialties() ([]Response, error)
	create(vet *repository.Vet) (*Response, error)
	update(vet *repository.Vet) (*Response, error)
}

type Service struct {
	logger     *zap.Logger
	repository repository.VetRepositorier
}

func NewService(logger *zap.Logger, repository repository.VetRepositorier) *Service {
	return &Service{logger: logger, repository: repository}
}

// getAllSpecialties - retrieve all specialties
func (service *Service) getAllSpecialties() (*specialtiesResponse, error) {
	specialties, err := service.repository.FindAllSpecialties()
	if err != nil {
		service.logger.Error("fail to retrieve all specialties.", zap.String("error", err.Error()))
		return nil, err
	}

	specialtyResponses := ToSpecialtyResponses(specialties)

	return &specialtiesResponse{
		Specialties: *specialtyResponses,
	}, nil
}

// getVetById - retrieve vet by id
func (service *Service) getVetById(id int) (*Response, error) {
	vet, err := service.repository.FindById(id)
	if err != nil {
		service.logger.Error("fail to retrieve vet by id.",
			zap.Int("id", id), zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromVet(vet)
	return response, nil
}

func (service *Service) getVetByIdWithSpecialties(id int) (*Response, error) {
	vet, err := service.repository.FindByIdWithSpecialties(id)
	if err != nil {
		service.logger.Error("fail to retrieve vet by id with specialties",
			zap.Int("id", id), zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromVet(vet)
	return response, nil
}

func (service *Service) getVetByLastName(lastName string) ([]Response, error) {
	vets, err := service.repository.FindByLastName(lastName)
	if err != nil {
		service.logger.Error("fail to retrieve the vets by last name.",
			zap.String("lastName", lastName), zap.String("error", err.Error()))
		return nil, err
	}

	return FromVets(vets), nil
}

// getAllVets - retrieve all vets
func (service *Service) getAllVets() ([]Response, error) {
	vets, err := service.repository.FindAll()

	if err != nil {
		service.logger.Error("fail to retrieve all vets.", zap.String("error", err.Error()))
		return nil, err
	}
	service.logger.Info("counts of all vets")
	return FromVets(vets), nil
}

func (service *Service) getAllVetsWithSpecialties() ([]Response, error) {
	vets, err := service.repository.FindAllPreload()

	if err != nil {
		service.logger.Error("fail to retrieve all vets.", zap.String("error", err.Error()))
		return nil, err
	}
	service.logger.Info("counts of all vets.", zap.Int("count", len(vets)))
	return FromVets(vets), nil
}

// create - create new vet
func (service *Service) create(vet *repository.Vet) (*Response, error) {
	service.logger.Info("Create new vet.", zap.Any("vet", vet))
	newVet, err := service.repository.Insert(vet)

	if err != nil {
		service.logger.Error("insert new vet failed.", zap.String("error", err.Error()))
		return &Response{}, err
	}

	response := &Response{}
	response.FromVet(newVet)
	return response, nil
}

// update - update vet
func (service *Service) update(vet *repository.Vet) (*Response, error) {
	service.logger.Info("update vet.", zap.Any("vet", vet))
	updatedVet, err := service.repository.Update(vet)

	if err != nil {
		service.logger.Error("Update vet failed.", zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromVet(updatedVet)
	return response, nil
}
