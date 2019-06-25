package orm

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	driverName = "postgres"
	dbName     = "cat"
)

var host, port, user, password string
var Engine *gorm.DB

func init() {
	initialEnv()
	//createDatabase()
	connectDBEngine()
	Engine.AutoMigrate(new(Users), new(Cat), new(CatThumbnails))
	Engine.LogMode(true)
}

func initialEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	} else {
		log.Println("loading .env file success")
	}
	host = os.Getenv("POSTGRES_HOST")
	port = os.Getenv("POSTGRES_PORT")
	user = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
}

func createDatabase() {
	db, err := sql.Open(driverName, fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, password, ))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)", dbName)
	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		panic(err)
	}
	if exists {
		log.Println("DATABASE EXIST,CONNECT DATABASE")
	} else {
		log.Println("DATABASE DOES NOT EXIST,CREATE DATABASE")
		_, err = db.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			panic(err)
		}
	}
}

func connectDBEngine() {
	var err error
	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	Engine, err = gorm.Open(driverName, dataSource)
	if err != nil {
		panic(err)
	}
}
