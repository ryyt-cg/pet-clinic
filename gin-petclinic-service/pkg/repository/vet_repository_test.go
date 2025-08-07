package repository

import (
	model2 "gin-petclinic-service/pkg/repository/model"
	"gin-petclinic-service/pkg/test"
	"gorm.io/gorm"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/stretchr/testify/suite"

	"gopkg.in/go-playground/assert.v1"
)

type VetRepoTestSuite struct {
	suite.Suite
	postgresql    *embeddedpostgres.EmbeddedPostgres
	vetRepository *VetRepository
}

// This will run before the tests in the suite are run
func (suite *VetRepoTestSuite) SetupSuite() {
	suite.postgresql = test.PgStart(suite.T(), "test/migrations")
	suite.vetRepository = getVetRepository(suite.T())
}

func (suite *VetRepoTestSuite) TearDownSuite() {
	err := suite.postgresql.Stop()
	if err != nil {
		suite.T().Fatal(err)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestVetRepoTestSuite(t *testing.T) {
	suite.Run(t, new(VetRepoTestSuite))
}

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
			Person: model2.Person{
				FirstName: "James",
				LastName:  "Carter",
			},
		}},
		{2, Vet{
			Model: gorm.Model{
				ID: 2,
			},
			Person: model2.Person{
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
			Person: model2.Person{
				FirstName: "James",
				LastName:  "Carter",
			},
		}},
		{"Leary", Vet{
			Person: model2.Person{
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
