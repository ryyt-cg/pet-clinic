package info

import (
	"encoding/json"
	"errors"
	"gin-petclinic-service/config/app"
	resterr "gin-petclinic-service/internal/errors"
	"gin-petclinic-service/internal/test"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
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

	val := intf.(*Info)
	return val, args.Error(1)
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

func Test_appInfo(t *testing.T) {
	info := &Info{
		AppName:     "Gin unit test",
		Description: "This is Gin unit test",
		Version:     "1.5.0",
	}

	// instantiate application config object
	app.Config = app.AppConfig{
		Server: app.ServerConfig{
			Host: "localhost",
		},
	}

	infoMock := infoServiceMock{}
	ipMock := ipServiceMock{}

	infoMock.On("getAppInfo").Return(info, nil)
	ipMock.On("lookupIP", "localhost").Return([]net.IP{net.ParseIP("1.2.3.4")}, nil)
	infoRouter := NewRouter(&infoMock, &ipMock)

	r := test.SetupRouter()
	infoRouter.Register(r.Group("/info"))
	expectedResponse := Info{
		AppName:     "Gin unit test",
		Description: "This is Gin unit test",
		Version:     "1.5.0",
		Ip:          "1.2.3.4",
	}
	jsonStr, _ := json.Marshal(expectedResponse)
	tests := []test.APITestCase{
		{"test appInfo", "GET", "/info", "", nil, http.StatusOK, string(jsonStr)},
	}
	for _, tc := range tests {
		test.Endpoint(t, r, tc)
	}
}

func Test_appInfoWithErrors(t *testing.T) {
	info := &Info{
		AppName:     "Gin unit test",
		Description: "This is Gin unit test",
		Version:     "1.5.0",
		Ip:          "Unknown host",
	}

	tests := []struct {
		mockResult     *Info
		mockError      error
		mockIPResult   []net.IP
		mockIPError    error
		statusCode     int
		expectedResult interface{}
	}{
		{
			mockResult:     info,
			mockError:      nil,
			mockIPResult:   nil,
			mockIPError:    errors.New("unable to fetch IPs"),
			statusCode:     http.StatusOK,
			expectedResult: info,
		},
		{
			mockResult:     nil,
			mockError:      errors.New("unable to fetch application info"),
			statusCode:     http.StatusInternalServerError,
			expectedResult: resterr.InternalServerError("unable to fetch application info"),
		},
	}

	// instantiate application config object
	app.Config = app.AppConfig{
		Server: app.ServerConfig{
			Host: "localhost",
		},
	}

	for _, tc := range tests {
		infoMock := infoServiceMock{}
		ipMock := ipServiceMock{}

		r := test.SetupRouter()
		infoRouter := NewRouter(&infoMock, &ipMock)
		infoRouter.Register(r.Group("/info"))
		infoMock.On("getAppInfo").Return(tc.mockResult, tc.mockError)
		if tc.statusCode == http.StatusOK {
			ipMock.On("lookupIP", "localhost").Return(tc.mockIPResult, tc.mockIPError)
		}

		jsonStr, _ := json.Marshal(tc.expectedResult)
		tests := []test.APITestCase{
			{"test appInfo", "GET", "/info", "", nil, tc.statusCode, string(jsonStr)},
		}
		for _, tc := range tests {
			test.Endpoint(t, r, tc)
		}
	}
}
