package owner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fiber-petclinic-service/api/pet"
	resterr "fiber-petclinic-service/internal/errors"
	"fiber-petclinic-service/internal/repository/model"
	"fiber-petclinic-service/internal/test"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_AllOwners(t *testing.T) {
	ownersResponses := &Responses{
		Context: model.Context{Count: 2},
		Owners: []Response{
			{
				ID:        35,
				FirstName: "Charles",
				LastName:  "Ward",
			},
			{
				ID:        26,
				FirstName: "John",
				LastName:  "Ward",
			},
			{
				ID:        28,
				FirstName: "John",
				LastName:  "Ward",
			},
		},
	}

	noOwnersFoundResponse := &Responses{
		Context: model.Context{Count: 0},
		Owners:  []Response{},
	}

	errorResponse := resterr.InternalServerError("fail to get all owners")

	testCases := []struct {
		name             string
		mockAllOwner     *Responses
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get all owners",
			mockAllOwner:     ownersResponses,
			mockError:        nil,
			route:            "/v1/owners/all",
			statusCode:       http.StatusOK,
			expectedResponse: ownersResponses,
		},
		{
			name:             "get no owner",
			mockAllOwner:     noOwnersFoundResponse,
			mockError:        nil,
			route:            "/v1/owners/all",
			statusCode:       http.StatusNotFound,
			expectedResponse: noOwnersFoundResponse,
		},
		{
			name:             "fail to get all owners",
			mockAllOwner:     nil,
			mockError:        gorm.ErrRecordNotFound,
			route:            "/v1/owners/all",
			statusCode:       http.StatusNotFound,
			expectedResponse: resterr.NotFound("Find no owners"),
		},
		{
			name:             "fail to get all owners",
			mockAllOwner:     nil,
			mockError:        errors.New("fail to get all owners"),
			route:            "/v1/owners/all",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: &errorResponse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := NewMockServicer(t)
			ownerMock.EXPECT().getAllOwners().Return(tc.mockAllOwner, tc.mockError)

			r := test.SetupRouter()
			v1 := r.Group("/v1")
			ownerRouter := NewRouter(ownerMock)
			ownerRouter.Register(v1.Group("/owners"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)

			switch tc.statusCode {
			case http.StatusOK | http.StatusNotFound:
				actualPetResponses := &Responses{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponses)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*Responses), actualPetResponses)
				assert.Equal(t, tc.expectedResponse.(*Responses).Context.Count, actualPetResponses.Context.Count)
				assert.Equal(t, tc.expectedResponse.(*Responses).Owners, actualPetResponses.Owners)
			case http.StatusInternalServerError:
				actualPetResponses := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponses)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}

				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*resterr.ErrorResponse), actualPetResponses)
			}

		})
	}
}

func Test_ownerById(t *testing.T) {
	ownerResponse := &Response{
		ID:        1,
		FirstName: "Nat",
		LastName:  "Cole",
		Address:   "1234 Elm St",
		City:      "New York",
		Telephone: "1234567890",
	}
	notFound := resterr.NotFound("record not found")
	internalError := resterr.InternalServerError("fail to get owner by id")

	testCases := []struct {
		name             string
		id               uint
		mockOwner        *Response
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get owner by ID",
			id:               1,
			mockOwner:        ownerResponse,
			mockError:        nil,
			route:            "/v1/owners/1",
			statusCode:       http.StatusOK,
			expectedResponse: ownerResponse,
		},
		{
			name:             "get no owner by ID",
			id:               1,
			mockOwner:        nil,
			mockError:        gorm.ErrRecordNotFound,
			route:            "/v1/owners/1",
			statusCode:       http.StatusNotFound,
			expectedResponse: &notFound,
		},
		{
			name:             "fail to get owner by ID",
			id:               1,
			mockOwner:        nil,
			mockError:        errors.New("fail to get owner by id"),
			route:            "/v1/owners/1",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: &internalError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := NewMockServicer(t)
			ownerMock.EXPECT().getOwnerById(tc.id).Return(tc.mockOwner, tc.mockError)

			r := test.SetupRouter()
			v1 := r.Group("/v1")
			ownerRouter := NewRouter(ownerMock)
			ownerRouter.Register(v1.Group("/owners"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)

			switch tc.statusCode {
			case http.StatusOK:
				actualOwnerResponse := &Response{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualOwnerResponse)

				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*Response), actualOwnerResponse)
				assert.Equal(t, tc.expectedResponse.(*Response).ID, actualOwnerResponse.ID)
				assert.Equal(t, tc.expectedResponse.(*Response).City, actualOwnerResponse.City)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualOwnerResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualOwnerResponse)

				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}

				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*resterr.ErrorResponse), actualOwnerResponse)
			}
		})
	}
}

