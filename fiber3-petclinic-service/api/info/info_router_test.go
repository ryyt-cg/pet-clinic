package info

import (
	"encoding/json"
	"errors"
	"fiber3-petclinic-service/config/app"
	resterr "fiber3-petclinic-service/pkg/errors"
	"fiber3-petclinic-service/pkg/test"
	"github.com/gofiber/fiber/v3"
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
		AppName:     "fiber unit test",
		Description: "This is fiber unit test",
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

	req := httptest.NewRequest("GET", "/info", nil)
	resp, _ := r.Test(req, fiber.TestConfig{})

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
	assert.Equal(t, "1.2.3.4", actualInfoResponse.Ip)
	assert.Equal(t, info.Version, actualInfoResponse.Version)
}

func Test_appInfoWithErrors(t *testing.T) {
	info := &Info{
		AppName:     "fiber error test",
		Description: "This is fiber error test",
		Version:     "1.5.0",
	}

	expectError := resterr.InternalServerError("unable to fetch application info")

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
			expectedResult: &expectError,
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

		switch tc.statusCode {
		case http.StatusOK:
			ipMock.On("lookupIP", "localhost").Return(tc.mockIPResult, tc.mockIPError)
		case http.StatusInternalServerError:

		}

		req := httptest.NewRequest("GET", "/info", nil)
		resp, _ := r.Test(req, fiber.TestConfig{})
		// Read the response body
		body, _ := io.ReadAll(resp.Body)

		switch resp.StatusCode {
		case http.StatusOK:
			actualInfoResponse := Info{}
			err := json.Unmarshal(body, &actualInfoResponse)
			if err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
				return
			}
			assert.Equal(t, tc.statusCode, resp.StatusCode)
			assert.Equal(t, tc.expectedResult, &actualInfoResponse)
		case http.StatusInternalServerError:
			actualInfoResponse := resterr.ErrorResponse{}
			err := json.Unmarshal(body, &actualInfoResponse)
			if err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
				return
			}
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			assert.Equal(t, tc.expectedResult, &actualInfoResponse)
		}
	}
}
