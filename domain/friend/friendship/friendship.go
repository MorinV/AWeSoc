package friendship

import "time"

type Friendship struct {
	Id       int
	PersonId int
	FriendId int
	State    State
	Created  time.Time
}

func New(id int, personalId int, friendId int, state State) *Friendship {
	return &Friendship{Id: id, PersonId: personalId, FriendId: friendId, State: state}
}

func (f *Friendship) AcceptFriendship() error {
	state, err := NewStateFromName("accepted")
	if err != nil {
		f.State = state
	}

	return err
}
