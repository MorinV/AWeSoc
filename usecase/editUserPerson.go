package usecase

import (
	"AWesomeSocial/domain/person"
	"AWesomeSocial/internal"
	"errors"
	"strconv"
	"time"
)

func (app *Application) EditUserPerson(firstname, secondname, surname, birthdate, gender, city, interests string, userId int) (int, error) {
	genderValue, err := strconv.Atoi(gender)
	if err != nil {
		return 0, err
	}
	p, err := app.repositoryRegistry.GetPersonReadRepo().GetByUserId(userId)
	if err != nil {
		if errors.Is(internal.ErrNoRecord, err) {
			p, err = makePerson(firstname, secondname, surname, birthdate, genderValue, city, interests, userId)
			if err != nil {
				return 0, err
			}
			id, err := app.repositoryRegistry.GetPersonWriteRepo().Add(p)
			return id, err
		} else {
			return 0, err
		}
	}

	birthdateTime, err := time.Parse(dateLayoutIso, birthdate)
	if err != nil {
		return 0, err
	}
	genderStruct, err := person.NewGenderFromValue(genderValue)
	if err != nil {
		return 0, err
	}
	err = p.Update(firstname, secondname, surname, birthdateTime, genderStruct, city, interests, userId)
	if err != nil {
		return 0, err
	}
	id, err := app.repositoryRegistry.GetPersonWriteRepo().Add(p)
	if err != nil {
		return 0, err
	}
	return id, err
}

func makePerson(firstname string, secondname string, surname string, birthdate string, gender int, city string, interests string, userId int) (*person.Person, error) {
	birthdateTime, err := time.Parse(dateLayoutIso, birthdate)
	if err != nil {
		return nil, err
	}
	genderStruct, err := person.NewGenderFromValue(gender)
	if err != nil {
		return nil, err
	}
	p := person.New(0, firstname, secondname, surname, birthdateTime, genderStruct, city, interests, userId)
	return p, nil
}
