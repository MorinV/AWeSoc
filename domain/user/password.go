package user

import "golang.org/x/crypto/bcrypt"

type Password string

func NewPassword(userPassword string) (Password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.MinCost)
	return Password(hash), err
}

func (p Password) Equals(userPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(userPassword))
	if err != nil {
		return false
	}

	return true
}
