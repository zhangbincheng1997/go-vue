package global

import (
	"main/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

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
