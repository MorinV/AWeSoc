package user

import (
	"AWesomeSocial/domain/user"
	"AWesomeSocial/internal"
	"database/sql"
	"errors"
	"time"
)

const dateLayoutIso = "2006-01-02"
const dateLayoutIsoTime = "2006-01-02 01:02:03"

type WriteRepo struct {
	DB *sql.DB
}

type ReadRepo struct {
	DB *sql.DB
}

func (rr *ReadRepo) GetByLogin(login user.Login) (*user.User, error) {
	stmt := `SELECT * FROM users WHERE login = ?`
	u := &user.User{}
	err := rr.DB.QueryRow(stmt, login.String()).
		Scan(&u.Id, &u.Login, &u.Password, &u.Email, &u.Created, &u.LastLogin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internal.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}

func (wr *WriteRepo) Add(u *user.User) (int, error) {
	if u.Id == 0 {
		id, err := wr.insert(u)
		if err != nil {
			return 0, err
		}

		return id, nil
	} else {
		err := wr.update(u)

		if err != nil {
			return 0, err
		}

		return u.Id, err
	}
}

func (wr *WriteRepo) insert(u *user.User) (int, error) {
	stmt := `INSERT INTO users (login, password, email, created, last_login) 
VALUES(?,?,?,?,?)`
	result, err := wr.DB.Exec(stmt, u.Login, u.Password, u.Email, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (wr *WriteRepo) update(u *user.User) error {
	stmt := `UPDATE users SET login = ?, password = ?, email = ?, created = ?, last_login = ? WHERE id = ?`
	_, err := wr.DB.Exec(stmt, u.Login, u.Password, u.Email, u.Created.Format(dateLayoutIsoTime), u.LastLogin.Format(dateLayoutIsoTime), u.Id)

	return err
}
