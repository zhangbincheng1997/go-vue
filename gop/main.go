package main

import (
	"main/constant"
	"main/controller"
	"main/core"
	"main/global"
	"main/middleware"
	"main/model"
	"main/utils"

	_ "main/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Swagger Example API
// @version 0.0.1
// @description roro.ishere
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	core.Viper()
	global.LOG = core.Zap()
	global.DB = core.MySQL()
	global.RDB = core.Redis()
	global.MGO = core.MongoDB()
	authMiddleware := core.JWT()

	initDB()

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	v1 := r.Group("/v1")

	userV1 := v1.Group("/user")
	{
		userV1.POST("/login", authMiddleware.LoginHandler)
		userV1.POST("/logout", authMiddleware.LogoutHandler)
		userV1.POST("/register", controller.Register)
		userV1.Use(authMiddleware.MiddlewareFunc())
		{
			userV1.GET("/info", controller.GetUserInfo)
			userV1.PUT("/password", controller.UpdatePassword)
			userV1.PUT("/info", controller.UpdateInfo)
			userV1.Use(core.GroupAuthorizator(constant.ADMIN))
			{
				userV1.GET("/list", controller.GetUserList)
				userV1.GET("/role", controller.GetRoleOptions)
				userV1.PUT("/role", controller.UpdateRole)
				userV1.DELETE("", controller.DeleteUser)
			}
		}
	}
	itemV1 := v1.Group("/item")
	{
		itemV1.Use(authMiddleware.MiddlewareFunc())
		{
			itemV1.GET("/list", controller.GetItemList)
			itemV1.GET("/status", controller.GetStatusOptions)
			itemV1.Use(core.GroupAuthorizator(constant.ADMIN))
			{
				itemV1.PUT("/text", controller.UpdateText)
				itemV1.PUT("/record/text", controller.UpdateRecordText)
				itemV1.PUT("/status", controller.UpdateStatus)
				itemV1.DELETE("", controller.DeleteItem)
				itemV1.POST("/import", controller.ImportData)
				itemV1.GET("/export", controller.ExportData)
			}
		}
	}
	uploadV1 := v1.Group("/upload")
	{
		uploadV1.Use(authMiddleware.MiddlewareFunc())
		{
			uploadV1.POST("", controller.UploadFile)
		}
	}

	// swag init
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // http://localhost:8080/swagger/index.html

	r.Run(":8080")
}

func initDB() {
	allModels := []interface{}{&model.User{}}

	if err := global.DB.AutoMigrate(allModels...); err != nil {
		global.LOG.Errorf("Failed to AutoMigrate: %v", err)
	}

	for _, m := range allModels {
		if !global.DB.Migrator().HasTable(m) {
			global.LOG.Errorf("Failed to HasTable %#v", m)
		}
	}

	users := []model.User{
		{Username: "admin", Password: utils.MD5("111111"), Role: constant.ADMIN},
		{Username: "111111", Password: utils.MD5("111111"), Role: constant.GUEST},
	}
	if err := global.DB.Create(&users).Error; err != nil {
		global.LOG.Errorf("Failed to Create: %v", err)
	}
}
