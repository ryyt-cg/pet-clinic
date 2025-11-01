package visit

import (
	"encoding/json"
	"errors"
	resterr "gin-petclinic-service/internal/errors"
	"gin-petclinic-service/internal/repository"
	"gin-petclinic-service/internal/repository/model"
	"gin-petclinic-service/internal/test"
	"net/http"
	"testing"

	"gorm.io/gorm"
)

func Test_visitById(t *testing.T) {
	visitResponse := &Response{
		ID:          1,
		VisitDate:   "2023-03-05",
		Description: "rabies shot",
		PetID:       101,
	}
	notFound := resterr.NotFound("record not found")
	internalError := resterr.InternalServerError("fail to get a visit by id")

	testCases := []struct {
		name             string
		id               uint
		mockVisit        *Response
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get owner by ID",
			id:               1,
			mockVisit:        visitResponse,
			mockError:        nil,
			route:            "/v1/visits/1",
			statusCode:       http.StatusOK,
			expectedResponse: visitResponse,
		},
		{
			name:             "get no owner by ID",
			id:               1,
			mockVisit:        nil,
			mockError:        gorm.ErrRecordNotFound,
			route:            "/v1/visits/1",
			statusCode:       http.StatusNotFound,
			expectedResponse: &notFound,
		},
		{
			name:             "fail to get owner by ID",
			id:               1,
			mockVisit:        nil,
			mockError:        errors.New("fail to get a visit by id"),
			route:            "/v1/visits/1",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: &internalError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			visitMock := NewMockServicer(t)
			visitMock.EXPECT().getVisitById(tc.id).Return(tc.mockVisit, tc.mockError)

			visitRouter := NewRouter(visitMock)
			visitRouter.Register(v1.Group("/visits"))

			jsonStr, _ := json.Marshal(tc.expectedResponse)
			tests := []test.APITestCase{
				{tc.name, "GET", tc.route, "", nil, tc.statusCode, string(jsonStr)},
			}
			for _, tc := range tests {
				test.Endpoint(t, r, tc)
			}
		})
	}
}

// Test_visitById_BadRequest
func Test_visitById_BadRequest(t *testing.T) {
	r := test.SetupRouter()
	v1 := r.Group("/v1")
	visitMock := NewMockServicer(t)

	petRouter := NewRouter(visitMock)
	petRouter.Register(v1.Group("/visits"))

	jsonStr, _ := json.Marshal(resterr.BadRequest("strconv.Atoi: parsing \"invalid-id\": invalid syntax"))
	tests := []test.APITestCase{
		{"Get visit by invalid it", "GET", "/v1/visits/invalid-id", "", nil, http.StatusBadRequest, string(jsonStr)},
	}
	for _, tc := range tests {
		test.Endpoint(t, r, tc)
	}
}

// Test_allVisits
func Test_allVisits(t *testing.T) {
	visitsResponses := &Responses{
		Context: model.Context{Count: 3},
		Visits: []Response{
			{
				ID:          35,
				VisitDate:   "2023-03-05",
				Description: "rabies shot",
				PetID:       201,
			},
			{
				ID:          26,
				VisitDate:   "2023-03-06",
				Description: "vaccination",
				PetID:       202,
			},
			{
				ID:          28,
				VisitDate:   "2023-03-07",
				Description: "check-up",
				PetID:       203,
			},
		},
	}

	noVisitsFoundResponse := resterr.NotFound("record not found")
	errorResponse := resterr.InternalServerError("fail to get all visits")

	testCases := []struct {
		name              string
		mockAllVisits     *Responses
		mockError         error
		route             string
		statusCode        int
		expectedResponses interface{}
	}{
		{
			name:              "get all visits",
			mockAllVisits:     visitsResponses,
			mockError:         nil,
			route:             "/v1/visits/all",
			statusCode:        http.StatusOK,
			expectedResponses: visitsResponses,
		},
		{
			name:              "get no owner",
			mockAllVisits:     nil,
			mockError:         gorm.ErrRecordNotFound,
			route:             "/v1/visits/all",
			statusCode:        http.StatusNotFound,
			expectedResponses: &noVisitsFoundResponse,
		},
		{
			name:              "fail to get all owners",
			mockAllVisits:     nil,
			mockError:         errors.New("fail to get all visits"),
			route:             "/v1/visits/all",
			statusCode:        http.StatusInternalServerError,
			expectedResponses: &errorResponse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			visitMock := NewMockServicer(t)
			visitMock.EXPECT().getAllVisits().Return(tc.mockAllVisits, tc.mockError)

			petRouter := NewRouter(visitMock)
			petRouter.Register(v1.Group("/visits"))

			jsonStr, _ := json.Marshal(tc.expectedResponses)
			tests := []test.APITestCase{
				{tc.name, "GET", tc.route, "", nil, tc.statusCode, string(jsonStr)},
			}
			for _, tc := range tests {
				test.Endpoint(t, r, tc)
			}

		})
	}
}

