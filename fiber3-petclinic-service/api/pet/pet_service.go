package pet

import (
	"errors"
	"fiber3-petclinic-service/internal/repository"
	"fiber3-petclinic-service/internal/repository/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Servicer interface {
	getAllPets() (*Responses, error)
	getPetById(id uint) (*Response, error)
	getPetWithVisitsById(id uint) (*Response, error)
	getPetsByName(name string) (*Responses, error)
	create(pet *repository.Pet) (*Response, error)
	update(pet *repository.Pet) (*Response, error)
}

type Service struct {
	repository repository.PetRepositorier
}

func NewService(repository repository.PetRepositorier) *Service {
	return &Service{repository: repository}
}

// getAllPets - retrieve all pets
func (service *Service) getAllPets() (*Responses, error) {
	log.Debug().Msg("Retrieve all pets")
	pets, err := service.repository.FindAll()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Msg("No pets found.")
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Msg("Fail to retrieve all pets.")
		return nil, err
	}

	petResponses := FromPets(pets)
	contextJson := model.Context{Count: len(petResponses)}
	return &Responses{Pets: petResponses, Context: contextJson}, nil
}

// getPetById - retrieve pet by id
func (service *Service) getPetById(id uint) (*Response, error) {
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

	response := ToResponse(petF)
	return response, nil
}

func (service *Service) getPetWithVisitsById(id uint) (*Response, error) {
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

	response := ToResponse(petF)
	return response, nil
}

// getPetByName - retrieve pet by name
func (service *Service) getPetsByName(name string) (*Responses, error) {
	log.Debug().Str("name", name).Msg("retrieve pets by name.")

	pets, err := service.repository.FindByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Str("name", name).Msg("No pets found with this name.")
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Str("name", name).Msg("fail to retrieve pets by name.")
		return nil, err
	}

	return ToResponses(pets), nil
}

// create - create new pet
func (service *Service) create(pet *repository.Pet) (*Response, error) {
	log.Debug().Str("name", pet.Name).Msg("Create new pet.")
	newPet, err := service.repository.Insert(pet)

	if err != nil {
		log.Error().Err(err).Msg("Fail new pet failed.")
		return nil, err
	}

	response := ToResponse(newPet)
	return response, nil
}

// update - update pet
func (service *Service) update(pet *repository.Pet) (*Response, error) {
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

	response := ToResponse(updatedPet)
	return response, nil
}
