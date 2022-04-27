package usecase

import (
	"AWesomeSocial/domain/friend"
	"AWesomeSocial/internal"
	"errors"
)

func (app *Application) GetFriend(personId, friendId int) (*friend.Friend, error) {
	f, err := app.repositoryRegistry.GetFriendReadRepo().GetFriend(personId, friendId)
	if err != nil {
		if errors.Is(err, internal.ErrNoRecord) {
			f = &friend.Friend{}
			f.PersonDetails = friend.PersonDetails{PersonId: personId}
			f.FriendDetails = friend.PersonDetails{PersonId: friendId}
		} else {
			return nil, err
		}
	}

	return f, nil
}
