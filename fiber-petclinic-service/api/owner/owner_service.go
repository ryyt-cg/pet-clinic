package owner

import (
	"errors"
	"fiber-petclinic-service/internal/repository"
	"fiber-petclinic-service/internal/repository/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Servicer - owner service interface
type Servicer interface {
	getOwnerById(id uint) (*Response, error)
	getAllOwners() (*Responses, error)
	getOwnerByIdWithPets(id uint) (*Response, error)
	getOwnerByLastName(lastName string) (*Responses, error)
	create(ownerRequest *AddRequest) (*Response, error)
	update(id uint, updateOwner *UpdateRequest) (*UpdateResponse, error)
}

type Service struct {
	repository repository.OwnerRepositorier
}

func NewService(repository repository.OwnerRepositorier) *Service {
	return &Service{repository: repository}
}

// getOwnerById - retrieve owner by id
func (service *Service) getOwnerById(id uint) (*Response, error) {
	log.Debug().Uint("id", id).Msg("Get owner by id.")
	owner, err := service.repository.FindById(id)
	if err != nil {
		log.Error().Err(err).Uint("id", id).Msg("Fail to retrieve owner by id.")
		return nil, err
	}

	response := ToResponse(owner)
	return response, nil
}

// getOwnerByLastName - retrieve owners with pets and visits by last name
func (service *Service) getOwnerByLastName(lastName string) (*Responses, error) {
	log.Debug().Str("lastName", lastName).Msg("Get owners by last name.")
	owners, err := service.repository.FindByLastName(lastName)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &Responses{
				Owners: []Response{},
			}, nil
		}
		log.Error().Err(err).Msg("Fail to retrieve owner by last name.")
		return nil, err
	}

	ownersJson := FromOwners(owners)
	contextJson := model.Context{Count: len(ownersJson)}
	return &Responses{Owners: ownersJson, Context: contextJson}, nil
}

// getAllOwners - retrieve all owners
func (service *Service) getAllOwners() (*Responses, error) {
	log.Debug().Msg("Get all owner")
	owners, err := service.repository.FindAll()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &Responses{Owners: []Response{}, Context: model.Context{}}, nil
		}
		log.Error().Err(err).Msg("Fail to retrieve all owner.")
		return nil, err
	}

	// convert to []Response & model.Context
	ownersJson := FromOwners(owners)
	contextJson := model.Context{Count: len(ownersJson)}
	return &Responses{Owners: ownersJson, Context: contextJson}, nil
}

// getOwnerByIdWithPets - retrieve by id with pets
func (service *Service) getOwnerByIdWithPets(id uint) (*Response, error) {
	log.Debug().Uint("id", id).Msg("Fail owner with pets by id.")
	owner, err := service.repository.FindByIdWithPets(id)
	if err != nil {
		log.Error().Err(err).Msg("Fail to retrieve all owners.")
		return nil, err
	}

	response := ToResponse(owner)
	return response, nil
}

// create - create new owner
func (service *Service) create(ownerRequest *AddRequest) (*Response, error) {
	log.Debug().Msg("Create new owner")
	owner := ToOwnerEntity(ownerRequest)
	newOwner, err := service.repository.Insert(owner)

	if err != nil {
		log.Error().Err(err).Msg("Fail new owner fail.")
		return nil, err
	}

	response := ToResponse(newOwner)
	return response, nil
}

// update - update owner
func (service *Service) update(id uint, request *UpdateRequest) (*UpdateResponse, error) {
	log.Debug().Msg("Update an owner.")
	ownerEntity := ToOwnerEntityFromUpdateRequest(request)
	ownerEntity.ID = id
	updatedOwner, err := service.repository.Update(ownerEntity)

	if err != nil {
		log.Error().Err(err).Msg("Fail to update an owner.")
		return nil, err
	}

	response := ToUpdateResponse(updatedOwner)
	return response, nil
}
