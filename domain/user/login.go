package user

type Login string

func NewLogin(userLogin string) (Login, error) {
	return Login(userLogin), nil
}

func (l Login) String() string {
	return string(l)
}
