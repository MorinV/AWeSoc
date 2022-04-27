package usecase

import (
	"AWesomeSocial/domain/friend"
	"AWesomeSocial/domain/friend/friendship"
	"AWesomeSocial/internal"
)

func (app *Application) AddFriend(personId, friendId int) error {
	f, err := makeNewFriend(personId, friendId)
	if err != nil {
		return err
	}
	incomingFriend, err := app.repositoryRegistry.GetFriendReadRepo().GetFriend(friendId, personId)
	switch err {
	case internal.ErrNoRecord:
		_, err := app.repositoryRegistry.GetFriendWriteRepo().Add(f)
		if err != nil {
			return err
		}
		break
	case nil:
		err = f.Friendship.AcceptFriendship()
		if err != nil {
			return err
		}
		_, err = app.repositoryRegistry.GetFriendWriteRepo().Add(f)
		if err != nil {
			return err
		}

		err := incomingFriend.Friendship.AcceptFriendship()
		if err != nil {
			return err
		}
		_, err = app.repositoryRegistry.GetFriendWriteRepo().Add(incomingFriend)
		if err != nil {
			return err
		}
	default:
		return err
	}

	return nil
}

func makeNewFriend(personId int, friendId int) (*friend.Friend, error) {
	state, err := friendship.NewStateFromName("pending")
	if err != nil {
		return nil, err
	}
	fs := friendship.New(0, personId, friendId, state)
	f := friend.New(0, fs, friend.PersonDetails{}, friend.PersonDetails{})
	return f, nil
}