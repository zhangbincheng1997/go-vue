package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

var logger = InitLogger()
var db *gorm.DB = InitMySQL()
var rdb *redis.Client = InitRedis()
var mgo *mongo.Database = InitMongoDB()

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "ok",
		"data":    data,
	})
}

func failure(c *gin.Context, message interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": message,
	})
}

func info(c *gin.Context) {
	user := getAuthUser(c)
	success(c, user)
}

func register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 50000, "message": err.Error()})
		return
	}
	if db.Where("username = ?", req.Username).Take(&User{}).Error != nil {
		user := User{
			Username:     req.Username,
			Password:     req.Password,
			Role:         "admin",
			Introduction: "I am a super administrator",
			Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			Name:         "Super Admin",
		}
		res := db.Create(&user)
		if err := res.Error; err != nil {
			logger.Error(fmt.Sprintf("%v", err))
			failure(c, err)
			return
		}
		success(c, res.RowsAffected)
		return
	}
	failure(c, "用户名已存在")
}

func getUserList(c *gin.Context) {
	var req UserPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 50000, "message": err.Error()})
		return
	}
	var count int64
	db.Model(User{}).Count(&count)
	var list []User
	offset := (req.Page - 1) * req.Limit
	if req.Sort {
		db.Order("id desc")
	}
	db.Limit(req.Limit).Offset(offset).Find(&list)
	success(c, gin.H{
		"list":  list,
		"count": count,
	})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	res := db.Where("username = ?", id).Delete(User{})
	if err := res.Error; err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		failure(c, err)
		return
	}
	success(c, res.RowsAffected)
}

func updateInfo(c *gin.Context) {
	var req UpdateInfoReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 50000, "message": err.Error()})
		return
	}
	user := getAuthUser(c)
	copier.Copy(&user, &req)
	res := db.Model(&user).Updates(&user) // 不会更新空值
	if err := res.Error; err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		failure(c, err)
		return
	}
	rdb.Del(user.Username)
	success(c, res.RowsAffected)
}

func updatePassword(c *gin.Context) {
	var req UpdatePasswordReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 50000, "message": err.Error()})
		return
	}
	user := getAuthUser(c)
	res := db.Model(User{}).Where("username = ? and password = ?", user.Username, req.OldPwd).Update("password", req.NewPwd)
	if err := res.Error; err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		failure(c, err)
		return
	}
	success(c, res.RowsAffected)
}

func getNextID(collection string) int {
	var generator Generator
	filter := bson.M{"collection": collection}
	update := bson.M{"$inc": bson.M{"id": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	err := mgo.Collection(IDGENERATOR).FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&generator)
	if err != nil {
		logger.Info(fmt.Sprint(err))
		panic(err)
	}
	return generator.ID
}

func getList(c *gin.Context) {
	var req PageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 50000, "message": err.Error()})
		return
	}
	logger.Info(fmt.Sprintf("%v", req))
	offset := (req.Page - 1) * req.Limit
	filter := bson.M{}
	if req.Status != 0 {
		filter[req.Language+".status"] = req.Status
	}
	if len(req.Keyword) > 0 {
		filter["text"] = bson.M{"$regex": req.Keyword}
	}
	findOptions := options.Find()
	if req.Sort != 0 {
		findOptions.SetSort(bson.M{"id": req.Sort})
	}
	findOptions.SetLimit(req.Limit)
	findOptions.SetSkip(offset)
	count, _ := mgo.Collection(req.Table).CountDocuments(context.TODO(), filter)
	res, _ := mgo.Collection(req.Table).Find(context.TODO(), filter, findOptions)

	list := make([]bson.M, 0)
	for res.Next(context.TODO()) {
		var item bson.M
		err := res.Decode(&item)
		if err != nil {
			panic(err)
		}
		list = append(list, item)
	}
	success(c, gin.H{
		"list":  list,
		"count": count,
	})
}

func getStatus(c *gin.Context) {
	success(c, statusOptions)
}

func updateText(c *gin.Context) { // 添加翻译
	var req UpdateTextReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 50000, "message": err.Error()})
		return
	}
	filter := bson.M{"id": req.ID} // 翻译
	update := bson.M{
		"$set": bson.M{"text": req.Text},
	}
	res, err := mgo.Collection(req.Table).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	modifiedCount := res.ModifiedCount
	success(c, modifiedCount)
}

func updateRecordText(c *gin.Context) { // 添加翻译
	var req UpdateTextReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 50000, "message": err.Error()})
		return
	}
	filter := bson.M{"id": req.ID} // 翻译
	update := bson.M{
		"$set": bson.M{
			req.Language + ".text":   req.Text,
			req.Language + ".status": WAITING,
		},
	}
	res, err := mgo.Collection(req.Table).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	modifiedCount := res.ModifiedCount
	success(c, modifiedCount)
}

