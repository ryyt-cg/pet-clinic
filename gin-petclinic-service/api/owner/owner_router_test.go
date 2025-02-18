package owner

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/rhtran/gin-petclinic-service/api/pet"
	resterr "github.com/rhtran/gin-petclinic-service/middleware/errors"
	"github.com/rhtran/gin-petclinic-service/pkg/model"
	"github.com/rhtran/gin-petclinic-service/pkg/test"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

// config the gin engine for testing purpose
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func Test_ownerById(t *testing.T) {
	logger, _ := zap.NewProduction()
	logger.Info("Owner by ID endpoint", zap.String("function", "Test_ownerById"))

	ownerResponse := &Response{
		ID:        1,
		FirstName: "Nat",
		LastName:  "Cole",
		Address:   "1234 Elm St",
		City:      "New York",
		Telephone: "1234567890",
	}

	_, convErr := strconv.Atoi("a1")

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
			expectedError: convErr,
			status:        http.StatusBadRequest,
			jsonResponse:  test.JsonString(resterr.BadRequestWithDetails(convErr)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ownerMock := MockServicer{}
			ownerMock.On("getOwnerById", uint(1)).Return(tc.expectedOwner, tc.expectedError)
			ownerRouter := NewRouter(logger, &ownerMock)

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
	//logger := log.New().With(nil, "function", "Test_ownerById")
	logger, _ := zap.NewProduction()
	ownerMock := MockServicer{}

	ownerResponse := &Response{
		ID:        1,
		FirstName: "Nat",
		LastName:  "Cole",
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

	ownerMock.On("getOwnerByIdWithPets", uint(1)).Return(ownerResponse, nil)
	ownerRouter := NewRouter(logger, &ownerMock)

	r := setupRouter()
	v1 := r.Group("/v1")
	ownerRouter.Register(v1.Group("/owners"))

	tc1, _ := jsoniter.Marshal(ownerResponse)

	tests := []test.APITestCase{
		{"Get Owner with Pets By ID ", "GET", "/v1/owners/1/pets", "", nil, http.StatusOK, string(tc1)},
	}
	for _, tc := range tests {
		test.Endpoint(t, r, tc)
	}
}

func Test_OwnerByLastName(t *testing.T) {
	//logger := log.New().With(nil, "function", "Test_OwnerByLastName")
	logger, _ := zap.NewProduction()
	ownerMock := MockServicer{}

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

	ownerMock.On("getOwnerByLastName", "Ward").Return(ownersResponse, nil)
	ownerRouter := NewRouter(logger, &ownerMock)

	r := setupRouter()
	v1 := r.Group("/v1")
	ownerRouter.Register(v1.Group("/owners"))

	tc1, _ := jsoniter.Marshal(ownersResponse)

	tests := []test.APITestCase{
		{"Get Owner By Last Name", "GET", "/v1/owners?last-name=Ward", "", nil, http.StatusOK, string(tc1)},
	}
	for _, tc := range tests {
		test.Endpoint(t, r, tc)
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
	//logger := log.New().With(nil, "function", "Test_CreateOwner")
	ownerMock := MockServicer{}
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

	ownerMock.On("create", ownerRequest).Return(owner, nil)
	ownerRouter := NewRouter(logger, &ownerMock)

	r := setupRouter()
	v1 := r.Group("/v1")
	ownerRouter.Register(v1.Group("/owners"))

	tc1, _ := jsoniter.Marshal(owner)

	tests := []test.APITestCase{
		{"Create Owner", "POST", "/v1/owners", `{"id":1,"firstName":"Nat","lastName":"Cole","address":"1234 Elm St","city":"New York","telephone":"1234567890"}`, nil, http.StatusCreated, string(tc1)},
	}

	for _, tc := range tests {
		test.Endpoint(t, r, tc)
	}
}

func Test_UpdateOwner(t *testing.T) {
	logger, _ := zap.NewProduction()
	//logger := log.New().With(nil, "function", "Test_UpdateOwner")
	ownerMock := MockServicer{}

	requestID := uint(1)
	ownerRequest := &UpdateRequest{
		FirstName: "Nat",
		LastName:  "Cole",
		Address:   "1234 Elm St",
		City:      "New York",
		Telephone: "1234567890",
	}

	requestBoday, _ := jsoniter.Marshal(ownerRequest)

	owner := &UpdateResponse{
		ID:        1,
		FirstName: "Nat",
		LastName:  "Cole",
		Address:   "1234 Elm St",
		City:      "New York",
		Telephone: "1234567890",
	}

	ownerMock.On("update", requestID, ownerRequest).Return(owner, nil)
	ownerRouter := NewRouter(logger, &ownerMock)

	r := setupRouter()
	v1 := r.Group("/v1")
	ownerRouter.Register(v1.Group("/owners"))

	tc1, _ := jsoniter.Marshal(owner)
	tests := []test.APITestCase{
		{"Update Owner", "PUT", "/v1/owners/1", string(requestBoday), nil, http.StatusOK, string(tc1)},
	}

	for _, tc := range tests {
		test.Endpoint(t, r, tc)
	}
}
