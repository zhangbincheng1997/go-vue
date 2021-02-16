package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *sql.DB = InitMySQL()
var collection *mongo.Collection = InitMongoDB()

func getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, statusOptions)
}

func getList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	rows, err := db.Query("SELECT * FROM t_user limit ?,?", offset, limit)
	if err != nil {
		panic(err)
	}
	list := make([]User, 0)
	var user User
	for rows.Next() {
		rows.Scan(&user.ID, &user.Username, &user.Password)
		list = append(list, user)
	}
	c.JSON(http.StatusOK, list)
}

func register(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	ret, _ := db.Exec("INSERT INTO t_user (username, password) VALUES (?, ?)", username, password)
	lastInsertID, _ := ret.LastInsertId()
	rowsAffected, _ := ret.RowsAffected()
	c.JSON(http.StatusOK, gin.H{
		"lastInsertID": lastInsertID,
		"rowsAffected": rowsAffected,
	})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	ret, _ := db.Exec("DELETE FROM t_user WHERE id = ?", id)
	rowsAffected, _ := ret.RowsAffected()
	c.JSON(http.StatusOK, gin.H{
		"rowsAffected": rowsAffected,
	})
}

func updatePassword(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	ret, _ := db.Exec("UPDATE t_user SET password = ? WHERE username = ?", password, username)
	rowsAffected, _ := ret.RowsAffected()
	c.JSON(http.StatusOK, gin.H{
		"rowsAffected": rowsAffected,
	})
}

func getItemList(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	sort, _ := strconv.ParseInt(c.DefaultQuery("sort", "1"), 10, 64)
	offset := (page - 1) * limit
	keyword := c.Query("keyword")
	filter := bson.M{}
	if len(keyword) > 0 {
		filter["text"] = bson.M{"$regex": keyword}
	}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"id": sort})
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	res, err := collection.Find(context.TODO(), filter, findOptions)

	list := make([]bson.M, 0)
	for res.Next(context.TODO()) {
		var item bson.M
		err := res.Decode(&item)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, item)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": list,
	})
}

func addItem(c *gin.Context) { // 添加条目
	id, _ := strconv.Atoi(c.Request.FormValue("id"))
	text := c.Request.FormValue("text")
	property := c.Request.FormValue("property")
	item := bson.M{ // 未翻译
		"id":       id,
		"text":     text,
		"property": property,
	}
	res, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		fmt.Println(err)
		return
	}
	insertedID := res.InsertedID
	c.JSON(http.StatusOK, gin.H{
		"insertedID": insertedID,
	})
}

func addResult(c *gin.Context) { // 添加翻译
	id := c.Request.FormValue("id")
	language := c.Request.FormValue("language")
	text := c.Request.FormValue("text")
	filter := bson.M{"id": id} // 翻译
	update := bson.M{
		"$set": bson.M{
			language + ".text":   text,
			language + ".status": WAITING,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
		return
	}
	modifiedCount := res.ModifiedCount
	c.JSON(http.StatusOK, gin.H{
		"modifiedCount": modifiedCount,
	})
}

func updateStatus(c *gin.Context) {
	// ids := c.PostFormArray("ids") // TODO
	ids := [...]int{1, 2}
	language := c.Request.FormValue("language")
	status := c.Request.FormValue("status")
	filter := bson.M{
		"id": bson.M{"$in": ids},
	}
	update := bson.M{
		"$set": bson.M{
			language + ".status": status,
		},
	}
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
		return
	}
	modifiedCount := res.ModifiedCount
	c.JSON(http.StatusOK, gin.H{
		"modifiedCount": modifiedCount,
	})
}

func deleteItem(c *gin.Context) {
	var ids []int
	c.ShouldBindJSON(&ids)
	filter := bson.M{
		"id": bson.M{"$in": ids},
	}
	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return
	}
	deletedCount := res.DeletedCount
	c.JSON(http.StatusOK, gin.H{
		"deletedCount": deletedCount,
	})
}

func main() {
	r := gin.Default()
	userV1 := r.Group("/v1/user")
	{
		userV1.GET("/list", getList)
		userV1.POST("/", register)
		userV1.DELETE("/:id", deleteUser)
		userV1.PUT("/", updatePassword)
	}
	itemV1 := r.Group("/v1/item")
	{
		itemV1.GET("/", getItemList)
		itemV1.POST("/", addItem)
		itemV1.POST("/result", addResult)
		itemV1.POST("/status", updateStatus)
		itemV1.DELETE("/", deleteItem)
		// importData
		// exportData
	}
	r.Run(":8080")
}
