package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
)

var pic []byte


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
