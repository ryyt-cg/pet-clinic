package vet

import (
	"testing"
)

func Test_allSpecialties(t *testing.T) {
}

func Test_allVets(t *testing.T) {}

func Test_vetById(t *testing.T) {
	//vetMock := vetServiceMock{}
	//
	//vetResponse := &Response{
	//	ID:        1,
	//	FirstName: "Nat",
	//	LastName:  "Cole",
	//}
	//
	//vetMock.On("getVetById", 1).Return(*vetResponse, nil)
	//vetRouter := NewRouter(&vetMock)
	//
	//r := setupRouter()
	//v1 := r.Group("/v1")
	//vetRouter.Register(v1.Group("/vets"))
	//
	//w := httptest.NewRecorder()
	//req, _ := http.NewRequest("GET", "/v1/vets/1", nil)
	//r.ServeHTTP(w, req)
	//
	//// Assert we encoded correctly,
	//// the request gives a 200
	//assert.Equal(t, http.StatusOK, w.Code)
	//
	//// unmarshal to Vet struct for asserts.
	//actualVet := &repository.Vet{}
	//jsoniter.Unmarshal(w.Body.Bytes(), actualVet)
	//assert.Equal(t, vetResponse.ID, actualVet.ID)
	//assert.Equal(t, vetResponse.FirstName, actualVet.FirstName)
	//assert.Equal(t, vetResponse.LastName, actualVet.LastName)
}

func Test_getVetByIdWithSpecialties(t *testing.T) {
	//vetMock := vetServiceMock{}
	//
	//vetResponse := &Response{
	//	ID:        1,
	//	FirstName: "Nat",
	//	LastName:  "Cole",
	//}
	//
	//vetMock.On("getVetByIdWithSpecialties", 1).Return(*vetResponse, nil)
	//vetRouter := NewRouter(&vetMock)
	//
	//r := setupRouter()
	//v1 := r.Group("/v1")
	//vetRouter.Register(v1.Group("/vets"))
	//
	//w := httptest.NewRecorder()
	//req, _ := http.NewRequest("GET", "/v1/vets/1/specialties", nil)
	//r.ServeHTTP(w, req)
	//
	//// Assert we encoded correctly,
	//// the request gives a 200
	//assert.Equal(t, http.StatusOK, w.Code)
}

func Test_VetByLastName(t *testing.T) {
	//vetMock := vetServiceMock{}
	//
	//var vetResponses = make([]Response, 1)
	//var vetResponse = &Response{
	//	ID:        15,
	//	FirstName: "Charles",
	//	LastName:  "Ward",
	//}
	//vetResponses[0] = *vetResponse
	//
	//vetMock.On("getVetByLastName", "Ward").Return(vetResponses, nil)
	//vetRouter := NewRouter(&vetMock)
	//
	//r := setupRouter()
	//v1 := r.Group("/v1")
	//vetRouter.Register(v1.Group("/vets"))
	//
	//w := httptest.NewRecorder()
	//req, _ := http.NewRequest("GET", "/v1/vets?last-name=Ward", nil)
	//r.ServeHTTP(w, req)
	//
	//// Assert we encoded correctly,
	//// the request gives a 200
	//assert.Equal(t, http.StatusOK, w.Code)
	//
	//// unmarshal to Vet struct for asserts.
	//actualVetResponses := make([]Response, 1)
	//jsoniter.Unmarshal(w.Body.Bytes(), &actualVetResponses)
	//
	//assert.Equal(t, vetResponse.ID, actualVetResponses[0].ID)
	//assert.Equal(t, vetResponse.FirstName, actualVetResponses[0].FirstName)
	//assert.Equal(t, vetResponse.LastName, actualVetResponses[0].LastName)
}

func Test_createVet(t *testing.T) {
	//vetMock := vetServiceMock{}
	//
	//vetRequest := &repository.Vet{
	//	FirstName: "Nat",
	//	LastName:  "Cole",
	//}
	//
	//vetResponse := &Response{
	//	ID:        1,
	//	FirstName: "Nat",
	//	LastName:  "Cole",
	//}
	//
	//vetMock.On("create", vetRequest).Return(*vetResponse, nil)
	//vetRouter := NewRouter(&vetMock)
	//
	//r := setupRouter()
	//v1 := r.Group("/v1")
	//vetRouter.Register(v1.Group("/vets"))
	//
	//w := httptest.NewRecorder()
	//req, _ := http.NewRequest("POST", "/v1/vets", bytes.NewBuffer(jsoniter.Marshal(vetRequest)))
	//req.Header.Set("Content-Type", "application/json")
	//r.ServeHTTP(w, req)
	//
	//// Assert we encoded correctly,
	//// the request gives a 200
	//assert.Equal(t, http.StatusCreated, w.Code)
}
func Test_updateVet(t *testing.T) {

}
