package core

import (
	"main/global"

	"github.com/go-redis/redis"
)

// Redis ...
func Redis() *redis.Client {
	cfg := global.CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		global.LOG.Errorf("Redis连接失败：%v", err)
		return nil
	}
	global.LOG.Infof("Redis连接成功：%v", cfg.Addr)
	return client
}
