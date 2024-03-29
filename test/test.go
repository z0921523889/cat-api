package main

import (
	"cat-api/src/app/conf"
	"cat-api/src/app/db"
	"cat-api/src/app/env"
	"cat-api/src/app/orm"
)

func main() {
	env.LoadEnv()
	db.InitDataBase()
	orm.ConnectDBEngine()
	defer orm.CloseDBEngine()
	conf.CheckDataBaseConfig()
}
