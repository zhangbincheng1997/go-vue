package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

var logger = InitLogger()
var db *gorm.DB = InitMySQL()
var database *mongo.Database = InitMongoDB()
var collection *mongo.Collection = database.Collection(TEXT)
var idGenerator *mongo.Collection = database.Collection(IDGENERATOR)

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
	user := make(map[string]interface{})
	user["roles"] = "admin"
	user["introduction"] = "I am a superadministrator"
	user["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	user["name"] = "Super Admin"
	success(c, user)
}

func login(c *gin.Context) {
	// username := c.Request.FormValue("username")
	// password := c.Request.FormValue("password")
	// ret, _ := db.Exec("INSERT INTO t_user (username, password) VALUES (?, ?)", username, password)
	success(c, "token...")
}

func logout(c *gin.Context) {
	success(c, "")
}

func register(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	ret, _ := db.Exec("INSERT INTO t_user (username, password) VALUES (?, ?)", username, password)
	rowsAffected, _ := ret.RowsAffected()
	success(c, rowsAffected)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	ret, _ := db.Exec("DELETE FROM t_user WHERE id = ?", id)
	rowsAffected, _ := ret.RowsAffected()
	success(c, rowsAffected)
}

func updatePassword(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	ret, _ := db.Exec("UPDATE t_user SET password = ? WHERE username = ?", password, username)
	rowsAffected, _ := ret.RowsAffected()
	success(c, rowsAffected)
}

func getNextID(collection string) int {
	var generator Generator
	filter := bson.M{"collection": collection}
	update := bson.M{"$inc": bson.M{"id": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	err := idGenerator.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&generator)
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
	count, _ := database.Collection(req.Table).CountDocuments(context.TODO(), filter)
	res, _ := database.Collection(req.Table).Find(context.TODO(), filter, findOptions)

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
	res, err := database.Collection(req.Table).UpdateOne(context.TODO(), filter, update)
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
	res, err := database.Collection(req.Table).UpdateOne(context.TODO(), filter, update)
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
	res, err := database.Collection(req.Table).UpdateMany(context.TODO(), filter, update)
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
	res, err := database.Collection(req.Table).DeleteMany(context.TODO(), filter)
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
	res, err := database.Collection(table).InsertMany(context.TODO(), list)
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
	res, _ := database.Collection(table).Find(context.TODO(), filter, findOptions)

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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))
	// r.Use(cors.Default())

	userV1 := r.Group("/v1/user")
	{
		userV1.GET("/info", info)
		userV1.POST("/login", login)
		userV1.POST("/logout", logout)
		userV1.POST("/register", register)
		userV1.DELETE("/:id", deleteUser)
		userV1.PUT("/password", updatePassword)
	}
	itemV1 := r.Group("/v1/item")
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

	r.Run(":8080")
}
