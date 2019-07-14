package controller

import (
	"cat-api/src/app/orm"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
	"strconv"
)

type CatController struct {
}

type PostCatRequest struct {
	Name             string `form:"name" json:"name"`
	Level            string `form:"level" json:"level"`
	Price            int64  `form:"price" json:"price"`
	Deposit          int64  `form:"deposit" json:"deposit"`
	PetCoin          int64  `form:"pet_coin" json:"pet_coin"`
	ReservationPrice int64  `form:"reservation_price" json:"reservation_price"`
	AdoptionPrice    int64  `form:"adoption_price" json:"adoption_price"`
	ContractDays     int64  `form:"contract_days" json:"contract_days"`
	ContractBenefit  int64  `form:"contract_benefit" json:"contract_benefit"`
}

// @Description Add a new cat
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "貓的名稱"
// @Param level formData string true "貓的級別"
// @Param price formData string true "貓的價格"
// @Param deposit formData string true "貓的押金"
// @Param pet_coin formData string true "貓的pet幣"
// @Param reservation_price formData string true "貓的預約價格"
// @Param adoption_price formData string true "貓的即搶價格"
// @Param contract_days formData string true "貓的合約時間"
// @Param contract_benefit formData string true "貓的合約增益"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/cat [post]
func (controller *CatController) PostCat(context *gin.Context) {
	var request PostCatRequest
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	cat := orm.Cat{
		Name:             request.Name,
		Level:            request.Level,
		Price:            request.Price,
		Deposit:          request.Deposit,
		PetCoin:          request.PetCoin,
		ReservationPrice: request.ReservationPrice,
		AdoptionPrice:    request.AdoptionPrice,
		Status:           1,
		ContractDays:     request.ContractDays,
		ContractBenefit:  request.ContractBenefit,
	}
	if err := orm.Engine.Create(&cat).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("insert cat complete catID=%d", cat.ID)})
}

type PutModifyCatRequest struct {
	Name             string `form:"name" json:"name"`
	Level            string `form:"level" json:"level"`
	Price            int64  `form:"price_min" json:"price"`
	Deposit          int64  `form:"deposit" json:"deposit"`
	PetCoin          int64  `form:"pet_coin" json:"pet_coin"`
	ReservationPrice int64  `form:"reservation_price" json:"reservation_price"`
	AdoptionPrice    int64  `form:"adoption_price" json:"adoption_price"`
	ContractDays     int64  `form:"contract_days" json:"contract_days"`
	ContractBenefit  int64  `form:"contract_benefit" json:"contract_benefit"`
	Status           int64  `form:"status" json:"status"`
}

// @Description modify cat information to the database
// @Accept multipart/form-data
// @Produce json
// @Param catId path int true "貓的ID"
// @Param name formData string true "貓的名稱"
// @Param level formData string true "貓的級別"
// @Param price formData string true "貓的價格"
// @Param deposit formData string true "貓的押金"
// @Param pet_coin formData string true "貓的pet幣"
// @Param reservation_price formData string true "貓的預約價格"
// @Param adoption_price formData string true "貓的即搶價格"
// @Param contract_days formData string true "貓的合約時間"
// @Param contract_benefit formData string true "貓的合約增益"
// @Param status formData string true "貓的狀態(系統掛售中 : 1 / 期約到期掛售中 : 2/轉讓中 : 3 /領養增值中 :4 / 等待裂變中 :5)"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/cats/{catId} [put]
func (controller *CatController) PutModifyCat(context *gin.Context) {
	catIdString := context.Param("catId")
	catId, err := strconv.ParseUint(catIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	var request PutModifyCatRequest
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	var cat orm.Cat
	if err := orm.Engine.First(&cat, catId).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	cat.Name = request.Name
	cat.Level = request.Level
	cat.Price = request.Price
	cat.Deposit = request.Deposit
	cat.PetCoin = request.PetCoin
	cat.ReservationPrice = request.ReservationPrice
	cat.AdoptionPrice = request.AdoptionPrice
	cat.Status = request.Status
	cat.ContractDays = request.ContractDays
	cat.ContractBenefit = request.ContractBenefit
	if err := orm.Engine.Save(&cat).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("update cat complete catID=%d", cat.ID)})
}

type GetCatListRequest struct {
	ListRequest
	Status int `form:"status" json:"status"`
}

type GetCatListResponse struct {
	ListResponse
	Cats []CatItem `form:"cats" json:"cats"`
}

type CatItem struct {
	Id               uint   `form:"id" json:"id"`
	Name             string `form:"name" json:"name"`
	Level            string `form:"level" json:"level"`
	Price            int64  `form:"price" json:"price"`
	Deposit          int64  `form:"deposit" json:"deposit"`
	PetCoin          int64  `form:"pet_coin" json:"pet_coin"`
	ReservationPrice int64  `form:"reservation_price" json:"reservation_price"`
	AdoptionPrice    int64  `form:"adoption_price" json:"adoption_price"`
	ContractDays     int64  `form:"contract_days" json:"contract_days"`
	ContractBenefit  int64  `form:"contract_benefit" json:"contract_benefit"`
	CatThumbnailPath string `form:"cat_thumbnail_path" json:"cat_thumbnail_path"`
	Status           int64  `form:"status" json:"status"`
}

// @Description get cat list with status from database
// @Accept json
// @Produce json
// @Param lower query int true "貓列表的lower"
// @Param upper query int true "貓列表的upper"
// @Param status query int true "貓列表的狀態(系統掛售中 : 1 / 期約到期掛售中 : 2/轉讓中 : 3 /領養增值中 :4 / 等待裂變中 :5)"
// @Success 200 {object} controller.GetCatListResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/cats [get]
func (controller *CatController) GetCatList(context *gin.Context) {
	var total int
	var cats []orm.Cat
	var request GetCatListRequest
	var response GetCatListResponse
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	command := orm.Engine.Table("cat")
	if request.Status > 0 {
		command = command.Where("Status = ?", request.Status)
	}
	if err := command.Count(&total).
		Limit(request.Upper - request.Lower + 1).Offset(request.Lower).
		Find(&cats).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	response.Cats = make([]CatItem,0)
	for _, cat := range cats {
		catThumbnailPath := ""
		if cat.CatThumbnailId != 0 {
			catThumbnailPath = fmt.Sprintf("/api/v1/cat/thumbnail/%d", cat.CatThumbnailId)
		}
		response.Cats = append(response.Cats, CatItem{
			Id:               cat.ID,
			Name:             cat.Name,
			Level:            cat.Level,
			Price:            cat.Price,
			PetCoin:          cat.PetCoin,
			ReservationPrice: cat.ReservationPrice,
			AdoptionPrice:    cat.AdoptionPrice,
			ContractDays:     cat.ContractDays,
			ContractBenefit:  cat.ContractBenefit,
			CatThumbnailPath: catThumbnailPath,
			Status:           cat.Status,
		})
	}
	response.Lower = request.Lower
	response.Upper = request.Upper
	response.Total = total
	context.JSON(http.StatusOK, &response)
}
