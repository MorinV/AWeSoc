package internal

import (
	"AWesomeSocial/domain/friend"
	"AWesomeSocial/domain/person"
	"AWesomeSocial/domain/user"
	"errors"
)

type RepositoryRegistry interface {
	GetFriendReadRepo() friend.ReadRepo
	GetFriendWriteRepo() friend.WriteRepo
	GetUserReadRepo() user.ReadRepo
	GetUserWriteRepo() user.WriteRepo
	GetPersonReadRepo() person.ReadRepo
	GetPersonWriteRepo() person.WriteRepo
}

var ErrNoRecord = errors.New("models: подходящей записи не найдено")
