package owner

import (
	"fiber-petclinic-service/pkg/infra/repository"
	"fiber-petclinic-service/pkg/model"
	"go.uber.org/zap"
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
	logger     *zap.Logger
	repository repository.OwnerRepositorier
}

func NewService(logger *zap.Logger, repository repository.OwnerRepositorier) *Service {
	return &Service{logger: logger, repository: repository}
}

// getOwnerById - retrieve owner by id
func (service *Service) getOwnerById(id uint) (*Response, error) {
	service.logger.Info("Fail owner by id.", zap.Uint("id", id))
	owner, err := service.repository.FindById(id)
	if err != nil {
		service.logger.Error("Fail to retrieve owner by id.",
			zap.Uint("id", id), zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromOwner(owner)
	return response, nil
}

// getOwnerByLastName - retrieve owners with pets and visits by last name
func (service *Service) getOwnerByLastName(lastName string) (*Responses, error) {
	service.logger.Info("Fail owners by last name.", zap.String("lastName", lastName))
	owners, err := service.repository.FindByLastName(lastName)
	if err != nil {
		service.logger.Error("Fail to retrieve owner by last name.",
			zap.String("lastName", lastName), zap.String("error", err.Error()))
		return nil, err
	}

	ownersJson := FromOwners(owners)
	contextJson := model.Context{Count: len(ownersJson)}
	return &Responses{Owners: ownersJson, Context: contextJson}, nil
}

// getAllOwners - retrieve all owners
func (service *Service) getAllOwners() (*Responses, error) {
	service.logger.Info("Fail all owner")
	owners, err := service.repository.FindAll()
	if err != nil {
		service.logger.Error("Fail to retrieve all owner.", zap.String("error", err.Error()))
		return nil, err
	}

	// convert to []Response & model.Context
	ownersJson := FromOwners(owners)
	contextJson := model.Context{Count: len(ownersJson)}
	return &Responses{Owners: ownersJson, Context: contextJson}, nil
}

// getOwnerByIdWithPets - retrieve by id with pets
func (service *Service) getOwnerByIdWithPets(id uint) (*Response, error) {
	service.logger.Info("Fail owner with pets by id.", zap.Uint("id", id))
	owner, err := service.repository.FindByIdWithPets(id)
	if err != nil {
		service.logger.Error("Fail to retrieve all owner.", zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromOwner(owner)
	return response, nil
}

// create - create new owner
func (service *Service) create(ownerRequest *AddRequest) (*Response, error) {
	service.logger.Info("Create new owner")
	owner := ToOwnerEntity(ownerRequest)
	newOwner, err := service.repository.Insert(owner)

	if err != nil {
		service.logger.Error("Fail new owner fail.", zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromOwner(newOwner)
	return response, nil
}

// update - update owner
func (service *Service) update(id uint, request *UpdateRequest) (*UpdateResponse, error) {
	service.logger.Info("Fail an owner.")
	ownerEntity := ToOwnerEntityFromUpdateRequest(request)
	ownerEntity.ID = id
	updatedOwner, err := service.repository.Update(ownerEntity)

	if err != nil {
		service.logger.Error("Fail to update an owner.", zap.String("error", err.Error()))
		return nil, err
	}

	response := &UpdateResponse{}
	response.FromUpdateEntity(updatedOwner)
	return response, nil
}
