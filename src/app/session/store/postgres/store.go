package postgres

import (
	"github.com/gin-contrib/sessions"
	gsessions "github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

type Store struct {
	*SessionPostgresStore
}

func NewPostgresStore(keyPairs ...[]byte) (*Store, error) {
	s, err := newPostgresStore(keyPairs...)
	if err != nil {
		return nil, err
	}
	return &Store{s}, nil
}

func (store *Store) Options(options sessions.Options) {
	store.SessionPostgresStore.Options = &gsessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
