package main

import (
	"AWesomeSocial/domain/friend"
	"AWesomeSocial/internal"
	"AWesomeSocial/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (s *service) home(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Path != "/" {
		s.notFound(w)
		return
	}

	user := s.getUser(session)
	personals, err := s.app.GetLatestRegisteredPersons()
	if err != nil {
		s.serverError(w, err)
		return
	}

	data := &templateData{Persons: personals, User: user}

	s.render(w, r, "home.page.tmpl", data)
}

func (s *service) showPersonalPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		s.notFound(w)
		return
	}

	session, err := store.Get(r, CookieName)
	if err != nil {
		s.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	sessionUser := s.getUser(session)

	person, err := s.app.GetPerson(id)
	if err != nil {
		if errors.Is(err, internal.ErrNoRecord) {
			s.notFound(w)
		} else {
			s.serverError(w, err)
		}
		return
	}

	f := &friend.Friend{}

	if sessionUser.Authenticated && sessionUser.UserId != person.UserId {
		f, err = s.app.GetFriend(sessionUser.PersonId, person.Id)
		if err != nil {
			s.serverError(w, err)
			return
		}
	}

	data := &templateData{Person: person, User: sessionUser, Friend: f}

	s.render(w, r, "person.page.tmpl", data)
}

func (s *service) showRegisterForm(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, CookieName)
	if err != nil {
		s.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	sessionUser := s.getUser(session)
	if sessionUser.Authenticated {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data := &templateData{User: sessionUser}
	s.render(w, r, "register.page.tmpl", data)
}

func (s *service) register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.render(w, r, "register.page.tmpl", nil)
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")
	email := r.FormValue("email")

	_, err := s.app.PersistUser(0, login, password, email)
	if err != nil {
		s.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/editPersonForm", http.StatusSeeOther)
}

func (s *service) showEditPersonForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	session, err := store.Get(r, CookieName)
	if err != nil {
		s.serverError(w, err)
		return
	}
	sessionUser := s.getUser(session)
	if !sessionUser.Authenticated {
		s.clientError(w, http.StatusUnauthorized)
		return
	}
	person, err := s.app.GetPerson(sessionUser.PersonId)
	if err != nil {
		if errors.Is(internal.ErrNoRecord, err) {
			person = nil
		} else {
			s.serverError(w, err)
			return
		}
	}

	data := &templateData{User: sessionUser, Person: person}
	s.render(w, r, "editPerson.page.tmpl", data)
}

func (s *service) editPerson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	session, err := store.Get(r, CookieName)
	if err != nil {
		s.serverError(w, err)
		return
	}
	sessionUser := s.getUser(session)
	if !sessionUser.Authenticated {
		s.clientError(w, http.StatusUnauthorized)
		return
	}

	firstname := r.FormValue("firstname")
	secondname := r.FormValue("secondname")
	surname := r.FormValue("surname")
	birthdate := r.FormValue("birthdate")
	gender := r.FormValue("gender")
	city := r.FormValue("city")
	interests := r.FormValue("interests")

	id, err := s.app.EditUserPerson(firstname, secondname, surname, birthdate, gender, city, interests, sessionUser.UserId)
	if err != nil {
		s.serverError(w, err)
		return
	}
	sessionUser.PersonId = id
	sessionUser.Fullname = surname + " " + firstname + " " + secondname
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/personal?id=%d", id), http.StatusSeeOther)
}

func (s *service) showLoginForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	session, err := store.Get(r, CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := s.getUser(session)
	if user.Authenticated {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data := &templateData{User: user, Flashes: session.Flashes()}
	s.render(w, r, "login.page.tmpl", data)
	return
}

func (s *service) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/loginForm", http.StatusSeeOther)
		return
	}

	session, err := store.Get(r, CookieName)
	if err != nil {
		s.serverError(w, err)
		return
	}

	sessionUser := s.getUser(session)
	if sessionUser.Authenticated {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	auth, err := s.app.AuthUser(login, password)
	if err != nil {
		s.serverError(w, err)
		return
	}

	sessionUser = User{
		Username:      auth.User.Login.String(),
		Authenticated: true,
		UserId:        auth.User.Id,
		PersonId:      auth.Person.Id,
		Fullname:      auth.Person.Fullname,
	}

	session.Values["user"] = sessionUser
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *service) logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = User{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *service) addFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, CookieName)
	if err != nil {
		s.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	sessionUser := s.getUser(session)
	if auth := sessionUser.Authenticated; !auth {
		s.clientError(w, http.StatusUnauthorized)
		return
	}
	t := struct {
		FriendId *string `json:"friend_id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&t)
	if err != nil {
		s.serverError(w, err)
		return
	}
	friendId, err := strconv.Atoi(*t.FriendId)
	if err != nil {
		s.serverError(w, err)
		return
	}
	err = s.app.AddFriend(sessionUser.PersonId, friendId)

	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "service/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		s.serverError(w, err)
		return
	}
}

func (s *service) showFriendsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, CookieName)
	if err != nil {
		s.serverError(w, err)
		return
	}

	sessionUser := s.getUser(session)
	if auth := sessionUser.Authenticated; !auth {
		s.clientError(w, http.StatusUnauthorized)
		return
	}

	friends, err := s.app.GetFriends(sessionUser.PersonId)
	if err != nil {
		s.serverError(w, err)
		return
	}

	friendsIncoming, err := s.app.GetIncomingFriends(sessionUser.PersonId)
	if err != nil {
		s.serverError(w, err)
		return
	}
	data := &templateData{User: sessionUser, Friends: friends, IncomingFriends: friendsIncoming}
	s.render(w, r, "friendslist.page.tmpl", data)
	return
}

func (s *service) search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, CookieName)
	if err != nil {
		s.serverError(w, err)
		return
	}

	sessionUser := s.getUser(session)

	params := usecase.SearchPersonParams{
		FirstNamePref: r.URL.Query().Get("firstNamePref"),
		SurNamePref:   r.URL.Query().Get("surNamePref"),
	}

	persons, err := s.app.SearchPersons(params)
	if err != nil {
		s.infoLog.Println(err)
	}

	data := &templateData{User: sessionUser, SearchPersonParams: params, Persons: persons}
	s.render(w, r, "search.page.tmpl", data)
	return
}
