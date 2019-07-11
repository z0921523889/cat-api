package controller

import (
	"cat-api/src/app/orm"
	"cat-api/src/app/session"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
	"time"
)

type CatWithUserController struct {
}

type PostCatReservationRequest struct {
	CatId          uint `form:"cat_id" json:"cat_id"`
	TimeScheduleId uint `form:"time_schedule_id" json:"time_schedule_id"`
}

// @Description Add a new cat_user_reservation
// @Accept multipart/form-data
// @Produce json
// @Param cat_id formData int true "貓的ID識別"
// @Param time_schedule_id formData int true "時段的ID識別"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/cat/reservation [post]
func (controller *CatWithUserController) PostCatReservations(context *gin.Context) {
	var request PostCatReservationRequest
	var user orm.User
	var adoptionTimePeriodCatPivot orm.AdoptionTimePeriodCatPivot
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	sessionData := session.Get(context, session.UserSessionKey)
	sessionValue := sessionData.(session.UserSessionValue)
	if orm.Engine.Where("cat_id = ?", request.CatId).
		Where("adoption_time_period_id = ?", request.TimeScheduleId).
		Preload("Cat").Preload("AdoptionTimePeriod").
		First(&adoptionTimePeriodCatPivot).RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, errors.New("cat record not found"))
		return
	}
	if adoptionTimePeriodCatPivot.Cat.Status != 1 {
		httputil.NewError(context, http.StatusInternalServerError, errors.New("cat is not on sale"))
		return
	}
	if !orm.Engine.Table("cat_user_reservation").
		Where("adoption_time_period_cat_pivot_id = ?", adoptionTimePeriodCatPivot.ID).
		Where("user_id = ?", sessionValue.User.ID).
		First(&orm.CatUserReservation{}).RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, errors.New("record already exist"))
		return
	}
	if orm.Engine.Preload("Wallet").First(&user, sessionValue.User.ID).RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, errors.New("user record not found"))
		return
	}
	timeLayout := "2006-01-02"
	if time.Now().Format(timeLayout) < adoptionTimePeriodCatPivot.AdoptionTimePeriod.StartTime.Format(timeLayout) {
		year, month, day := adoptionTimePeriodCatPivot.AdoptionTimePeriod.StartTime.Date()
		httputil.NewError(context, http.StatusInternalServerError, errors.New(
			fmt.Sprintf("cat not reservation until %d-%s-%d", year, month.String(), day)))
		return
	} else if time.Now().Format(timeLayout) > adoptionTimePeriodCatPivot.AdoptionTimePeriod.StartTime.Format(timeLayout) {
		year, month, day := adoptionTimePeriodCatPivot.AdoptionTimePeriod.StartTime.Date()
		httputil.NewError(context, http.StatusInternalServerError, errors.New(
			fmt.Sprintf("cat not reservation after %d-%s-%d", year, month.String(), day)))
		return
	}
	var coin int64
	if user.Wallet.Coin >= 0 {
		if adoptionTimePeriodCatPivot.AdoptionTimePeriod.StartTime.After(time.Now()) {
			coin = adoptionTimePeriodCatPivot.Cat.ReservationPrice + adoptionTimePeriodCatPivot.Cat.Deposit
		} else {
			coin = adoptionTimePeriodCatPivot.Cat.AdoptionPrice + adoptionTimePeriodCatPivot.Cat.Deposit
		}
		if user.Wallet.Coin < coin {
			httputil.NewError(context, http.StatusInternalServerError, errors.New("user wallet does not have enough coin"))
			return
		}
		user.Wallet.Coin -= coin
	}
	ormSession := orm.Engine.Begin()
	catUserReservation := orm.CatUserReservation{
		AdoptionTimePeriodCatPivotsId: adoptionTimePeriodCatPivot.ID,
		UserId:                        sessionValue.User.ID,
		Coin:                          coin,
	}
	if err := ormSession.Create(&catUserReservation).Error; err != nil {
		ormSession.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	if err := ormSession.Save(&user).Error; err != nil {
		ormSession.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	ormSession.Commit()
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("insert catUserReservation complete catUserReservationID=%d", catUserReservation.ID)})
}

type PutAdoptionCatOwnerRequest struct {
	CatId          uint `form:"cat_id" json:"cat_id"`
	TimeScheduleId uint `form:"time_schedule_id" json:"time_schedule_id"`
	UserId         uint `form:"user_id" json:"user_id"`
}

// @Description modify adoption cat owner
// @Accept multipart/form-data
// @Produce json
// @Param cat_id formData int true "貓的ID識別"
// @Param time_schedule_id formData int true "時段的ID識別"
// @Param user_id formData int true "用戶的ID識別"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/cat/adoption/owner [put]
func (controller *CatWithUserController) PutAdoptionCatOwner(context *gin.Context) {
	var request PutAdoptionCatOwnerRequest
	var adoptionTimePeriodCatPivot orm.AdoptionTimePeriodCatPivot
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if orm.Engine.Where("cat_id = ?", request.CatId).
		Where("adoption_time_period_id = ?", request.TimeScheduleId).
		Preload("Cat").Preload("AdoptionTimePeriod").
		First(&adoptionTimePeriodCatPivot).RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, errors.New("cat record not found"))
		return
	}
	if orm.Engine.First(&orm.User{}).RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, errors.New("user record not found"))
		return
	}
	adoptionTimePeriodCatPivot.UserId = request.UserId
	if err := orm.Engine.Save(&adoptionTimePeriodCatPivot).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("adoption cat will be assigned to user after end of timeline,userID=%d", request.UserId)})
}
