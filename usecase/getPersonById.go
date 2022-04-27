package usecase

import "AWesomeSocial/domain/person"

func (app *Application) GetPerson(id int) (*person.Person, error) {
	persons, err := app.repositoryRegistry.GetPersonReadRepo().Get(id)
	if err != nil {
		return nil, err
	}

	return persons, nil
}
