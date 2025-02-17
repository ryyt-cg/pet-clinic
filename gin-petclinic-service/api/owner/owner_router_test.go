package owner

import (
	"encoding/json"
	"github.com/rhtran/gin-petclinic-service/api/pet"
	"github.com/rhtran/gin-petclinic-service/pkg/infra/repository/test"
	"github.com/rhtran/gin-petclinic-service/pkg/model"
	"go.uber.org/zap"
	"net/http"
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
	//logger := log.New().With(nil, "function", "Test_ownerById")
	logger, _ := zap.NewProduction()
	ownerMock := MockServicer{}

	ownerResponse := &Response{
		ID:        1,
		FirstName: "Nat",
		LastName:  "Cole",
	}

	ownerMock.On("getOwnerById", uint(1)).Return(ownerResponse, nil)
	ownerRouter := NewRouter(logger, &ownerMock)

	r := setupRouter()
	v1 := r.Group("/v1")
	ownerRouter.Register(v1.Group("/owners"))

	tc1, _ := json.Marshal(ownerResponse)

	tests := []test.APITestCase{
		{"Get Owner By ID", "GET", "/v1/owners/1", "", nil, http.StatusOK, string(tc1)},
	}
	for _, tc := range tests {
		test.Endpoint(t, r, tc)
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

	tc1, _ := json.Marshal(ownerResponse)

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

	tc1, _ := json.Marshal(ownersResponse)

	tests := []test.APITestCase{
		{"Get Owner By Last Name", "GET", "/v1/owners?last-name=Ward", "", nil, http.StatusOK, string(tc1)},
	}
	for _, tc := range tests {
		test.Endpoint(t, r, tc)
	}
}

func Test_AllOwners(t *testing.T) {
	logger, _ := zap.NewProduction()
	//logger := log.New().With(nil, "function", "TestAllOwners")
	ownerMock := MockServicer{}

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

	ownerMock.On("getAllOwners").Return(ownersResponse, nil)
	ownerRouter := NewRouter(logger, &ownerMock)

	r := setupRouter()
	v1 := r.Group("/v1")
	ownerRouter.Register(v1.Group("/owners"))

	tc1, _ := json.Marshal(ownersResponse)

	tests := []test.APITestCase{
		{"Get All Owners", "GET", "/v1/owners/all", "", nil, http.StatusOK, string(tc1)},
	}
	for _, tc := range tests {
		test.Endpoint(t, r, tc)
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

	tc1, _ := json.Marshal(owner)

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

	requestBoday, _ := json.Marshal(ownerRequest)

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

	tc1, _ := json.Marshal(owner)
	tests := []test.APITestCase{
		{"Update Owner", "PUT", "/v1/owners/1", string(requestBoday), nil, http.StatusOK, string(tc1)},
	}

	for _, tc := range tests {
		test.Endpoint(t, r, tc)
	}
}
