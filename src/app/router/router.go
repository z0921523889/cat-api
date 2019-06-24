package router

import (
	"cat-api/src/app/controller"
	"cat-api/src/app/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine
var userController = &controller.UserController{}
//middleware
var authMiddleware = &middleware.AuthMiddleware{}

func init() {
	Router = gin.Default()
	Router.MaxMultipartMemory = 1
	// set up session using cookie
	store := cookie.NewStore([]byte("secret"))
	Router.Use(sessions.Sessions("cat_session", store))
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	Router.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	Router.Use(gin.Recovery())
	//Group v1
	v1 := Router.Group("api/v1")
	// Group v1 non auth
	v1.POST("/user/login", userController.PostUserLogin)
	// Group v1 auth
	auth := v1.Use(middleware.GetHandlerFunc(authMiddleware))
	auth.GET("/user/info", userController.GetUserInfo)
	auth.GET("/user/avatar", userController.GetUserAvatar)
	auth.POST("/user/avatar", userController.PostUserAvatar)
	//router.POST("/somePost", posting)
	//router.PUT("/somePut", putting)
	//router.DELETE("/someDelete", deleting)
	//router.PATCH("/somePatch", patching)
	//router.HEAD("/someHead", head)
	//router.OPTIONS("/someOptions", options)
}
