package usecase

import (
	"AWesomeSocial/domain/friend"
)

func (app *Application) GetFriends(personId int) ([]*friend.Friend, error) {
	friends, err := app.repositoryRegistry.GetFriendReadRepo().GetFriendsList(personId)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

func (app *Application) GetIncomingFriends(personId int) ([]*friend.Friend, error) {
	incomingFriends, err := app.repositoryRegistry.GetFriendReadRepo().GetIncomingList(personId)
	if err != nil {
		return nil, err
	}

	return incomingFriends, nil
}
