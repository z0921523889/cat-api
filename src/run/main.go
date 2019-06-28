package main

import (
	_ "cat-api/docs"
	"cat-api/src/app/db"
	"cat-api/src/app/env"
	"cat-api/src/app/orm"
	"cat-api/src/app/router"
	"cat-api/src/app/session"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
)

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://github.com/EDDYCJY/go-gin-example
// @license.name MIT
// @license.url https://github.com/EDDYCJY/go-gin-example/blob/master/LICENSE
func main() {
	env.LoadEnv()
	db.InitDataBase()
	orm.ConnectDBEngine()
	defer orm.CloseDBEngine()
	session.ConnectSessionEngine()
	defer session.CloseSessionEngine()
	routerEngine := router.Router
	routerEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := routerEngine.Run(":8085")
	if err != nil {
		log.Println(err)
	}
}

