package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/gorilla/securecookie"
	"time"
	"net/http"
	"errors"
)

type PGStore struct {
	Codecs []securecookie.Codec
	Options *Options
	Path string
	DBPool *sqlx.DB
}

type PGSession struct {
	ID int64
	Key string
	Data string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

func NewPGStoreFromPool(db *sqlx.DB, keyPairs ...[]byte) (*PGStore, error) {
	store := &PGStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &Options{
			Path: "/",
			MaxAge: 86400 * 30,
		},
		DBPool: db,
	}
	store.createSessionTable()
	return store, nil
}

func (store *PGStore) createSessionTable() error {
	stmt := `
		CREATE TABLE IF NOT EXISTS http_sessions (
			id serial,
			key bytea,
			data bytea,
			created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
			updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
			expires_at timestamp without time zone 
		)
	`
	store.DBPool.MustExec(stmt)
	return nil
}

func (store *PGStore) New(r *http.Request, name string) (*Session, error) {
	session := &Session{
		Values: make(map[string]interface{}),
		name: name,
		Options: new(Options),
	}
	options := *store.Options
	session.Options = &(options)
	session.IsNew = true

	var err error
	if cookie, cookieErr := r.Cookie(name); cookieErr == nil {
		err = securecookie.DecodeMulti(name, cookie.Value, &session.ID, store.Codecs...)
		if err == nil {
			err = store.load(session)
			if err == nil {
				session.IsNew = false
			} 
		}
	}

	store.MaxAge(store.Options.MaxAge)

	return session, err
}

func (store *PGStore) load(session *Session) error {
	var s PGSession
	err := store.selectOne(&s, session.ID)
	if err != nil {
		return err
	}
	return securecookie.DecodeMulti(session.Name(), string(s.Data), &session.Values, store.Codecs...)
}

func (store *PGStore) selectOne(session *PGSession, key string) error {
	err := store.DBPool.Get(session, "SELECT * FROM http_sessions WHERE key = $1", key)
	if err != nil {
		return errors.New("Unable to find session in the database")
	}
	return nil
}

func (store *PGStore) MaxAge(age int) {
	store.Options.MaxAge = age
	for _, codec := range store.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

func (store *PGStore) Get(r *http.Request, name string) (*Session, error) {
	return GetRegistry(r).Get(*store, name)
}

func (store *PGStore) Save(r *http.Request, w http.ResponseWriter, session *Session) error {
	
}