package info

import (
	"fiber3-petclinic-service/config/app"

	"github.com/rs/zerolog/log"
)

type Servicer interface {
	getAppInfo() (*Info, error)
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (service *Service) getAppInfo() (*Info, error) {
	// this block make GetAppInfo() not testable.
	info := &Info{
		AppName:     app.Config.AppInfo.Name,
		Description: app.Config.AppInfo.Description,
		Version:     app.Config.AppInfo.Version,
	}

	log.Debug().Any("appInfo", info).Msg("App info")
	return info, nil
}
