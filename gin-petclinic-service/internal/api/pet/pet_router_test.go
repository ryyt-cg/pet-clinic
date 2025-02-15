package pet

import (
	"encoding/json"
	"github.com/rhtran/gin-petclinic-service/pkg/infra/repository"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"

	"github.com/stretchr/testify/mock"
)

// petServiceMock
// require to mock all functions; otherwise, compilation errors.
type petServiceMock struct {
	mock.Mock
}

func (petM *petServiceMock) getAllPets() (*Responses, error) {
	args := petM.Called()
	intf := args.Get(0)
	val := intf.(Responses)
	ptr := &val

	return ptr, args.Error(1)
}

func (petM *petServiceMock) getPetById(id int) (*Response, error) {
	args := petM.Called(id)
	intf := args.Get(0)
	val := intf.(Response)
	ptr := &val

	return ptr, args.Error(1)
}

func (petM *petServiceMock) getPetWithVisitsById(id int) (*Response, error) {
	args := petM.Called(id)
	intf := args.Get(0)
	val := intf.(Response)
	ptr := &val

	return ptr, args.Error(1)
}

func (petM *petServiceMock) getPetsByName(name string) ([]Response, error) {
	args := petM.Called(name)
	intf := args.Get(0)
	val := intf.([]Response)
	ptr := val

	return ptr, args.Error(1)
}

func (petM *petServiceMock) create(pet *repository.Pet) (*Response, error) {
	args := petM.Called(pet)
	intf := args.Get(0)
	val := intf.(*Response)
	ptr := val

	return ptr, args.Error(1)
}

func (petM *petServiceMock) update(pet *repository.Pet) (*Response, error) {
	args := petM.Called(pet)
	intf := args.Get(0)
	val := intf.(*Response)
	ptr := val

	return ptr, args.Error(1)
}

// config the gin engine for testing purpose
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func Test_PetById(t *testing.T) {
	//logger := log.New().With(nil, "function", "Test_PetById")
	logger := zap.NewNop()
	petMock := petServiceMock{}
	petResponse := &Response{
		ID:   1,
		Name: "Nash",
	}
	petMock.On("getPetById", 1).Return(*petResponse, nil)
	petRouter := NewRouter(logger, &petMock)

	r := setupRouter()
	v1 := r.Group("/v1")
	petRouter.Register(v1.Group("/pets"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/pets/1", nil)
	r.ServeHTTP(w, req)

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)

	// unmarshal to Pet struct for asserts.
	actualPetResponse := &Response{}
	json.Unmarshal(w.Body.Bytes(), actualPetResponse)
	assert.Equal(t, petResponse.ID, actualPetResponse.ID)
	assert.Equal(t, petResponse.Name, actualPetResponse.Name)
}

func Test_PetsByName(t *testing.T) {
	//logger := log.New().With(nil, "function", "Test_PetByName")
	logger := zap.NewNop()
	petMock := petServiceMock{}

	petResponses := make([]Response, 1)
	petRespopnse := &Response{
		ID:   15,
		Name: "Charles",
	}
	petResponses[0] = *petRespopnse

	petMock.On("GetPetByName", "Charles").Return(petResponses, nil)
	petRouter := NewRouter(logger, &petMock)

	r := setupRouter()
	v1 := r.Group("/v1")
	petRouter.Register(v1.Group("/pets"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/pets?name=Charles", nil)
	r.ServeHTTP(w, req)

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)

	// unmarshal to Pet struct for asserts.
	actualPetResponses := make([]Response, 1)
	json.Unmarshal(w.Body.Bytes(), &actualPetResponses)
	assert.Equal(t, petRespopnse.ID, actualPetResponses[0].ID)
	assert.Equal(t, petRespopnse.Name, actualPetResponses[0].Name)
}
