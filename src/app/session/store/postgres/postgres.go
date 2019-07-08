package postgres

import (
	"cat-api/src/app/orm"
	"encoding/base32"
	"errors"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"net/http"
	"strings"
	"time"
)

// Amount of time for cookies/redis keys to expire.
var sessionExpire = 86400 * 30

type SessionPostgresStore struct {
	Codecs        []securecookie.Codec
	Options       *sessions.Options // default configuration
	DefaultMaxAge int               // default Redis TTL for a MaxAge == 0 session
	maxLength     int
	keyPrefix     string
	serializer    SessionSerializer
}

func (store *SessionPostgresStore) SetMaxLength(l int) {
	if l >= 0 {
		store.maxLength = l
	}
}

func (store *SessionPostgresStore) SetKeyPrefix(p string) {
	store.keyPrefix = p
}

func (store *SessionPostgresStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(store, name)
}

func (store *SessionPostgresStore) New(r *http.Request, name string) (*sessions.Session, error) {
	var (
		err error
		ok  bool
	)
	session := sessions.NewSession(store, name)
	// make a copy
	options := *store.Options
	session.Options = &options
	session.IsNew = true
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, store.Codecs...)
		if err == nil {
			ok, err = store.load(session)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	return session, err
}

func (store *SessionPostgresStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Marked for deletion.
	if session.Options.MaxAge <= 0 {
		//clean
		if err := store.delete(session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		// Build an alphanumeric key for the redis store.
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		if err := store.save(session); err != nil {
			return err
		}
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, store.Codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}
	return nil
}

// load reads the session from postgres.
// returns true if there is a sessoin data in DB
func (store *SessionPostgresStore) load(session *sessions.Session) (bool, error) {
	var count int
	var sessionData orm.Session
	if err := orm.Engine.
		Where("token = ?", store.keyPrefix+session.ID).
		Where("expiry > ?", time.Now()).
		Find(&sessionData).Count(&count).Error; err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil // no data was associated with this key
	}
	return true, store.serializer.Deserialize(sessionData.Data, session)
}

// save stores the session in redis.
func (store *SessionPostgresStore) save(session *sessions.Session) error {
	b, err := store.serializer.Serialize(session)
	if err != nil {
		return err
	}
	if store.maxLength != 0 && len(b) > store.maxLength {
		return errors.New("SessionStore: the value to store is too big")
	}
	age := session.Options.MaxAge
	if age == 0 {
		age = store.DefaultMaxAge
	}
	token := store.keyPrefix + session.ID
	expiry := time.Now().Add(time.Second * time.Duration(age))
	var sessionData orm.Session
	if orm.Engine.First(&sessionData, "token = ?", token).RecordNotFound() {
		sessionData = orm.Session{
			Token:  store.keyPrefix + session.ID,
			Data:   b,
			Expiry: expiry,
		}
		if err := orm.Engine.Create(&sessionData).Error; err != nil {
			return err
		}
	} else {
		sessionData.Data = b
		sessionData.Expiry = expiry
		if err := orm.Engine.Save(&sessionData).Error; err != nil {
			return err
		}
	}
	return nil
}

// delete removes keys from redis if MaxAge<0
func (store *SessionPostgresStore) delete(session *sessions.Session) error {
	if err := orm.Engine.Where("token LIKE ?", store.keyPrefix+session.ID).Delete(orm.Session{}).Error; err != nil {
		return err
	}
	return nil
}

// ping does an internal ping against a server to check if it is alive.
func (store *SessionPostgresStore) ping() (bool, error) {
	err := orm.Engine.DB().Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}

func newPostgresStore(keyPairs ...[]byte) (*SessionPostgresStore, error) {
	rs := &SessionPostgresStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: sessionExpire,
		},
		DefaultMaxAge: 60 * 20, // 20 minutes seems like a reasonable default
		maxLength:     4096,
		keyPrefix:     "session_",
		serializer:    GobSerializer{},
	}
	_, err := rs.ping()
	return rs, err
}
