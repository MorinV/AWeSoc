package person

import (
	"AWesomeSocial/domain/person"
	"AWesomeSocial/internal"
	"database/sql"
	"errors"
)

const dateLayoutIso = "2006-01-02"

type WriteRepo struct {
	DB *sql.DB
}

type ReadRepo struct {
	DB *sql.DB
}

func (rr *ReadRepo) Search(firstNamePref, surNamePref string) ([]*person.Person, error) {
	stmt := `SELECT * FROM personal WHERE firstname like ? AND surname like ? order by id`
	rows, err := rr.DB.Query(stmt, firstNamePref+"%", surNamePref+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var persons []*person.Person

	for rows.Next() {
		pdb := &DbValue{}
		err = rows.Scan(&pdb.Id, &pdb.Firstname, &pdb.Secondname, &pdb.Surname, &pdb.Fullname, &pdb.Birthdate, &pdb.Gender, &pdb.City, &pdb.Interests, &pdb.UserId)
		if err != nil {
			return nil, err
		}
		p, err := pdb.Hydrate()
		if err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}

func (rr *ReadRepo) Get(id int) (*person.Person, error) {
	stmt := `SELECT * FROM personal WHERE id = ?`
	pdb := &DbValue{}
	err := rr.DB.QueryRow(stmt, id).
		Scan(&pdb.Id, &pdb.Firstname, &pdb.Secondname, &pdb.Surname, &pdb.Fullname, &pdb.Birthdate, &pdb.Gender, &pdb.City, &pdb.Interests, &pdb.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internal.ErrNoRecord
		} else {
			return nil, err
		}
	}
	p, err := pdb.Hydrate()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (rr *ReadRepo) GetByUserId(id int) (*person.Person, error) {
	stmt := `SELECT * FROM personal WHERE user_id = ?`
	pdb := &DbValue{}
	err := rr.DB.QueryRow(stmt, id).
		Scan(&pdb.Id, &pdb.Firstname, &pdb.Secondname, &pdb.Surname, &pdb.Fullname, &pdb.Birthdate, &pdb.Gender, &pdb.City, &pdb.Interests, &pdb.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internal.ErrNoRecord
		} else {
			return nil, err
		}
	}
	p, err := pdb.Hydrate()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (rr *ReadRepo) Latest() ([]*person.Person, error) {
	stmt := `SELECT * FROM personal order by id DESC LIMIT 10`
	rows, err := rr.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var persons []*person.Person

	for rows.Next() {
		pdb := &DbValue{}
		err = rows.Scan(&pdb.Id, &pdb.Firstname, &pdb.Secondname, &pdb.Surname, &pdb.Fullname, &pdb.Birthdate, &pdb.Gender, &pdb.City, &pdb.Interests, &pdb.UserId)
		if err != nil {
			return nil, err
		}
		p, err := pdb.Hydrate()
		if err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}

func (wr *WriteRepo) Add(p *person.Person) (int, error) {
	if p.Id == 0 {
		id, err := wr.insert(p)
		if err != nil {
			return 0, err
		}

		return id, nil
	} else {
		err := wr.update(p)

		if err != nil {
			return 0, err
		}

		return p.Id, err
	}
}

func (wr *WriteRepo) insert(p *person.Person) (int, error) {
	stmt := `INSERT INTO personal (firstname, secondname, surname, fullname, birthdate, gender, city, interests, user_id) 
VALUES(?,?,?,?,?,?,?,?,?)`
	result, err := wr.DB.Exec(stmt, p.Firstname, p.Secondname, p.Surname, p.Fullname, p.Birthdate.Format(dateLayoutIso), p.Gender.Int(), p.City, p.Interests, p.UserId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (wr *WriteRepo) update(p *person.Person) error {
	stmt := `UPDATE personal 
SET firstname = ?, secondname = ?, surname = ?, fullname = ?, birthdate = ?, gender = ?, city = ?, interests = ?, user_id = ?
WHERE id = ?`
	_, err := wr.DB.Exec(stmt, p.Firstname, p.Secondname, p.Surname, p.Fullname, p.Birthdate.Format(dateLayoutIso), p.Gender.Int(), p.City, p.Interests, p.Id)

	return err
}
