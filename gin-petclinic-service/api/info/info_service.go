package info

import (
	"github.com/rhtran/gin-petclinic-service/config/app"
	"go.uber.org/zap"
)

type Servicer interface {
	getAppInfo() (*Info, error)
}

type Service struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) *Service {
	return &Service{logger}
}

func (service *Service) getAppInfo() (*Info, error) {
	// this block make GetAppInfo() not testable.
	info := &Info{
		AppName:     app.Config.AppInfo.Name,
		Description: app.Config.AppInfo.Description,
		Version:     app.Config.AppInfo.Version,
	}

	service.logger.Info("app info", zap.Any("info", info))
	return info, nil
}
