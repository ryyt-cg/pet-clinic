package pet

/*
func Test_getAllPets(t *testing.T) {
	mockEntityPets := []repository.Pet{
		{Model: gorm.Model{ID: 1}, Name: "Tom",
			Birthdate: time.Date(2015, 11, 19, 0, 0, 0, 00, time.UTC),
			TypeID:    19, OwnerID: 7},
		{Model: gorm.Model{ID: 2}, Name: "Mike",
			Birthdate: time.Date(2018, 4, 17, 0, 0, 0, 0, time.UTC),
			TypeID:    20, OwnerID: 7},
	}

	mockPets := &Responses{
		Context: model.Context{
			Count: 2,
		},
		Pets: FromPets(mockEntityPets),
	}

	tests := []struct {
		name             string
		mockPets         *Responses
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
		//{
		//	name:             "get no pet",
		//	mockPets:         nil,
		//	mockError:        gorm.ErrRecordNotFound,
		//	route:            "/v1/pets/all",
		//	statusCode:       http.StatusNotFound,
		//	expectedResponse: resterr.NotFound("Pet not found"),
		//},
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
				actualPetResponse := &Responses{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(*Responses).Context.Count, actualPetResponse.Context.Count)
				assert.Equal(t, tc.expectedResponse.(*Responses).Pets, actualPetResponse.Pets)
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
	mockPet := &Response{
		ID:        1,
		Name:      "Running Water",
		Birthdate: "2018-02-26",
	}

	tests := []struct {
		name             string
		id               interface{}
		mockPet          *Response
		mockError        error
		route            string
		statusCode       int
		expectedResponse interface{}
	}{
		{
			name:             "get pet by id",
			id:               int(1),
			mockPet:          mockPet,
			mockError:        nil,
			route:            "/v1/pets/1",
			statusCode:       http.StatusOK,
			expectedResponse: mockPet,
		},
		{
			name:             "get no pet by id",
			id:               int(1),
			mockPet:          nil,
			mockError:        gorm.ErrRecordNotFound,
			route:            "/v1/pets/1",
			statusCode:       http.StatusNotFound,
			expectedResponse: resterr.NotFound("Pet not found"),
		},
		{
			name:             "fail to get pet by id",
			id:               int(1),
			mockPet:          nil,
			mockError:        errors.New("unable to get pet by id"),
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
			petMock.EXPECT().getPetById(1).Return(tc.mockPet, tc.mockError)

			petRouter := NewRouter(petMock)
			petRouter.Register(v1.Group("/pets"))

			req := httptest.NewRequest("GET", tc.route, nil)
			resp, _ := r.Test(req, 5)
			assert.Equal(t, tc.statusCode, resp.StatusCode)

			switch resp.StatusCode {
			case http.StatusOK:
				actualPetResponse := &Response{}
				// Read the response body
				body, _ := io.ReadAll(resp.Body)
				err := json.Unmarshal(body, actualPetResponse)
				if err != nil {
					t.Errorf("Error unmarshalling response body: %v", err)
					return
				}
				assert.Equal(t, tc.expectedResponse.(*Response).ID, actualPetResponse.ID)
				assert.Equal(t, tc.expectedResponse.(*Response).Name, actualPetResponse.Name)
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

func Test_getPetByInvalidId(t *testing.T) {
	mockError := errors.New("unable to get pet by id")
	route := "/v1/pets/a1"
	expectedResponse := resterr.BadRequest("failed to convert: strconv.Atoi: parsing \"a1\": invalid syntax")

	r := test.SetupRouter()
	v1 := r.Group("/v1")

	petMock := NewMockServicer(t)
	petMock.EXPECT().getPetById("a1").Return(nil, mockError)

	petRouter := NewRouter(petMock)
	petRouter.Register(v1.Group("/pets"))

	req := httptest.NewRequest("GET", route, nil)
	resp, _ := r.Test(req, 5)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	actualPetResponse := &resterr.ErrorResponse{}
	// Read the response body
	body, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal(body, actualPetResponse)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
		return
	}
	assert.Equal(t, expectedResponse.Message, actualPetResponse.Message)
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
*/
