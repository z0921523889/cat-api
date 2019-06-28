package controller

import (
	"cat-api/src/app/format"
	"cat-api/src/app/orm"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type TimePeriodController struct {
}

type PostTimePeriodRequest struct {
	StartAt string `form:"start_at" json:"start_at"`
	EndAt   string `form:"end_at" json:"end_at"`
}

func (controller *TimePeriodController) PostTimePeriod(context *gin.Context) {
	var startAt, endAt time.Time
	var err error
	var request PostTimePeriodRequest
	if err := context.Bind(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if startAt, err = format.ParseTime(request.StartAt); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if endAt, err = format.ParseTime(request.StartAt); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var count int
	err = orm.Engine.Table("AdoptionTimePeriod").
		Where("start_at BETWEEN ? AND ?", startAt, endAt).
		Where("end_at BETWEEN ? AND ?", startAt, endAt).Count(&count).Error
	if err != nil && count > 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	timePeriod := orm.AdoptionTimePeriod{
		StartAt: startAt,
		EndAt:   endAt,
	}
	err = orm.Engine.Create(&timePeriod).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.Status(http.StatusOK)
}
