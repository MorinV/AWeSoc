package mysql

import (
	"AWesomeSocial/domain/friend"
	"AWesomeSocial/domain/person"
	"AWesomeSocial/domain/user"
	friendmysql "AWesomeSocial/internal/mysql/friend"
	personmysql "AWesomeSocial/internal/mysql/person"
	usermysql "AWesomeSocial/internal/mysql/user"
	"database/sql"
)

const dateLayoutIso = "2006-01-02"

type RepositoryRegistry struct {
	masterDb *sql.DB
}

func New(masterDb *sql.DB) *RepositoryRegistry {
	return &RepositoryRegistry{masterDb: masterDb}
}

func (r *RepositoryRegistry) GetFriendReadRepo() friend.ReadRepo {
	return &friendmysql.ReadRepo{DB: r.masterDb}
}

func (r *RepositoryRegistry) GetFriendWriteRepo() friend.WriteRepo {
	return &friendmysql.WriteRepo{DB: r.masterDb}
}

func (r *RepositoryRegistry) GetUserReadRepo() user.ReadRepo {
	return &usermysql.ReadRepo{DB: r.masterDb}
}

func (r *RepositoryRegistry) GetUserWriteRepo() user.WriteRepo {
	return &usermysql.WriteRepo{DB: r.masterDb}
}

func (r *RepositoryRegistry) GetPersonReadRepo() person.ReadRepo {
	return &personmysql.ReadRepo{DB: r.masterDb}
}

func (r *RepositoryRegistry) GetPersonWriteRepo() person.WriteRepo {
	return &personmysql.WriteRepo{DB: r.masterDb}
}
