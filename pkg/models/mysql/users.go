package mysql

import (
	"AWesomeSocial/pkg/models"
	"database/sql"
	"errors"
	"time"
)

type UsersRepository struct {
	DB *sql.DB
}

func (m *UsersRepository) Insert(login, password, email string) (int, error) {
	stmt := `INSERT INTO users (login, password, email, created, last_login) 
VALUES(?,?,?,?,?)`
	result, err := m.DB.Exec(stmt, login, password, email, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UsersRepository) GetByLogin(login string) (*models.User, error) {
	stmt := `SELECT * FROM users WHERE login = ?`
	user := &models.User{}
	err := m.DB.QueryRow(stmt, login).
		Scan(&user.Id, &user.Login, &user.Password, &user.Email, &user.Created, &user.LastLogin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return user, nil
}
