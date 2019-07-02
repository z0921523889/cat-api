package middleware

import (
"github.com/gin-gonic/gin"
)

type UserAuthMiddleware struct{}

func (middleware *UserAuthMiddleware) Execute(context *gin.Context) {

}

type AdminAuthMiddleware struct{}

func (middleware *AdminAuthMiddleware) Execute(context *gin.Context) {

}