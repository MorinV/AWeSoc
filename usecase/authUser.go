package usecase

import (
	"AWesomeSocial/domain/person"
	"AWesomeSocial/domain/user"
	"AWesomeSocial/internal"
	"errors"
)

type AuthenticatedUser struct {
	User   *user.User
	Person *person.Person
}

func (app *Application) AuthUser(userLogin, userPassword string) (*AuthenticatedUser, error) {
	if userLogin == "" || userPassword == "" {
		return nil, errors.New("не введены учетные данные")
	}

	login, err := user.NewLogin(userLogin)
	if err != nil {
		return nil, err
	}
	u, err := app.repositoryRegistry.GetUserReadRepo().GetByLogin(login)
	if err != nil {
		return nil, err
	}
	if !u.Password.Equals(userPassword) {
		return nil, errors.New("пароль не совпадает")
	}

	p, err := app.repositoryRegistry.GetPersonReadRepo().GetByUserId(u.Id)
	if err != nil {
		if errors.Is(internal.ErrNoRecord, err) {
			p = &person.Person{}
		} else {
			return nil, err
		}
	}
	_, err = app.repositoryRegistry.GetUserWriteRepo().Add(u)
	if err != nil {
		return nil, err
	}

	return &AuthenticatedUser{User: u, Person: p}, nil
}
