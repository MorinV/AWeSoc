package main

import (
	"AWesomeSocial/pkg/models"
	"AWesomeSocial/pkg/models/mysql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	user := app.getUser(session)
	personals, err := app.rm.Repositories.PersonalRepository.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Personals: personals, User: user}

	app.render(w, r, "home.page.tmpl", data)
}

func (app *application) showPersonalPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	session, err := store.Get(r, CookieName)
	if err != nil {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	sessionUser := app.getUser(session)

	personal, err := app.rm.Repositories.PersonalRepository.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	friendPersonal := &mysql.FriendPersonal{}

	if sessionUser.Authenticated && sessionUser.UserId != personal.UserId {
		friendPersonal, err = app.rm.Repositories.FriendsRepository.GetFriend(sessionUser.UserId, personal.UserId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				friendPersonal.Personal = personal
			} else {
				app.serverError(w, err)
				return
			}
		}
	}

	data := &templateData{Personal: personal, User: sessionUser, FriendPersonal: friendPersonal}

	app.render(w, r, "personal.page.tmpl", data)
}

func (app *application) showRegisterForm(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, CookieName)
	if err != nil {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	sessionUser := app.getUser(session)
	if sessionUser.Authenticated {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data := &templateData{User: sessionUser}
	app.render(w, r, "register.page.tmpl", data)
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.render(w, r, "register.page.tmpl", nil)
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")
	email := r.FormValue("email")
	hashedPassword, err := app.hashPassword(password)
	if err != nil {
		app.serverError(w, err)
		return
	}
	tr, err := app.rm.Db.Begin()
	if err != nil {
		app.serverError(w, err)
		return
	}
	userId, err := app.rm.Repositories.UsersRepository.Insert(login, hashedPassword, email)
	if err != nil {
		tr.Rollback()
		app.serverError(w, err)
		return
	}

	firstname := r.FormValue("firstname")
	secondname := r.FormValue("secondname")
	surname := r.FormValue("surname")
	birthdate := r.FormValue("birthdate")
	gender := r.FormValue("gender")
	city := r.FormValue("city")
	interests := r.FormValue("interests")

	id, err := app.rm.Repositories.PersonalRepository.Insert(firstname, secondname, surname, birthdate, gender, city, interests, userId)
	if err != nil {
		tr.Rollback()
		app.serverError(w, err)
		return
	}
	if err = tr.Commit(); err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/personal?id=%d", id), http.StatusSeeOther)
}

func (app *application) showLoginForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	session, err := store.Get(r, CookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := app.getUser(session)
	if user.Authenticated {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data := &templateData{User: user, Flashes: session.Flashes()}
	app.render(w, r, "login.page.tmpl", data)
	return
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/loginForm", http.StatusSeeOther)
		return
	}

	session, err := store.Get(r, CookieName)
	if err != nil {
		app.serverError(w, err)
		return
	}

	sessionUser := app.getUser(session)
	if sessionUser.Authenticated {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	if password == "" || login == "" {
		session.AddFlash("Must enter a code")
		err = session.Save(r, w)
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.clientError(w, http.StatusForbidden)
		return
	}

	user, err := app.rm.Repositories.UsersRepository.GetByLogin(login)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if auth := app.verifyPassword(password, user.Password); !auth {
		app.clientError(w, http.StatusForbidden)
		return
	}

	personal, err := app.rm.Repositories.PersonalRepository.GetByUserId(user.Id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	sessionUser = User{
		Username:      user.Login,
		Authenticated: true,
		UserId:        user.Id,
		PersonalId:    personal.Id,
		Fullname:      personal.Fullname,
	}
	session.Values["user"] = sessionUser
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
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

func (app *application) addFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, CookieName)
	if err != nil {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	sessionUser := app.getUser(session)
	if auth := sessionUser.Authenticated; !auth {
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	t := struct {
		Friend_id *string `json:"friend_id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&t)
	if err != nil {
		app.serverError(w, err)
		return
	}
	friendId, err := strconv.Atoi(*t.Friend_id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	friend, err := app.rm.Repositories.FriendsRepository.GetFriend(friendId, sessionUser.PersonalId)
	if err != nil {
		if errors.Is(models.ErrNoRecord, err) {
			err = app.rm.Repositories.FriendsRepository.Insert(sessionUser.PersonalId, friendId, mysql.FriendsState["pending"])
			if err != nil {
				app.serverError(w, err)
				return
			}
		} else {
			app.serverError(w, err)
			return
		}
	} else {
		err = app.rm.Repositories.FriendsRepository.Insert(sessionUser.PersonalId, friendId, mysql.FriendsState["accepted"])
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.rm.Repositories.FriendsRepository.UpdateState(friend.Personal_id, friend.Friend_id, mysql.FriendsState["accepted"])
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) showFriendsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, CookieName)
	if err != nil {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	sessionUser := app.getUser(session)
	if auth := sessionUser.Authenticated; !auth {
		app.clientError(w, http.StatusUnauthorized)
		return
	}

	friends, err := app.rm.Repositories.FriendsRepository.GetFriendsList(sessionUser.PersonalId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	friendsIncoming, err := app.rm.Repositories.FriendsRepository.GetIncomingList(sessionUser.PersonalId)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := &templateData{User: sessionUser, Friends: friends, IncomingFriends: friendsIncoming}
	app.render(w, r, "friendslist.page.tmpl", data)
	return
}
