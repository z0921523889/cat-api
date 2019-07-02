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

type TimeScheduleController struct {
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
// @Router /api/v1/time/schedule [post]
func (controller *TimeScheduleController) PostTimeSchedule(context *gin.Context) {
	var startAt, endAt time.Time
	var err error
	var request PostTimePeriodRequest
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if startAt, err = format.ParseTime(request.StartAt, format.TimeFormatter); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if endAt, err = format.ParseTime(request.EndAt, format.TimeFormatter); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if !orm.Engine.
		Where("start_time = ?", startAt).
		Where("end_time = ?", endAt).
		First(&orm.AdoptionTimePeriods{}).
		RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	timePeriod := orm.AdoptionTimePeriods{
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

type GetTimePeriodRequest struct {
	ListRequest
	Date string `form:"date" json:"date"`
}

type GetTimePeriodResponse struct {
	ListResponse
	TimePeriods []TimePeriodItem `form:"time_periods" json:"time_periods"`
}
type TimePeriodItem struct {
	Id      uint      `form:"id" json:"id"`
	StartAt time.Time `form:"start_at" json:"start_at"`
	EndAt   time.Time `form:"end_at" json:"end_at"`
	Cats    []CatItem `form:"cats" json:"cats"`
}

// @Description get timePeriod list with date
// @Accept json
// @Produce json
// @Param lower query int true "時段列表的lower"
// @Param upper query int true "時段列表的upper"
// @Param date query string true "時段列表的日期(yyyy-MM-dd)"
// @Success 200 {object} controller.GetTimePeriodResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/time/schedules [get]
func (controller *TimeScheduleController) GetTimeScheduleList(context *gin.Context) {
	var total int
	var timePeriods []orm.AdoptionTimePeriods
	var request GetTimePeriodRequest
	var timePeriodResponse []TimePeriodItem
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	time, err := format.ParseTime(request.Date, format.DateFormatter)
	if err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	before := time
	after := time.AddDate(0, 0, 1)
	if err := orm.Engine.
		Table("adoption_time_periods").Preload("Cats").
		Where("start_time BETWEEN ? AND ?", before, after).
		Count(&total).Limit(request.Upper - request.Lower + 1).Offset(request.Lower).
		Find(&timePeriods).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	for _, timePeriod := range timePeriods {
		var catList []CatItem
		for _, cat := range timePeriod.Cats {
			catList = append(catList, CatItem{
				Id:               cat.ID,
				Name:             cat.Name,
				Level:            cat.Level,
				Price:            cat.Price,
				PetCoin:          cat.PetCoin,
				ReservationPrice: cat.ReservationPrice,
				AdoptionPrice:    cat.AdoptionPrice,
				ContractDays:     cat.ContractDays,
				ContractBenefit:  cat.ContractBenefit,
				CatThumbnailPath: fmt.Sprintf("/api/v1/cat/%d/thumbnail", cat.ID),
			})
		}
		timePeriodResponse = append(timePeriodResponse, TimePeriodItem{
			Id:      timePeriod.ID,
			StartAt: timePeriod.EndAt,
			EndAt:   timePeriod.EndAt,
			Cats:    catList,
		})
	}
	context.JSON(http.StatusOK, &GetTimePeriodResponse{
		ListResponse: ListResponse{
			Lower: request.Lower,
			Upper: request.Upper,
			Total: total,
		},
		TimePeriods: timePeriodResponse,
	})
}
