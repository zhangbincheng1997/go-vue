package global

import (
	"main/config"
	"main/model"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// 全局变量
var (
	DB     *gorm.DB
	RDB    *redis.Client
	MGO    *mongo.Database
	LOG    *zap.Logger
	CONFIG config.Config
)

// ConfigFile ...
var ConfigFile = "config.yaml"

// GetAuthUser ...
func GetAuthUser(c *gin.Context) *model.User {
	claims, exists := c.Get(jwt.IdentityKey)
	if !exists {
		LOG.Error("JWT解析错误！！！")
		return nil
	}
	return claims.(*model.User)
}
