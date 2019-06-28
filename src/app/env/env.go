package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	Port     string
	User     string
	Password string
	Host     string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	} else {
		log.Println("loading .env file success")
	}
	Host = os.Getenv("POSTGRES_HOST")
	Port = os.Getenv("POSTGRES_PORT")
	User = os.Getenv("POSTGRES_USER")
	Password = os.Getenv("POSTGRES_PASSWORD")
}
