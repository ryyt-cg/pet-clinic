package pet

import (
	"github.com/rhtran/gin-petclinic-service/pkg/infra/repository"
	"github.com/rhtran/gin-petclinic-service/pkg/model"
	"go.uber.org/zap"
)

type Servicer interface {
	getAllPets() (*Responses, error)
	getPetById(id int) (*Response, error)
	getPetWithVisitsById(id int) (*Response, error)
	getPetsByName(name string) ([]Response, error)
	create(pet *repository.Pet) (*Response, error)
	update(pet *repository.Pet) (*Response, error)
}

type Service struct {
	logger     *zap.Logger
	repository repository.PetRepositorier
}

func NewService(logger *zap.Logger, repository repository.PetRepositorier) *Service {
	return &Service{logger: logger, repository: repository}
}

// getAllPets - retrieve all pets
func (service *Service) getAllPets() (*Responses, error) {
	service.logger.Info("Fail all pets")
	pets, err := service.repository.FindAll()
	if err != nil {
		service.logger.Error("Fail to retrieve all pets.", zap.String("error", err.Error()))
		return nil, err
	}

	petResponses := FromPets(pets)

	contextJson := model.Context{Count: len(petResponses)}
	return &Responses{Pets: petResponses, Context: contextJson}, nil
}

// getPetById - retrieve pet by id
func (service *Service) getPetById(id int) (*Response, error) {
	service.logger.Info("Fail pet by id.", zap.Int("id", id))
	petF, err := service.repository.FindById(id)
	if err != nil {
		service.logger.Error("Fail to retrieve pet by id.",
			zap.Int("id", id), zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromPet(petF)
	return response, nil
}

func (service *Service) getPetWithVisitsById(id int) (*Response, error) {
	service.logger.Info("Fail pet by id.", zap.Int("id", id))
	petF, err := service.repository.FindByIdWithVisits(id)
	if err != nil {
		service.logger.Error("Fail to retrieve pet by id.",
			zap.Int("id", id), zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromPet(petF)
	return response, nil
}

// getPetByName - retrieve pet by name
func (service *Service) getPetsByName(name string) ([]Response, error) {
	service.logger.Info("Fail pet by name.", zap.String("name", name))

	pets, err := service.repository.FindByName(name)
	if err != nil {
		service.logger.Error("Fail to retrieve pet by name.",
			zap.String("name", name), zap.String("error", err.Error()))
		return nil, err
	}

	return FromPets(pets), nil
}

// create - create new pet
func (service *Service) create(pet *repository.Pet) (*Response, error) {
	service.logger.Info("Create new pet.", zap.String("name", pet.Name))
	newPet, err := service.repository.Insert(pet)

	if err != nil {
		service.logger.Error("Fail new pet failed.", zap.String("error", err.Error()))
		return &Response{}, err
	}

	response := &Response{}
	response.FromPet(newPet)
	return response, nil
}

// update - update pet
func (service *Service) update(pet *repository.Pet) (*Response, error) {
	service.logger.Info("Fail vet.", zap.Any("pet", pet))
	updatedPet, err := service.repository.Update(pet)

	if err != nil {
		service.logger.Error("Update pet failed.", zap.String("error", err.Error()))
		return nil, err
	}

	response := &Response{}
	response.FromPet(updatedPet)
	return response, nil
}
