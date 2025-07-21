package owner

import (
	"errors"
	resterr "fiber-petclinic-service/Pkg/errors"
	"fiber-petclinic-service/api/pet"
	"fiber-petclinic-service/pkg/repository/model"
	"fiber-petclinic-service/pkg/test"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

// config the gin engine for testing purpose
func setupRouter() *fiber.App {
	gin.SetMode(gin.TestMode)
	r := fiber.New()
	return r
}

func Test_ownerById(t *testing.T) {
	log.Info().Str("function", "Test_ownerById").Msg("Owner by ID endpoint")

	ownerResponse := &Response{
		ID:        1,
		FirstName: "Nat",
		LastName:  "Cole",
		Address:   "1234 Elm St",
		City:      "New York",
		Telephone: "1234567890",
	}

	testCases := []struct {
		url           string
		name          string
		expectedOwner *Response
		expectedError error
		status        int
		jsonResponse  string
	}{
		{
			url:           "/v1/owners/1",
			name:          "Test getting owner by ID",
			expectedOwner: ownerResponse,
			expectedError: nil,
			status:        http.StatusOK,
			jsonResponse:  test.JsonString(ownerResponse),
		},
		{
			url:           "/v1/owners/1",
			name:          "Test finding no owner",
			expectedOwner: nil,
			expectedError: gorm.ErrRecordNotFound,
			status:        http.StatusNotFound,
			jsonResponse:  test.JsonString(resterr.NotFound(gorm.ErrRecordNotFound.Error())),
		},
		{
			url:           "/v1/owners/1",
			name:          "Test get owner by ID with error",
			expectedOwner: nil,
			expectedError: errors.New("unexpected error occurred"),
			status:        http.StatusInternalServerError,
			jsonResponse:  test.JsonString(resterr.InternalServerError("unexpected error occurred")),
		},
		{
			url:           "/v1/owners/a1",
			name:          "Test get owner by ID with invalid ID",
			expectedOwner: nil,
			expectedError: strconv.ErrSyntax,
			status:        http.StatusBadRequest,
			jsonResponse:  test.JsonString(resterr.BadRequest(errors.New("strconv.Atoi: parsing \"a1\": invalid syntax"))),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := MockServicer{}
			ownerMock.On("getOwnerById", uint(1)).Return(tc.expectedOwner, tc.expectedError)
			ownerRouter := NewRouter(&ownerMock)

			r := setupRouter()
			v1 := r.Group("/v1")
			ownerRouter.Register(v1.Group("/owners"))

			tests := []test.APITestCase{
				{"Get Owner By ID", "GET", tc.url, "", nil, tc.status, tc.jsonResponse},
			}
			for _, tc := range tests {
				test.Endpoint(t, r, tc)
			}
		})
	}
}

func Test_ownerByIdWithPets(t *testing.T) {
	logger, _ := zap.NewProduction()
	logger.Info("Owner by ID with Pets endpoint", zap.String("function", "Test_ownerByIdWithPets"))

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
				Type:      "Dog",
			},
			{
				ID:        2,
				Name:      "Lucky",
				Birthdate: "2015-09-07",
				Type:      "Dog",
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

	_, convErr := strconv.Atoi("a1")

	testCases := []struct {
		url           string
		pathParam     string
		name          string
		expectedOwner *Response
		expectedError error
		status        int
		jsonResponse  string
	}{
		{
			url:           "/v1/owners/1/pets",
			pathParam:     "1",
			name:          "Test getting owner with pets by ID",
			expectedOwner: ownerWithPetsResponse,
			expectedError: nil,
			status:        http.StatusOK,
			jsonResponse:  test.JsonString(ownerWithPetsResponse),
		},
		{
			url:           "/v1/owners/2/pets",
			pathParam:     "2",
			name:          "Test getting owner with no pets by ID",
			expectedOwner: ownerResponse,
			expectedError: nil,
			status:        http.StatusOK,
			jsonResponse:  test.JsonString(ownerResponse),
		},
		{
			url:           "/v1/owners/5/pets",
			pathParam:     "5",
			name:          "Test finding no owner",
			expectedOwner: nil,
			expectedError: gorm.ErrRecordNotFound,
			status:        http.StatusNotFound,
			jsonResponse:  test.JsonString(resterr.NotFound(gorm.ErrRecordNotFound.Error())),
		},
		{
			url:           "/v1/owners/6/pets",
			pathParam:     "6",
			name:          "Test get owner by ID with error",
			expectedOwner: nil,
			expectedError: errors.New("unexpected error occurred"),
			status:        http.StatusInternalServerError,
			jsonResponse:  test.JsonString(resterr.InternalServerError("unexpected error occurred")),
		},
		{
			url:           "/v1/owners/a1/pets",
			pathParam:     "a1",
			name:          "Test get owner by ID with invalid ID",
			expectedOwner: nil,
			expectedError: convErr,
			status:        http.StatusBadRequest,
			jsonResponse:  test.JsonString(resterr.BadRequest(convErr.Error())),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := MockServicer{}
			id, _ := strconv.Atoi(tc.pathParam)
			ownerMock.On("getOwnerByIdWithPets", uint(id)).Return(tc.expectedOwner, tc.expectedError)
			ownerRouter := NewRouter(logger, &ownerMock)

			r := setupRouter()
			v1 := r.Group("/v1")
			ownerRouter.Register(v1.Group("/owners"))

			tests := []test.APITestCase{
				{"Get Owner with Pets By ID", "GET", tc.url, "", nil, tc.status, tc.jsonResponse},
			}
			for _, tc := range tests {
				test.Endpoint(t, r, tc)
			}
		})
	}
}

