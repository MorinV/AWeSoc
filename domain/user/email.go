package user

type Email string

func NewEmail(userEmail string) (Email, error) {
	return Email(userEmail), nil
}
