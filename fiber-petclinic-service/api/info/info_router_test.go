package info

import (
	"encoding/json"
	"fiber-petclinic-service/config/app"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

type infoServiceMock struct {
	mock.Mock
}

type ipServiceMock struct {
	mock.Mock
}

func (infoM *infoServiceMock) getAppInfo() (*Info, error) {
	args := infoM.Called()
	intf := args.Get(0)

	if intf == nil {
		return nil, args.Error(1)
	}

	val := intf.(Info)
	return &val, args.Error(1)
}

func (ipM *ipServiceMock) lookupIP(host string) ([]net.IP, error) {
	args := ipM.Called(host)
	intf := args.Get(0)

	if intf == nil {
		return nil, args.Error(1)
	}

	val := intf.([]net.IP)
	return val, args.Error(1)
}

// config the gin engine for testing purpose
func setupRouter() *fiber.App {
	r := fiber.New()
	return r
}

func Test_appInfo(t *testing.T) {
	infoMock := infoServiceMock{}
	info := Info{
		AppName:     "fiber unit test",
		Description: "This is fiber unit test",
		Version:     "1.5.0",
		Ip:          "1.2.3.4",
	}

	// instantiate application config object
	app.Config = app.AppConfig{
		Server: app.ServerConfig{
			Host: "localhost",
		},
	}

	ipMock := ipServiceMock{}

	infoMock.On("getAppInfo").Return(info, nil)
	ipMock.On("lookupIP", "localhost").Return([]net.IP{net.ParseIP("1.2.3.4")}, nil)
	infoRouter := NewRouter(&infoMock, &ipMock)

	r := setupRouter()
	infoRouter.Register(r.Group("/info"))

	req := httptest.NewRequest("GET", "/info", nil)
	resp, _ := r.Test(req, 5)

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, _ := io.ReadAll(resp.Body)
	// unmarshal to Pet struct for asserts.
	actualInfoResponse := Info{}
	err := json.Unmarshal(body, &actualInfoResponse)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
		return
	}

	assert.Equal(t, info.AppName, actualInfoResponse.AppName)
	assert.Equal(t, info.Ip, actualInfoResponse.Ip)
	assert.Equal(t, info.Version, actualInfoResponse.Version)
}
