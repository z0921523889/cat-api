package middleware

import (
	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Execute(context *gin.Context)
}

func GetHandlerFunc(middleware Middleware) gin.HandlerFunc {
	return middleware.Execute
}
