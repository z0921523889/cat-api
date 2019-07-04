package session

import (
	"cat-api/src/app/orm"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var UserSessionKey = "session_user"
var AdminSessionKey = "session_admin"

type UserSessionValue struct {
	IsLogin bool
	User    orm.Users
}

type AdminSessionValue struct {
	IsLogin bool
	Admin   orm.Admins
}

func init()  {
	gob.Register(UserSessionValue{})
	gob.Register(AdminSessionValue{})
}

func Get(context *gin.Context, key string) interface{} {
	session := sessions.Default(context)
	return session.Get(key)
}

func Set(context *gin.Context, key string, value interface{}) error {
	session := sessions.Default(context)
	session.Set(key, value)
	return session.Save()
}
