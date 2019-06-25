package controller

import (
	"cat-api/src/app/orm"
	"github.com/gin-gonic/gin"
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
	Price            int64  `form:"price_min" json:"price"`
	PetCoin          int64  `form:"pet_coin" json:"pet_coin"`
	ReservationPrice int64  `form:"reservation_price" json:"reservation_price"`
	AdoptionPrice    int64  `form:"adoption_price" json:"adoption_price"`
	ContractDays     int64  `form:"contract_days" json:"contract_days"`
	ContractBenefit  int64  `form:"contract_benefit" json:"contract_benefit"`
}

func (controller *CatController) PostCat(context *gin.Context) {
	var request PostCatRequest
	if err := context.Bind(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat := orm.Cat{
		Name:             request.Name,
		Level:            request.Level,
		Price:            request.Price,
		PetCoin:          request.PetCoin,
		ReservationPrice: request.ReservationPrice,
		AdoptionPrice:    request.AdoptionPrice,
		Status:           1,
		ContractDays:     request.ContractDays,
		ContractBenefit:  request.ContractBenefit,
	}
	log.Println(cat)
	err := orm.Engine.Create(&cat).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, nil)
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
	context.JSON(http.StatusOK, nil)
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
	err = orm.Engine.Model(&cat).Related(&catThumbnail,"CatThumbnailId").Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	controller.serveBinaryFile(context,catThumbnail.Data)
}
