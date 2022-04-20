package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Personal struct {
	Id         int
	Firstname  string
	Secondname string
	Surname    string
	Fullname   string
	Birthdate  time.Time
	Gender     string
	City       string
	Interests  string
	UserId     int
}

type User struct {
	Id        int
	Login     string
	Password  string
	Email     string
	Created   time.Time
	LastLogin time.Time
}

type Friend struct {
	Id          int
	Personal_id int
	Friend_id   int
	State       string
	Created     time.Time
}
