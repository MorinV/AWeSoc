package main

import (
	"net/http"
	"path/filepath"
)

func (s *service) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.home)
	mux.HandleFunc("/personal", s.showPersonalPage)
	mux.HandleFunc("/register", s.register)
	mux.HandleFunc("/registerForm", s.showRegisterForm)
	mux.HandleFunc("/editPersonForm", s.showEditPersonForm)
	mux.HandleFunc("/editPerson", s.editPerson)
	mux.HandleFunc("/loginForm", s.showLoginForm)
	mux.HandleFunc("/login", s.login)
	mux.HandleFunc("/logout", s.logout)
	mux.HandleFunc("/addFriend", s.addFriend)
	mux.HandleFunc("/friendlist", s.showFriendsList)
	mux.HandleFunc("/search", s.search)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