func Test_ownerWithBadId(t *testing.T) {
	badRequest := resterr.BadRequest("failed to convert: strconv.Atoi: parsing \"invalid-id\": invalid syntax")
	ownerMock := NewMockServicer(t)

	r := test.SetupRouter()
	v1 := r.Group("/v1")
	ownerRouter := NewRouter(ownerMock)
	ownerRouter.Register(v1.Group("/owners"))

	req := httptest.NewRequest("GET", "/v1/owners/invalid-id", nil)
	resp, _ := r.Test(req, 5)

	actualOwnerResponse := &resterr.ErrorResponse{}
	// Read the response body
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal(body, actualOwnerResponse)

	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, &badRequest, actualOwnerResponse)
}

func Test_ownerByIdWithPets(t *testing.T) {
	ownerWithPetsResponse := &Response{
		ID:        1,
		FirstName: "Nat",
		LastName:  "Cole",
		Address:   "1234 Elm St",
		City:      "New York",
		Telephone: "1234567890",
		Pets: []pet.Response{
			{
				ID:        1,
				Name:      "Max",
				Birthdate: "2010-09-07",
				Species:   "Dog",
			},
			{
				ID:        2,
				Name:      "Lucky",
				Birthdate: "2015-09-07",
				Species:   "Dog",
			},
		},
	}

	ownerResponse := &Response{
		ID:        2,
		FirstName: "John",
		LastName:  "Smith",
		Address:   "5321 Oak St",
		City:      "New Jersey",
		Telephone: "1234567890",
	}
	notFoundResponse := resterr.NotFound(gorm.ErrRecordNotFound.Error())
	internalErrorResponse := resterr.InternalServerError("fail to get owner")

	testCases := []struct {
		name             string
		id               uint
		mockOwner        *Response
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get owner with pets by ID",
			id:               1,
			mockOwner:        ownerWithPetsResponse,
			mockError:        nil,
			route:            "/v1/owners/1/pets",
			statusCode:       http.StatusOK,
			expectedResponse: ownerWithPetsResponse,
		},
		{
			name:             "get owner with no pets by ID",
			id:               2,
			mockOwner:        ownerResponse,
			mockError:        nil,
			route:            "/v1/owners/2/pets",
			statusCode:       http.StatusOK,
			expectedResponse: ownerResponse,
		},
		{
			name:             "get no owner",
			id:               5,
			mockOwner:        nil,
			mockError:        gorm.ErrRecordNotFound,
			route:            "/v1/owners/5/pets",
			statusCode:       http.StatusNotFound,
			expectedResponse: notFoundResponse,
		},
		{
			name:             "fail to get owner by id",
			id:               6,
			mockOwner:        nil,
			mockError:        resterr.InternalServerError("fail to get owner"),
			route:            "/v1/owners/6/pets",
			statusCode:       http.StatusNotFound,
			expectedResponse: internalErrorResponse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := NewMockServicer(t)
			ownerMock.EXPECT().getOwnerByIdWithPets(tc.id).Return(tc.mockOwner, tc.mockError)

			r := test.SetupRouter()
			v1 := r.Group("/v1")
			ownerRouter := NewRouter(ownerMock)
			ownerRouter.Register(v1.Group("/owners"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)

			switch tc.statusCode {
			case http.StatusOK:
				actualOwnerResponse := &Response{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualOwnerResponse)

				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*Response), actualOwnerResponse)
				assert.Equal(t, tc.expectedResponse.(*Response).ID, actualOwnerResponse.ID)
				assert.Equal(t, tc.expectedResponse.(*Response).City, actualOwnerResponse.City)
				assert.Equal(t, tc.expectedResponse.(*Response).Pets, actualOwnerResponse.Pets)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualOwnerResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualOwnerResponse)

				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}

				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*resterr.ErrorResponse), actualOwnerResponse)
			}
		})
	}
}

func Test_ownerByIdWithPetsBadID(t *testing.T) {
	response := resterr.BadRequest("failed to convert: strconv.Atoi: parsing \"notID\": invalid syntax")
	ownerMock := NewMockServicer(t)

	r := test.SetupRouter()
	v1 := r.Group("/v1")
	ownerRouter := NewRouter(ownerMock)
	ownerRouter.Register(v1.Group("/owners"))

	req := httptest.NewRequest("GET", "/v1/owners/notID/pets", nil)
	resp, _ := r.Test(req, 5)

	actualOwnerResponse := &resterr.ErrorResponse{}
	// Read the response body
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal(body, actualOwnerResponse)

	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, &response, actualOwnerResponse)
}

