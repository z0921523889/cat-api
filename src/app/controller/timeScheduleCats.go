package controller

import (
	"cat-api/src/app/orm"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
	"strconv"
)

// @Description Add a new timePeriodCatPivot
// @Accept json
// @Produce json
// @Param catId path int true "貓的ID"
// @Param scheduleId path int true "時段的ID"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/time/schedules/{scheduleId}/cat/{catId} [post]
func (controller *TimeScheduleController) PostTimeScheduleCat(context *gin.Context) {
	catIdString := context.Param("catId")
	catId, err := strconv.ParseUint(catIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	scheduleIdString := context.Param("scheduleId")
	scheduleId, err := strconv.ParseUint(scheduleIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	var cat orm.Cats
	var timePeriod orm.AdoptionTimePeriods
	if orm.Engine.First(&cat, catId).RecordNotFound() ||
		orm.Engine.First(&timePeriod, scheduleId).RecordNotFound() ||
		!orm.Engine.Where("cats_id = ?", catId).Where("adoption_time_periods_id = ?", scheduleId).
			First(&orm.AdoptionTimePeriodCatPivots{}).RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	timePeriodCatPivot := orm.AdoptionTimePeriodCatPivots{
		CatsId:                cat.ID,
		AdoptionTimePeriodsId: timePeriod.ID,
	}
	if err := orm.Engine.Create(&timePeriodCatPivot).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("insert timePeriodCatPivot complete timePeriodCatPivotID=%d", timePeriodCatPivot.ID)})
}

type GetTimeScheduleCatRequest struct {
	ListRequest
}

type GetTimeScheduleCatResponse struct {
	ListResponse
	Cats []CatItem `form:"cats" json:"cats"`
}

// @Description get cat list with timeSchedule from database
// @Accept json
// @Produce json
// @Param scheduleId path int true "時段的ID"
// @Param lower query int true "貓列表的lower"
// @Param upper query int true "貓列表的upper"
// @Success 200 {object} controller.GetTimeScheduleCatResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/time/schedules/{scheduleId} [get]
func (controller *TimeScheduleController) GetTimeScheduleCat(context *gin.Context) {
	var total int
	var cats []orm.Cats
	var catIds []uint
	var timePeriod orm.AdoptionTimePeriods
	var request GetTimeScheduleCatRequest
	var catsResponse []CatItem
	scheduleIdString := context.Param("scheduleId")
	scheduleId, err := strconv.ParseUint(scheduleIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	orm.Engine.Preload("Cats").First(&timePeriod, scheduleId)
	for _, cat := range timePeriod.Cats {
		catIds = append(catIds, cat.ID)
	}
	if err := orm.Engine.Table("cats").
		Where("Id in (?)", catIds).
		Count(&total).Limit(request.Upper - request.Lower + 1).Offset(request.Lower).
		Find(&cats).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	for _, cat := range cats {
		catsResponse = append(catsResponse, CatItem{
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
	context.JSON(http.StatusOK, &GetTimeScheduleCatResponse{
		ListResponse: ListResponse{
			Lower: request.Lower,
			Upper: request.Upper,
			Total: total,
		},
		Cats: catsResponse,
	})
}
