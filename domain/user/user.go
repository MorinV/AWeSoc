package user

import "time"

type User struct {
	Id        int
	Login     Login
	Password  Password
	Email     Email
	Created   time.Time
	LastLogin time.Time
}

func New(id int, login Login, password Password, email Email) *User {
	return &User{Id: id, Login: login, Password: password, Email: email, Created: time.Now(), LastLogin: time.Now()}
}

func (u *User) UpdateLastLogin() {
	u.LastLogin = time.Now()
}
