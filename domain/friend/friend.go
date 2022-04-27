package friend

import "AWesomeSocial/domain/friend/friendship"

type Friend struct {
	Id            int
	Friendship    *friendship.Friendship
	PersonDetails PersonDetails
	FriendDetails PersonDetails
}

func New(id int, friendship *friendship.Friendship, personDetails PersonDetails, friendDetails PersonDetails) *Friend {
	return &Friend{Id: id, Friendship: friendship, PersonDetails: personDetails, FriendDetails: friendDetails}
}
