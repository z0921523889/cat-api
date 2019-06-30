package controller

import (
	"cat-api/src/app/format"
	"cat-api/src/app/orm"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
	"time"
)

type TimePeriodController struct {
}

type PostTimePeriodRequest struct {
	StartAt string `form:"start_at" json:"start_at"`
	EndAt   string `form:"end_at" json:"end_at"`
}

// @Description Add a new timePeriod
// @Accept multipart/form-data
// @Produce json
// @Param start_at formData string true "起始時間(yyyy-MM-dd'T'HH:mm:ss)"
// @Param end_at formData string true "結束時間(yyyy-MM-dd'T'HH:mm:ss)"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/time/period [post]
func (controller *TimePeriodController) PostTimePeriod(context *gin.Context) {
	var startAt, endAt time.Time
	var err error
	var request PostTimePeriodRequest
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if startAt, err = format.ParseTime(request.StartAt); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if endAt, err = format.ParseTime(request.EndAt); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	var count int
	err = orm.Engine.Table("AdoptionTimePeriod").
		Where("start_at BETWEEN ? AND ?", startAt, endAt).
		Where("end_at BETWEEN ? AND ?", startAt, endAt).Count(&count).Error
	if err != nil && count > 0 {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	timePeriod := orm.AdoptionTimePeriod{
		StartAt: startAt,
		EndAt:   endAt,
	}
	err = orm.Engine.Create(&timePeriod).Error
	if err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("insert time period complete timePeriodID=%d", timePeriod.ID)})
}
