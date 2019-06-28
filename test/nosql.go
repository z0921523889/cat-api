package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"

	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/postgresstore"

	_ "github.com/lib/pq"
)

var sessionManager *scs.SessionManager

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/cat?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize a new session manager and configure it to use PostgreSQL as
	// the session store.
	sessionManager = scs.New()
	sessionManager.Store = postgresstore.New(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/put", putHandler)
	mux.HandleFunc("/get", getHandler)

	http.ListenAndServe(":8000", sessionManager.LoadAndSave(mux))
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager.Put(r.Context(), "message", "Hello from a session!")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	msg := sessionManager.GetString(r.Context(), "message")
	io.WriteString(w, msg)
}
