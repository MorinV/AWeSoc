package usecase

import "AWesomeSocial/domain/person"

func (app *Application) GetLatestRegisteredPersons() ([]*person.Person, error) {
	persons, err := app.repositoryRegistry.GetPersonReadRepo().Latest()
	if err != nil {
		return nil, err
	}

	return persons, nil
}
