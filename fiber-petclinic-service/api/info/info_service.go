package info

import (
	"fiber-petclinic-service/config/app"

	"github.com/rs/zerolog/log"
)

type Servicer interface {
	getAppInfo() (*response, error)
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (service *Service) getAppInfo() (*response, error) {
	// this block make GetAppInfo() not testable.
	info := &response{
		AppName:     app.Config.AppInfo.Name,
		Description: app.Config.AppInfo.Description,
		Version:     app.Config.AppInfo.Version,
	}

	log.Debug().Any("appInfo", info).Msg("App info")
	return info, nil
}
