package pet

import (
	"bytes"
	"encoding/json"
	"errors"
	resterr "fiber-petclinic-service/internal/errors"
	"fiber-petclinic-service/internal/repository"
	"fiber-petclinic-service/internal/repository/model"
	"fiber-petclinic-service/internal/test"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_getAllPets(t *testing.T) {
	mockEntityPets := []repository.Pet{
		{Model: gorm.Model{ID: 1}, Name: "Tom",
			Birthdate: test.ToDate("2015-11-19"),
			Visits:    nil,
			SpeciesID: 19, OwnerID: 7},
		{Model: gorm.Model{ID: 2}, Name: "Mike",
			Birthdate: test.ToDate("2018-04-17"),
			Visits:    nil,
			SpeciesID: 20, OwnerID: 7},
	}

	mockPets := &responses{
		Context: model.Context{
			Count: 2,
		},
		Pets: fromPets(mockEntityPets),
	}

	tests := []struct {
		name             string
		mockPets         *responses
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get all pets",
			mockPets:         mockPets,
			mockError:        nil,
			route:            "/v1/pets/all",
			statusCode:       http.StatusOK,
			expectedResponse: mockPets,
		},
		{
			name:             "get no pet",
			mockPets:         nil,
			mockError:        gorm.ErrRecordNotFound,
			route:            "/v1/pets/all",
			statusCode:       http.StatusNotFound,
			expectedResponse: resterr.NotFound("Pet not found"),
		},
		{
			name:             "fail to get pets",
			mockPets:         nil,
			mockError:        errors.New("fail to get pets"),
			route:            "/v1/pets/all",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: resterr.InternalServerError("fail to get pets"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			petMock := NewMockServicer(t)
			petMock.EXPECT().getAllPets().Return(tc.mockPets, tc.mockError)

			petRouter := NewRouter(petMock)
			petRouter.Register(v1.Group("/pets"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)
			assert.Equal(t, tc.statusCode, resp.StatusCode)

			switch resp.StatusCode {
			case http.StatusOK:
				actualPetResponse := &responses{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(*responses).Context.Count, actualPetResponse.Context.Count)
				assert.Equal(t, tc.expectedResponse.(*responses).Pets, actualPetResponse.Pets)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualPetResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Status, actualPetResponse.Status)
				assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Message, actualPetResponse.Message)
			}
		})
	}
}

func Test_getPetById(t *testing.T) {
	mockPet := &response{
		ID:        1,
		Name:      "Running Water",
		Birthdate: "2018-02-26",
	}

	tests := []struct {
		name             string
		id               uint
		mockPet          *response
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get pet by id",
			id:               1,
			mockPet:          mockPet,
			mockError:        nil,
			route:            "/v1/pets/1",
			statusCode:       http.StatusOK,
			expectedResponse: mockPet,
		},
		{
			name:             "get no pet by id",
			id:               1,
			mockPet:          nil,
			mockError:        gorm.ErrRecordNotFound,
			route:            "/v1/pets/1",
			statusCode:       http.StatusNotFound,
			expectedResponse: resterr.NotFound("Pet not found"),
		},
		{
			name:             "fail to get pet by id",
			id:               1,
			mockPet:          nil,
			mockError:        resterr.InternalServerError("unable to get pet by id"),
			route:            "/v1/pets/1",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: resterr.InternalServerError("unable to get pet by id"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			petMock := NewMockServicer(t)
			petMock.EXPECT().getPetById(tc.id).Return(tc.mockPet, tc.mockError)

			petRouter := NewRouter(petMock)
			petRouter.Register(v1.Group("/pets"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)
			assert.Equal(t, tc.statusCode, resp.StatusCode)

			switch resp.StatusCode {
			case http.StatusOK:
				actualPetResponse := &response{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(*response).ID, actualPetResponse.ID)
				assert.Equal(t, tc.expectedResponse.(*response).Name, actualPetResponse.Name)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualPetResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Status, actualPetResponse.Status)
				assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Message, actualPetResponse.Message)
			}
		})
	}
}

