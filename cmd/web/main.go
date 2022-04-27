package main

import (
	"AWesomeSocial/internal/mysql"
	"AWesomeSocial/usecase"
	"database/sql"
	"encoding/gob"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
)

var store *sessions.FilesystemStore

type User struct {
	Username      string
	Authenticated bool
	UserId        int
	PersonId      int
	Fullname      string
}

const CookieName = "AWeSoc"

type service struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
	app           *usecase.Application
}

func main() {
	authKeyDefault := []byte("pQlGHWaQHYB6CNPmXm4lb+j4XRwpy+d2eAJa77Q8KRBS4UAm/4hMsFNiUfJRQNKs8ks4nGQCxuYTYUVz7lGgZQ==")
	encryptionKeyDefault := []byte("7dNa2RUGZYn+3RSqU0KZrCQu0Q1Pvxg=")
	addr := flag.String("addr", ":4000", "Сетевой аддрес")
	dsn := flag.String("dsn", "web:qwerty@/AWeSoc?parseTime=true", "DSN для MySQL")
	authKeyOne := flag.String("authKeyOne", string(authKeyDefault), "Ключ авторизации")
	encryptionKeyOne := flag.String("encryptionKey", string(encryptionKeyDefault), "Ключ расшифровки")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	store = sessions.NewFilesystemStore(
		"",
		[]byte(*authKeyOne),
		[]byte(*encryptionKeyOne),
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(User{})

	repoRegistry := mysql.New(db)

	s := &service{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
		app:           usecase.New(repoRegistry),
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  s.routes(),
	}

	infoLog.Printf("Запуск веб-сервера на http://localhost%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
