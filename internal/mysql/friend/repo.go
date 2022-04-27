package friend

import (
	"AWesomeSocial/domain/friend"
	"AWesomeSocial/internal"
	"database/sql"
	"errors"
	"time"
)

type WriteRepo struct {
	DB *sql.DB
}

type ReadRepo struct {
	DB *sql.DB
}

func (rr *ReadRepo) GetFriendsList(personalId int) ([]*friend.Friend, error) {
	stmt := `SELECT 
	f.id,
	f.state,
	f.created,
    person.id,
	person.fullname,
	friend.id,
	friend.fullname
FROM friends f 
    INNER JOIN personal person on f.personal_id = person.id 
    INNER JOIN personal friend on f.friend_id = friend.id
WHERE personal_id = ? ORDER BY state`
	rows, err := rr.DB.Query(stmt, personalId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*friend.Friend
	for rows.Next() {
		fdb := &DbValue{}
		err := rows.Scan(
			&fdb.Id, &fdb.State, &fdb.Created, &fdb.PersonId, &fdb.PersonFullname, &fdb.FriendId, &fdb.FriendFullname)
		if err != nil {
			return nil, err
		}
		f, err := fdb.Hydrate()
		if err != nil {
			return nil, err
		}
		friends = append(friends, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}

func (rr *ReadRepo) GetIncomingList(personalId int) ([]*friend.Friend, error) {
	stmt := `SELECT 
	f.id,
	f.state,
	f.created,
    person.id,
	person.fullname,
	friend.id,
	friend.fullname
FROM friends f 
    INNER JOIN personal friend on f.personal_id = friend.id 
    INNER JOIN personal person on f.friend_id = person.id
WHERE friend_id = ? AND state='pending'`
	rows, err := rr.DB.Query(stmt, personalId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*friend.Friend
	for rows.Next() {
		fdb := &DbValue{}
		err := rows.Scan(
			&fdb.Id, &fdb.State, &fdb.Created, &fdb.PersonId, &fdb.PersonFullname, &fdb.FriendId, &fdb.FriendFullname)
		if err != nil {
			return nil, err
		}
		f, err := fdb.Hydrate()
		if err != nil {
			return nil, err
		}
		friends = append(friends, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}

func (rr *ReadRepo) GetFriend(personalId, friendId int) (*friend.Friend, error) {
	stmt := `SELECT 
	f.id,
	f.state,
	f.created,
    person.id,
	person.fullname,
	friend.id,
	friend.fullname
FROM friends f 
    INNER JOIN personal person on f.personal_id = person.id 
    INNER JOIN personal friend on f.friend_id = friend.id
WHERE personal_id = ? AND friend_id = ?`

	fdb := &DbValue{}
	err := rr.DB.QueryRow(stmt, personalId, friendId).
		Scan(&fdb.Id, &fdb.State, &fdb.Created, &fdb.PersonId, &fdb.PersonFullname, &fdb.FriendId, &fdb.FriendFullname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internal.ErrNoRecord
		}
		return nil, err
	}
	result, err := fdb.Hydrate()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rr *ReadRepo) FindFriend(personalId, friendId int) (*friend.Friend, error) {
	result, err := rr.GetFriend(personalId, friendId)
	if errors.Is(err, internal.ErrNoRecord) {
		return nil, nil
	}
	return result, err
}

func (wr *WriteRepo) Add(f *friend.Friend) (int, error) {
	if f.Id == 0 {
		id, err := wr.insert(f)
		if err != nil {
			return 0, err
		}

		return id, nil
	} else {
		err := wr.update(f)

		if err != nil {
			return 0, err
		}

		return f.Id, err
	}
}

func (wr *WriteRepo) insert(f *friend.Friend) (int, error) {
	stmt := `INSERT INTO friends (personal_id, friend_id, state, created) VALUES (?, ?, ?, ?)`
	result, err := wr.DB.Exec(stmt, f.Friendship.PersonId, f.Friendship.FriendId, f.Friendship.State.Int(), time.Now().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (wr *WriteRepo) update(f *friend.Friend) error {
	stmt := `UPDATE friends SET state = ?, personal_id = ?, friend_id = ? WHERE id = ?`
	_, err := wr.DB.Exec(stmt, f.Friendship.State.Int(), f.Friendship.PersonId, f.Friendship.FriendId, f.Id)

	return err
}

func (wr *WriteRepo) UpdateState(personalId, friendId, state int) error {
	stmt := `UPDATE friends SET state = ? WHERE personal_id = ? AND friend_id = ?`
	_, err := wr.DB.Exec(stmt, state, personalId, friendId)

	return err
}
