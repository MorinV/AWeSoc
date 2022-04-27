package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"runtime/debug"
)

func (s *service) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s/n%s", err.Error(), debug.Stack())
	s.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (s *service) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (s *service) notFound(w http.ResponseWriter) {
	s.clientError(w, http.StatusNotFound)
}

func (s *service) render(w http.ResponseWriter, r *http.Request, name string, templateData *templateData) {
	ts, ok := s.templateCache[name]
	if !ok {
		s.serverError(w, fmt.Errorf("Шаблон %s не существует!", name))
		return
	}

	err := ts.Execute(w, templateData)
	if err != nil {
		s.serverError(w, err)
	}
}

func (s *service) getUser(sessions *sessions.Session) User {
	val := sessions.Values["user"]
	var user = User{}
	user, ok := val.(User)
	if !ok {
		return User{Authenticated: false}
	}
	return user
}

func (s *service) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash), err
}

func (s *service) verifyPassword(userPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))
	if err != nil {
		s.errorLog.Printf("Ошибка авторизации: %s", err)
		return false
	}

	return true
}