func Test_OwnerByLastName(t *testing.T) {
	logger, _ := zap.NewProduction()
	logger.Info("Owner by Last Name endpoint", zap.String("function", "Test_OwnerByLastName"))

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

	testCases := []struct {
		name           string
		queryParam     string
		url            string
		expectedOwners *Responses
		expectedError  error
		status         int
		jsonResponse   string
	}{
		{
			name:           "Test getting owners by last name",
			queryParam:     "Ward",
			url:            "/v1/owners?last-name=Ward",
			expectedOwners: ownersResponse,
			expectedError:  nil,
			status:         http.StatusOK,
			jsonResponse:   test.JsonString(ownersResponse),
		},
		{
			name:           "Test finding no owner",
			queryParam:     "Jackson",
			url:            "/v1/owners?last-name=Jackson",
			expectedOwners: noOwnersFoundResponse,
			expectedError:  nil,
			status:         http.StatusNotFound,
			jsonResponse:   test.JsonString(resterr.NotFound("Find no owner with last name: Jackson")),
		},
		{
			name:           "Test get owner by last name with error",
			queryParam:     "Ward",
			url:            "/v1/owners?last-name=Ward",
			expectedOwners: nil,
			expectedError:  errors.New("unexpected error occurred"),
			status:         http.StatusInternalServerError,
			jsonResponse:   test.JsonString(resterr.InternalServerError("unexpected error occurred")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := MockServicer{}
			ownerMock.On("getOwnerByLastName", tc.queryParam).Return(tc.expectedOwners, tc.expectedError)
			ownerRouter := NewRouter(&ownerMock)

			r := setupRouter()
			v1 := r.Group("/v1")
			ownerRouter.Register(v1.Group("/owners"))

			tests := []test.APITestCase{
				{"Get Owner By Last Name", "GET", tc.url, "", nil, tc.status, tc.jsonResponse},
			}
			for _, tc := range tests {
				test.Endpoint(t, r, tc)
			}
		})
	}
}

