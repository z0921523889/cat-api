package controller

import (
	"cat-api/src/app/orm"
	"cat-api/src/app/session"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
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
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	sessionData := session.Get(context, session.UserSessionKey)
	sessionValue := sessionData.(session.UserSessionValue)
	var adoptionTimePeriodCatPivot orm.AdoptionTimePeriodCatPivots
	if orm.Engine.Where("cats_id = ?", request.CatId).
		Where("adoption_time_periods_id = ?", request.TimeScheduleId).
		Preload("Cats").
		First(&adoptionTimePeriodCatPivot).RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, errors.New("cat record not found"))
		return
	}
	if !orm.Engine.Table("cat_user_reservations").
		Where("adoption_time_period_cat_pivots_id = ?", adoptionTimePeriodCatPivot.ID).
		Where("users_id = ?", sessionValue.User.ID).
		First(&orm.CatUserReservations{}).RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, errors.New("record already exist"))
		return
	}
	if sessionValue.User.Wallets.Coin > 0 && sessionValue.User.Wallets.Coin < adoptionTimePeriodCatPivot.Cat.Deposit {
		httputil.NewError(context, http.StatusInternalServerError, errors.New("user wallet does not have enough coin"))
		return
	}
	ormSession := orm.Engine.Begin()
	catUserReservation := orm.CatUserReservations{
		AdoptionTimePeriodCatPivotsId: adoptionTimePeriodCatPivot.ID,
		UsersId:                       sessionValue.User.ID,
	}
	if err := ormSession.Create(&catUserReservation).Error; err != nil {
		ormSession.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	wallet := sessionValue.User.Wallets
	wallet.Coin -= adoptionTimePeriodCatPivot.Cat.Deposit
	if err := ormSession.Save(&wallet).Error; err != nil {
		ormSession.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	ormSession.Commit()
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("insert catUserReservation complete catUserReservationID=%d", catUserReservation.ID)})
}
