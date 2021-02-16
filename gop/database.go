package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// config for MySQL
const (
	USERNAME = "root"
	PROTOCOL = "tcp"
	HOST     = "www.littleredhat1997.com"
	PORT     = "3306"
	PASSWORD = "Zhangbincheng0"
	DATABASE = "test"
	CHARSET  = "utf8"
)

// config for MongoDB
const (
	HOSTMONGO       = "120.79.157.49"
	PORTMONGO       = "27017"
	DATABASEMONGO   = "test"
	COLLECTIONMONGO = "item"
)

// InitMySQL init
func InitMySQL() *sql.DB {
	dbDSN := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s", USERNAME, PASSWORD, PROTOCOL, HOST, PORT, DATABASE, CHARSET)
	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		log.Println("dbDSN: " + dbDSN)
		panic("MySQL数据源配置不正确: " + err.Error())
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

// InitMongoDB init
func InitMongoDB() *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := fmt.Sprintf("mongodb://%s:%s", HOSTMONGO, PORTMONGO)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("uri: " + uri)
		panic("Redis数据源配置不正确: " + err.Error())
	}
	database := client.Database(DATABASEMONGO)
	collection := database.Collection(COLLECTIONMONGO)
	return collection
}
