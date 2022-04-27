package person

import "time"

type Person struct {
	Id         int
	Firstname  string
	Secondname string
	Surname    string
	Fullname   string
	Birthdate  time.Time
	Gender     Gender
	City       string
	Interests  string
	UserId     int
}

func New(id int, firstname string, secondname string, surname string, birthdate time.Time, gender Gender, city string, interests string, userId int) *Person {
	fullname := surname + " " + firstname + " " + secondname
	return &Person{Id: id, Firstname: firstname, Secondname: secondname, Surname: surname, Fullname: fullname, Birthdate: birthdate, Gender: gender, City: city, Interests: interests, UserId: userId}
}

func (p *Person) Update(firstname string, secondname string, surname string, birthdate time.Time, gender Gender, city string, interests string, userId int) error {
	p.Firstname = firstname
	p.Secondname = secondname
	p.Surname = surname
	p.Fullname = surname + " " + firstname + " " + secondname
	p.Birthdate = birthdate
	p.Gender = gender
	p.City = city
	p.Interests = interests
	p.UserId = userId

	return nil
}
