package orm

import (
	"cat-api/src/app/format"
	"log"
)

var list = [][]string{
	{"2006-01-02T09:30:00", "2006-01-02T10:30:00"},
	{"2006-01-02T10:30:00", "2006-01-02T11:30:00"},
	{"2006-01-02T13:30:00", "2006-01-02T14:30:00"},
	{"2006-01-02T14:30:00", "2006-01-02T15:30:00"},
	{"2006-01-02T15:30:00", "2006-01-02T16:30:00"},
	{"2006-01-02T16:30:00", "2006-01-02T17:30:00"},
}

func CheckDefaultAdmin() {
	if Engine.Where("account = ?", "system").First(&Admin{}).RecordNotFound() {
		session := Engine.Begin()
		wallet := Wallet{
			Coin:    -1,
			PetCoin: -1,
		}
		if err := session.Create(&wallet).Error; err != nil {
			log.Fatal("fail on create default wallet err : " + err.Error())
		}
		user := User{
			Phone:            "system",
			UserName:         "system",
			Password:         "system",
			SecurityPassword: "system",
			IdentifiedCode:   "system",
			Status:           1,
			Wallet:           wallet,
		}
		if err := session.Create(&user).Error; err != nil {
			log.Fatal("fail on create default user err : " + err.Error())
		}
		admin := Admin{
			Account:  "system",
			Password: "system",
			User:     user,
		}
		if err := session.Create(&admin).Error; err != nil {
			log.Fatal("fail on create default admin err : " + err.Error())
		}
		if err := session.Commit().Error; err != nil {
			log.Fatal("fail on commit default admin")
		}
	}
}

func CheckDefaultAdminTimePeriodTemplate() {
	var admin Admin
	if err := Engine.Where("account = ?", "system").Preload("AdminTimePeriodTemplates").First(&admin).Error; err != nil {
		log.Fatal("fail on find default admin")
	}
	if len(admin.AdminTimePeriodTemplates) == 0 {
		session := Engine.Begin()
		for _, value := range list {
			start, _ := format.ParseTime(value[0], format.TimeFormatter)
			end, _ := format.ParseTime(value[1], format.TimeFormatter)
			adminTimePeriodTemplate := AdminTimePeriodTemplate{
				AdminId: admin.ID,
				StartAt: start,
				EndAt:   end,
			}
			if err := session.Create(&adminTimePeriodTemplate).Error; err != nil {
				session.Rollback()
				log.Fatal("fail on create default admin time period template")
			}
		}
		if err := session.Commit().Error; err != nil {
			log.Fatal("fail on commit default admin time period template")
		}
	}
}
