package user

type WriteRepo interface {
	Add(u *User) (int, error)
}
type ReadRepo interface {
	GetByLogin(login Login) (*User, error)
}
