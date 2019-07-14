package schedule

import (
	"cat-api/src/app/conf"
	"cat-api/src/app/format"
	"cat-api/src/app/orm"
	"log"
	"strconv"
	"time"
)

func generateTimePeriod() {
	admin := orm.Admin{}
	var templates []orm.AdminTimePeriodTemplate
	id, err := strconv.Atoi(conf.DefaultConfig["TimePeriodTemplateAdminId"])
	if err != nil {
		id = 1
	}
	if err := orm.Engine.First(&admin, id).Related(&templates, "AdminId").Error; err != nil {
		log.Println("can not find default admin for time period template")
		return
	}
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0).Add(time.Nanosecond * -1)
	ormSession := orm.Engine.Begin()
	for _, value := range templates {
		day := time.Now()
		for end.After(day) {
			start := format.ReplaceDateOfTime(value.StartAt, day)
			end := format.ReplaceDateOfTime(value.EndAt, day)
			day = day.AddDate(0, 0, 1)
			if orm.Engine.
				Where("start_time = ?", start).
				Where("end_time = ?", end).
				First(&orm.AdoptionTimePeriod{}).
				RecordNotFound() {
				if err := ormSession.Create(&orm.AdoptionTimePeriod{
					StartTime: start,
					EndTime:   end,
				}).Error; err != nil {
					ormSession.Rollback()
				}
				log.Println("data not exist --> create")
			} else {
				log.Println("data exist --> break")
				continue
			}
		}
	}
	ormSession.Commit()
}
