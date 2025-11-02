package vet

import (
	"errors"
	"fiber-petclinic-service/internal/repository"
	"fiber-petclinic-service/internal/repository/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Servicer interface {
	getAllSpecialties() (*specialtiesResponse, error)
	getVetById(id uint) (*response, error)
	getVetByIdWithSpecialties(id uint) (*response, error)
	getVetByLastName(lastName string) (*responses, error)
	getAllVets() (*responses, error)
	getAllVetsWithSpecialties() (*responses, error)
	create(vet *repository.Vet) (*response, error)
	update(vet *repository.Vet) (*response, error)
}

type Service struct {
	repository repository.VetRepositorier
}

func NewService(repository repository.VetRepositorier) *Service {
	return &Service{repository: repository}
}

// getAllSpecialties - retrieve all specialties
func (service *Service) getAllSpecialties() (*specialtiesResponse, error) {
	specialties, err := service.repository.FindAllSpecialties()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &specialtiesResponse{
				Context:     model.Context{},
				Specialties: []specialtyResponse{},
			}, nil
		}

		log.Error().Err(err).Msg("Fail to retrieve all specialties.")
		return nil, err
	}

	specialtiesResponses := toSpecialtiesResponses(specialties)
	return specialtiesResponses, nil
}

// getVetById - retrieve vet by id
func (service *Service) getVetById(id uint) (*response, error) {
	log.Debug().Uint("id", id).Msg("GET vet by id")
	vet, err := service.repository.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Uint("id", id).Msg("Vet by id not found.")
			return nil, gorm.ErrRecordNotFound
		}

		log.Error().Err(err).Uint("id", id).Msg("Fail to retrieve vet by id.")
		return nil, err
	}

	response := &response{}
	response.fromVet(vet)
	return response, nil
}

func (service *Service) getVetByIdWithSpecialties(id uint) (*response, error) {
	vet, err := service.repository.FindByIdWithSpecialties(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Uint("id", id).Msg("Vet by id not found with specialties.")
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Uint("id", id).Msg("Fail to retrieve vet by id with specialties")
		return nil, err
	}

	response := &response{}
	response.fromVet(vet)
	return response, nil
}

func (service *Service) getVetByLastName(lastName string) (*responses, error) {
	vets, err := service.repository.FindByLastName(lastName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &responses{
				Context: model.Context{},
				Vets:    []response{},
			}, nil
		}
		log.Error().Err(err).Str("lastName", lastName).Msg("Fail to retrieve the vets by last name.")
		return nil, err
	}

	vetsJson := fromVets(vets)
	contextJson := model.Context{Count: len(vetsJson)}
	return &responses{Vets: vetsJson, Context: contextJson}, nil
}

// getAllVets - retrieve all vets
func (service *Service) getAllVets() (*responses, error) {
	vets, err := service.repository.FindAll()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &responses{
				Context: model.Context{},
				Vets:    []response{},
			}, nil
		}
		log.Error().Err(err).Msg("Fail to retrieve all vets.")
		return nil, err
	}
	log.Debug().Msg("counts of all vets")

	vetsJson := fromVets(vets)
	contextJson := model.Context{Count: len(vetsJson)}
	return &responses{Vets: vetsJson, Context: contextJson}, nil
}

func (service *Service) getAllVetsWithSpecialties() (*responses, error) {
	vets, err := service.repository.FindAllPreload()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &responses{
				Context: model.Context{},
				Vets:    []response{},
			}, nil
		}
		log.Error().Err(err).Msg("Fail to retrieve all vets.")
		return nil, err
	}

	return toResponses(vets), nil
}

// create - create new vet
func (service *Service) create(vet *repository.Vet) (*response, error) {
	log.Info().Any("vet", vet).Msg("Create new vet.")
	newVet, err := service.repository.Insert(vet)

	if err != nil {
		log.Error().Err(err).Msg("Fail new vet failed.")
		return nil, err
	}

	response := &response{}
	response.fromVet(newVet)
	return response, nil
}

// update - update vet
func (service *Service) update(vet *repository.Vet) (*response, error) {
	log.Info().Any("vet", vet).Msg("Update a vet.")
	updatedVet, err := service.repository.Update(vet)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("No vet found for update.")
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Msg("Update vet failed.")
		return nil, err
	}

	response := &response{}
	response.fromVet(updatedVet)
	return response, nil
}
