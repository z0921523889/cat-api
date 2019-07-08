package schedule

import (
	"cat-api/src/app/conf"
	"cat-api/src/app/orm"
	"github.com/robfig/cron"
	"log"
	"strconv"
	"time"
)

func StartScheduleJobs() {
	c := cron.New()
	c.AddFunc("@monthly", GenerateTimePeriod)
	c.Start()
}

func GenerateTimePeriod() {
	log.Println("run generateTimePeriod")
	admin := orm.Admin{}
	var templates []orm.AdminTimePeriodTemplate
	id, err := strconv.Atoi(conf.DefaultConfig["TimePeriodTemplateAdminId"])
	if err != nil {
		id = 1
	}
	if err := orm.Engine.First(&admin, id).Related(&templates, "AdminId").Error; err != nil {
		return
	}
	monthLater := time.Now().AddDate(0, 1, 0)
	session := orm.Engine.Begin()
	for _, value := range templates {
		day := time.Now()
		for monthLater.After(day) {
			start := time.Date(
				day.Year(),
				day.Month(),
				day.Day(),
				value.StartAt.Hour(),
				value.StartAt.Minute(),
				value.StartAt.Second(),
				value.StartAt.Nanosecond(),
				value.StartAt.Location(),
			)
			end := time.Date(
				day.Year(),
				day.Month(),
				day.Day(),
				value.EndAt.Hour(),
				value.EndAt.Minute(),
				value.EndAt.Second(),
				value.EndAt.Nanosecond(),
				value.EndAt.Location(),
			)
			day = day.AddDate(0, 0, 1)
			if orm.Engine.
				Where("start_time = ?", start).
				Where("end_time = ?", end).
				First(&orm.AdoptionTimePeriod{}).
				RecordNotFound() {
				if err := session.Create(&orm.AdoptionTimePeriod{
					StartAt: start,
					EndAt:   end,
				}).Error; err != nil {
					session.Rollback()
				}
				log.Println("data not exist --> create")
			} else {
				log.Println("data exist --> break")
				continue
			}
		}
	}
	session.Commit()
}
