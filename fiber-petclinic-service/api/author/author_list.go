package author

import (
	"fiber-petclinic-service/exception"
)

var (
	authorList []Author
)

func init() {
	authorList = []Author{
		{ID: 1, Name: "Lucas Boyce", Bio: "Lucas Boyce is a software engineer with a passion for building scalable applications.", Website: "https://lucasboyce.dev"},
		{ID: 2, Name: "Jane Doe", Bio: "Jane Doe is a tech enthusiast and writer.", Website: "https://janedoe.com"},
		{ID: 3, Name: "John Smith", Bio: "John Smith is a developer and open-source contributor.", Website: "https://johnsmith.io"},
	}
}

func getAuthorByID(ID int) (*Author, error) {
	// Simulate a database lookup
	for _, author := range authorList {
		if author.ID == ID {
			return &author, nil
		}
	}

	return nil, exception.ErrorNotFound // Return nil if not found
}

func getAuthorByName(name string) (*Author, error) {
	for _, author := range authorList {
		if author.Name == name {
			return &author, nil
		}
	}

	return nil, exception.ErrorNotFound
}
