package controller

import (
	"cat-api/src/app/orm"
	"cat-api/src/app/session"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
	"time"
)

type UserController struct {
	FileController
}

type PostUserLoginRequest struct {
	Phone    string `form:"phone" json:"phone"`
	Password string `form:"password" json:"password"`
}

// @Description user account login
// @Accept multipart/form-data
// @Produce json
// @Param phone formData string true "用戶手機號碼"
// @Param password formData string true "用戶密碼"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/user/login [post]
func (controller *UserController) PostUserLogin(context *gin.Context) {
	var request PostUserLoginRequest
	var user orm.User
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if orm.Engine.Where("phone = ?", request.Phone).
		Where("password = ?", request.Password).
		Preload("Wallet").
		First(&user).RecordNotFound() {
		httputil.NewError(context, http.StatusUnauthorized, errors.New("account or password incorrect"))
		return
	}
	if err := session.Set(context, session.UserSessionKey, &session.UserSessionValue{
		IsLogin: true,
		User:    user,
	}); err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	if err := session.Set(context, session.AdminSessionKey, &session.AdminSessionValue{
		IsLogin: false,
		Admin:   orm.Admin{},
	}); err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("login success")})
}

// @Description user account logout
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/user/logout [post]
func (controller *UserController) PostUserLogout(context *gin.Context) {
	if err := session.Clear(context); err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("logout success")})
}

type PostUserRegisterRequest struct {
	Phone            string `form:"phone" json:"phone"`
	ValidValue       string `form:"valid_value" json:"valid_value"`
	UserName         string `form:"user_name" json:"user_name"`
	Password         string `form:"password" json:"password"`
	SecurityPassword string `form:"security_password" json:"security_password"`
	IdentifiedCode   string `form:"identified_code" json:"identified_code"`
}

// @Description user account register
// @Accept multipart/form-data
// @Produce json
// @Param phone formData string true "用戶手機號碼"
// @Param valid_value formData string true "用戶手機認證碼"
// @Param user_name formData string true "用戶名稱"
// @Param password formData string true "用戶密碼"
// @Param security_password formData string true "用戶二級密碼"
// @Param identified_code formData string false "推薦用戶識別碼"
// @Success 200 {object} controller.Message
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/user/register [post]
func (controller *UserController) PostUserRegister(context *gin.Context) {
	var request PostUserRegisterRequest
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	sessionData := session.Get(context, session.UserSessionKey)
	if sessionData != nil {
		userSessionValue := sessionData.(session.UserSessionValue)
		if userSessionValue.IsLogin {
			httputil.NewError(context, http.StatusInternalServerError,
				errors.New(
					fmt.Sprintf("already login with another user userID=%d", userSessionValue.User.ID)))
			return
		}
	}
	if !orm.Engine.Where("phone = ?", request.Phone).
		First(&orm.User{}).RecordNotFound() {
		httputil.NewError(context, http.StatusBadRequest, errors.New("phone already exist"))
		return
	}
	hash := sha1.New()
	hash.Write([]byte(request.Phone + time.Now().String()))
	code := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	transition := orm.Engine.Begin()
	wallet := orm.Wallet{}
	if err := transition.Create(&wallet).Error; err != nil {
		transition.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	user := orm.User{
		Phone:            request.Phone,
		UserName:         request.UserName,
		Password:         request.Password,
		SecurityPassword: request.SecurityPassword,
		IdentifiedCode:   code[:15],
		Status:           1,
		Wallet:           wallet,
	}
	if err := transition.Create(&user).Error; err != nil {
		transition.Rollback()
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	if err := transition.Commit().Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("insert user complete userID=%d", user.ID)})
}

type GetUserInfoResponse struct {
	UserName       string `form:"user_name" json:"user_name"`
	Phone          string `form:"phone" json:"phone"`
	IdentifiedCode string `form:"identified_code" json:"identified_code"`
}

// @Description get user info
// @Accept json
// @Produce json
// @Success 200 {object} controller.GetUserInfoResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/user/info [get]
func (controller *UserController) GetUserInfo(context *gin.Context) {
	sessionData := session.Get(context, session.UserSessionKey)
	sessionValue := sessionData.(session.UserSessionValue)
	context.JSON(http.StatusOK, GetUserInfoResponse{
		UserName:       sessionValue.User.UserName,
		Phone:          sessionValue.User.Phone,
		IdentifiedCode: sessionValue.User.IdentifiedCode,
	})
}

type GetUserListRequest struct {
	ListRequest
	Status int `form:"status" json:"status"`
}

type GetUserListResponse struct {
	ListResponse
	Users []UserItem `form:"users" json:"users"`
}

type UserItem struct {
	Id             uint   `form:"id" json:"id"`
	UserName       string `form:"user_name" json:"user_name"`
	Phone          string `form:"phone" json:"phone"`
	IdentifiedCode string `form:"identified_code" json:"identified_code"`
}

// @Description get user list
// @Accept json
// @Produce json
// @Param lower query int true "用戶列表的lower"
// @Param upper query int true "用戶列表的upper"
// @Success 200 {object} controller.GetUserListResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/v1/users [get]
func (controller *UserController) GetUserList(context *gin.Context) {
	var total int
	var users []orm.User
	var request GetUserListRequest
	var response GetUserListResponse
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	command := orm.Engine.Table("user")
	if request.Status > 0 {
		command = command.Where("Status = ?", request.Status)
	}
	command = command.Not("phone = ?", "system")
	if err := command.Count(&total).
		Limit(request.Upper - request.Lower + 1).Offset(request.Lower).
		Find(&users).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	response.Users = make([]UserItem, 0)
	for _, user := range users {
		response.Users = append(response.Users, UserItem{
			Id:             user.ID,
			UserName:       user.UserName,
			Phone:          user.Phone,
			IdentifiedCode: user.IdentifiedCode,
		})
	}
	response.Lower = request.Lower
	response.Upper = request.Upper
	response.Total = total
	context.JSON(http.StatusOK, &response)
}
