package schedule

import (
	"cat-api/src/app/conf"
	"cat-api/src/app/orm"
	"github.com/robfig/cron"
	"strconv"
	"time"
)

func StartScheduleJobs() {
	c := cron.New()
	c.AddFunc("@daily", generateTimePeriod)
	c.Start()
}

func generateTimePeriod() {
	admin := orm.Admins{}
	adminTimePeriodTemplate := orm.AdminTimePeriodTemplate{}
	id, err := strconv.Atoi(conf.DefaultConfig["TimePeriodTemplateAdminId"])
	if err != nil {
		id = 1
	}
	var count int
	orm.Engine.Table("adoption_time_period").Where("start_at > ?", time.Now()).Count(&count)
	if count > 0 {
		return
	}
	if err := orm.Engine.First(&admin, id).Related(&adminTimePeriodTemplate, "AdminId").Error; err != nil {
		monthLater := time.Now().AddDate(0, 1, 0)
		day := time.Now()
		session := orm.Engine.Begin()
		for monthLater.After(day) {
			start := time.Date(
				day.Year(),
				day.Month(),
				day.Day(),
				adminTimePeriodTemplate.StartAt.Hour(),
				adminTimePeriodTemplate.StartAt.Minute(),
				adminTimePeriodTemplate.StartAt.Second(),
				adminTimePeriodTemplate.StartAt.Nanosecond(),
				adminTimePeriodTemplate.StartAt.Location(),
			)
			end := time.Date(
				day.Year(),
				day.Month(),
				day.Day(),
				adminTimePeriodTemplate.EndAt.Hour(),
				adminTimePeriodTemplate.EndAt.Minute(),
				adminTimePeriodTemplate.EndAt.Second(),
				adminTimePeriodTemplate.EndAt.Nanosecond(),
				adminTimePeriodTemplate.EndAt.Location(),
			)
			if err := session.Create(&orm.AdoptionTimePeriod{
				StartAt: start,
				EndAt:   end,
			}).Error; err != nil {
				session.Rollback()
			}
			day.AddDate(0, 0, 1)
		}
		session.Commit()
	}
}
