package app

import (
	"cat-api/src/app/conf"
	"cat-api/src/app/db"
	_ "cat-api/src/app/docs"
	"cat-api/src/app/env"
	"cat-api/src/app/orm"
	"cat-api/src/app/router"
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
func StartServer() {
	env.LoadEnv()
	db.InitDataBase()
	orm.ConnectDBEngine()
	defer orm.CloseDBEngine()
	conf.CheckDataBaseConfig()
	routerEngine := router.InitialRouterEngine()
	if env.SwaggerDocEnable {
		routerEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	err := routerEngine.Run(":" + env.ApplicationPort)
	if err != nil {
		log.Println(err)
	}
}
