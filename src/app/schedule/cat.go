package schedule

import (
	"cat-api/src/app/conf"
	"cat-api/src/app/format"
	"cat-api/src/app/orm"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
)

func assignedMarketCat(data AssignedMarketCatTaskData) {
	var adoptionTimePeriodCatPivot orm.AdoptionTimePeriodCatPivot
	var catUserReservations []orm.CatUserReservation
	var cat orm.Cat
	var ownerId uint
	var ownerReservation orm.CatUserReservation
	if err := orm.Engine.Where("cat_id = ?", data.CatId).
		Where("adoption_time_period_id = ?", data.AdoptionTimePeriodId).
		Preload("Cat").Preload("CatUserReservations").
		First(&adoptionTimePeriodCatPivot).Error; err != nil {
		return
	}
	cat = adoptionTimePeriodCatPivot.Cat
	catUserReservations = adoptionTimePeriodCatPivot.CatUserReservations
	if len(catUserReservations) == 0 {
		return
	}
	if adoptionTimePeriodCatPivot.UserId == 0 {
		rand.Seed(time.Now().Unix())
		ownerId = catUserReservations[rand.Intn(len(catUserReservations))].UserId
	} else {
		ownerId = adoptionTimePeriodCatPivot.UserId
	}
	ormSession := orm.Engine.Begin()
	for _, catUserReservation := range catUserReservations {
		if ownerId == catUserReservation.UserId {
			ownerReservation = catUserReservation
			continue
		}
		var user orm.User
		var wallet orm.Wallet
		if orm.Engine.First(&user, catUserReservation.UserId).Related(&wallet).RecordNotFound() {
			continue
		}
		wallet.Coin += catUserReservation.Coin
		if err := ormSession.Set("gorm:association_autoupdate", false).
			Save(&wallet).Error; err != nil {
			continue
		}
	}
	switch cat.Status {
	case 1:
		catUserAdoption := orm.CatUserAdoption{
			CatId:     cat.ID,
			UserId:    ownerId,
			StartTime: time.Now(),
			EndTime:   time.Now().AddDate(0, 0, int(cat.ContractDays)),
			Status:    1,
		}
		cat.Status = 4
		if err := ormSession.Create(&catUserAdoption).Save(&cat).Error; err != nil {
			ormSession.Rollback()
			return
		}
		endingCatContractTaskData := EndingCatContractTaskData{
			AdoptionTimePeriodCatPivotId: int(adoptionTimePeriodCatPivot.ID),
			CatUserAdoptionId:            int(catUserAdoption.ID),
		}
		data, err := json.Marshal(endingCatContractTaskData)
		if err != nil {
			ormSession.Rollback()
			return
		}
		executeTask := orm.ExecuteTask{
			ExecuteType: EndingCatContract,
			ExecuteTime: catUserAdoption.EndTime,
			Data:        data,
			Done:        false,
		}
		if err := ormSession.Create(&executeTask).Error; err != nil {
			ormSession.Rollback()
			return
		}
	case 2:
		var catUserAdoption orm.CatUserAdoption
		if err := orm.Engine.Where("cat_id = ?", cat.ID).
			Where("status = ?", 2).
			Find(&catUserAdoption).Error; err != nil {
			ormSession.Rollback()
			return
		}
		catUserTransfer := orm.CatUserTransfer{
			AdoptionTimePeriodCatPivotId: adoptionTimePeriodCatPivot.ID,
			CatUserAdoptionId:            catUserAdoption.ID,
			UserId:                       ownerId,
			Status:                       1,
			StartTime:                    time.Now(),
		}
		cat.Status = 3
		if err := ormSession.Create(&catUserTransfer).Save(&cat).Error; err != nil {
			ormSession.Rollback()
			return
		}

		checkCatTransferCancelTaskData := CheckCatTransferCancelTaskData{
			CatUserTransferId: int(catUserTransfer.ID),
			ReturnCoin:        int(ownerReservation.Coin - cat.Deposit),
		}
		data, err := json.Marshal(checkCatTransferCancelTaskData)
		if err != nil {
			ormSession.Rollback()
			return
		}
		transferDuration, err := strconv.Atoi(conf.DefaultConfig["CatTransferDuration"])
		if err != nil {
			ormSession.Rollback()
			return
		}
		executeTask := orm.ExecuteTask{
			ExecuteType: CheckCatTransferCancel,
			ExecuteTime: time.Now().Add(time.Hour * time.Duration(transferDuration)),
			Data:        data,
			Done:        false,
		}
		if err := ormSession.Create(&executeTask).Error; err != nil {
			ormSession.Rollback()
			return
		}
	}
	ormSession.Commit()
}

