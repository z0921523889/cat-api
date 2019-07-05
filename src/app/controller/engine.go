package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)

type FileController struct {
}

func (controller *FileController) serveBinaryFile(context *gin.Context, data []byte) {
	context.Data(http.StatusOK, "image/jpeg", data)
}

func (controller *FileController) getBinaryDataFromBody(context *gin.Context) ([]byte, error) {
	body := context.Request.Body
	return ioutil.ReadAll(body)
}

func (controller *FileController) getBinaryDataFromForm(context *gin.Context, fileName string) ([]byte, error) {
	file, _, err := context.Request.FormFile(fileName)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type ListRequest struct {
	Lower int `form:"lower" json:"lower"`
	Upper int `form:"upper" json:"upper"`
}

type ListResponse struct {
	Lower int `form:"lower" json:"lower"`
	Upper int `form:"upper" json:"upper"`
	Total int `form:"total" json:"total"`
}

type Message struct {
	Message string `json:"message" example:"message"`
}

type EngineController struct {
}

func (controller *EngineController) GetTest(context *gin.Context) {
	context.JSON(http.StatusOK, Message{Message: fmt.Sprintf("test on api engine...")})
}
