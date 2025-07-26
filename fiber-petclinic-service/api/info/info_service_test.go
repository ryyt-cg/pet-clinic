package info

import (
	"fiber-petclinic-service/config/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getAppInfo(t *testing.T) {
	// instantiate application config object
	app.Config = app.AppConfig{
		AppInfo: app.AppInfoConfig{
			Name:        "Fiber App",
			Description: "App using Go Fiber",
			Version:     "1.5.0",
		},
		Server: app.ServerConfig{
			Host: "localhost",
		},
	}

	appInfo := NewService()
	result, _ := appInfo.getAppInfo()

	assert.Equal(t, "Fiber App", result.AppName)
	assert.Equal(t, "1.5.0", result.Version)
	assert.Equal(t, "App using Go Fiber", result.Description)
}
