package middleware

import (
"github.com/gin-gonic/gin"
)

type AuthMiddleware struct{}

func (middleware *AuthMiddleware) Execute(context *gin.Context) {

}
