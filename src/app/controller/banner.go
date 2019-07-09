package controller

import (
	"cat-api/src/app/orm"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"mime/multipart"
	"net/http"
	"strconv"
)

type BannerController struct {
	FileController
}

type PostBannerRequest struct {
	Banner *multipart.FileHeader `form:"banner" json:"banner" binding:"required"`
	Order  int                   `form:"order" json:"order" binding:"required"`
}

// @Description insert banner to the database
// @Accept multipart/form-data
// @Produce json
// @Param banner formData file true "banner的縮圖檔案"
// @Param order formData int true "banner的排序"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/banner [post]
func (controller *BannerController) PostBanner(context *gin.Context) {
	var request PostBannerRequest
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	data, err := controller.getBinaryDataFromMultipartFile(request.Banner)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	banner := orm.Banner{
		Data: data,
		Sort: int64(request.Order),
	}
	if err := orm.Engine.Create(&banner).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("insert banner complete bannerID=%d", banner.ID)})
}

type PutModifyBannerRequest struct {
	Banner *multipart.FileHeader `form:"banner" json:"banner" binding:"required"`
	Order  int                   `form:"order" json:"order" binding:"required"`
}

// @Description modify banner information to the database
// @Accept multipart/form-data
// @Produce json
// @Param bannerId path int true "banner的ID"
// @Param banner formData file true "banner的縮圖檔案"
// @Param order formData int true "banner的排序"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/banner/{bannerId} [put]
func (controller *BannerController) PutModifyBanner(context *gin.Context) {
	bannerIdString := context.Param("bannerId")
	bannerId, err := strconv.ParseUint(bannerIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	var request PutModifyBannerRequest
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	data, err := controller.getBinaryDataFromMultipartFile(request.Banner)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	var banner orm.Banner
	if err := orm.Engine.First(&banner, bannerId).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	banner.Data = data
	banner.Sort = int64(request.Order)
	if err := orm.Engine.Save(&banner).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("update cat complete bannerID=%d", banner.ID)})
}

type GetBannerListRequest struct {
	ListRequest
}

type GetBannerListResponse struct {
	ListResponse
	Banners []BannerItem `form:"banners" json:"banners"`
}

type BannerItem struct {
	Id            uint   `form:"id" json:"id"`
	Order         int    `form:"order" json:"order"`
	ThumbnailPath string `form:"cat_thumbnail_path" json:"cat_thumbnail_path"`
}

// @Description get banner list from database
// @Accept json
// @Produce json
// @Param lower query int true "banner列表的lower"
// @Param upper query int true "banner列表的upper"
// @Success 200 {object} controller.GetBannerListResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/banners [get]
func (controller *BannerController) GetBannerList(context *gin.Context) {
	var total int
	var banners []orm.Banner
	var request GetBannerListRequest
	var response GetBannerListResponse
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if err := orm.Engine.Table("banner").Count(&total).
		Limit(request.Upper - request.Lower + 1).Offset(request.Lower).
		Order("sort asc").Find(&banners).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	response.Banners = make([]BannerItem, 0)
	for _, banner := range banners {
		bannerThumbnailPath := fmt.Sprintf("/api/v1/banner/%d", banner.ID)
		response.Banners = append(response.Banners, BannerItem{
			Id:            banner.ID,
			Order:         int(banner.Sort),
			ThumbnailPath: bannerThumbnailPath,
		})
	}
	response.Lower = request.Lower
	response.Upper = request.Upper
	response.Total = total
	context.JSON(http.StatusOK, &response)
}

// @Description get banner thumbnail from the database
// @Accept json
// @Produce image/jpeg
// @Param bannerId path int true "banner的ID"
// @Success 200 {string} binary
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/banner/{bannerId} [get]
func (controller *BannerController) GetBannerThumbnail(context *gin.Context) {
	bannerIdString := context.Param("bannerId")
	bannerId, err := strconv.ParseUint(bannerIdString, 10, 32)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	var banner orm.Banner
	err = orm.Engine.First(&banner, bannerId).Error
	if err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	controller.serveBinaryFile(context, banner.Data)
}
