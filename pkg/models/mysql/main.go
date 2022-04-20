package mysql

import (
	"database/sql"
)

type RepositoryManager struct {
	Db           *sql.DB
	Repositories *Repositories
}

type Repositories struct {
	PersonalRepository *PersonalRepository
	UsersRepository    *UsersRepository
	FriendsRepository  *FriendsRepository
}

func (rm *RepositoryManager) CreateRepositories() {
	rm.Repositories = &Repositories{
		PersonalRepository: &PersonalRepository{DB: rm.Db},
		UsersRepository:    &UsersRepository{DB: rm.Db},
		FriendsRepository:  &FriendsRepository{DB: rm.Db},
	}
}
