package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

var pic []byte

type UserController struct {
	FileController
}

type response struct {
	Error    bool   `json:"error"`
	Message  string `json:"message"`
	FuncName string `json:"func_name"`
	Data     interface{}
}

func (controller *UserController) GetUserInfo(context *gin.Context) {
	re := response{
		Error:    false,
		Message:  "nothing",
		FuncName: "GetUserInfo",
	}
	context.JSON(http.StatusOK, re)
}

func (controller *UserController) PostUserLogin(context *gin.Context) {
	var re response
	session := sessions.Default(context)
	session.Set("user_login", true)
	err := session.Save()
	if err != nil {
		re = response{
			Error:   true,
			Message: "session save fail",
		}
	} else {
		re = response{
			Error:   false,
			Message: "nothing",
		}
	}
	context.JSON(http.StatusOK, re)
}

func (controller *UserController) GetUserAvatar(context *gin.Context) {
	controller.serveBinaryFile(context, pic)
}

func (controller *UserController) PostUserAvatar(context *gin.Context) {
	data, err := controller.getBinaryFile(context)
	pic = data
	if err != nil {
		context.AbortWithError(500, err)
	}
	re := response{
		Error:   false,
		Message: "nothing",
	}
	context.JSON(http.StatusOK, re)
}
