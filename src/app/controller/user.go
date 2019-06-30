package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
)

var pic []byte

type UserController struct {
	FileController
}

func (controller *UserController) GetUserInfo(context *gin.Context) {
	session := sessions.Default(context)
	login := session.Get("user_login")
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("login state : %t", login)})
}

func (controller *UserController) PostUserLogin(context *gin.Context) {
	session := sessions.Default(context)
	session.Set("user_login", true)
	session.Save()
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("login success")})
}

func (controller *UserController) GetUserAvatar(context *gin.Context) {
	controller.serveBinaryFile(context, pic)
}

func (controller *UserController) PostUserAvatar(context *gin.Context) {
	data, err := controller.getBinaryDataFromBody(context)
	pic = data
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("post user avatar success")})
}
