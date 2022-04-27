package usecase

import "AWesomeSocial/domain/user"

func (app *Application) PersistUser(id int, login, password, email string) (int, error) {
	l, err := user.NewLogin(login)
	if err != nil {
		return 0, err
	}
	p, err := user.NewPassword(password)
	if err != nil {
		return 0, err
	}
	e, err := user.NewEmail(email)
	if err != nil {
		return 0, err
	}

	u := user.New(id, l, p, e)
	id, err = app.repositoryRegistry.GetUserWriteRepo().Add(u)
	if err != nil {
		return 0, err
	}

	return id, nil
}
