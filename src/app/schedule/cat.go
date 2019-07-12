package schedule

import (
	"cat-api/src/app/orm"
	"time"
)

func AssignedMarketCat() {
	var adoptionTimePeriods []orm.AdoptionTimePeriod
	if err := orm.Engine.Where("done = ?", false).
		Where("end_time < ?", time.Now()).
		Preload("Cats").
		Find(&adoptionTimePeriods).Error; err != nil {
		return
	}
	ormSession := orm.Engine.Begin()
	for _, adoptionTimePeriod := range adoptionTimePeriods {
		if len(adoptionTimePeriod.Cats) == 0 {
			continue
		}
		for _, cat := range adoptionTimePeriod.Cats {
			var adoptionTimePeriodCatPivot orm.AdoptionTimePeriodCatPivot
			if err := orm.Engine.Where("cat_id = ?", cat.ID).
				Where("adoption_time_period_id = ?", adoptionTimePeriod.ID).
				Find(&adoptionTimePeriodCatPivot).Error; err != nil {
				continue
			}
			switch cat.Status {
			case 1:
				catUserAdoption := orm.CatUserAdoption{
					CatId:     cat.ID,
					UserId:    adoptionTimePeriodCatPivot.UserId,
					StartTime: time.Now(),
					EndTime:   time.Now().AddDate(0, 0, int(cat.ContractDays)),
					Status:    1,
				}
				cat.Status = 4
				if err := ormSession.Create(&catUserAdoption).Save(&cat).Error; err != nil {
					ormSession.Rollback()
					continue
				}
			case 2:
				var catUserAdoption orm.CatUserAdoption
				if err := orm.Engine.Where("cat_id = ?", cat.ID).
					Where("status = ?", 2).
					Find(&catUserAdoption).Error; err != nil {
					return
				}
				catUserTransfer := orm.CatUserTransfer{
					CatId:          cat.ID,
					OriginalUserId: catUserAdoption.UserId,
					NewUserId:      adoptionTimePeriodCatPivot.UserId,
					Status:         1,
					StartTime:      time.Now(),
				}
				cat.Status = 3
				if err := ormSession.Create(&catUserTransfer).Save(&cat).
					Set("gorm:association_autoupdate", false).
					Save(&catUserAdoption).Error; err != nil {
					ormSession.Rollback()
					continue
				}
				if err := ormSession.Set("gorm:association_autoupdate", false).Save(&catUserAdoption).Error; err != nil {
					ormSession.Rollback()
					continue
				}
			}
		}
		adoptionTimePeriod.Done = true
		if err := ormSession.
			Set("gorm:association_autoupdate", false).
			Save(&adoptionTimePeriod).Error; err != nil {
			ormSession.Rollback()
			continue
		}
		ormSession.Commit()
	}
}
