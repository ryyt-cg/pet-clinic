package owner

import (
	"encoding/json"
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
	logger := zap.NewNop()
	ownerMock := MockServicer{}

	ownerResponse := &Response{
		ID:        1,
		FirstName: "Nat",
		LastName:  "Cole",
	}

	ownerMock.On("getOwnerById", 1).Return(ownerResponse, nil)
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

func Test_OwnerByLastName(t *testing.T) {
	//logger := log.New().With(nil, "function", "Test_OwnerByLastName")
	logger := zap.NewNop()
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

func Test_GetAllOwners(t *testing.T) {
	logger := zap.NewNop()
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
