package middleware

import (
	"cat-api/src/app/session"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
)

type UserAuthMiddleware struct{}

func (middleware *UserAuthMiddleware) Execute(context *gin.Context) {
	userSessionData := session.Get(context, session.UserSessionKey)
	if userSessionData != nil {
		userSessionValue := userSessionData.(session.UserSessionValue)
		if userSessionValue.IsLogin && userSessionValue.User.Status == 1 {
			return
		}
		httputil.NewError(context, http.StatusUnauthorized, errors.New("user has been blocked"))
	}else{
		httputil.NewError(context, http.StatusUnauthorized, errors.New("user does not authorized"))
	}
	context.Abort()
}

type AdminAuthMiddleware struct{}

func (middleware *AdminAuthMiddleware) Execute(context *gin.Context) {
	sessionData := session.Get(context, session.AdminSessionKey)
	if sessionData != nil {
		sessionValue := sessionData.(session.AdminSessionValue)
		if sessionValue.IsLogin {
			return
		}
	}
	httputil.NewError(context, http.StatusUnauthorized, errors.New("admin does not authorized"))
	context.Abort()
}
