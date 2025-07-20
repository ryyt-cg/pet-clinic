package author

import (
	"errors"
	"fiber-petclinic-service/exception"
	"github.com/rs/zerolog/log"
)

type Servicer interface {
	getAuthorByID(int) (*Author, error)
	getAuthorByName(string) (*Author, error)
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (service *Service) getAuthorByID(ID int) (*Author, error) {
	author, err := getAuthorByID(ID)
	if err != nil {
		if errors.Is(err, exception.ErrorNotFound) {
			log.Warn().Int("ID", ID).Msg("Author not found by ID")
			return nil, exception.ErrorNotFound // Return a more specific error
		}
		log.Error().Err(err).Int("ID", ID).Msg("Failed to fetch author by ID")
		return nil, err // Propagate other errors
	}

	log.Debug().Any("author", author).Msg("getAuthorByID")
	return author, nil
}

func (service *Service) getAuthorByName(name string) (*Author, error) {
	author, err := getAuthorByName(name)
	if err != nil {
		if errors.Is(err, exception.ErrorNotFound) {
			log.Warn().Str("name", name).Msg("Author not found by name")
			return nil, exception.ErrorNotFound // Return a more specific error
		}
		log.Error().Err(err).Str("name", name).Msg("Failed to fetch author by name")
		return nil, err // Propagate other errors
	}

	log.Debug().Any("author", author).Msg("getAuthorByName")
	return author, nil
}