func updateStatus(c *gin.Context) {
	var req StatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 50000, "message": err.Error()})
		return
	}
	filter := bson.M{
		"id": bson.M{"$in": req.Ids},
	}
	update := bson.M{
		"$set": bson.M{
			req.Language + ".status": req.Status,
		},
	}
	res, err := mgo.Collection(req.Table).UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	modifiedCount := res.ModifiedCount
	success(c, modifiedCount)
}

func deleteItem(c *gin.Context) {
	var req DeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 50000, "message": err.Error()})
		return
	}
	filter := bson.M{
		"id": bson.M{"$in": req.Ids},
	}
	res, err := mgo.Collection(req.Table).DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	deletedCount := res.DeletedCount
	success(c, deletedCount)
}

func importData(c *gin.Context) {
	file, _ := c.FormFile("file")
	table := c.Request.FormValue("table")
	csvFile, _ := file.Open()
	defer csvFile.Close()
	r := csv.NewReader(csvFile)

	// 解析CSV
	var list []interface{}
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		text := line[0]
		property := line[1]
		item := bson.M{
			"id":       getNextID(table),
			"text":     text,
			"property": property,
		}
		list = append(list, item)
	}
	res, err := mgo.Collection(table).InsertMany(context.TODO(), list)
	if err != nil {
		panic(err)
	}
	insertedIDs := res.InsertedIDs
	success(c, insertedIDs)
}

func exportData(c *gin.Context) {
	table := c.Query("table")
	language := c.Query("language")
	logger.Info(table + " " + language)
	dir := "data"
	_filename := "data.csv"
	filename := path.Join(dir, _filename)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(file)

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"id": 1})
	res, _ := mgo.Collection(table).Find(context.TODO(), filter, findOptions)

	data := [][]string{}
	for res.Next(context.TODO()) {
		var item map[string]interface{}
		err := res.Decode(&item)
		if err != nil {
			panic(err)
		}
		content := []string{}
		content = append(content, fmt.Sprint(item["id"].(int32)))
		content = append(content, item["text"].(string))
		content = append(content, item["property"].(string))
		l := item[language].(map[string]interface{})
		if l["text"] != nil {
			content = append(content, l["text"].(string))
		} else {
			content = append(content, "")
		}
		if l["status"] != nil {
			content = append(content, statusMap[l["status"].(int32)])
		} else {
			content = append(content, "")
		}
		data = append(data, content)
	}
	w.WriteAll(data)
	w.Flush()
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", _filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(filename)
}

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           time.Hour * 24,
	}))

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "go-vue",
		Key:         []byte("roro"),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			username := claims[identityKey].(string)

			val, err := rdb.Get(username).Result()
			if err != nil {
				var user User
				res := db.Where("username = ?", username).First(&user) // MySQL
				if err := res.Error; err != nil {
					return nil
				}
				user.Password = "" // 敏感数据
				jsonStr, _ := json.Marshal(user)
				res2 := rdb.Set(username, string(jsonStr), 0) // Redis
				if err := res2.Err(); err != nil {
					return nil
				}
				return &user
			}
			var user User
			_ = json.Unmarshal([]byte(val), &user)
			return &user
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var req LoginReq
			if err := c.ShouldBind(&req); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			username := req.Username
			password := req.Password
			var user User
			res := db.Where("username = ? and password = ?", username, password).First(&user)
			if err := res.Error; err != nil {
				logger.Error(fmt.Sprintf("%v", err))
				return nil, jwt.ErrFailedAuthentication
			}
			return &User{Username: username}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			logger.Info(fmt.Sprint(data.(*User)))
			if v, ok := data.(*User); ok && v.Role == "admin" {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	userV1 := r.Group("/v1/user")
	{
		userV1.POST("/login", authMiddleware.LoginHandler)
		userV1.POST("/logout", authMiddleware.LogoutHandler)
		userV1.POST("/register", register)
		userV1.Use(authMiddleware.MiddlewareFunc())
		{
			userV1.GET("/info", info)
			userV1.GET("/list", getUserList)
			userV1.DELETE("/:id", deleteUser)
			userV1.PUT("/password", updatePassword)
			userV1.PUT("/info", updateInfo)
		}
	}
	itemV1 := r.Group("/v1/item")
	{
		itemV1.Use(authMiddleware.MiddlewareFunc())
		{
			itemV1.GET("/list", getList)
			itemV1.GET("/status", getStatus)
			itemV1.PUT("/text", updateText)
			itemV1.PUT("/record/text", updateRecordText)
			itemV1.PUT("/status", updateStatus)
			itemV1.DELETE("", deleteItem)
			itemV1.POST("/import", importData)
			itemV1.GET("/export", exportData)
		}
	}

	r.Run(":8080")
}
