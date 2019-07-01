package orm

import (
	"cat-api/src/app/format"
	"github.com/jinzhu/gorm"
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
	admin := Admins{}
	err := Engine.First(&admin, 1).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			admin.ID = 1
			admin.Name = "system"
			admin.Password = "system"
			session := Engine.Begin()
			if err := session.Create(&admin).Error; err != nil {
				session.Rollback()
				log.Fatal("fail on create default admin")
			}
			for _, value := range list {
				start, _ := format.ParseTime(value[0])
				end, _ := format.ParseTime(value[1])
				adminTimePeriodTemplate := AdminTimePeriodTemplate{
					AdminId: 1,
					StartAt: start,
					EndAt:   end,
				}
				if err := session.Create(&adminTimePeriodTemplate).Error; err != nil {
					session.Rollback()
					log.Fatal("fail on create default adminTimePeriodTemplate")
				}
			}
			if err := session.Commit().Error; err != nil {
				log.Fatal("fail on commit default admin")
			}
		}else{
			log.Fatal("fail on find default admin")
		}
	}
}