package schedule

import (
	"github.com/robfig/cron"
)

func StartScheduleJobs() {
	c := cron.New()
	c.AddFunc("@hourly", AssignedMarketCat)
	c.AddFunc("@monthly", GenerateTimePeriod)
	c.Start()
}


