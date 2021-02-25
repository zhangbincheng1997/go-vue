package main

import (
	"main/constant"
	"main/controller"
	"main/core"
	"main/global"
	"main/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "main/docs"
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

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	authMiddleware := core.JWT()

	userV1 := r.Group("/v1/user")
	{
		userV1.POST("/login", authMiddleware.LoginHandler)
		userV1.POST("/logout", authMiddleware.LogoutHandler)
		userV1.POST("/register", controller.Register)
		userV1.Use(authMiddleware.MiddlewareFunc())
		{
			userV1.GET("/info", controller.UserInfo)
			userV1.PUT("/password", controller.UpdatePassword)
			userV1.PUT("/info", controller.UpdateInfo)
			userV1.Use(core.GroupAuthorizator(constant.ADMIN))
			{
				userV1.GET("/list", controller.GetUserList)
				userV1.DELETE("/", controller.DeleteUser)
			}
		}
	}
	itemV1 := r.Group("/v1/item")
	{
		itemV1.Use(authMiddleware.MiddlewareFunc())
		{
			itemV1.GET("/list", controller.GetList)
			itemV1.GET("/status", controller.GetStatus)
			itemV1.Use(core.GroupAuthorizator(constant.ADMIN))
			{
				itemV1.PUT("/text", controller.UpdateText)
				itemV1.PUT("/record/text", controller.UpdateRecordText)
				itemV1.PUT("/status", controller.UpdateStatus)
				itemV1.DELETE("/", controller.DeleteItem)
				itemV1.POST("/import", controller.ImportData)
				itemV1.GET("/export", controller.ExportData)
			}
		}
	}

	uploadV1 := r.Group("/v1/upload")
	{
		uploadV1.POST("", controller.Upload)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // http://localhost:8080/swagger/index.html

	r.Run(":8080")
}
