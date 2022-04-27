package usecase

import (
	"AWesomeSocial/domain/person"
	"errors"
)

type SearchPersonParams struct {
	FirstNamePref string
	SurNamePref   string
}

func (app *Application) SearchPersons(params SearchPersonParams) ([]*person.Person, error) {
	if len(params.FirstNamePref) < 2 || len(params.SurNamePref) < 2 {
		return nil, errors.New("слишком короткая строка поиска")
	}
	persons, err := app.repositoryRegistry.GetPersonReadRepo().Search(params.FirstNamePref, params.SurNamePref)
	if err != nil {
		return nil, err
	}

	return persons, nil
}