func Test_updateVisit(t *testing.T) {
	visit := &repository.Visit{
		Model: gorm.Model{
			ID: 1,
		},
		VisitDate:   test.ToDate("2023-03-05"),
		Description: "rabies shot",
		PetID:       101,
	}
	visitRequest := &UpdateRequest{
		ID:          1,
		VisitDate:   "2023-03-05",
		Description: "rabies shot",
		PetID:       101,
	}

	visitResponse := &Response{
		ID:          1,
		VisitDate:   "2023-03-05",
		Description: "rabies shot",
		PetID:       101,
	}
	notFound := resterr.NotFound("record not found")
	internalError := resterr.InternalServerError("fail to update a visit")

	testCases := []struct {
		name             string
		id               uint
		visitEntity      *repository.Visit
		request          *UpdateRequest
		mockVisit        *Response
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "update visit by ID",
			id:               1,
			visitEntity:      visit,
			request:          visitRequest,
			mockVisit:        visitResponse,
			mockError:        nil,
			route:            "/v1/visits/1",
			statusCode:       http.StatusOK,
			expectedResponse: visitResponse,
		},
		{
			name:             "no visit found by ID to update",
			id:               1,
			visitEntity:      visit,
			request:          visitRequest,
			mockVisit:        nil,
			mockError:        gorm.ErrRecordNotFound,
			route:            "/v1/visits/1",
			statusCode:       http.StatusNotFound,
			expectedResponse: &notFound,
		},
		{
			name:             "fail to update visit by ID",
			id:               1,
			visitEntity:      visit,
			request:          visitRequest,
			mockVisit:        nil,
			mockError:        errors.New("fail to update a visit"),
			route:            "/v1/visits/1",
			statusCode:       http.StatusInternalServerError,
			expectedResponse: &internalError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			visitMock := NewMockServicer(t)
			visitMock.EXPECT().update(tc.visitEntity).Return(tc.mockVisit, tc.mockError)

			petRouter := NewRouter(visitMock)
			petRouter.Register(v1.Group("/visits"))

			jsonStr, _ := json.Marshal(tc.expectedResponse)
			requestStr, _ := json.Marshal(tc.request)
			tests := []test.APITestCase{
				{tc.name, "PUT", tc.route, string(requestStr), nil, tc.statusCode, string(jsonStr)},
			}
			for _, tc := range tests {
				test.Endpoint(t, r, tc)
			}

		})
	}
}

