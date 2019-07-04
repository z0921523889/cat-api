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
	var user orm.Users
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if orm.Engine.Where("phone = ?", request.Phone).
		Where("password = ?", request.Password).
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
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("login success")})
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
	if !orm.Engine.Where("phone = ?", request.Phone).
		First(&orm.Users{}).RecordNotFound() {
		httputil.NewError(context, http.StatusBadRequest, errors.New("phone already exist"))
		return
	}
	hash := sha1.New()
	hash.Write([]byte(request.Phone + time.Now().String()))
	code := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	user := orm.Users{
		Phone:            request.Phone,
		UserName:         request.UserName,
		Password:         request.Password,
		SecurityPassword: request.SecurityPassword,
		IdentifiedCode:   code[:15],
		Status:           1,
	}
	if err := orm.Engine.Create(&user).Error; err != nil {
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
	var users []orm.Users
	var request GetUserListRequest
	var usersResponse []UserItem
	if err := context.Bind(&request); err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	if err := orm.Engine.Table("users").Count(&total).
		Limit(request.Upper - request.Lower + 1).Offset(request.Lower).
		Find(&users).Error; err != nil {
		httputil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	for _, user := range users {
		usersResponse = append(usersResponse, UserItem{
			Id:             user.ID,
			UserName:       user.UserName,
			Phone:          user.Phone,
			IdentifiedCode: user.IdentifiedCode,
		})
	}
	context.JSON(http.StatusOK, &GetUserListResponse{
		ListResponse: ListResponse{
			Lower: request.Lower,
			Upper: request.Upper,
			Total: total,
		},
		Users: usersResponse,
	})
}
