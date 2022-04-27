package main

import (
	"AWesomeSocial/domain/friend"
	"AWesomeSocial/domain/person"
	"AWesomeSocial/usecase"
	"html/template"
	"path/filepath"
)

type templateData struct {
	Person             *person.Person
	Persons            []*person.Person
	Friend             *friend.Friend
	Friends            []*friend.Friend
	IncomingFriends    []*friend.Friend
	User               User
	Flashes            []interface{}
	SearchPersonParams usecase.SearchPersonParams
	ErrorMessage       string
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
