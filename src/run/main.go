package main

import (
	_ "cat-api/src/app/orm"
	router2 "cat-api/src/app/router"
	"log"
)

func main() {
	router := router2.Router
	err := router.Run(":8085")
	if err != nil {
		log.Println(err)
	}
}
