package person

type WriteRepo interface {
	Add(p *Person) (int, error)
}
type ReadRepo interface {
	Get(id int) (*Person, error)
	GetByUserId(id int) (*Person, error)
	Latest() ([]*Person, error)
	Search(firstNamePref, surNamePref string) ([]*Person, error)
}
