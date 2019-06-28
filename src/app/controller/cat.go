package controller

import (
	"cat-api/src/app/orm"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"log"
	"net/http"
	"strconv"
)

type CatController struct {
	FileController
}

type PostCatRequest struct {
	Name             string `form:"name" json:"name"`
	Level            string `form:"level" json:"level"`
	Price            int64  `form:"price" json:"price"`
	PetCoin          int64  `form:"pet_coin" json:"pet_coin"`
	ReservationPrice int64  `form:"reservation_price" json:"reservation_price"`
	AdoptionPrice    int64  `form:"adoption_price" json:"adoption_price"`
	ContractDays     int64  `form:"contract_days" json:"contract_days"`
	ContractBenefit  int64  `form:"contract_benefit" json:"contract_benefit"`
}

//@Summary Add new cat to the database
// @Description Add a new cat
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "貓的名稱"
// @Param level formData string true "貓的級別"
// @Param price formData string true "貓的價格"
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
	log.Println("hello")
	var request PostCatRequest
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	cat := orm.Cat{
		Name:             request.Name,
		Level:            request.Level,
		Price:            request.Price,
		PetCoin:          request.PetCoin,
		ReservationPrice: request.ReservationPrice,
		AdoptionPrice:    request.AdoptionPrice,
		Status:           0,
		ContractDays:     request.ContractDays,
		ContractBenefit:  request.ContractBenefit,
	}
	if err := orm.Engine.Create(&cat).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("insert cat complete catID=%d", cat.ID)})
}

type GetCatRequest struct {
	ListRequest
	Status int `form:"status" json:"status"`
}

type GetCatResponse struct {
	ListResponse
	Cats []CatItem `form:"cats" json:"cats"`
}

type CatItem struct {
	Name             string `form:"pet_coin" json:"name"`
	Level            string `form:"pet_coin" json:"level"`
	Price            int64  `form:"pet_coin" json:"price"`
	PetCoin          int64  `form:"pet_coin" json:"pet_coin"`
	ReservationPrice int64  `form:"pet_coin" json:"reservation_price"`
	AdoptionPrice    int64  `form:"pet_coin" json:"adoption_price"`
	ContractDays     int64  `form:"pet_coin" json:"contract_days"`
	ContractBenefit  int64  `form:"pet_coin" json:"contract_benefit"`
	CatThumbnailPath string `form:"pet_coin" json:"cat_thumbnail_path"`
}

func (controller *CatController) GetCat(context *gin.Context) {
	var total int
	var cats []orm.Cat
	var request GetCatRequest
	var catsResponse []CatItem
	if err := context.Bind(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := orm.Engine.Table("cats").Where("Status = ?", request.Status).Count(&total).Limit(request.Upper - request.Lower).Offset(request.Lower).Find(&cats).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, cat := range cats {
		catsResponse = append(catsResponse, CatItem{
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
		log.Println(cat.CreatedAt)
	}
	context.JSON(http.StatusOK, &GetCatResponse{
		ListResponse: ListResponse{
			Lower: request.Lower,
			Upper: request.Upper,
			Total: total,
		},
		Cats: catsResponse,
	})
}

func (controller *CatController) PostCatThumbnail(context *gin.Context) {
	catIdString := context.Param("catId")
	catId, err := strconv.ParseUint(catIdString, 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	data, err := controller.getBinaryFile(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": &err})
	}

	session := orm.Engine.Begin()
	catThumbnail := orm.CatThumbnails{
		Data: data,
	}
	err = session.Create(&catThumbnail).Error
	if err != nil {
		session.Rollback()
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var cat orm.Cat
	err = session.Where("id = ?", catId).First(&cat).Error
	if err != nil {
		session.Rollback()
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.CatThumbnailId = catThumbnail.ID
	err = session.Save(&cat).Error
	if err != nil {
		session.Rollback()
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = session.Commit().Error
	if err != nil {
		session.Rollback()
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.Status(http.StatusOK)
}

func (controller *CatController) GetCatThumbnail(context *gin.Context) {
	catIdString := context.Param("catId")
	catId, err := strconv.ParseUint(catIdString, 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var cat orm.Cat
	err = orm.Engine.Where("id = ?", catId).First(&cat).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var catThumbnail orm.CatThumbnails
	err = orm.Engine.Model(&cat).Related(&catThumbnail, "CatThumbnailId").Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.serveBinaryFile(context, catThumbnail.Data)
}

type PutModifyCatRequest struct {
	Name             string `form:"name" json:"name"`
	Level            string `form:"level" json:"level"`
	Price            int64  `form:"price_min" json:"price"`
	PetCoin          int64  `form:"pet_coin" json:"pet_coin"`
	ReservationPrice int64  `form:"reservation_price" json:"reservation_price"`
	AdoptionPrice    int64  `form:"adoption_price" json:"adoption_price"`
	ContractDays     int64  `form:"contract_days" json:"contract_days"`
	ContractBenefit  int64  `form:"contract_benefit" json:"contract_benefit"`
	Status           int64  `form:"status" json:"status"`
}

func (controller *CatController) PutModifyCat(context *gin.Context) {
	catIdString := context.Param("catId")
	catId, err := strconv.ParseUint(catIdString, 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var request PutModifyCatRequest
	if err := context.Bind(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var cat orm.Cat
	err = orm.Engine.Where("id = ?", catId).First(&cat).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.Name = request.Name
	cat.Level = request.Level
	cat.Price = request.Price
	cat.PetCoin = request.PetCoin
	cat.ReservationPrice = request.ReservationPrice
	cat.AdoptionPrice = request.AdoptionPrice
	cat.Status = request.Status
	cat.ContractDays = request.ContractDays
	cat.ContractBenefit = request.ContractBenefit
	err = orm.Engine.Save(&cat).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.Status(http.StatusOK)
}
