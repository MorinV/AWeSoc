package friend

import (
	"AWesomeSocial/domain/friend"
	"AWesomeSocial/domain/friend/friendship"
	"time"
)

type DbValue struct {
	Id             int
	State          string
	Created        time.Time
	PersonFullname string
	PersonId       int
	FriendFullname string
	FriendId       int
}

func (f DbValue) Hydrate() (*friend.Friend, error) {
	personDetails := friend.PersonDetails{PersonId: f.PersonId, Fullname: f.PersonFullname}
	friendDetails := friend.PersonDetails{PersonId: f.FriendId, Fullname: f.FriendFullname}
	state, err := friendship.NewStateFromName(f.State)
	if err != nil {
		return nil, err
	}
	fs := &friendship.Friendship{Id: f.Id, PersonId: f.PersonId, FriendId: f.FriendId, State: state, Created: f.Created}

	return &friend.Friend{Id: f.Id, Friendship: fs, PersonDetails: personDetails, FriendDetails: friendDetails}, nil
}
