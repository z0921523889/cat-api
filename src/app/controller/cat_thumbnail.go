package controller

import (
	"cat-api/src/app/orm"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
	"strconv"
)

type CatThumbnailController struct {
	FileController
}

// @Description insert or update cat thumbnail to the database
// @Accept image/jpeg
// @Produce json
// @Param catId path int true "貓的ID"
// @Param thumbnail formData File true "貓的縮圖檔案"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/cat/{catId}/thumbnail [post]
func (controller *CatThumbnailController) PostCatThumbnail(context *gin.Context) {
	catIdString := context.Param("catId")
	catId, err := strconv.ParseUint(catIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	data, err := controller.getBinaryDataFromForm(context,"thumbnail")
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	session := orm.Engine.Begin()
	var cat orm.Cat
	var catThumbnail orm.CatThumbnails
	if err := orm.Engine.First(&cat,catId).Related(&catThumbnail, "CatThumbnailId").Error; err != nil {
		catThumbnail.Data = data
		if err == gorm.ErrRecordNotFound {
			if err := session.Create(&catThumbnail).Error; err != nil {
				session.Rollback()
				httputil.NewError(context, http.StatusInternalServerError, err)
				return
			}
			cat.CatThumbnailId = catThumbnail.ID
			if err := session.Save(&cat).Error; err != nil {
				session.Rollback()
				httputil.NewError(context, http.StatusInternalServerError, err)
				return
			}
		} else {
			session.Rollback()
			httputil.NewError(context, http.StatusInternalServerError, err)
			return
		}
	} else {
		catThumbnail.Data = data
		if err := session.Save(&catThumbnail).Error; err != nil {
			session.Rollback()
			httputil.NewError(context, http.StatusInternalServerError, err)
			return
		}
	}
	err = session.Commit().Error
	if err != nil {
		session.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("insert or update cat thumbnail complete catID=%d", cat.ID)})
}

// @Description get cat thumbnail from the database
// @Accept json
// @Produce image/jpeg
// @Param catId path int true "貓的ID"
// @Success 200 {string} binary
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/cat/{catId}/thumbnail [get]
func (controller *CatThumbnailController) GetCatThumbnail(context *gin.Context) {
	catIdString := context.Param("catId")
	catId, err := strconv.ParseUint(catIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	var cat orm.Cat
	var catThumbnail orm.CatThumbnails
	err = orm.Engine.First(&cat,catId).Related(&catThumbnail,"CatThumbnailId").Error
	if err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	controller.serveBinaryFile(context, catThumbnail.Data)
}