func endingCatContract(data EndingCatContractTaskData) {
	var catUserAdoption orm.CatUserAdoption
	var adoptionTimePeriodCatPivot orm.AdoptionTimePeriodCatPivot
	var userWallet orm.Wallet
	var timePeriod orm.AdoptionTimePeriod
	if err := orm.Engine.Preload("Cat").Preload("User").First(&catUserAdoption, data.CatUserAdoptionId).Error; err != nil {
		return
	}
	if err := orm.Engine.Preload("AdoptionTimePeriod").First(&adoptionTimePeriodCatPivot, data.AdoptionTimePeriodCatPivotId).Error; err != nil {
		return
	}
	if err := orm.Engine.Model(&catUserAdoption.User).Related(&userWallet).Error; err != nil {
		return
	}
	catUserAdoption.Cat.Price = catUserAdoption.Cat.Price * (100 + catUserAdoption.Cat.ContractBenefit) / 100
	catUserAdoption.Cat.Status = 2
	catUserAdoption.Status = 2
	userWallet.PetCoin += catUserAdoption.Cat.PetCoin
	maxPriceOfCat, err := strconv.Atoi(conf.DefaultConfig["MaxPriceOfCat"])
	if err != nil {
		return
	}
	if int(catUserAdoption.Cat.Price) >= maxPriceOfCat {
		catUserAdoption.Cat.Status = 5
		if err := orm.Engine.Save(&catUserAdoption).Save(&userWallet).Error; err != nil {
			return
		}
		return
	}
	ormSession := orm.Engine.Begin()
	if err := ormSession.Save(&catUserAdoption).Save(&userWallet).Error; err != nil {
		ormSession.Rollback()
		return
	}
	date := time.Now().AddDate(0, 0, 1)
	start := format.ReplaceDateOfTime(adoptionTimePeriodCatPivot.AdoptionTimePeriod.StartTime, date)
	end := format.ReplaceDateOfTime(adoptionTimePeriodCatPivot.AdoptionTimePeriod.EndTime, date)
	if err := orm.Engine.Where("start_time = ?", start).Where("end_time = ?", end).
		First(&timePeriod).Error; err != nil {
		return
	}
	timePeriodCatPivot := orm.AdoptionTimePeriodCatPivot{
		CatId:                catUserAdoption.Cat.ID,
		AdoptionTimePeriodId: timePeriod.ID,
	}
	if err := ormSession.Create(&timePeriodCatPivot).Error; err != nil {
		ormSession.Rollback()
		return
	}
	assignedMarketCatTaskData := AssignedMarketCatTaskData{
		AdoptionTimePeriodId: int(timePeriod.ID),
		CatId:                int(catUserAdoption.Cat.ID),
	}
	taskData, err := json.Marshal(assignedMarketCatTaskData)
	if err != nil {
		ormSession.Rollback()
		return
	}
	executeTask := orm.ExecuteTask{
		ExecuteType: AssignedMarketCat,
		ExecuteTime: timePeriod.EndTime,
		Data:        taskData,
		Done:        false,
	}
	if err := ormSession.Create(&executeTask).Error; err != nil {
		ormSession.Rollback()
		return
	}
	ormSession.Commit()
}

