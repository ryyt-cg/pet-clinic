package info

import (
	"fiber-petclinic-service/config/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_lookupIP(t *testing.T) {
	// instantiate application config object
	app.Config = app.AppConfig{
		Server: app.ServerConfig{
			Host: "localhost",
		},
	}

	expected := []string{"::1", "127.0.0.1"}

	ipService := NewIPService()
	result, _ := ipService.lookupIP(app.Config.Server.Host)

	for i, ip := range result {
		assert.Equal(t, expected[i], ip.String())
	}
}
