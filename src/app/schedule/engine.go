package schedule

import (
	"cat-api/src/app/orm"
	"encoding/json"
	"github.com/robfig/cron"
	"log"
	"time"
)

const (
	AssignedMarketCat      = 0
	EndingCatContract      = 1
	EndingCatTransfer      = 2
	CheckCatTransferCancel = 3
)

type AssignedMarketCatTaskData struct {
	AdoptionTimePeriodId int
	CatId                int
}

type EndingCatContractTaskData struct {
	AdoptionTimePeriodCatPivotId int
	CatUserAdoptionId            int
}

type EndingCatTransferTaskData struct {
	CatUserTransferId            int
}
type CheckCatTransferCancelTaskData struct {
	CatUserTransferId int
	ReturnCoin int
}

func StartScheduleJobs() {
	generateTimePeriod()
	c := cron.New()
	c.AddFunc("@every 1m", RunningExecuteTask)
	c.AddFunc("@monthly", generateTimePeriod)
	c.Start()
}

func RunningExecuteTask() {
	var executeTasks []orm.ExecuteTask
	if orm.Engine.Where("done = ?", false).
		Where("execute_time < ?", time.Now()).
		Find(&executeTasks).RecordNotFound() {
		return
	}
	for _, executeTask := range executeTasks {
		switch executeTask.ExecuteType {
		case AssignedMarketCat:
			var data AssignedMarketCatTaskData
			err := json.Unmarshal(executeTask.Data, &data)
			if err != nil {
				log.Println("RunningExecuteTask : AssignedMarketCat fail err : ", err)
				continue
			}
			assignedMarketCat(data)
		case EndingCatContract:
			var data EndingCatContractTaskData
			err := json.Unmarshal(executeTask.Data, &data)
			if err != nil {
				log.Println("RunningExecuteTask : EndingCatContract fail err : ", err)
				continue
			}
			endingCatContract(data)
		case EndingCatTransfer:
			var data EndingCatTransferTaskData
			err := json.Unmarshal(executeTask.Data, &data)
			if err != nil {
				log.Println("RunningExecuteTask : EndingCatTransfer fail err : ", err)
				continue
			}
			endingCatTransfer(data)
		case CheckCatTransferCancel:
			var data CheckCatTransferCancelTaskData
			err := json.Unmarshal(executeTask.Data, &data)
			if err != nil {
				log.Println("RunningExecuteTask : CheckCatTransferCancel fail err : ", err)
				continue
			}
			checkCatTransferCancel(data)
		}
		executeTask.Done = true
		if err := orm.Engine.Save(&executeTask).Error; err != nil {
			log.Println("executeTask update done status fail err : ", err)
			continue
		}
	}
}