func endingCatTransfer(data EndingCatTransferTaskData) {
	var catUserTransfer orm.CatUserTransfer
	var timePeriod orm.AdoptionTimePeriod
	var cat orm.Cat
	if err := orm.Engine.Not("status = ?", 3).
		Preload("CatUserAdoption").Preload("AdoptionTimePeriodCatPivot").
		First(&catUserTransfer, data.CatUserTransferId).Error; err != nil {
		return
	}
	if err := orm.Engine.Model(&catUserTransfer.AdoptionTimePeriodCatPivot).
		Related(&catUserTransfer.AdoptionTimePeriodCatPivot.AdoptionTimePeriod).Error; err != nil {
		return
	}
	date := time.Now().AddDate(0, 0, 1)
	start := format.ReplaceDateOfTime(catUserTransfer.AdoptionTimePeriodCatPivot.AdoptionTimePeriod.StartTime, date)
	end := format.ReplaceDateOfTime(catUserTransfer.AdoptionTimePeriodCatPivot.AdoptionTimePeriod.EndTime, date)
	if err := orm.Engine.Where("start_time = ?", start).Where("end_time = ?", end).
		First(&timePeriod).Error; err != nil {
		return
	}
	if err := orm.Engine.Model(&catUserTransfer.CatUserAdoption).Related(&cat).Error; err != nil {
		return
	}
	catUserAdoption := orm.CatUserAdoption{
		CatId:     catUserTransfer.CatUserAdoption.CatId,
		UserId:    catUserTransfer.UserId,
		StartTime: time.Now(),
		EndTime:   time.Now().AddDate(0, 0, int(cat.ContractDays)),
		Status:    1,
	}

	ormSession := orm.Engine.Begin()
	if err := ormSession.Create(&catUserAdoption).Error; err != nil {
		ormSession.Rollback()
		return
	}
	timePeriodCatPivot := orm.AdoptionTimePeriodCatPivot{
		CatId:                catUserTransfer.CatUserAdoption.CatId,
		AdoptionTimePeriodId: timePeriod.ID,
	}
	if err := ormSession.Create(&timePeriodCatPivot).Error; err != nil {
		ormSession.Rollback()
		return
	}
	cat.Status = 4
	catUserTransfer.Status = 3
	catUserTransfer.CatUserAdoption.Status = 3
	if err := ormSession.Save(&catUserTransfer).Save(&cat).Error; err != nil {
		ormSession.Rollback()
		return
	}
	endingCatContractTaskData := EndingCatContractTaskData{
		AdoptionTimePeriodCatPivotId: int(timePeriodCatPivot.ID),
		CatUserAdoptionId:            int(catUserAdoption.ID),
	}
	taskData, err := json.Marshal(endingCatContractTaskData)
	if err != nil {
		ormSession.Rollback()
		return
	}
	executeTask := orm.ExecuteTask{
		ExecuteType: EndingCatContract,
		ExecuteTime: catUserAdoption.EndTime,
		Data:        taskData,
		Done:        false,
	}
	if err := ormSession.Create(&executeTask).Error; err != nil {
		ormSession.Rollback()
		return
	}
	ormSession.Commit()
}

func checkCatTransferCancel(data CheckCatTransferCancelTaskData) {
	var catUserTransfer orm.CatUserTransfer
	var timePeriod orm.AdoptionTimePeriod
	var cat orm.Cat
	var user orm.User
	var userWallet orm.Wallet
	if err := orm.Engine.Not("status = ?", 3).
		Preload("CatUserAdoption").Preload("AdoptionTimePeriodCatPivot").
		First(&catUserTransfer, data.CatUserTransferId).Error; err != nil {
		return
	}
	if len(catUserTransfer.Certificate) > 0 {
		return
	}
	if err := orm.Engine.Model(&catUserTransfer.AdoptionTimePeriodCatPivot).
		Related(&catUserTransfer.AdoptionTimePeriodCatPivot.AdoptionTimePeriod).Error; err != nil {
		return
	}
	if err := orm.Engine.Model(&catUserTransfer.CatUserAdoption).Related(&cat).Related(&user).Error; err != nil {
		return
	}
	if err := orm.Engine.Model(&user).Related(&userWallet).Error; err != nil {
		return
	}
	date := time.Now().AddDate(0, 0, 1)
	start := format.ReplaceDateOfTime(catUserTransfer.AdoptionTimePeriodCatPivot.AdoptionTimePeriod.StartTime, date)
	end := format.ReplaceDateOfTime(catUserTransfer.AdoptionTimePeriodCatPivot.AdoptionTimePeriod.EndTime, date)
	if err := orm.Engine.Where("start_time = ?", start).Where("end_time = ?", end).
		First(&timePeriod).Error; err != nil {
		return
	}
	timePeriodCatPivot := orm.AdoptionTimePeriodCatPivot{
		CatId:                cat.ID,
		AdoptionTimePeriodId: timePeriod.ID,
	}
	cat.Status = 2
	catUserTransfer.Status = 4
	userWallet.Coin += int64(data.ReturnCoin)
	ormSession := orm.Engine.Begin()
	if err := ormSession.Create(&timePeriodCatPivot).Save(&catUserTransfer).Save(&cat).Save(&userWallet).Error; err != nil {
		ormSession.Rollback()
		return
	}
	assignedMarketCatTaskData := AssignedMarketCatTaskData{
		AdoptionTimePeriodId: int(timePeriod.ID),
		CatId:                int(cat.ID),
	}
	taskData, err := json.Marshal(assignedMarketCatTaskData)
	if err != nil {
		ormSession.Rollback()
		return
	}
	executeTask := orm.ExecuteTask{
		ExecuteType: AssignedMarketCat,
		ExecuteTime: timePeriod.EndTime,
		Data:        taskData,
		Done:        false,
	}
	if err := ormSession.Create(&executeTask).Error; err != nil {
		ormSession.Rollback()
		return
	}
	ormSession.Commit()
}