func Test_OwnerByLastName(t *testing.T) {
	ownersResponse := &Responses{
		Context: model.Context{Count: 2},
		Owners: []Response{
			{
				ID:        15,
				FirstName: "Charles",
				LastName:  "Ward",
			},
			{
				ID:        11,
				FirstName: "John",
				LastName:  "Ward",
			},
		},
	}

	noOwnersFoundResponse := &Responses{
		Context: model.Context{Count: 0},
		Owners:  []Response{},
	}

	errorResponse := resterr.InternalServerError("unexpected error occurred")

	testCases := []struct {
		name             string
		lastName         string
		mockOwners       *Responses
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get owners by last name",
			lastName:         "Ward",
			mockOwners:       ownersResponse,
			mockError:        nil,
			route:            "/v1/owners?last-name=Ward",
			statusCode:       http.StatusOK,
			expectedResponse: ownersResponse,
		},
		{
			name:             "get no owners by last name",
			lastName:         "Jackson",
			mockOwners:       noOwnersFoundResponse,
			mockError:        nil,
			route:            "/v1/owners?last-name=Jackson",
			statusCode:       http.StatusNotFound,
			expectedResponse: noOwnersFoundResponse,
		},
		{
			name:             "Test get owner by last name with error",
			lastName:         "Ward",
			route:            "/v1/owners?last-name=Ward",
			mockOwners:       nil,
			mockError:        errors.New("unexpected error occurred"),
			statusCode:       http.StatusInternalServerError,
			expectedResponse: &errorResponse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := NewMockServicer(t)
			ownerMock.EXPECT().getOwnerByLastName(tc.lastName).Return(tc.mockOwners, tc.mockError)

			r := test.SetupRouter()
			v1 := r.Group("/v1")
			ownerRouter := NewRouter(ownerMock)
			ownerRouter.Register(v1.Group("/owners"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)

			switch tc.statusCode {
			case http.StatusOK:
				actualOwnersResponse := &Responses{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualOwnersResponse)

				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*Responses), actualOwnersResponse)
				assert.Equal(t, tc.expectedResponse.(*Responses).Context.Count, actualOwnersResponse.Context.Count)
			case http.StatusNotFound | http.StatusInternalServerError:
				actualOwnerResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualOwnerResponse)

				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}

				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*resterr.ErrorResponse), actualOwnerResponse)
			}
		})
	}
}

func Test_addNewOwner(t *testing.T) {
	ownerRequest := &AddRequest{
		FirstName: "Nat",
		LastName:  "Cole",
		Address:   "1234 Elm St",
		City:      "New York",
		Telephone: "1234567890",
	}

	owner := &Response{
		ID:        1,
		FirstName: "Nat",
		LastName:  "Cole",
		Address:   "1234 Elm St",
		City:      "New York",
		Telephone: "1234567890",
	}

	internalError := resterr.InternalServerError("fail to add a new owner")
	testCases := []struct {
		name             string
		request          *AddRequest
		mockOwner        *Response
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "add a new owner",
			request:          ownerRequest,
			mockOwner:        owner,
			mockError:        nil,
			route:            "/v1/owners",
			statusCode:       http.StatusCreated,
			expectedResponse: owner,
		},
		{
			name:             "fail to create a new owner",
			request:          ownerRequest,
			mockOwner:        nil,
			mockError:        errors.New("fail to add a new owner"),
			route:            "/v1/owners",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: &internalError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := NewMockServicer(t)
			ownerMock.EXPECT().create(tc.request).Return(tc.mockOwner, tc.mockError)

			r := test.SetupRouter()
			v1 := r.Group("/v1")
			ownerRouter := NewRouter(ownerMock)
			ownerRouter.Register(v1.Group("/owners"))

			requestBody, _ := json.Marshal(tc.request)

			req := httptest.NewRequest("POST", tc.route, bytes.NewReader(requestBody))
			// Must se content-type header.
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, 10)

			switch resp.StatusCode {
			case http.StatusCreated:
				actualOwnerResponse := &Response{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualOwnerResponse)

				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*Response), actualOwnerResponse)
			case http.StatusInternalServerError:
				actualOwnerResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualOwnerResponse)

				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}

				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*resterr.ErrorResponse), actualOwnerResponse)
			}
		})
	}
}

