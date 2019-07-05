package orm

import (
	"cat-api/src/app/env"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	driverName = "postgres"
	dbName     = "cat"
)

var Engine *gorm.DB

func ConnectDBEngine() {
	var err error
	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", env.PostgresHost, env.PostgresPort, env.PostgresUser, env.PostgresPassword, dbName)
	Engine, err = gorm.Open(driverName, dataSource)
	if err != nil {
		panic(err)
	}
	Engine.DB().SetMaxIdleConns(10)
	Engine.AutoMigrate(
		new(Sessions), new(ApplicationConfigs),
		new(Admins),new(AdminProfiles),
		new(Users),
		new(Cats), new(CatThumbnails),
		new(AdminTimePeriodTemplates), new(AdoptionTimePeriods), new(AdoptionTimePeriodCatPivots),
		new(CatUserReservations),
		new(Wallets),
		)
	Engine.LogMode(true)
	CheckDefaultAdmin()
}

func CloseDBEngine() error {
	return Engine.Close()
}
