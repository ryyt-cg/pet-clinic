package vet

import (
	"bytes"
	"encoding/json"
	resterr "fiber3-petclinic-service/internal/errors"
	"fiber3-petclinic-service/internal/repository"
	"fiber3-petclinic-service/internal/repository/model"
	"fiber3-petclinic-service/internal/test"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_allSpecialties(t *testing.T) {
	specialties := &specialtiesResponse{
		Context: model.Context{
			Count: 2,
		},
		Specialties: []specialtyResponse{
			{
				ID:   1,
				Name: "Radiology",
			},
			{
				ID:   2,
				Name: "Surgery",
			},
		},
	}

	noSpecialtiesResponse := &specialtiesResponse{
		Context:     model.Context{},
		Specialties: []specialtyResponse{},
	}

	testCases := []struct {
		name            string
		mockSpecialties *specialtiesResponse
		mockErr         error
		statusCode      int
		expectedResult  interface{}
	}{
		{
			name:            "get all specialties",
			mockSpecialties: specialties,
			mockErr:         nil,
			statusCode:      http.StatusOK,
			expectedResult:  specialties,
		},
		{
			name:            "get no specialties",
			mockSpecialties: noSpecialtiesResponse,
			mockErr:         nil,
			statusCode:      http.StatusOK,
			expectedResult:  noSpecialtiesResponse,
		},
		{
			name:            "fail to get all specialties",
			mockSpecialties: nil,
			mockErr:         resterr.InternalServerError("failed to get specialties"),
			statusCode:      http.StatusInternalServerError,
			expectedResult:  resterr.InternalServerError("failed to get specialties"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock service
			vetMock := NewMockServicer(t)
			vetMock.EXPECT().getAllSpecialties().Return(tc.mockSpecialties, tc.mockErr)

			// Create a new router with the mock service
			vetRouter := NewRouter(vetMock)

			// Setup the router
			r := test.SetupRouter()
			v1 := r.Group("/v1")
			vetRouter.Register(v1.Group("/vets"))

			// Create a request to the endpoint
			req, _ := http.NewRequest("GET", "/v1/vets/specialties", nil)
			resp, _ := r.Test(req, fiber.TestConfig{})

			switch tc.statusCode {
			case http.StatusOK:
				actualSpecialtiesResponse := &specialtiesResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualSpecialtiesResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResult.(*specialtiesResponse), actualSpecialtiesResponse)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualPetResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResult.(resterr.ErrorResponse).Status, actualPetResponse.Status)
				assert.Equal(t, tc.expectedResult.(resterr.ErrorResponse).Message, actualPetResponse.Message)
			}
		})
	}
}

func Test_allVets(t *testing.T) {
	vetsResponse := &responses{
		Context: model.Context{
			Count: 2,
		},
		Vets: []response{
			{
				ID:        1,
				FirstName: "Nat",
				LastName:  "Cole",
			},
			{
				ID:        2,
				FirstName: "John",
				LastName:  "Smith",
			},
		},
	}

	noVetsResponse := &responses{
		Context: model.Context{},
		Vets:    []response{},
	}

	testCases := []struct {
		name           string
		mockVets       *responses
		mockErr        error
		statusCode     int
		expectedResult interface{}
	}{
		{
			name:           "get all vets",
			mockVets:       vetsResponse,
			mockErr:        nil,
			statusCode:     http.StatusOK,
			expectedResult: vetsResponse,
		},
		{
			name:           "get no vets",
			mockVets:       noVetsResponse,
			mockErr:        nil,
			statusCode:     http.StatusOK,
			expectedResult: noVetsResponse,
		},
		{
			name:           "fail to get all vets",
			mockVets:       nil,
			mockErr:        resterr.InternalServerError("failed to get specialties"),
			statusCode:     http.StatusInternalServerError,
			expectedResult: resterr.InternalServerError("failed to get specialties"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock service
			vetMock := NewMockServicer(t)
			vetMock.EXPECT().getAllVets().Return(tc.mockVets, tc.mockErr)

			// Create a new router with the mock service
			vetRouter := NewRouter(vetMock)

			// Setup the router
			r := test.SetupRouter()
			v1 := r.Group("/v1")
			vetRouter.Register(v1.Group("/vets"))

			// Create a request to the endpoint
			req, _ := http.NewRequest("GET", "/v1/vets/all", nil)
			resp, _ := r.Test(req, fiber.TestConfig{})

			switch tc.statusCode {
			case http.StatusOK:
				actualVetsResponse := &responses{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualVetsResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResult.(*responses), actualVetsResponse)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualPetResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResult.(resterr.ErrorResponse).Status, actualPetResponse.Status)
				assert.Equal(t, tc.expectedResult.(resterr.ErrorResponse).Message, actualPetResponse.Message)
			}
		})
	}
}

