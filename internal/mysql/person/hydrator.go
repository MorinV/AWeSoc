package person

import (
	"AWesomeSocial/domain/person"
	"time"
)

type DbValue struct {
	Id         int
	Firstname  string
	Secondname string
	Surname    string
	Fullname   string
	Birthdate  time.Time
	Gender     string
	City       string
	Interests  string
	UserId     int
}

func (p *DbValue) Hydrate() (*person.Person, error) {
	gender, err := person.NewGenderFromName(p.Gender)
	if err != nil {
		return nil, err
	}

	return &person.Person{
		Id:         p.Id,
		Firstname:  p.Firstname,
		Secondname: p.Secondname,
		Surname:    p.Surname,
		Fullname:   p.Fullname,
		Birthdate:  p.Birthdate,
		Gender:     gender,
		City:       p.City,
		Interests:  p.Interests,
		UserId:     p.UserId,
	}, nil
}
