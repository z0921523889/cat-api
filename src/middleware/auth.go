package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

type AuthMiddleware struct{}

func (middleware *AuthMiddleware) Execute(context *gin.Context) {
	log.Println("AuthMiddleware Execute")
	session := sessions.Default(context)
	success := session.Get("user_login").(bool)
	if !success {
		panic("login auth fail")
	}
}
