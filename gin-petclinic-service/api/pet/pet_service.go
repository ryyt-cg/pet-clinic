package pet

import (
	"gin-petclinic-service/pkg/model"
	"gin-petclinic-service/pkg/repository"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("Fail to retrieve all pets.")
		return nil, err
	}

	petResponses := FromPets(pets)

	contextJson := model.Context{Count: len(petResponses)}
	return &Responses{Pets: petResponses, Context: contextJson}, nil
}

// getPetById - retrieve pet by id
func (service *Service) getPetById(id int) (*Response, error) {
	log.Debug().Int("id", id).Msg("Fail pet by id.")
	petF, err := service.repository.FindById(id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Fail to retrieve pet by id.")
		return nil, err
	}

	response := &Response{}
	response.FromPet(petF)
	return response, nil
}

func (service *Service) getPetWithVisitsById(id int) (*Response, error) {
	log.Debug().Int("id", id).Msg("Retrieve pet by id.")
	petF, err := service.repository.FindByIdWithVisits(id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Fail to retrieve pet by id.")
		return nil, err
	}

	response := &Response{}
	response.FromPet(petF)
	return response, nil
}

// getPetByName - retrieve pet by name
func (service *Service) getPetsByName(name string) ([]Response, error) {
	log.Debug().Str("name", name).Msg("Retrieve pet by name.")

	pets, err := service.repository.FindByName(name)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Fail to retrieve pet by name.")
		return nil, err
	}

	return FromPets(pets), nil
}

// create - create new pet
func (service *Service) create(pet *repository.Pet) (*Response, error) {
	log.Debug().Str("name", pet.Name).Msg("Create new pet.")
	newPet, err := service.repository.Insert(pet)

	if err != nil {
		log.Error().Err(err).Msg("Create a new pet failed.")
		return &Response{}, err
	}

	response := &Response{}
	response.FromPet(newPet)
	return response, nil
}

// update - update pet
func (service *Service) update(pet *repository.Pet) (*Response, error) {
	log.Debug().Any("pet", pet).Msg("Update a vet.")
	updatedPet, err := service.repository.Update(pet)

	if err != nil {
		log.Error().Err(err).Msg("Update pet failed.")
		return nil, err
	}

	response := &Response{}
	response.FromPet(updatedPet)
	return response, nil
}
