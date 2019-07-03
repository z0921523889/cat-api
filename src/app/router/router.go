package router

import (
	"cat-api/src/app/controller"
	"cat-api/src/app/middleware"
	"cat-api/src/app/session/store/postgres"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var adminController = &controller.AdminController{}
var userController = &controller.UserController{}
var catController = &controller.CatController{}
var catThumbnailController = &controller.CatThumbnailController{}
var timePeriodController = &controller.TimeScheduleController{}
//middleware
var userAuthMiddleware = &middleware.UserAuthMiddleware{}
var adminMiddleware = &middleware.AdminAuthMiddleware{}

func InitialRouterEngine() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 1
	// set up session using cookie
	store, _ := postgres.NewPostgresStore([]byte("secret"))
	router.Use(sessions.Sessions("cat_session", store))
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())
	//Group v1
	v1 := router.Group("api/v1")
	// Group v1 non auth
	v1.POST("/user/register", userController.PostUserRegister)
	v1.POST("/user/login", userController.PostUserLogin)
	v1.POST("/admin/login", adminController.PostAdminLogin)
	// Group admin auth
	adminAuth := v1.Use(middleware.GetHandlerFunc(adminMiddleware))
	//cat
	adminAuth.POST("/cat", catController.PostCat)
	adminAuth.PUT("/cats/:catId", catController.PutModifyCat)
	adminAuth.POST("/cat/thumbnail", catThumbnailController.PostCatThumbnail)
	adminAuth.PUT("/thumbnails/:thumbnailId/cats/:catId", catThumbnailController.PutModifyCatThumbnail)
	//schedules
	adminAuth.POST("/time/schedule", timePeriodController.PostTimeSchedule)
	adminAuth.POST("/time/schedules/:scheduleId/cat/:catId", timePeriodController.PostTimeScheduleCat)
	// Group user auth
	userAuth := v1.Use(middleware.GetHandlerFunc(userAuthMiddleware))
	userAuth.GET("/user/info", userController.GetUserInfo)
	//cat
	userAuth.GET("/cats", catController.GetCatList)
	userAuth.GET("/cat/thumbnails", catThumbnailController.GetCatThumbnailList)
	userAuth.GET("/thumbnail/cats/:catId", catThumbnailController.GetCatThumbnail)
	//schedules
	userAuth.GET("/time/schedules", timePeriodController.GetTimeScheduleList)
	userAuth.GET("/time/schedules/:scheduleId", timePeriodController.GetTimeScheduleCat)
	return router
}