func Test_addNewOwnerWithBadRequest(t *testing.T) {
	r := test.SetupRouter()
	v1 := r.Group("/v1")

	ownerMock := NewMockServicer(t)
	ownerRouter := NewRouter(ownerMock)
	ownerRouter.Register(v1.Group("/owners"))

	// assign number to lastName to purposely fail JSON unmarshalling
	addRequest := map[string]interface{}{
		"firstName": "James",
		"lastName":  5,
		"address":   "123 Main St.",
	}
	addRequestJSON, _ := json.Marshal(addRequest)

	req := httptest.NewRequest("POST", "/v1/owners", bytes.NewReader(addRequestJSON))
	// Must se content-type header.
	req.Header.Set("Content-Type", "application/json")
	resp, _ := r.Test(req, 10)

	actualOwnerResponse := &resterr.ErrorResponse{}
	// Read the response body
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal(body, actualOwnerResponse)

	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
		return
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

func Test_UpdateOwner(t *testing.T) {
	ownerRequest := &UpdateRequest{
		FirstName: "Nat",
		LastName:  "Cole",
		Address:   "1234 Elm St",
		City:      "New York",
		Telephone: "1234567890",
	}

	owner := &UpdateResponse{
		ID:        15,
		FirstName: "Nat",
		LastName:  "Cole",
		Address:   "1234 Elm St",
		City:      "New York",
		Telephone: "1234567890",
	}

	internalError := resterr.InternalServerError("update: unable to update owner")

	testCases := []struct {
		name             string
		id               uint
		request          *UpdateRequest
		mockOwner        *UpdateResponse
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "update an owner",
			id:               15,
			request:          ownerRequest,
			mockOwner:        owner,
			mockError:        nil,
			route:            "/v1/owners/15",
			statusCode:       http.StatusOK,
			expectedResponse: owner,
		},
		{
			name:             "fail to update owner",
			id:               39,
			request:          ownerRequest,
			mockOwner:        nil,
			mockError:        errors.New("update: unable to update owner"),
			route:            "/v1/owners/39",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: &internalError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := NewMockServicer(t)
			ownerMock.EXPECT().update(tc.id, tc.request).Return(tc.mockOwner, tc.mockError)

			r := test.SetupRouter()
			v1 := r.Group("/v1")
			ownerRouter := NewRouter(ownerMock)
			ownerRouter.Register(v1.Group("/owners"))

			requestBody, _ := json.Marshal(tc.request)
			req := httptest.NewRequest("PUT", tc.route, bytes.NewReader(requestBody))
			// Must se content-type header.
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, 5)

			switch resp.StatusCode {
			case http.StatusCreated:
				actualOwnerResponse := &UpdateResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualOwnerResponse)

				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*UpdateResponse), actualOwnerResponse)
			case http.StatusInternalServerError:
				actualOwnerResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualOwnerResponse)

				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}

				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*resterr.ErrorResponse), actualOwnerResponse)
			}

		})
	}
}

func Test_updateOwnerWithBadRequest(t *testing.T) {
	ownerRequest := &UpdateRequest{
		FirstName: "Nat",
		LastName:  "Cole",
		Address:   "1234 Elm St",
		City:      "New York",
		Telephone: "1234567890",
	}
	// assign number to lastName to purposely fail JSON unmarshalling
	updateRequest := map[string]interface{}{
		"firstName": "James",
		"lastName":  5,
		"address":   "123 Main St.",
	}
	updateRequestJSON, _ := json.Marshal(updateRequest)

	badResponse := resterr.BadRequest("failed to convert: strconv.Atoi: parsing \"a1\": invalid syntax")
	badResponse2 := resterr.BadRequest("json: cannot unmarshal number into Go struct field UpdateRequest.lastName of type string")

	testCases := []struct {
		name             string
		id               interface{}
		request          interface{}
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "update owner request with invalid id",
			id:               "a1",
			request:          ownerRequest,
			route:            "/v1/owners/a1",
			statusCode:       http.StatusBadRequest,
			expectedResponse: &badResponse,
		},
		{
			name:             "bad update owner request",
			id:               int(17),
			request:          updateRequestJSON,
			route:            "/v1/owners/17",
			statusCode:       http.StatusBadRequest,
			expectedResponse: &badResponse2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := NewMockServicer(t)

			r := test.SetupRouter()
			v1 := r.Group("/v1")
			ownerRouter := NewRouter(ownerMock)
			ownerRouter.Register(v1.Group("/owners"))
			var req *http.Request

			switch tc.request.(type) {
			case []byte:
				req = httptest.NewRequest("PUT", tc.route, bytes.NewReader(tc.request.([]byte)))
			default:
				requestBody, _ := json.Marshal(tc.request)
				req = httptest.NewRequest("PUT", tc.route, bytes.NewReader(requestBody))
			}

			// Must se content-type header.
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, 5)
			actualOwnerResponse := &resterr.ErrorResponse{}
			// Read the response body
			body, _ := io.ReadAll(resp.Body)
			err := json.Unmarshal(body, actualOwnerResponse)

			if err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
				return
			}

			assert.Equal(t, tc.statusCode, resp.StatusCode)
			assert.Equal(t, tc.expectedResponse.(*resterr.ErrorResponse), actualOwnerResponse)
		})
	}
}
