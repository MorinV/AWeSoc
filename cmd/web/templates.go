package main

import (
	"AWesomeSocial/pkg/models"
	"AWesomeSocial/pkg/models/mysql"
	"html/template"
	"path/filepath"
)

type templateData struct {
	Personal        *models.Personal
	Personals       []*models.Personal
	FriendPersonal  *mysql.FriendPersonal
	Friends         []*mysql.FriendPersonal
	IncomingFriends []*mysql.FriendPersonal
	User            User
	Flashes         []interface{}
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
