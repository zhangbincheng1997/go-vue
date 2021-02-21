package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// config for MySQL
const (
	USERNAME  = "root"
	PROTOCOL  = "tcp"
	HOST      = "www.littleredhat1997.com"
	PORT      = "3306"
	PASSWORD  = "Zhangbincheng0"
	DATABASE  = "test"
	CHARSET   = "utf8"
	PARSETIME = "TRUE"
	LOC       = "Local"
)

// config for Redis
const (
	ADDRREDIS     = "www.littleredhat1997.com"
	PORTREDIS     = "6379"
	PASSWORDREDIS = "Zhangbincheng0"
	DBREDIS       = 0
)

// config for MongoDB
const (
	HOSTMONGO     = "120.79.157.49"
	PORTMONGO     = "27017"
	DATABASEMONGO = "test"
	ITEM          = "item"
	TEXT          = "text"
	IMAGE         = "image"
	IDGENERATOR   = "id_generator"
)

// InitMySQL init
func InitMySQL() *gorm.DB {
	dbDSN := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s", USERNAME, PASSWORD, PROTOCOL, HOST, PORT, DATABASE, CHARSET, PARSETIME, LOC)
	dialector := mysql.New(mysql.Config{
		DSN:                       dbDSN, // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,   // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,  // disable datetime precision support, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // smart configure based on used version
	})
	config := &gorm.Config{}
	db, err := gorm.Open(dialector, config)
	if err != nil {
		log.Println("dbDSN: " + dbDSN)
		panic("MySQL数据源配置不正确: " + err.Error())
	}
	// db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, &Role{}, &UserRole{})
	return db
}

// InitRedis init
func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", ADDRREDIS, PORTREDIS),
		Password: PASSWORDREDIS,
		DB:       0,
	})
	return rdb
}

// InitMongoDB init
func InitMongoDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := fmt.Sprintf("mongodb://%s:%s", HOSTMONGO, PORTMONGO)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("uri: " + uri)
		panic("Redis数据源配置不正确: " + err.Error())
	}
	mgo := client.Database(DATABASEMONGO)
	return mgo
}