func Test_getById_BadRequest(t *testing.T) {
	testCases := []struct {
		name             string
		id               interface{}
		mockPet          *response
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get pet by invalid id",
			id:               "invalid_id",
			mockPet:          nil,
			mockError:        resterr.BadRequest("strconv.Atoi: parsing \"invalid_id\": invalid syntax"),
			route:            "/v1/pets/invalid_id",
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("strconv.Atoi: parsing \"invalid_id\": invalid syntax"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			petMock := NewMockServicer(t)
			petRouter := NewRouter(petMock)
			petRouter.Register(v1.Group("/pets"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)
			assert.Equal(t, tc.statusCode, resp.StatusCode)

			actualPetResponse := &resterr.ErrorResponse{}
			// Read the response body
			body, _ := io.ReadAll(resp.Body)
			err := json.Unmarshal(body, actualPetResponse)
			if err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
				return
			}
			assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Status, actualPetResponse.Status)
			assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Message, actualPetResponse.Message)
		})
	}
}

func Test_getPetByIdWithVisits(t *testing.T) {
	mockPet := &response{
		ID:        1,
		Name:      "Running Water",
		Birthdate: "2018-02-26",
		Species:   "Dog",
		Visits: []visitResponse{
			{
				ID:          1,
				VisitDate:   "2023-10-01",
				Description: "Regular check-up",
				PetID:       1,
			},
			{
				ID:          2,
				VisitDate:   "2023-10-15",
				Description: "Vaccination",
				PetID:       1,
			},
		},
	}

	tests := []struct {
		name             string
		id               uint
		mockPet          *response
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get pet by id with visits",
			id:               1,
			mockPet:          mockPet,
			mockError:        nil,
			route:            "/v1/pets/1/visits",
			statusCode:       http.StatusOK,
			expectedResponse: mockPet,
		},
		{
			name:             "get no pet by id",
			id:               1,
			mockPet:          nil,
			mockError:        gorm.ErrRecordNotFound,
			route:            "/v1/pets/1/visits",
			statusCode:       http.StatusNotFound,
			expectedResponse: resterr.NotFound("Pet not found"),
		},
		{
			name:             "fail to get pet by id",
			id:               1,
			mockPet:          nil,
			mockError:        resterr.InternalServerError("unable to get pet by id"),
			route:            "/v1/pets/1/visits",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: resterr.InternalServerError("unable to get pet by id"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			petMock := NewMockServicer(t)
			petMock.EXPECT().getPetWithVisitsById(tc.id).Return(tc.mockPet, tc.mockError)

			petRouter := NewRouter(petMock)
			petRouter.Register(v1.Group("/pets"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)
			assert.Equal(t, tc.statusCode, resp.StatusCode)

			switch resp.StatusCode {
			case http.StatusOK:
				actualPetResponse := &response{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(*response).ID, actualPetResponse.ID)
				assert.Equal(t, tc.expectedResponse.(*response).Name, actualPetResponse.Name)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualPetResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Status, actualPetResponse.Status)
				assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Message, actualPetResponse.Message)
			}
		})
	}
}

func Test_getByIdWithVisits_BadRequest(t *testing.T) {
	testCases := []struct {
		name             string
		id               interface{}
		mockPet          *response
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get pet by invalid id",
			id:               "invalid_id",
			mockPet:          nil,
			mockError:        resterr.BadRequest("strconv.Atoi: parsing \"invalid_id\": invalid syntax"),
			route:            "/v1/pets/invalid_id/visits",
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("strconv.Atoi: parsing \"invalid_id\": invalid syntax"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			petMock := NewMockServicer(t)
			petRouter := NewRouter(petMock)
			petRouter.Register(v1.Group("/pets"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)
			assert.Equal(t, tc.statusCode, resp.StatusCode)

			actualPetResponse := &resterr.ErrorResponse{}
			// Read the response body
			body, _ := io.ReadAll(resp.Body)
			err := json.Unmarshal(body, actualPetResponse)
			if err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
				return
			}
			assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Status, actualPetResponse.Status)
			assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Message, actualPetResponse.Message)
		})
	}
}

