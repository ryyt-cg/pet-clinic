package repository

import "testing"

func TestVisitRepository_FindAll(t *testing.T) {

}

func TestVisitRepository_FindByID(t *testing.T) {}

func TestVisitRepository_Insert(t *testing.T) {

}

func TestVisitRepository_Update(t *testing.T) {}

/*
func (suite *VisitRepoTestSuite) Test_FindById() {
	var testCases = []struct {
		input    uint
		expected Visit
	}{
		{input: 2, expected: Visit{
			Model: gorm.Model{
				ID: 2,
			},
			PetID:       8,
			Description: "rabies shot",
		}},
		{input: 1, expected: Visit{
			Model: gorm.Model{
				ID: 1,
			},
			PetID:       7,
			Description: "rabies shot",
		}},
	}

	for _, testCase := range testCases {
		visit, _ := suite.visitRepository.FindById(testCase.input)
		assert.Equal(suite.T(), testCase.expected.ID, visit.ID)
		assert.Equal(suite.T(), testCase.expected.PetID, visit.PetID)
		assert.Equal(suite.T(), testCase.expected.Description, visit.Description)
		//assert.Equal(suite.T(), testCase.expected.VisitDate, visit.VisitDate)
	}
}
*/
