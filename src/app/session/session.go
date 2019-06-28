package session

import (
	"cat-api/src/app/env"
	"context"
	"database/sql"
	"fmt"
	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/postgresstore"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var session *scs.SessionManager
var db *sql.DB

const dbName = "cat"

func ConnectSessionEngine() {
	var err error
	dataSource := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", env.User, env.Password, env.Host, env.Port, dbName)
	db, err = sql.Open("postgres", dataSource)
	if err != nil {
		log.Fatal(err)
	}
	session = scs.New()
	session.Store = postgresstore.New(db)
	session.Cookie.Name = "cat_session"
	session.Lifetime = 24 * time.Hour
	//Session.IdleTimeout = 20 * time.Minute
	//Session.Cookie.Domain = "example.com"
	//Session.Cookie.HttpOnly = true
	//Session.Cookie.Path = "/example/"
	//Session.Cookie.Persist = true
	//Session.Cookie.SameSite = http.SameSiteStrictMode
	//Session.Cookie.Secure = true
}

func CloseSessionEngine() error {
	return db.Close()
}

func Put(ctx context.Context, key string, val interface{}) {
	session.Put(ctx, key, val)
}

func Get(ctx context.Context, key string) interface{} {
	return session.Get(ctx, key)
}
