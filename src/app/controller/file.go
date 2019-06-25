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