func Test_AllOwners(t *testing.T) {
	logger, _ := zap.NewProduction()
	logger.Info("All owner endpoint", zap.String("function", "TestAllOwners"))

	ownersResponse := &Responses{
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

	testCases := []struct {
		name           string
		expectedOwners *Responses
		expectedError  error
		status         int
		jsonResponse   string
	}{
		{
			name:           "Test getting all owners",
			expectedOwners: ownersResponse,
			expectedError:  nil,
			status:         http.StatusOK,
			jsonResponse:   test.JsonString(ownersResponse),
		},
		{
			name:           "Test finding no owner",
			expectedOwners: noOwnersFoundResponse,
			expectedError:  nil,
			status:         http.StatusNotFound,
			jsonResponse:   test.JsonString(resterr.NotFound("Find no owner")),
		},
		{
			name:           "Test get all owners with error",
			expectedOwners: nil,
			expectedError:  errors.New("unexpected error occurred"),
			status:         http.StatusInternalServerError,
			jsonResponse:   test.JsonString(resterr.InternalServerError("unexpected error occurred")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := MockServicer{}
			ownerMock.On("getAllOwners").Return(tc.expectedOwners, tc.expectedError)
			ownerRouter := NewRouter(logger, &ownerMock)

			r := setupRouter()
			v1 := r.Group("/v1")
			ownerRouter.Register(v1.Group("/owners"))

			tests := []test.APITestCase{
				{"Get All Owners", "GET", "/v1/owners/all", "", nil, tc.status, tc.jsonResponse},
			}
			for _, tc := range tests {
				test.Endpoint(t, r, tc)
			}
		})
	}
}

func Test_CreateOwner(t *testing.T) {
	logger, _ := zap.NewProduction()
	logger.Info("Create Owner endpoint", zap.String("function", "Test_CreateOwner"))
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

	testCases := []struct {
		name          string
		request       *AddRequest
		expectedOwner *Response
		expectedError error
		status        int
		jsonResponse  string
	}{
		{
			name:          "Test create owner",
			request:       ownerRequest,
			expectedOwner: owner,
			expectedError: nil,
			status:        http.StatusCreated,
			jsonResponse:  test.JsonString(owner),
		},
		{
			name:          "Test create owner with error",
			request:       ownerRequest,
			expectedOwner: nil,
			expectedError: errors.New("unexpected error occurred"),
			status:        http.StatusInternalServerError,
			jsonResponse:  test.JsonString(resterr.InternalServerError("unexpected error occurred")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := MockServicer{}
			ownerMock.On("create", tc.request).Return(tc.expectedOwner, tc.expectedError)
			ownerRouter := NewRouter(&ownerMock)

			r := setupRouter()
			v1 := r.Group("/v1")
			ownerRouter.Register(v1.Group("/owners"))

			tests := []test.APITestCase{
				{"Create Owner", "POST", "/v1/owners", test.JsonString(tc.request), nil, tc.status, tc.jsonResponse},
			}
			for _, tc := range tests {
				test.Endpoint(t, r, tc)
			}
		})
	}
}

func Test_UpdateOwner(t *testing.T) {
	logger, _ := zap.NewProduction()
	logger.Info("Update Owner endpoint", zap.String("function", "Test_UpdateOwner"))

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

	testCases := []struct {
		name          string
		requestID     string
		request       *UpdateRequest
		url           string
		expectedOwner *UpdateResponse
		expectedError error
		status        int
		jsonResponse  string
	}{
		{
			name:          "Test update owner",
			requestID:     "15",
			request:       ownerRequest,
			url:           "/v1/owners/15",
			expectedOwner: owner,
			expectedError: nil,
			status:        http.StatusOK,
			jsonResponse:  test.JsonString(owner),
		},
		{
			name:          "Test update owner",
			requestID:     "a1",
			request:       ownerRequest,
			url:           "/v1/owners/a1",
			expectedOwner: nil,
			expectedError: strconv.ErrSyntax,
			status:        http.StatusBadRequest,
			jsonResponse:  test.JsonString(resterr.BadRequest(errors.New("strconv.Atoi: parsing \"a1\": invalid syntax"))),
		},
		{
			name:          "Test update owner with error",
			requestID:     "17",
			request:       ownerRequest,
			url:           "/v1/owners/17",
			expectedOwner: nil,
			expectedError: errors.New("unexpected error occurred"),
			status:        http.StatusInternalServerError,
			jsonResponse:  test.JsonString(resterr.InternalServerError("unexpected error occurred")),
		},
		{
			name:          "Test update owner fail",
			requestID:     "39",
			request:       ownerRequest,
			url:           "/v1/owners/39",
			expectedOwner: nil,
			expectedError: errors.New("update: unable to update owner"),
			status:        http.StatusInternalServerError,
			jsonResponse:  test.JsonString(resterr.InternalServerError("update: unable to update owner")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := MockServicer{}
			id, _ := strconv.Atoi(tc.requestID)
			ownerMock.On("update", uint(id), tc.request).Return(tc.expectedOwner, tc.expectedError)
			ownerRouter := NewRouter(&ownerMock)

			r := setupRouter()
			v1 := r.Group("/v1")
			ownerRouter.Register(v1.Group("/owners"))

			tests := []test.APITestCase{
				{"Update Owner", "PUT", tc.url, test.JsonString(tc.request), nil, tc.status, tc.jsonResponse},
			}
			for _, tc := range tests {
				test.Endpoint(t, r, tc)
			}
		})
	}
}
