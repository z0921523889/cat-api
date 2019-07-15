package router

import (
	"cat-api/src/app/controller"
	"cat-api/src/app/env"
	"cat-api/src/app/middleware"
	"cat-api/src/app/session/store/postgres"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	//controller
	adminController          = &controller.AdminController{}
	userController           = &controller.UserController{}
	catController            = &controller.CatController{}
	catThumbnailController   = &controller.CatThumbnailController{}
	catReservationController = &controller.CatReservationController{}
	catTransferController    = &controller.CatTransferController{}
	catAdoptionController    = &controller.CatAdoptionController{}
	timePeriodController     = &controller.TimeScheduleController{}
	bannerController         = &controller.BannerController{}
	//middleware
	userAuthMiddleware = &middleware.UserAuthMiddleware{}
	adminMiddleware    = &middleware.AdminAuthMiddleware{}
)

func InitialRouterEngine() *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 1
	// set up cors
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Accept", "Content-Type", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// set up session using cookie & postgres
	store, _ := postgres.NewPostgresStore([]byte("secret"))
	store.Options(sessions.Options{
		Path:   "/",
		Domain: env.DomainName,
		MaxAge: 86400 * 30,
	})
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
	v1.GET("/identified/code/:code/check", userController.GetCheckIdentifiedCode)
	v1.POST("/user/register", userController.PostUserRegister)
	v1.POST("/user/login", userController.PostUserLogin)
	v1.POST("/admin/login", adminController.PostAdminLogin)
	// Group admin auth
	adminAuth := router.Group("api/v1").Use(middleware.GetHandlerFunc(adminMiddleware))
	//user
	adminAuth.GET("/users", userController.GetUserList)
	//cat
	adminAuth.POST("/cat", catController.PostCat)
	adminAuth.PUT("/cats/:catId", catController.PutModifyCat)
	adminAuth.POST("/cat/thumbnail", catThumbnailController.PostCatThumbnail)
	adminAuth.POST("/cat/thumbnails/:thumbnailId/cats/:catId", catThumbnailController.PostCatThumbnailBind)
	//schedules
	adminAuth.POST("/time/schedule", timePeriodController.PostTimeSchedule)
	adminAuth.POST("/time/schedules/:scheduleId/cat/:catId", timePeriodController.PostTimeScheduleCat)
	//cat_reservation
	adminAuth.PUT("/cat/adoption/owner", catReservationController.PutAdoptionCatOwner)
	//banner
	adminAuth.POST("/banner", bannerController.PostBanner)
	adminAuth.PUT("/banner/:bannerId", bannerController.PutModifyBanner)
	// Group user auth
	userAuth := router.Group("api/v1").Use(middleware.GetHandlerFunc(userAuthMiddleware))
	//user
	userAuth.GET("/user/info", userController.GetUserInfo)
	userAuth.POST("/user/logout", userController.PostUserLogout)
	//cat
	userAuth.GET("/cats", catController.GetCatList)
	userAuth.GET("/cat/thumbnails", catThumbnailController.GetCatThumbnailList)
	userAuth.GET("/cat/thumbnail/:thumbnailId", catThumbnailController.GetCatThumbnail)
	//schedules
	userAuth.GET("/time/schedules", timePeriodController.GetTimeScheduleList)
	userAuth.GET("/cats/time/schedules/:scheduleId", timePeriodController.GetTimeScheduleCat)
	//cat_reservation
	userAuth.POST("/cat/reservation", catReservationController.PostCatReservations)
	//cat_transfer
	userAuth.GET("/transfer/cats", catTransferController.GetTransferCatList)
	userAuth.GET("/certificate/transfer/:transferId", catTransferController.GetTransferCatCertificateThumbnail)
	userAuth.POST("/transfer/:transferId/certificate", catTransferController.PostTransferCatCertificate)
	userAuth.POST("/transfer/:transferId/migrate", catTransferController.PostTransferMigrate)
	//cat_adoption
	userAuth.GET("/adoption/cats", catAdoptionController.GetAdoptionCatList)
	//banner
	userAuth.GET("/banners", bannerController.GetBannerList)
	userAuth.GET("/banner/:bannerId", bannerController.GetBannerThumbnail)
	return router
}
