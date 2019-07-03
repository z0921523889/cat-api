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

type AdminController struct {
}

type PostAdminLoginRequest struct {
	Account  string `form:"account" json:"account"`
	Password string `form:"password" json:"password"`
}

// @Description admin account login
// @Accept multipart/form-data
// @Produce json
// @Param account formData string true "管理員帳號"
// @Param password formData string true "管理員密碼"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/admin/login [post]
func (controller *AdminController) PostAdminLogin(context *gin.Context) {
	var request PostAdminLoginRequest
	var admin orm.Admins
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if orm.Engine.Where("account = ?", request.Account).
		Where("password = ?", request.Password).
		First(&admin).RecordNotFound() {
		httputil.NewError(context, http.StatusUnauthorized, errors.New("account or password incorrect"))
		return
	}
	if err := session.Set(context, session.AdminSessionKey, session.AdminSessionValue{
		IsLogin: true,
		Admin:   admin,
	}); err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("admin login success")})
}