/*
func Test_updateVisit_BadRequest(t *testing.T) {
	visitRequest := `{
	  "id": 1,
	  "visitDate": "2023-03-31",
	  "description": "rabies shot",
	  "petID": "101"
	}`

	visitRequest2 := `{
	  "id": 1,
	  "visitDate": "2023-03-32",
	  "description": "rabies shot",
	  "petID": 101
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
			name:             "update visit found by ID with invalid ID",
			id:               "invalid-id",
			route:            "/v1/visits/invalid-id",
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("failed to convert: strconv.Atoi: parsing \"invalid-id\": invalid syntax"),
		},
		{
			name:             "update visit by ID with bad request",
			id:               1,
			request:          visitRequest,
			route:            "/v1/visits/1",
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("json: cannot unmarshal string into Go struct field UpdateRequest.petID of type uint"),
		},
		{
			name:             "update visit by ID with bad request",
			id:               1,
			request:          visitRequest2,
			route:            "/v1/visits/1",
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("parsing time \"2023-03-32\": day out of range"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			visitMock := NewMockServicer(t)
			visitRouter := NewRouter(visitMock)
			visitRouter.Register(v1.Group("/visits"))

			req := httptest.NewRequest("PUT", tc.route, bytes.NewReader([]byte(tc.request)))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, 5)

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

func Test_addNewVisit(t *testing.T) {
	visit := &repository.Visit{
		VisitDate:   test.ToDate("2023-03-05"),
		Description: "rabies shot",
		PetID:       101,
	}

	visitRequest := &AddRequest{
		VisitDate:   "2023-03-05",
		Description: "rabies shot",
		PetID:       101,
	}

	visitResponse := &Response{
		ID:          1,
		VisitDate:   "2023-03-05",
		Description: "rabies shot",
		PetID:       101,
	}
	internalError := resterr.InternalServerError("fail to create a visit")

	testCases := []struct {
		name             string
		newVisit         *repository.Visit
		request          *AddRequest
		mockVisit        *Response
		mockError        error
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "add new visit",
			newVisit:         visit,
			request:          visitRequest,
			mockVisit:        visitResponse,
			mockError:        nil,
			statusCode:       http.StatusCreated,
			expectedResponse: visitResponse,
		},
		{
			name:             "fail to add new visit",
			newVisit:         visit,
			request:          visitRequest,
			mockVisit:        nil,
			mockError:        errors.New("fail to create a visit"),
			statusCode:       http.StatusInternalServerError,
			expectedResponse: &internalError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			visitMock := NewMockServicer(t)
			visitMock.EXPECT().create(tc.newVisit).Return(tc.mockVisit, tc.mockError)
			visitRouter := NewRouter(visitMock)
			visitRouter.Register(v1.Group("/visits"))
			requestBody, _ := json.Marshal(tc.request)
			req := httptest.NewRequest("POST", "/v1/visits", bytes.NewReader(requestBody))
			// Must se content-type header.
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, 5)

			switch tc.statusCode {
			case http.StatusCreated:
				actualVisitResponse := &Response{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualVisitResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.statusCode, resp.StatusCode)
				assert.Equal(t, tc.expectedResponse.(*Response), actualVisitResponse)
			case http.StatusInternalServerError:
				actualVisitResponse := &resterr.ErrorResponse{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualVisitResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(*resterr.ErrorResponse).Status, actualVisitResponse.Status)
				assert.Equal(t, tc.expectedResponse.(*resterr.ErrorResponse).Message, actualVisitResponse.Message)
			}
		})
	}
}

func Test_addNewVisit_BadRequest(t *testing.T) {
	visitRequest := `{
	  "visitDate": "2023-03-31",
	  "description": "rabies shot",
	  "petID": "101"
	}`

	visitRequest2 := `{
	  "visitDate": "2023-03-32",
	  "description": "rabies shot",
	  "petID": 101
	}`

	testCases := []struct {
		name             string
		request          string
		statusCode       int
		expectedResponse resterr.ErrorResponse
	}{
		{
			name:             "add new visit with bad request",
			request:          visitRequest,
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("json: cannot unmarshal string into Go struct field AddRequest.petID of type uint"),
		},
		{
			name:             "update visit by ID with bad request",
			request:          visitRequest2,
			statusCode:       http.StatusBadRequest,
			expectedResponse: resterr.BadRequest("parsing time \"2023-03-32\": day out of range"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := test.SetupRouter()
			v1 := r.Group("/v1")

			visitMock := NewMockServicer(t)
			visitRouter := NewRouter(visitMock)
			visitRouter.Register(v1.Group("/visits"))

			req := httptest.NewRequest("POST", "/v1/visits", bytes.NewReader([]byte(tc.request)))
			// Must se content-type header.
			req.Header.Set("Content-Type", "application/json")
			resp, _ := r.Test(req, 5)

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


*/
