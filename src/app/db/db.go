package db

import (
	"cat-api/src/app/env"
	"database/sql"
	"fmt"
	"log"
)

const (
	driverName        = "postgres"
	dbName            = "cat"
)

func InitDataBase() {
	db, err := sql.Open(driverName, fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", env.Host, env.Port, env.User, env.Password))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	createDatabase(db)
}

func createDatabase(db *sql.DB) {
	row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_database WHERE datname = $1)", dbName)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		panic(err)
	}
	if exists {
		log.Println("DATABASE EXIST")
	} else {
		log.Println("DATABASE DOES NOT EXIST,CREATE DATABASE")
		_, err = db.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			panic(err)
		}
	}
}