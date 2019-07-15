package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	Type               string
	ProjectName        string
	Version            string
	FullPathDomainName string
	DomainName         string
	ApplicationHost    string
	ApplicationPort    string
	PostgresHost       string
	PostgresPort       string
	PostgresUser       string
	PostgresPassword   string
	ScheduleJobEnable  bool
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	} else {
		log.Println("loading .env file success")
	}
	Type = os.Getenv("TYPE")
	ProjectName = os.Getenv("PROJECT_NAME")
	Version = os.Getenv("VERSION")
	FullPathDomainName = os.Getenv("FULL_PATH_DOMAIN_NAME")
	DomainName = os.Getenv("DOMAIN_NAME")
	ApplicationHost = os.Getenv("APPLICATION_HOST")
	ApplicationPort = os.Getenv("APPLICATION_PORT")
	PostgresHost = os.Getenv("POSTGRES_HOST")
	PostgresPort = os.Getenv("POSTGRES_PORT")
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	ScheduleJobEnable, _ = strconv.ParseBool(os.Getenv("SCHEDULE_JOB_ENABLE"))
}
