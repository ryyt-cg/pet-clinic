package pet

import (
	"errors"
	"gin-petclinic-service/internal/repository"
	"gin-petclinic-service/internal/repository/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Servicer interface {
	getAllPets() (*responses, error)
	getPetById(id uint) (*response, error)
	getPetWithVisitsById(id uint) (*response, error)
	getPetsByName(name string) (*responses, error)
	create(pet *repository.Pet) (*addResponse, error)
	update(pet *repository.Pet) (*updateResponse, error)
}

type Service struct {
	repository repository.PetRepositorier
}

func NewService(repository repository.PetRepositorier) *Service {
	return &Service{repository: repository}
}

// getAllPets - retrieve all pets
func (service *Service) getAllPets() (*responses, error) {
	log.Debug().Msg("Retrieve all pets")
	pets, err := service.repository.FindAll()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &responses{
				Context: model.Context{},
				Pets:    []response{},
			}, nil
		}
		log.Error().Err(err).Msg("Fail to retrieve all pets.")
		return nil, err
	}

	petResponses := fromPets(pets)
	contextJson := model.Context{Count: len(petResponses)}
	return &responses{Pets: petResponses, Context: contextJson}, nil
}

// getPetById - retrieve pet by id
func (service *Service) getPetById(id uint) (*response, error) {
	log.Debug().Uint("id", id).Msg("Retrieve pet by id.")
	petF, err := service.repository.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Uint("id", id).Msg("Pet by id not found.")
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Uint("id", id).Msg("Fail to retrieve pet by id.")
		return nil, err
	}

	response := toResponse(petF)
	return response, nil
}

func (service *Service) getPetWithVisitsById(id uint) (*response, error) {
	log.Debug().Uint("id", id).Msg("Retrieve pet by id.")
	petF, err := service.repository.FindByIdWithVisits(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Uint("id", id).Msg("Pet by id not found.")
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Uint("id", id).Msg("Fail to retrieve pet by id.")
		return nil, err
	}

	response := toResponse(petF)
	return response, nil
}

// getPetByName - retrieve pet by name
func (service *Service) getPetsByName(name string) (*responses, error) {
	log.Debug().Str("name", name).Msg("retrieve pets by name.")

	pets, err := service.repository.FindByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &responses{
				Context: model.Context{},
				Pets:    []response{},
			}, nil
		}
		log.Error().Err(err).Str("name", name).Msg("fail to retrieve pets by name.")
		return nil, err
	}

	return toResponses(pets), nil
}

// create - create new pet
func (service *Service) create(pet *repository.Pet) (*addResponse, error) {
	log.Debug().Str("name", pet.Name).Msg("Create new pet.")
	newPet, err := service.repository.Insert(pet)

	if err != nil {
		log.Error().Err(err).Msg("Fail new pet failed.")
		return nil, err
	}

	response := toAddResponse(newPet)
	return response, nil
}

// update - update pet
func (service *Service) update(pet *repository.Pet) (*updateResponse, error) {
	log.Info().Any("pet", pet).Msg("Update a vet.")
	updatedPet, err := service.repository.Update(pet)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("No pet found for update.")
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Msg("Fail to Update pet.")
		return nil, err
	}

	response := toUpdateResponse(updatedPet)
	return response, nil
}
