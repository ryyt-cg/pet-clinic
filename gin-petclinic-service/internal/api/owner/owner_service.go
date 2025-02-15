package owner

import (
	"github.com/rhtran/gin-petclinic-service/pkg/infra/repository"
	"github.com/rhtran/gin-petclinic-service/pkg/model"
	"go.uber.org/zap"
)

type Servicer interface {
	getOwnerById(id uint) (*Response, error)
	getAllOwners() (*Responses, error)
	getOwnerByIdWithPets(id uint) (*Responses, error)
	getOwnerByLastName(lastName string) (*Responses, error)
	create(owner *repository.Owner) (*Response, error)
	update(owner *repository.Owner) (*Response, error)
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
	service.logger.Info("retrieve owner by id.", zap.Uint("id", id))
	owner, err := service.repository.FindById(id)
	if err != nil {
		service.logger.Error("fails to retrieve owner by id.",
			zap.Uint("id", id), zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromOwner(owner)
	return response, nil
}

// getOwnerByLastName - retrieve owners with pets and visits by last name
func (service *Service) getOwnerByLastName(lastName string) (*Responses, error) {
	service.logger.Info("retrieve owners by last name.", zap.String("lastName", lastName))
	owners, err := service.repository.FindByLastName(lastName)
	if err != nil {
		service.logger.Error("fail to retrieve owner by last name.",
			zap.String("lastName", lastName), zap.String("error", err.Error()))
		return nil, err
	}

	ownersJson := FromOwners(owners)
	contextJson := model.Context{Count: len(ownersJson)}
	return &Responses{Owners: ownersJson, Context: contextJson}, nil
}

// getAllOwners - retrieve all owners
func (service *Service) getAllOwners() (*Responses, error) {
	service.logger.Info("retrieve all owner")
	owners, err := service.repository.FindAll()
	if err != nil {
		service.logger.Error("fail to retrieve all owner.", zap.String("error", err.Error()))
		return nil, err
	}

	ownersJson := FromOwners(owners)
	contextJson := model.Context{Count: len(ownersJson)}
	return &Responses{Owners: ownersJson, Context: contextJson}, nil
}

// getAllOwnersWithPets - retrieve all owners with pets
func (service *Service) getOwnerByIdWithPets(id uint) (*Responses, error) {
	service.logger.Info("retrieve all owner")
	owners, err := service.repository.FindByIdWithPets(id)
	if err != nil {
		service.logger.Error("fail to retrieve all owner.", zap.String("error", err.Error()))
		return nil, err
	}

	ownersJson := FromOwners(owners)
	contextJson := model.Context{Count: len(ownersJson)}
	return &Responses{Owners: ownersJson, Context: contextJson}, nil
}

// create - create new owner
func (service *Service) create(owner *repository.Owner) (*Response, error) {
	service.logger.Info("Create new owner")
	newOwner, err := service.repository.Insert(owner)

	if err != nil {
		service.logger.Error("insert new owner failed.", zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromOwner(newOwner)
	return response, nil
}

// update - update owner
func (service *Service) update(owner *repository.Owner) (*Response, error) {
	service.logger.Info("update an owner.")
	updatedOwner, err := service.repository.Update(owner)

	if err != nil {
		service.logger.Error("fail to update an owner.", zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromOwner(updatedOwner)
	return response, nil
}
