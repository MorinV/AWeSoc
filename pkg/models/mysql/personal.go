package mysql

import (
	"AWesomeSocial/pkg/models"
	"database/sql"
	"errors"
)

type PersonalRepository struct {
	DB *sql.DB
}

var Genders = map[string]int{
	"Мужчина": 1,
	"Женщина": 2,
}

func (m *PersonalRepository) Insert(firstname, secondname, surname, birthdate, gender, city, interests string, userId int) (int, error) {
	stmt := `INSERT INTO personal (firstname, secondname, surname, fullname, birthdate, gender, city, interests, user_id) 
VALUES(?,?,?,?,?,?,?,?,?)`
	fullname := surname + " " + firstname + " " + secondname
	result, err := m.DB.Exec(stmt, firstname, secondname, surname, fullname, birthdate, gender, city, interests, userId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PersonalRepository) Get(id int) (*models.Personal, error) {
	stmt := `SELECT * FROM personal WHERE id = ?`
	p := &models.Personal{}
	err := m.DB.QueryRow(stmt, id).
		Scan(&p.Id, &p.Firstname, &p.Secondname, &p.Surname, &p.Fullname, &p.Birthdate, &p.Gender, &p.City, &p.Interests, &p.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return p, nil
}

func (m *PersonalRepository) GetByUserId(id int) (*models.Personal, error) {
	stmt := `SELECT * FROM personal WHERE user_id = ?`
	p := &models.Personal{}
	err := m.DB.QueryRow(stmt, id).
		Scan(&p.Id, &p.Firstname, &p.Secondname, &p.Surname, &p.Fullname, &p.Birthdate, &p.Gender, &p.City, &p.Interests, &p.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return p, nil
}

func (m *PersonalRepository) Latest() ([]*models.Personal, error) {
	stmt := `SELECT * FROM personal order by id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var personals []*models.Personal

	for rows.Next() {
		p := &models.Personal{}
		err = rows.Scan(&p.Id, &p.Firstname, &p.Secondname, &p.Surname, &p.Fullname, &p.Birthdate, &p.Gender, &p.City, &p.Interests, &p.UserId)
		if err != nil {
			return nil, err
		}
		personals = append(personals, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return personals, nil
}
