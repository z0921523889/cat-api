package controller

import (
	"cat-api/src/app/conf"
	"cat-api/src/app/orm"
	"cat-api/src/app/schedule"
	"cat-api/src/app/session"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type CatTransferController struct {
	FileController
}

type GetTransferCatListRequest struct {
	ListRequest
	Status int `form:"status" json:"status"`
}

type GetTransferCatListResponse struct {
	ListResponse
	TransferCatItems []TransferCatItem `form:"transfer_cats" json:"transfer_cats"`
}

type TransferCatItem struct {
	Id                    uint      `form:"id" json:"id"`
	CatName               string    `form:"cat_name" json:"cat_name"`
	Price                 int64     `form:"price" json:"price"`
	ContractDays          int64     `form:"contract_days" json:"contract_days"`
	ContractBenefit       int64     `form:"contract_benefit" json:"contract_benefit"`
	StartAt               time.Time `form:"start_at" json:"start_at"`
	EndAt                 time.Time `form:"end_at" json:"end_at"`
	SellerName            string    `form:"seller_name" json:"seller_name"`
	BuyerName             string    `form:"buyer_name" json:"buyer_name"`
	SellerPhone           string    `form:"seller_phone" json:"seller_phone"`
	BuyerPhone            string    `form:"buyer_phone" json:"buyer_phone"`
	SellerCertificatePath string    `form:"seller_certificate_path" json:"seller_certificate_path" binding:"required"`
	BuyerCertificatePath  string    `form:"buyer_certificate_path" json:"buyer_certificate_path" binding:"required"`
}

// @Description get transfer cat list with status from database
// @Accept json
// @Produce json
// @Param lower query int true "貓列表的lower"
// @Param upper query int true "貓列表的upper"
// @Param status query int true "貓列表的交易狀態( 待交易: 1 / 買家已上傳憑證 : 2/已完成 :3 /已取消)"
// @Success 200 {object} controller.GetCatListResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/transfer/cats [get]
func (controller *CatTransferController) GetTransferCatList(context *gin.Context) {
	var total int
	var catUserTransfers []orm.CatUserTransfer
	var request GetTransferCatListRequest
	var response GetTransferCatListResponse
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	sessionData := session.Get(context, session.UserSessionKey)
	sessionValue := sessionData.(session.UserSessionValue)
	command := orm.Engine.Table("cat_user_transfer").Joins("JOIN cat_user_adoption ON cat_user_adoption.id = cat_user_transfer.cat_user_adoption_id")
	command = command.Where("cat_user_adoption.user_id = ? OR new_user_id >= ?", sessionValue.User.ID, sessionValue.User.ID)
	if request.Status > 0 {
		command = command.Where("status = ?", request.Status)
	}
	command = command.Preload("CatUserAdoption").Preload("User")
	if err := command.Count(&total).
		Limit(request.Upper - request.Lower + 1).Offset(request.Lower).
		Find(&catUserTransfers).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	response.TransferCatItems = make([]TransferCatItem, 0)
	for _, catUserTransfer := range catUserTransfers {
		var cat orm.Cat
		var user orm.User
		if err := orm.Engine.Model(&catUserTransfer.CatUserAdoption).Related(&cat).Related(&user).Error; err != nil {
			httputil.NewError(context, http.StatusInternalServerError, err)
			return
		}
		response.TransferCatItems = append(response.TransferCatItems, TransferCatItem{
			Id:                   catUserTransfer.ID,
			CatName:              cat.Name,
			Price:                cat.Price,
			ContractDays:         cat.ContractDays,
			ContractBenefit:      cat.ContractBenefit,
			StartAt:              catUserTransfer.StartTime,
			EndAt:                catUserTransfer.EndTime,
			SellerName:           user.UserName,
			BuyerName:            catUserTransfer.User.UserName,
			SellerPhone:          user.Phone,
			BuyerPhone:           catUserTransfer.User.Phone,
			BuyerCertificatePath: fmt.Sprintf("/api/v1/transfer/%d/certificate", catUserTransfer.ID),
		})
	}
	response.Lower = request.Lower
	response.Upper = request.Upper
	response.Total = total
	context.JSON(http.StatusOK, &response)
}

type PostTransferCatCertificateRequest struct {
	CatCertificate *multipart.FileHeader `form:"certificate" json:"certificate" binding:"required"`
}

// @Description post transfer cat certificate to database
// @Accept multipart/form-data
// @Produce json
// @Param transferId path int true "transfer的ID"
// @Param certificate formData file true "certificate的縮圖檔案"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/transfer/{transferId}/certificate [post]
func (controller *CatTransferController) PostTransferCatCertificate(context *gin.Context) {
	var request PostTransferCatCertificateRequest
	transferIdString := context.Param("transferId")
	transferId, err := strconv.ParseUint(transferIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	sessionData := session.Get(context, session.UserSessionKey)
	sessionValue := sessionData.(session.UserSessionValue)
	data, err := controller.getBinaryDataFromMultipartFile(request.CatCertificate)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	var catUserTransfer orm.CatUserTransfer
	if orm.Engine.Where("user_id = ?", sessionValue.User.ID).
		First(&catUserTransfer, transferId).RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, errors.New("transfer record not found"))
		return
	}
	catUserTransfer.Status = 2
	catUserTransfer.Certificate = data
	ormSession := orm.Engine.Begin()
	if err := ormSession.Save(&catUserTransfer).Error; err != nil {
		ormSession.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	endingCatTransferTaskData := schedule.EndingCatTransferTaskData{
		CatUserTransferId:            int(catUserTransfer.ID),
	}
	taskData, err := json.Marshal(endingCatTransferTaskData)
	if err != nil {
		ormSession.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	transferDuration, err := strconv.Atoi(conf.DefaultConfig["CatTransferDuration"])
	if err != nil {
		ormSession.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	executeTask := orm.ExecuteTask{
		ExecuteType: schedule.EndingCatTransfer,
		ExecuteTime: time.Now().Add(time.Hour * time.Duration(transferDuration)),
		Data:        taskData,
		Done:        false,
	}
	if err := ormSession.Create(&executeTask).Error; err != nil {
		ormSession.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	ormSession.Commit()
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("update transfer with certificate transferID=%d", catUserTransfer.ID)})
}

// @Description get transfer cat certificate thumbnail from the database
// @Accept json
// @Produce image/jpeg
// @Param transferId path int true "transfer的ID"
// @Success 200 {string} binary
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/certificate/transfer/{transferId} [get]
func (controller *CatTransferController) GetTransferCatCertificateThumbnail(context *gin.Context) {
	transferIdString := context.Param("transferId")
	transferId, err := strconv.ParseUint(transferIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	sessionData := session.Get(context, session.UserSessionKey)
	sessionValue := sessionData.(session.UserSessionValue)
	var catUserTransfer orm.CatUserTransfer
	if orm.Engine.Preload("CatUserAdoption").First(&catUserTransfer, transferId).RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, errors.New("transfer record not found"))
		return
	}
	if sessionValue.User.ID != catUserTransfer.UserId && sessionValue.User.ID != catUserTransfer.CatUserAdoption.UserId {
		httputil.NewError(context, http.StatusUnauthorized, errors.New("use are unauthorized for this record"))
		return
	}
	controller.serveBinaryFile(context, catUserTransfer.Certificate)
}

// @Description post transfer cat migrate to the database
// @Accept json
// @Produce json
// @Param transferId path int true "transfer的ID"
// @Success 200 {string} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/transfer/{transferId}/migrate [post]
func (controller *CatTransferController) PostTransferMigrate(context *gin.Context) {
	transferIdString := context.Param("transferId")
	transferId, err := strconv.ParseUint(transferIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	sessionData := session.Get(context, session.UserSessionKey)
	sessionValue := sessionData.(session.UserSessionValue)
	var catUserTransfer orm.CatUserTransfer
	if orm.Engine.Preload("CatUserAdoption").First(&catUserTransfer, transferId).RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, errors.New("transfer record not found"))
		return
	}
	if sessionValue.User.ID != catUserTransfer.CatUserAdoption.UserId {
		httputil.NewError(context, http.StatusUnauthorized, errors.New("user are unauthorized for this record"))
		return
	}
	if len(catUserTransfer.Certificate) == 0 {
		httputil.NewError(context, http.StatusUnauthorized, errors.New("user has not post certificate yet"))
		return
	}
	ormSession := orm.Engine.Begin()
	endingCatTransferTaskData := schedule.EndingCatTransferTaskData{
		CatUserTransferId:            int(catUserTransfer.ID),
	}
	taskData, err := json.Marshal(endingCatTransferTaskData)
	if err != nil {
		ormSession.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	executeTask := orm.ExecuteTask{
		ExecuteType: schedule.EndingCatTransfer,
		ExecuteTime: time.Now(),
		Data:        taskData,
		Done:        false,
	}
	if err := ormSession.Create(&executeTask).Error; err != nil {
		ormSession.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	ormSession.Commit()
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("migrate cat transfer, transferID=%d", catUserTransfer.ID)})
}
