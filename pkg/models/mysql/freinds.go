package mysql

import (
	"AWesomeSocial/pkg/models"
	"database/sql"
	"errors"
	"time"
)

type FriendsRepository struct {
	DB *sql.DB
}

var FriendsState = map[string]int{
	"pending":  1,
	"accepted": 2,
}

type FriendPersonal struct {
	*models.Friend
	*models.Personal
}

func (fr *FriendsRepository) GetFriendsList(personalId int) ([]*FriendPersonal, error) {
	stmt := `SELECT 
	f.id,
	f.personal_id,
	f.friend_id,
	f.state,
	f.created,
    p.id,
	p.fullname
FROM friends f INNER JOIN personal p on f.personal_id = p.id WHERE personal_id = ? ORDER BY state`
	rows, err := fr.DB.Query(stmt, personalId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*FriendPersonal
	for rows.Next() {
		friend := &FriendPersonal{&models.Friend{}, &models.Personal{}}
		err := rows.Scan(
			&friend.Friend.Id, &friend.Friend.Personal_id, &friend.Friend.Friend_id, &friend.Friend.State, &friend.Friend.Created, &friend.Personal.Id, &friend.Personal.Fullname)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}

func (fr *FriendsRepository) GetIncomingList(personalId int) ([]*FriendPersonal, error) {
	stmt := `SELECT 
	f.id,
	f.personal_id,
	f.friend_id,
	f.state,
	f.created,
    p.id,
	p.fullname
FROM friends f INNER JOIN personal p on f.personal_id = p.id WHERE friend_id = ? AND state='pending'`
	rows, err := fr.DB.Query(stmt, personalId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*FriendPersonal
	for rows.Next() {
		friend := &FriendPersonal{&models.Friend{}, &models.Personal{}}
		err := rows.Scan(
			&friend.Friend.Id, &friend.Friend.Personal_id, &friend.Friend.Friend_id, &friend.Friend.State, &friend.Friend.Created, &friend.Personal.Id, &friend.Personal.Fullname)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}

func (fr *FriendsRepository) GetFriend(personalId, friendId int) (*FriendPersonal, error) {
	stmt := `
SELECT 
	f.id,
	f.personal_id,
	f.friend_id,
	f.state,
	f.created,
    p.id,
	p.fullname
FROM friends f INNER JOIN personal p on f.personal_id = p.id WHERE personal_id = ? AND friend_id = ?`

	result := &FriendPersonal{&models.Friend{}, &models.Personal{}}
	err := fr.DB.QueryRow(stmt, personalId, friendId).
		Scan(&result.Friend.Id, &result.Friend.Personal_id, &result.Friend.Friend_id, &result.Friend.State, &result.Friend.Created, &result.Personal.Id, &result.Personal.Fullname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, models.ErrNoRecord
		}
		return result, err
	}
	return result, nil
}

func (fr *FriendsRepository) Insert(personalId, friendId, state int) error {
	stmt := `INSERT INTO friends (personal_id, friend_id, state, created) VALUES (?, ?, ?, ?)`
	_, err := fr.DB.Exec(stmt, personalId, friendId, state, time.Now().Format(time.RFC3339))

	return err
}

func (fr *FriendsRepository) UpdateState(personalId, friendId, state int) error {
	stmt := `UPDATE friends SET state = ? WHERE personal_id = ? AND friend_id = ?`
	_, err := fr.DB.Exec(stmt, state, personalId, friendId)

	return err
}