func Test_vetById(t *testing.T) {
	vetResponse := &response{
		ID:        1,
		FirstName: "Code",
		LastName:  "Ninjas",
	}

	testCases := []struct {
		name           string
		id             uint
		mockVet        *response
		mockErr        error
		route          string
		statusCode     int
		expectedResult interface{}
	}{
		{
			name:           "get vet by id",
			id:             1,
			mockVet:        vetResponse,
			mockErr:        nil,
			route:          "/v1/vets/1",
			statusCode:     http.StatusOK,
			expectedResult: vetResponse,
		},
		{
			name:           "found no vet by id",
			id:             2,
			mockVet:        nil,
			mockErr:        gorm.ErrRecordNotFound,
			route:          "/v1/vets/2",
			statusCode:     http.StatusNotFound,
			expectedResult: gorm.ErrRecordNotFound,
		},
		{
			name:           "fail to get vet by id",
			id:             1,
			mockVet:        nil,
			mockErr:        resterr.InternalServerError("failed to get specialties"),
			route:          "/v1/vets/1",
			statusCode:     http.StatusInternalServerError,
			expectedResult: resterr.InternalServerError("failed to get specialties"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock service
			vetMock := NewMockServicer(t)
			vetMock.EXPECT().getVetById(tc.id).Return(tc.mockVet, tc.mockErr)

			// Create a new router with the mock service
			vetRouter := NewRouter(vetMock)

			// Setup the router
			r := test.SetupRouter()
			v1 := r.Group("/v1")
			vetRouter.Register(v1.Group("/vets"))

			// Create a request to the endpoint
			req, _ := http.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, fiber.TestConfig{})

			switch tc.statusCode {
			case http.StatusOK:
				actualVetsResponse := &response{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualVetsResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResult.(*response), actualVetsResponse)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualPetResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResult.(resterr.ErrorResponse).Status, actualPetResponse.Status)
				assert.Equal(t, tc.expectedResult.(resterr.ErrorResponse).Message, actualPetResponse.Message)
			}
		})
	}
}

func Test_vetByID_BadRequest(t *testing.T) {
	testCases := []struct {
		name           string
		id             interface{}
		route          string
		mockVet        *response
		expectedResult resterr.ErrorResponse
	}{
		{
			name:           "bad request with invalid id",
			id:             "invalid",
			route:          "/v1/vets/invalid",
			mockVet:        nil,
			expectedResult: resterr.BadRequest("strconv.Atoi: parsing \"invalid\": invalid syntax"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock service
			vetMock := NewMockServicer(t)

			// Create a new router with the mock service
			vetRouter := NewRouter(vetMock)
			r := test.SetupRouter()
			v1 := r.Group("/v1")
			vetRouter.Register(v1.Group("/vets"))

			req, _ := http.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, fiber.TestConfig{})
			actualVetsResponse := &resterr.ErrorResponse{}

			// Read the response body
			body, _ := io.ReadAll(resp.Body)
			err := json.Unmarshal(body, actualVetsResponse)
			if err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
				return
			}
			assert.Equal(t, tc.expectedResult.Status, resp.StatusCode)
			assert.Equal(t, tc.expectedResult.Message, actualVetsResponse.Message)
		})
	}
}

