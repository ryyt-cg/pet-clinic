package repository

import "testing"

func TestVetRepository_FindAll(t *testing.T) {

}

func TestVetRepository_FindById(t *testing.T)           {}
func TestVetRepository_FindAllWithVisits(t *testing.T)  {}
func TestVetRepository_FindByIdWithVisits(t *testing.T) {}
func TestVetRepository_FindAllSpecialties(t *testing.T) {

}

func TestVetRepository_FindAllSpecialtiesWithVisits(t *testing.T) {}
func TestVetRepository_FindByLastName(t *testing.T) {

}

func TestVetRepository_FindByLastNameWithVisits(t *testing.T) {}
func TestVetRepository_Insert(t *testing.T) {

}

func TestVetRepository_Update(t *testing.T) {}

/*
func getVetRepository(t *testing.T) *VetRepository {
	db, err := test.Connect()
	if err != nil {
		t.Fatal(err)
	}

	return NewVetRepository(db)
}

func (suite *VetRepoTestSuite) Test_FindById() {
	var testCases = []struct {
		input    uint
		expected Vet
	}{
		{1, Vet{
			Model: gorm.Model{
				ID: 1,
			},
			Person: model.Person{
				FirstName: "James",
				LastName:  "Carter",
			},
		}},
		{2, Vet{
			Model: gorm.Model{
				ID: 2,
			},
			Person: model.Person{
				FirstName: "Helen",
				LastName:  "Leary",
			},
		}},
	}

	for _, testCase := range testCases {
		vet, _ := suite.vetRepository.FindById(testCase.input)
		assert.Equal(suite.T(), testCase.expected.LastName, vet.LastName)
		assert.Equal(suite.T(), testCase.expected.FirstName, vet.FirstName)
	}
}

func (suite *VetRepoTestSuite) Test_FindByLastName() {
	var testCases = []struct {
		input    string
		expected Vet
	}{
		{"Carter", Vet{
			Person: model.Person{
				FirstName: "James",
				LastName:  "Carter",
			},
		}},
		{"Leary", Vet{
			Person: model.Person{
				FirstName: "Helen",
				LastName:  "Leary",
			},
		}},
	}

	for _, testCase := range testCases {
		vets, _ := suite.vetRepository.FindByLastName(testCase.input)
		assert.Equal(suite.T(), testCase.expected.LastName, vets[0].LastName)
		assert.Equal(suite.T(), testCase.expected.FirstName, vets[0].FirstName)
	}
}
*/
