package controller

import (
	"cat-api/src/app/orm"
	"cat-api/src/app/session"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
	"time"
)

type CatAdoptionController struct {
	FileController
}

type GetAdoptionCatListRequest struct {
	ListRequest
	Status int `form:"status" json:"status"`
}

type GetAdoptionCatListResponse struct {
	ListResponse
	AdoptionCatItem []AdoptionCatItem `form:"adoption_cats" json:"adoption_cats"`
}

type AdoptionCatItem struct {
	Id                 uint      `form:"id" json:"id"`
	Status             int       `form:"status" json:"status"`
	CatName            string    `form:"cat_name" json:"cat_name"`
	CatLevel           string    `form:"cat_level" json:"cat_level"`
	CatContractDays    int64     `form:"cat_contract_days" json:"cat_contract_days"`
	CatContractBenefit int64     `form:"cat_contract_benefit" json:"cat_contract_benefit"`
	CatThumbnailPath   string    `form:"cat_thumbnail_path" json:"cat_thumbnail_path"`
	StartAt            time.Time `form:"start_at" json:"start_at"`
	EndAt              time.Time `form:"end_at" json:"end_at"`
}

// @Description get adoption cat list with status from database
// @Accept json
// @Produce json
// @Param lower query int true "貓列表的lower"
// @Param upper query int true "貓列表的upper"
// @Param status query int true "貓列表的領養狀態( 增值中: 1 / 已完成 :2 / 已售出 : 3 / 等待裂變 : 4)"
// @Success 200 {object} controller.GetCatListResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/adoption/cats [get]
func (controller *CatAdoptionController) GetAdoptionCatList(context *gin.Context) {
	var total int
	var catUserAdoptions []orm.CatUserAdoption
	var request GetAdoptionCatListRequest
	var response GetAdoptionCatListResponse
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	sessionData := session.Get(context, session.UserSessionKey)
	sessionValue := sessionData.(session.UserSessionValue)

	command := orm.Engine.Table("cat_user_adoption")
	command = command.Where("cat_user_adoption.user_id = ?", sessionValue.User.ID)
	if request.Status > 0 {
		command = command.Where("status = ?", request.Status)
	}
	command = command.Preload("Cat")
	if err := command.Count(&total).
		Limit(request.Upper - request.Lower + 1).Offset(request.Lower).
		Find(&catUserAdoptions).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	response.AdoptionCatItem = make([]AdoptionCatItem, 0)
	for _, catUserAdoption := range catUserAdoptions {
		catThumbnailPath := ""
		if catUserAdoption.Cat.CatThumbnailId != 0 {
			catThumbnailPath = fmt.Sprintf("/api/v1/cat/thumbnail/%d", catUserAdoption.Cat.CatThumbnailId)
		}
		response.AdoptionCatItem = append(response.AdoptionCatItem, AdoptionCatItem{
			Id:                 catUserAdoption.ID,
			Status:             int(catUserAdoption.Status),
			CatName:            catUserAdoption.Cat.Name,
			CatLevel:           catUserAdoption.Cat.Level,
			CatContractDays:    catUserAdoption.Cat.ContractDays,
			CatContractBenefit: catUserAdoption.Cat.ContractBenefit,
			CatThumbnailPath:   catThumbnailPath,
			StartAt:            catUserAdoption.StartTime,
			EndAt:              catUserAdoption.EndTime,
		})
	}
	response.Lower = request.Lower
	response.Upper = request.Upper
	response.Total = total
	context.JSON(http.StatusOK, &response)
}
