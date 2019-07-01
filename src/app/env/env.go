package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	ApplicationHost   string
	ApplicationPort   string
	PostgresHost      string
	PostgresPort      string
	PostgresUser      string
	PostgresPassword  string
	ScheduleJobEnable bool
	SwaggerDocEnable  bool
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	} else {
		log.Println("loading .env file success")
	}
	ApplicationHost = os.Getenv("APPLICATION_HOST")
	ApplicationPort = os.Getenv("APPLICATION_PORT")
	PostgresHost = os.Getenv("POSTGRES_HOST")
	PostgresPort = os.Getenv("POSTGRES_PORT")
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	ScheduleJobEnable, _ = strconv.ParseBool(os.Getenv("SCHEDULE_JOB_ENABLE"))
	SwaggerDocEnable, _ = strconv.ParseBool(os.Getenv("SWAGGER_DOC_ENABLE"))
}
