package pet

import (
	resterr "fiber-petclinic-service/pkg/errors"
	"fiber-petclinic-service/pkg/repository"
	"fiber-petclinic-service/pkg/test"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// petServiceMock
// require to mock all functions; otherwise, compilation errors.
type petServiceMock struct {
	mock.Mock
}

func (petM *petServiceMock) getAllPets() (*Responses, error) {
	args := petM.Called()
	intf := args.Get(0)

	if intf == nil {
		return nil, args.Error(1)
	}

	val := intf.(Responses)
	return &val, args.Error(1)
}

func (petM *petServiceMock) getPetById(id int) (*Response, error) {
	args := petM.Called(id)
	intf := args.Get(0)

	if intf == nil {
		return nil, args.Error(1)
	}

	val := intf.(Response)
	return &val, args.Error(1)
}

func (petM *petServiceMock) getPetWithVisitsById(id int) (*Response, error) {
	args := petM.Called(id)
	intf := args.Get(0)

	if intf == nil {
		return nil, args.Error(1)
	}

	val := intf.(Response)
	ptr := &val
	return ptr, args.Error(1)
}

func (petM *petServiceMock) getPetsByName(name string) ([]Response, error) {
	args := petM.Called(name)
	intf := args.Get(0)

	if intf == nil {
		return nil, args.Error(1)
	}

	val := intf.([]Response)
	return val, args.Error(1)
}

func (petM *petServiceMock) create(pet *repository.Pet) (*Response, error) {
	args := petM.Called(pet)
	intf := args.Get(0)

	if intf == nil {
		return nil, args.Error(1)
	}

	val := intf.(*Response)
	return val, args.Error(1)
}

func (petM *petServiceMock) update(pet *repository.Pet) (*Response, error) {
	args := petM.Called(pet)
	intf := args.Get(0)

	if intf == nil {
		return nil, args.Error(1)
	}

	val := intf.(*Response)
	return val, args.Error(1)
}

func Test_PetById(t *testing.T) {
	tests := []struct {
		description      string
		id               string
		route            string
		expectedError    bool
		error            interface{}
		expectedCode     int
		expectedResponse interface{}
	}{
		{
			description:      "Invalid pet id",
			route:            "/v1/pets/4js",
			expectedError:    true,
			expectedCode:     http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("Invalid pet ID"),
		},
		{
			description:   "Found pet by id",
			id:            "1",
			route:         "/v1/pets/1",
			expectedError: false,
			expectedCode:  http.StatusOK,
			expectedResponse: Response{
				ID:   1,
				Name: "Nash",
			},
		},
		{
			description:      "Found no pet by id",
			route:            "/v1/pets/1",
			expectedError:    true,
			error:            gorm.ErrRecordNotFound,
			expectedCode:     http.StatusNotFound,
			expectedResponse: resterr.NotFound("Pet not found"),
		},
		{
			description:      "Internal Error",
			route:            "/v1/pets/1",
			expectedError:    true,
			error:            error(fiber.ErrInternalServerError),
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: resterr.InternalServerError("Something went wrong"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			petMock := petServiceMock{}
			if tc.expectedError {
				if tc.expectedCode == http.StatusBadRequest {
					// Do not expect petMock return anything
				} else {
					petMock.On("getPetById", 1).Return(nil, tc.error)
				}
			} else {
				petMock.On("getPetById", 1).Return(tc.expectedResponse, nil)
			}

			petRouter := NewRouter(&petMock)
			petRouter.Register(v1.Group("/pets"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)
			assert.Equal(t, tc.expectedCode, resp.StatusCode, tc.description)

			switch tc.expectedCode {
			case http.StatusOK:
				actualPetResponse := &Response{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(Response).ID, actualPetResponse.ID)
				assert.Equal(t, tc.expectedResponse.(Response).Name, actualPetResponse.Name)
			case http.StatusNotFound:
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

//func Test_PetsByName(t *testing.T) {
//	petMock := petServiceMock{}
//
//	petResponses := make([]Response, 1)
//	petRespopnse := &Response{
//		ID:   15,
//		Name: "Charles",
//	}
//	petResponses[0] = *petRespopnse
//
//	petMock.On("getPetsByName", "Charles").Return(petResponses, nil)
//	petRouter := NewRouter(&petMock)
//
//	r := setupRouter()
//	v1 := r.Group("/v1")
//	petRouter.Register(v1.Group("/pets"))
//
//	req := httptest.NewRequest("GET", "/v1/pets?name=Charles", nil)
//	resp, _ := r.Test(req, 5)
//
//	// Assert we encoded correctly,
//	// the request gives a 200
//	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200, got %d", resp.StatusCode)
//
//	// Read the response body
//	body, _ := io.ReadAll(resp.Body)
//	// unmarshal to Pet struct for asserts.
//	actualPetResponses := make([]Response, 1)
//	err := json.Unmarshal(body, &actualPetResponses)
//	if err != nil {
//		t.Errorf("Error unmarshalling response body: %v", err)
//		return
//	}
//
//	//assert.Equal(t, petRespopnse.ID, actualPetResponses[0].ID)
//	//assert.Equal(t, petRespopnse.Name, actualPetResponses[0].Name)
//}
