package controller

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type FileController struct {
}

func (controller *FileController) serveBinaryFile(context *gin.Context, data []byte) {
	context.Data(http.StatusOK, "image/jpeg", data)
}

func (controller *FileController) getBinaryFile(context *gin.Context) ([]byte, error) {
	body := context.Request.Body
	return ioutil.ReadAll(body)
}

type ListRequest struct {
	Lower int `form:"lower" json:"lower"`
	Upper int `form:"upper" json:"upper"`
}

type ListResponse struct {
	Lower int         `form:"lower" json:"lower"`
	Upper int         `form:"upper" json:"upper"`
	Total int         `form:"total" json:"total"`
}

type Message struct {
	Message string `json:"message" example:"message"`
}