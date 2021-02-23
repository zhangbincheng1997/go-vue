package core

import (
	"context"
	"fmt"
	"main/global"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// MongoDB ...
func MongoDB() *mongo.Database {
	cfg := global.CONFIG.MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := fmt.Sprintf("mongodb://%s", cfg.Addr)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		global.LOG.Error("MongoDB连接失败", zap.Any("err", err))
		return nil
	}
	mgo := client.Database(cfg.Database)
	global.LOG.Info("MongoDB连接成功", zap.String("uri", uri))
	return mgo
}
