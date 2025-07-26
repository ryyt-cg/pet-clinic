package vet

import (
	repository2 "gin-petclinic-service/pkg/repository"
	"github.com/rs/zerolog/log"
)

type Servicer interface {
	getAllSpecialties() (*specialtiesResponse, error)
	getVetById(id int) (*Response, error)
	getVetByIdWithSpecialties(id int) (*Response, error)
	getVetByLastName(lastName string) ([]Response, error)
	getAllVets() ([]Response, error)
	getAllVetsWithSpecialties() ([]Response, error)
	create(vet *repository2.Vet) (*Response, error)
	update(vet *repository2.Vet) (*Response, error)
}

type Service struct {
	repository repository2.VetRepositorier
}

func NewService(repository repository2.VetRepositorier) *Service {
	return &Service{repository: repository}
}

// getAllSpecialties - retrieve all specialties
func (service *Service) getAllSpecialties() (*specialtiesResponse, error) {
	specialties, err := service.repository.FindAllSpecialties()
	if err != nil {
		log.Error().Err(err).Msg("Fail to retrieve all specialties.")
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
		log.Error().Err(err).Int("id", id).Msg("Fail to retrieve vet by id.")
		return nil, err
	}

	response := &Response{}
	response.FromVet(vet)
	return response, nil
}

func (service *Service) getVetByIdWithSpecialties(id int) (*Response, error) {
	vet, err := service.repository.FindByIdWithSpecialties(id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Fail to retrieve vet by id with specialties")
		return nil, err
	}

	response := &Response{}
	response.FromVet(vet)
	return response, nil
}

func (service *Service) getVetByLastName(lastName string) ([]Response, error) {
	vets, err := service.repository.FindByLastName(lastName)
	if err != nil {
		log.Error().Err(err).Str("lastName", lastName).Msg("Fail to retrieve the vets by last name.")
		return nil, err
	}

	return FromVets(vets), nil
}

// getAllVets - retrieve all vets
func (service *Service) getAllVets() ([]Response, error) {
	vets, err := service.repository.FindAll()

	if err != nil {
		log.Error().Err(err).Msg("Fail to retrieve all vets.")
		return nil, err
	}
	log.Debug().Msg("counts of all vets")
	return FromVets(vets), nil
}

func (service *Service) getAllVetsWithSpecialties() ([]Response, error) {
	vets, err := service.repository.FindAllPreload()

	if err != nil {
		log.Error().Err(err).Msg("Fail to retrieve all vets.")
		return nil, err
	}
	log.Debug().Int("count", len(vets)).Msg("counts of all vets.")
	return FromVets(vets), nil
}

// create - create new vet
func (service *Service) create(vet *repository2.Vet) (*Response, error) {
	log.Debug().Any("vet", vet).Msg("Create new vet.")
	newVet, err := service.repository.Insert(vet)

	if err != nil {
		log.Error().Err(err).Msg("Fail new vet failed.")
		return &Response{}, err
	}

	response := &Response{}
	response.FromVet(newVet)
	return response, nil
}

// update - update vet
func (service *Service) update(vet *repository2.Vet) (*Response, error) {
	log.Debug().Any("vet", vet).Msg("Fail vet.")
	updatedVet, err := service.repository.Update(vet)

	if err != nil {
		log.Error().Err(err).Msg("Update vet failed.")
		return nil, err
	}

	response := &Response{}
	response.FromVet(updatedVet)
	return response, nil
}
