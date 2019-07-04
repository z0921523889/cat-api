package controller

import (
	"cat-api/src/app/orm"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
	"strconv"
)

type CatThumbnailController struct {
	FileController
}

// @Description insert or update cat thumbnail to the database
// @Accept multipart/form-data
// @Produce json
// @Param thumbnail formData file true "貓的縮圖檔案"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/cat/thumbnail [post]
func (controller *CatThumbnailController) PostCatThumbnail(context *gin.Context) {
	data, err := controller.getBinaryDataFromForm(context, "thumbnail")
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	catThumbnail := orm.CatThumbnails{
		Data: data,
	}
	if err := orm.Engine.Create(&catThumbnail).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("insert cat thumbnail complete catThumbnailID=%d", catThumbnail.ID)})
}

type GetCatThumbnailListRequest struct {
	ListRequest
}

type GetCatThumbnailListResponse struct {
	ListResponse
	CatThumbnails []CatThumbnailItem `form:"cat_thumbnails" json:"cat_thumbnails"`
}

type CatThumbnailItem struct {
	Id               uint   `form:"id" json:"id"`
	CatThumbnailPath string `form:"cat_thumbnail_path" json:"cat_thumbnail_path"`
}

// @Description get cat thumbnail list from the database
// @Accept json
// @Produce json
// @Param lower query int true "貓縮圖列表的lower"
// @Param upper query int true "貓縮圖列表的upper"
// @Success 200 {object} controller.GetCatThumbnailListResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/cat/thumbnails [get]
func (controller *CatThumbnailController) GetCatThumbnailList(context *gin.Context) {
	var total int
	var catThumbnails []orm.CatThumbnails
	var request GetCatThumbnailListRequest
	var catThumbnailResponse []CatThumbnailItem
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if err := orm.Engine.Table("cat_thumbnails").
		Count(&total).Limit(request.Upper - request.Lower + 1).Offset(request.Lower).
		Find(&catThumbnails).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	for _, catThumbnail := range catThumbnails {
		catThumbnailResponse = append(catThumbnailResponse, CatThumbnailItem{
			Id:               catThumbnail.ID,
			CatThumbnailPath: fmt.Sprintf("/api/v1/thumbnail/cats/%d", catThumbnail.ID),
		})
	}
	context.JSON(http.StatusOK, &GetCatThumbnailListResponse{
		ListResponse: ListResponse{
			Lower: request.Lower,
			Upper: request.Upper,
			Total: total,
		},
		CatThumbnails: catThumbnailResponse,
	})

}


// @Description get cat thumbnail from the database
// @Accept json
// @Produce image/jpeg
// @Param catId path int true "貓的ID"
// @Success 200 {string} binary
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/thumbnail/cats/{catId} [get]
func (controller *CatThumbnailController) GetCatThumbnail(context *gin.Context) {
	catIdString := context.Param("catId")
	catId, err := strconv.ParseUint(catIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	var cat orm.Cats
	var catThumbnail orm.CatThumbnails
	err = orm.Engine.First(&cat, catId).Related(&catThumbnail, "CatThumbnailId").Error
	if err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	controller.serveBinaryFile(context, catThumbnail.Data)
}

// @Description link cat thumbnail with cat to the database
// @Accept json
// @Produce json
// @Param catId path int true "貓的ID"
// @Param thumbnailId path int true "貓的縮圖ID"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/thumbnails/{thumbnailId}/cats/{catId} [put]
func (controller *CatThumbnailController) PostCatThumbnailBind(context *gin.Context) {
	catIdString := context.Param("catId")
	catId, err := strconv.ParseUint(catIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	thumbnailIdString := context.Param("thumbnailId")
	thumbnailId, err := strconv.ParseUint(thumbnailIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	var cat orm.Cats
	var catThumbnail orm.CatThumbnails
	if orm.Engine.First(&cat, catId).RecordNotFound() || orm.Engine.First(&catThumbnail, thumbnailId).RecordNotFound() {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	cat.CatThumbnailId = catThumbnail.ID
	if err := orm.Engine.Save(&cat).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("update cat with thumbnail complete catID=%d thumbnailID=%d", cat.ID, thumbnailId)})
}

