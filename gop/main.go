package main

import (
	"main/controller"
	"main/core"
	"main/global"
	"main/middleware"

	"github.com/gin-gonic/gin"
)

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

	// Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	userV1 := r.Group("/v1/user")
	{
		userV1.POST("/login", authMiddleware.LoginHandler)
		userV1.POST("/logout", authMiddleware.LogoutHandler)
		userV1.POST("/register", controller.Register)
		userV1.Use(authMiddleware.MiddlewareFunc())
		{
			userV1.GET("/info", controller.Info)
			userV1.GET("/list", controller.GetUserList)
			userV1.DELETE("/:id", controller.DeleteUser)
			userV1.PUT("/password", controller.UpdatePassword)
			userV1.PUT("/info", controller.UpdateInfo)
		}
	}
	itemV1 := r.Group("/v1/item")
	{
		itemV1.Use(authMiddleware.MiddlewareFunc())
		{
			itemV1.GET("/list", controller.GetList)
			itemV1.GET("/status", controller.GetStatus)
			itemV1.PUT("/text", controller.UpdateText)
			itemV1.PUT("/record/text", controller.UpdateRecordText)
			itemV1.PUT("/status", controller.UpdateStatus)
			itemV1.DELETE("", controller.DeleteItem)
			itemV1.POST("/import", controller.ImportData)
			itemV1.GET("/export", controller.ExportData)
		}
	}

	r.Run(":8080")
}