func Test_getVetByIdWithSpecialties(t *testing.T) {
	vetResponse := &response{
		ID:        1,
		FirstName: "Code",
		LastName:  "Ninjas",
		Specialties: &[]specialtyResponse{
			{
				ID:   1,
				Name: "Radiology",
			},
			{
				ID:   2,
				Name: "Surgery",
			},
		},
	}

	testCases := []struct {
		name           string
		id             uint
		mockVet        *response
		mockErr        error
		route          string
		statusCode     int
		expectedResult interface{}
	}{
		{
			name:           "get vet by id",
			id:             1,
			mockVet:        vetResponse,
			mockErr:        nil,
			route:          "/v1/vets/1/specialties",
			statusCode:     http.StatusOK,
			expectedResult: vetResponse,
		},
		{
			name:           "found no vet by id",
			id:             2,
			mockVet:        nil,
			mockErr:        gorm.ErrRecordNotFound,
			route:          "/v1/vets/2/specialties",
			statusCode:     http.StatusNotFound,
			expectedResult: gorm.ErrRecordNotFound,
		},
		{
			name:           "fail to get vet by id",
			id:             1,
			mockVet:        nil,
			mockErr:        resterr.InternalServerError("failed to get specialties"),
			route:          "/v1/vets/1/specialties",
			statusCode:     http.StatusInternalServerError,
			expectedResult: resterr.InternalServerError("failed to get specialties"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock service
			vetMock := NewMockServicer(t)
			vetMock.EXPECT().getVetByIdWithSpecialties(tc.id).Return(tc.mockVet, tc.mockErr)

			// Create a new router with the mock service
			vetRouter := NewRouter(vetMock)

			// Setup the router
			r := test.SetupRouter()
			v1 := r.Group("/v1")
			vetRouter.Register(v1.Group("/vets"))

			// Create a request to the endpoint
			req, _ := http.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, fiber.TestConfig{})

			switch tc.statusCode {
			case http.StatusOK:
				actualVetsResponse := &response{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualVetsResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResult.(*response), actualVetsResponse)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualPetResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResult.(resterr.ErrorResponse).Status, actualPetResponse.Status)
				assert.Equal(t, tc.expectedResult.(resterr.ErrorResponse).Message, actualPetResponse.Message)
			}
		})
	}
}

func Test_vetByIDWithSpecialties_BadRequest(t *testing.T) {
	testCases := []struct {
		name           string
		id             interface{}
		route          string
		expectedResult resterr.ErrorResponse
	}{
		{
			name:           "bad request with invalid id",
			id:             "invalid",
			route:          "/v1/vets/invalid/specialties",
			expectedResult: resterr.BadRequest("strconv.Atoi: parsing \"invalid\": invalid syntax"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock service
			vetMock := NewMockServicer(t)
			vetRouter := NewRouter(vetMock)
			r := test.SetupRouter()
			v1 := r.Group("/v1")
			vetRouter.Register(v1.Group("/vets"))

			req, _ := http.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, fiber.TestConfig{})

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			actualVetResponse := &resterr.ErrorResponse{}
			// Read the response body
			body, _ := io.ReadAll(resp.Body)
			err := json.Unmarshal(body, actualVetResponse)
			if err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
				return
			}
			assert.Equal(t, tc.expectedResult.Status, actualVetResponse.Status)
			assert.Equal(t, tc.expectedResult.Message, actualVetResponse.Message)
		})
	}
}

func Test_VetByLastName(t *testing.T) {
}

