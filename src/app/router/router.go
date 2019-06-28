package router

import (
	"cat-api/src/app/controller"
	"cat-api/src/app/middleware"
	"cat-api/src/app/session/store/postgres"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine
var userController = &controller.UserController{}
var catController = &controller.CatController{}
var timePeriodController = &controller.TimePeriodController{}
//middleware
var authMiddleware = &middleware.AuthMiddleware{}

func init() {
	Router = gin.Default()
	Router.MaxMultipartMemory = 1
	// set up session using cookie
	store := postgres.NewPostgresStore([]byte("secret"))
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

	auth.GET("/cat", catController.GetCat)
	auth.POST("/cat", catController.PostCat)
	auth.PUT("/cat/:catId", catController.PutModifyCat)
	auth.GET("/cat/:catId/thumbnail", catController.GetCatThumbnail)
	auth.POST("/cat/:catId/thumbnail", catController.PostCatThumbnail)

	auth.POST("/time/period", timePeriodController.PostTimePeriod)

}
