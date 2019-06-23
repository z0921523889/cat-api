package main

import (
	_ "cat-api/src/orm"
	"cat-api/src/router"
	"log"
)

func main() {
	router := router.Router
	err := router.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}
