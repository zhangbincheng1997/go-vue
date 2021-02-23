package core

import (
	"main/global"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// Redis ...
func Redis() *redis.Client {
	cfg := global.CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})
	pong, err := client.Ping().Result()
	if err != nil {
		global.LOG.Error("Redis连接失败", zap.Any("err", err))
		return nil
	}
	global.LOG.Info("Redis连接成功", zap.String("pong", pong))
	return client
}