func Test_getPetsByName(t *testing.T) {
	mockPets := &responses{
		Context: model.Context{
			Count: 2,
		},
		Pets: []response{
			{ID: 1, Name: "Tom", Birthdate: "2015-11-19"},
			{ID: 2, Name: "Tom", Birthdate: "2018-04-17"},
		},
	}

	tests := []struct {
		name             string
		nameQuery        string
		mockPets         *responses
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get pets by name",
			nameQuery:        "Tom",
			mockPets:         mockPets,
			mockError:        nil,
			route:            "/v1/pets?name=Tom",
			statusCode:       http.StatusOK,
			expectedResponse: mockPets,
		},
		{
			name:             "get no pets by name",
			nameQuery:        "Unknown",
			mockPets:         nil,
			mockError:        gorm.ErrRecordNotFound,
			route:            "/v1/pets?name=Unknown",
			statusCode:       http.StatusNotFound,
			expectedResponse: resterr.NotFound("No pets found with this name."),
		},
		{
			name:             "fail to get pets by name",
			nameQuery:        "Unknown",
			mockPets:         nil,
			mockError:        errors.New("failed to fetch pets by name"),
			route:            "/v1/pets?name=Unknown",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: resterr.InternalServerError("failed to fetch pets by name"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			petMock := NewMockServicer(t)
			petMock.EXPECT().getPetsByName(tc.nameQuery).Return(tc.mockPets, tc.mockError)

			petRouter := NewRouter(petMock)
			petRouter.Register(v1.Group("/pets"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)
			assert.Equal(t, tc.statusCode, resp.StatusCode)

			switch resp.StatusCode {
			case http.StatusOK:
				actualPetResponse := &responses{}
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(*responses).Context.Count, actualPetResponse.Context.Count)
				assert.Equal(t, tc.expectedResponse.(*responses).Pets, actualPetResponse.Pets)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualPetResponse := &resterr.ErrorResponse{}
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Status, actualPetResponse.Status)
				assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Message, actualPetResponse.Message)
			}
		})
	}
}

func Test_getByName_BadRequest(t *testing.T) {
	testCases := []struct {
		name             string
		nameQuery        string
		mockPets         *responses
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get pets by empty name",
			nameQuery:        "",
			mockPets:         nil,
			mockError:        resterr.BadRequest("pet name is empty"),
			route:            "/v1/pets?name=",
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("pet name is empty"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			petMock := NewMockServicer(t)
			petRouter := NewRouter(petMock)
			petRouter.Register(v1.Group("/pets"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)
			assert.Equal(t, tc.statusCode, resp.StatusCode)

			actualPetResponse := &resterr.ErrorResponse{}
			body, _ := io.ReadAll(resp.Body)
			err := json.Unmarshal(body, actualPetResponse)
			if err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
				return
			}
			assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Status, actualPetResponse.Status)
			assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Message, actualPetResponse.Message)
		})
	}
}

func Test_createNewPet(t *testing.T) {
	addPetRequest := &addRequest{
		Name:      "Tom",
		Birthdate: "2015-11-19",
		SpeciesID: 19,
		OwnerID:   7,
	}

	mockPet := &addResponse{
		ID:        1,
		Name:      "Tom",
		Birthdate: "2015-11-19",
		SpeciesID: 19,
		OwnerID:   7,
	}

	petEntity := &repository.Pet{
		Name:      "Tom",
		Birthdate: test.ToDate("2015-11-19"),
		SpeciesID: 19,
		OwnerID:   7,
	}

	tests := []struct {
		name             string
		request          *addRequest
		mockPet          *addResponse
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "create new pet",
			request:          addPetRequest,
			mockPet:          mockPet,
			mockError:        nil,
			route:            "/v1/pets",
			statusCode:       http.StatusCreated,
			expectedResponse: mockPet,
		},
		{
			name:             "fail to create new pet",
			request:          addPetRequest,
			mockPet:          nil,
			mockError:        errors.New("failed to create pet"),
			route:            "/v1/pets",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: resterr.InternalServerError("failed to create pet"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			petMock := NewMockServicer(t)
			petMock.EXPECT().create(petEntity).Return(tc.mockPet, tc.mockError)

			petRouter := NewRouter(petMock)
			petRouter.Register(v1.Group("/pets"))

			reqBody, _ := json.Marshal(tc.request)
			req := httptest.NewRequest("POST", tc.route, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, 5)

			switch resp.StatusCode {
			case http.StatusCreated:
				actualPetResponse := &addResponse{}
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(*addResponse), actualPetResponse)
			case http.StatusInternalServerError:
				actualPetResponse := &resterr.ErrorResponse{}
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Status, actualPetResponse.Status)
				assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Message, actualPetResponse.Message)
			}
		})
	}
}

func Test_createNewPet_BadRequest(t *testing.T) {
	testCases := []struct {
		name             string
		request          *addRequest
		mockPet          *response
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "create new pet with invalid birthdate",
			request:          &addRequest{Name: "Tom", Birthdate: "2015-11-31"},
			mockPet:          nil,
			mockError:        resterr.BadRequest("parsing time \"2015-11-31\": day out of range"),
			route:            "/v1/pets",
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("parsing time \"2015-11-31\": day out of range"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			petMock := NewMockServicer(t)
			petRouter := NewRouter(petMock)
			petRouter.Register(v1.Group("/pets"))

			reqBody, _ := json.Marshal(tc.request)
			req := httptest.NewRequest("POST", tc.route, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, 5)

			actualPetResponse := &resterr.ErrorResponse{}
			body, _ := io.ReadAll(resp.Body)
			err := json.Unmarshal(body, actualPetResponse)
			if err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
				return
			}
			assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Status, actualPetResponse.Status)
			assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Message, actualPetResponse.Message)
		})
	}
}

