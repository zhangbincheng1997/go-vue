package global

import (
	"main/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// 全局变量
var (
	CONFIG config.Config
	LOG    *zap.SugaredLogger
	DB     *gorm.DB
	RDB    *redis.Client
	MGO    *mongo.Database
)

// ConfigFile ...
var ConfigFile = "config.yaml"

// DataDir ...
var DataDir = "data"

// DataFile ...
var DataFile = "data.csv"

// GetAuthUser ...
func GetAuthUser(c *gin.Context) string {
	// func GetAuthUser(c *gin.Context) *model.User {
	claims, exists := c.Get(jwt.IdentityKey)
	if !exists {
		LOG.Error("JWT解析错误！！！")
		return ""
	}
	// return claims.(*model.User)
	return claims.(string)
}
