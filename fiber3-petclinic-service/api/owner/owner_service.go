package owner

import (
	"errors"
	"fiber3-petclinic-service/internal/repository"
	"fiber3-petclinic-service/internal/repository/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Servicer - owner service interface
type Servicer interface {
	getOwnerById(id uint) (*response, error)
	getAllOwners() (*responses, error)
	getOwnerByIdWithPets(id uint) (*response, error)
	getOwnerByLastName(lastName string) (*responses, error)
	create(ownerRequest *addRequest) (*response, error)
	update(id uint, updateOwner *updateRequest) (*updateResponse, error)
}

type Service struct {
	repository repository.OwnerRepositorier
}

func NewService(repository repository.OwnerRepositorier) *Service {
	return &Service{repository: repository}
}

// getOwnerById - retrieve owner by id
func (service *Service) getOwnerById(id uint) (*response, error) {
	log.Debug().Uint("id", id).Msg("Get owner by id.")
	owner, err := service.repository.FindById(id)
	if err != nil {
		log.Error().Err(err).Uint("id", id).Msg("Fail to retrieve owner by id.")
		return nil, err
	}

	response := toResponse(owner)
	return response, nil
}

// getOwnerByLastName - retrieve owners with pets and visits by last name
func (service *Service) getOwnerByLastName(lastName string) (*responses, error) {
	log.Debug().Str("lastName", lastName).Msg("Get owners by last name.")
	owners, err := service.repository.FindByLastName(lastName)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &responses{
				Owners: []response{},
			}, nil
		}
		log.Error().Err(err).Msg("Fail to retrieve owner by last name.")
		return nil, err
	}

	ownersJson := fromOwners(owners)
	contextJson := model.Context{Count: len(ownersJson)}
	return &responses{Owners: ownersJson, Context: contextJson}, nil
}

// getAllOwners - retrieve all owners
func (service *Service) getAllOwners() (*responses, error) {
	log.Debug().Msg("Get all owner")
	owners, err := service.repository.FindAll()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &responses{Owners: []response{}, Context: model.Context{}}, nil
		}
		log.Error().Err(err).Msg("Fail to retrieve all owner.")
		return nil, err
	}

	// convert to []Response & model.Context
	ownersJson := fromOwners(owners)
	return &responses{Owners: ownersJson, Context: model.Context{Count: len(ownersJson)}}, nil
}

// getOwnerByIdWithPets - retrieve by id with pets
func (service *Service) getOwnerByIdWithPets(id uint) (*response, error) {
	log.Debug().Uint("id", id).Msg("Fail owner with pets by id.")
	owner, err := service.repository.FindByIdWithPets(id)
	if err != nil {
		log.Error().Err(err).Msg("Fail to retrieve all owners.")
		return nil, err
	}

	response := toResponse(owner)
	return response, nil
}

// create - create new owner
func (service *Service) create(ownerRequest *addRequest) (*response, error) {
	log.Debug().Msg("Create new owner")
	owner := toOwnerEntity(ownerRequest)
	newOwner, err := service.repository.Insert(owner)

	if err != nil {
		log.Error().Err(err).Msg("Fail new owner fail.")
		return nil, err
	}

	response := toResponse(newOwner)
	return response, nil
}

// update - update owner
func (service *Service) update(id uint, request *updateRequest) (*updateResponse, error) {
	log.Debug().Msg("Update an owner.")
	ownerEntity := toOwnerEntityFromUpdateRequest(request)
	ownerEntity.ID = id
	updatedOwner, err := service.repository.Update(ownerEntity)

	if err != nil {
		log.Error().Err(err).Msg("Fail to update an owner.")
		return nil, err
	}

	response := toUpdateResponse(updatedOwner)
	return response, nil
}