func Test_updatePet(t *testing.T) {
	updatePetRequest := &updateRequest{
		Name:      "Tom",
		Birthdate: "2015-11-19",
		SpeciesID: 19,
		OwnerID:   7,
	}

	mockPet := &updateResponse{
		ID:        1,
		Name:      "Tom",
		Birthdate: "2015-11-19",
		SpeciesID: 19,
		OwnerID:   7,
	}

	petEntity := &repository.Pet{
		Model:     gorm.Model{ID: 1},
		Name:      "Tom",
		Birthdate: test.ToDate("2015-11-19"),
		SpeciesID: 19,
		OwnerID:   7,
	}

	tests := []struct {
		name             string
		id               uint
		request          *updateRequest
		mockPet          *updateResponse
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "update pet",
			id:               1,
			request:          updatePetRequest,
			mockPet:          mockPet,
			mockError:        nil,
			route:            "/v1/pets/1",
			statusCode:       http.StatusOK,
			expectedResponse: mockPet,
		},
		{
			name:             "fail to update pet",
			id:               1,
			request:          updatePetRequest,
			mockPet:          nil,
			mockError:        errors.New("failed to update pet"),
			route:            "/v1/pets/1",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: resterr.InternalServerError("failed to update pet"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			petMock := NewMockServicer(t)
			petMock.EXPECT().update(petEntity).Return(tc.mockPet, tc.mockError)

			r := test.SetupRouter()
			v1 := r.Group("/v1")
			petRouter := NewRouter(petMock)
			petRouter.Register(v1.Group("/pets"))

			reqBody, _ := json.Marshal(tc.request)
			req := httptest.NewRequest("PUT", tc.route, bytes.NewBuffer(reqBody))
			// Must se content-type header.
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, 5)

			switch resp.StatusCode {
			case http.StatusOK:
				actualPetResponse := &updateResponse{}
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(*updateResponse).ID, actualPetResponse.ID)
				assert.Equal(t, tc.expectedResponse.(*updateResponse).Name, actualPetResponse.Name)
			case http.StatusInternalServerError:
				actualPetResponse := &resterr.ErrorResponse{}
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Status, actualPetResponse.Status)
				assert.Equal(t, tc.expectedResponse.(resterr.ErrorResponse).Message, actualPetResponse.Message)
			}
		})
	}
}

func Test_updatePet_BadRequest(t *testing.T) {
	testCases := []struct {
		name             string
		id               interface{}
		request          *updateRequest
		mockPet          *response
		mockError        error
		route            string
		statusCode       int
		expectedResponse resterr.ErrorResponse
	}{
		{
			name:             "update pet with invalid id",
			id:               "invalid_id",
			request:          &updateRequest{Name: "Tom", Birthdate: "2015-11-19"},
			mockPet:          nil,
			mockError:        resterr.BadRequest("strconv.Atoi: parsing \"invalid_id\": invalid syntax"),
			route:            "/v1/pets/invalid_id",
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("strconv.Atoi: parsing \"invalid_id\": invalid syntax"),
		},
		{
			name:             "update pet with invalid birthdate",
			id:               1,
			request:          &updateRequest{Name: "Tom", Birthdate: "2015-11-31"}, // Invalid date
			mockPet:          nil,
			mockError:        resterr.BadRequest("parsing time \"2015-11-31\": day out of range"),
			route:            "/v1/pets/1",
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("parsing time \"2015-11-31\": day out of range"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			petMock := NewMockServicer(t)
			petRouter := NewRouter(petMock)
			petRouter.Register(v1.Group("/pets"))

			reqBody, _ := json.Marshal(tc.request)
			req := httptest.NewRequest("PUT", tc.route, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, 5)

			actualPetResponse := &resterr.ErrorResponse{}
			body, _ := io.ReadAll(resp.Body)
			err := json.Unmarshal(body, actualPetResponse)
			if err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
				return
			}
			assert.Equal(t, tc.expectedResponse.Status, actualPetResponse.Status)
			assert.Equal(t, tc.expectedResponse.Message, actualPetResponse.Message)
		})
	}
}