func Test_createVet(t *testing.T) {
	vetEntity := &repository.Vet{
		Person: model.Person{
			FirstName: "New",
			LastName:  "Vet",
		},
	}

	mockVet := &response{
		ID:        1,
		FirstName: "New",
		LastName:  "Vet",
	}

	createRequest := &addRequest{
		FirstName: "New",
		LastName:  "Vet",
	}

	testCases := []struct {
		name           string
		request        *addRequest
		vetEntity      *repository.Vet
		mockVet        *response
		mockErr        error
		statusCode     int
		expectedResult interface{}
	}{
		{
			name:           "create new vet",
			request:        createRequest,
			vetEntity:      vetEntity,
			mockVet:        mockVet,
			mockErr:        nil,
			statusCode:     http.StatusCreated,
			expectedResult: &response{ID: 1, FirstName: "New", LastName: "Vet"},
		},
		{
			name:           "fail to create new vet",
			request:        createRequest,
			vetEntity:      vetEntity,
			mockVet:        nil,
			mockErr:        gorm.ErrRecordNotFound,
			statusCode:     http.StatusNotFound,
			expectedResult: gorm.ErrRecordNotFound,
		},
		{
			name:           "fail to create new vet",
			request:        createRequest,
			vetEntity:      vetEntity,
			mockVet:        nil,
			mockErr:        resterr.InternalServerError("failed to create vet"),
			statusCode:     http.StatusInternalServerError,
			expectedResult: resterr.InternalServerError("failed to create vet"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock service
			vetMock := NewMockServicer(t)
			vetMock.EXPECT().create(tc.vetEntity).Return(tc.mockVet, tc.mockErr)

			// Create a new router with the mock service
			vetRouter := NewRouter(vetMock)
			r := test.SetupRouter()
			v1 := r.Group("/v1")
			vetRouter.Register(v1.Group("/vets"))

			requestBody, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/v1/vets", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, fiber.TestConfig{})

			switch tc.statusCode {
			case http.StatusCreated:
				actualVetsResponse := &response{}
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualVetsResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResult.(*response), actualVetsResponse)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualPetResponse := &resterr.ErrorResponse{}
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResult.(resterr.ErrorResponse).Status, actualPetResponse.Status)
				assert.Equal(t, tc.expectedResult.(resterr.ErrorResponse).Message, actualPetResponse.Message)
			}
		})
	}
}

func Test_createVet_BadRequest(t *testing.T) {
	vetRequest := `{
	 "firstName": "Sammy",
	 "lastName": 6
	}`

	testCases := []struct {
		name             string
		request          string
		route            string
		statusCode       int
		expectedResponse resterr.ErrorResponse
	}{
		{
			name:             "create vet with bad request",
			request:          vetRequest,
			route:            "/v1/vets",
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("json: cannot unmarshal number into Go struct field addRequest.lastName of type string"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			vetMock := NewMockServicer(t)
			vetRouter := NewRouter(vetMock)
			vetRouter.Register(v1.Group("/vets"))

			req := httptest.NewRequest("POST", tc.route, bytes.NewReader([]byte(tc.request)))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, fiber.TestConfig{})

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			actualErrorResponse := resterr.ErrorResponse{}
			body, _ := io.ReadAll(resp.Body)
			err := json.Unmarshal(body, &actualErrorResponse)
			if err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
				return
			}
			assert.Equal(t, tc.statusCode, actualErrorResponse.Status)
			assert.Equal(t, tc.expectedResponse, actualErrorResponse)
		})
	}
}

