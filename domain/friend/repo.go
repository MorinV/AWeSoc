package friend

type WriteRepo interface {
	Add(f *Friend) (int, error)
	UpdateState(personalId, friendId, state int) error
}
type ReadRepo interface {
	GetFriendsList(personalId int) ([]*Friend, error)
	GetIncomingList(personalId int) ([]*Friend, error)
	GetFriend(personalId, friendId int) (*Friend, error)
	FindFriend(personalId, friendId int) (*Friend, error)
}
