package orm

import (
	"database/sql"
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"log"
)

const (
	driverName = "postgres"
	dbName     = "cat"
)

var (
	host     = "192.168.99.100"
	port     = "5432"
	user     = "postgres"
	password = "password"
	//host     = os.Getenv("POSTGRES_HOST")
	//port     = os.Getenv("POSTGRES_PORT")
	//user     = os.Getenv("POSTGRES_USER")
	//password = os.Getenv("POSTGRES_PASSWORD")
)

func init() {
	createDatabase()
	dataSource := getDataSource()
	engine := getDBEngine(dataSource)
	err := engine.Sync2(new(User))
	if err != nil {
		log.Fatalln(err)
	}
}

func createDatabase() {
	db, err := sql.Open(driverName, fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, password, ))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("select count(*) from pg_catalog.pg_database where datname = '" + dbName+"'")
	if err != nil {
		log.Println(err)
		_, err = db.Exec("CREATE DATABASE '" + dbName+"'")
		if err != nil {
			panic(err)
		}
	}
}
func getDataSource() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
}

func getDBEngine(dataSource string) *xorm.Engine {
	engine, err := xorm.NewEngine(driverName, dataSource)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	err = engine.Ping()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println("connect postgresql success")
	return engine
}