func Test_updateVet(t *testing.T) {
	vetEntity := &repository.Vet{
		Model: gorm.Model{
			ID: 1,
		},
		Person: model.Person{
			FirstName: "Updated",
			LastName:  "Vet",
		},
	}

	mockVet := &response{
		ID:        1,
		FirstName: "Updated",
		LastName:  "Vet",
	}

	updateVetRequest := &updateRequest{
		ID:        1,
		FirstName: "Updated",
		LastName:  "Vet",
	}

	testCases := []struct {
		name           string
		id             uint
		request        *updateRequest
		vetEntity      *repository.Vet
		mockVet        *response
		mockErr        error
		route          string
		statusCode     int
		expectedResult interface{}
	}{
		{
			name:           "update vet by id",
			id:             1,
			request:        updateVetRequest,
			vetEntity:      vetEntity,
			mockVet:        mockVet,
			mockErr:        nil,
			route:          "/v1/vets/1",
			statusCode:     http.StatusOK,
			expectedResult: &response{ID: 1, FirstName: "Updated", LastName: "Vet"},
		},
		{
			name:           "fail to update vet by id",
			id:             1,
			request:        updateVetRequest,
			vetEntity:      vetEntity,
			mockVet:        nil,
			mockErr:        gorm.ErrRecordNotFound,
			route:          "/v1/vets/1",
			statusCode:     http.StatusNotFound,
			expectedResult: gorm.ErrRecordNotFound,
		},
		{
			name:           "fail to update vet by id with internal error",
			id:             1,
			request:        updateVetRequest,
			vetEntity:      vetEntity,
			mockVet:        nil,
			mockErr:        resterr.InternalServerError("failed to update vet"),
			route:          "/v1/vets/1",
			statusCode:     http.StatusInternalServerError,
			expectedResult: resterr.InternalServerError("failed to update vet"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock service
			vetMock := NewMockServicer(t)
			vetMock.EXPECT().update(tc.vetEntity).Return(tc.mockVet, tc.mockErr)

			// Create a new router with the mock service
			vetRouter := NewRouter(vetMock)
			r := test.SetupRouter()
			v1 := r.Group("/v1")
			vetRouter.Register(v1.Group("/vets"))

			requestBody, _ := json.Marshal(tc.request)
			// Create a request to the endpoint
			req, _ := http.NewRequest("PUT", tc.route, bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, fiber.TestConfig{})

			switch tc.statusCode {
			case http.StatusOK:
				actualVetsResponse := &response{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualVetsResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResult.(*response), actualVetsResponse)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualPetResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResult.(resterr.ErrorResponse).Status, actualPetResponse.Status)
				assert.Equal(t, tc.expectedResult.(resterr.ErrorResponse).Message, actualPetResponse.Message)
			}
		})
	}
}

func Test_updateVet_BadRequest(t *testing.T) {
	vetRequest := `{
	  "id": 1,
	  "firstName": "Sammy",
	  "lastName": "Sosa"
	}`

	vetRequest2 := `{
	 "id": "1",
	 "firstName": "Sammy",
	 "lastName": "Sosa"
	}`

	testCases := []struct {
		name             string
		id               interface{}
		request          string
		route            string
		statusCode       int
		expectedResponse resterr.ErrorResponse
	}{
		{
			name:             "update vet by ID with invalid ID",
			id:               "invalid-id",
			request:          vetRequest,
			route:            "/v1/vets/invalid-id",
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("strconv.Atoi: parsing \"invalid-id\": invalid syntax"),
		},
		{
			name:             "update vet by ID with bad request",
			id:               1,
			request:          vetRequest2,
			route:            "/v1/vets/1",
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("json: cannot unmarshal string into Go struct field updateRequest.id of type uint"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			vetMock := NewMockServicer(t)
			vetRouter := NewRouter(vetMock)
			vetRouter.Register(v1.Group("/vets"))

			req := httptest.NewRequest("PUT", tc.route, bytes.NewReader([]byte(tc.request)))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, fiber.TestConfig{})

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			actualErrorResponse := resterr.ErrorResponse{}
			body, _ := io.ReadAll(resp.Body)
			err := json.Unmarshal(body, &actualErrorResponse)
			if err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
				return
			}
			assert.Equal(t, tc.statusCode, actualErrorResponse.Status)
			assert.Equal(t, tc.expectedResponse, actualErrorResponse)
		})
	}
}
