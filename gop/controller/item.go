package controller

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"main/global"
	"main/model"
	"main/model/request"
	"main/model/response"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const idGenerator = "id_generator"

var statusOptions = []model.Status{
	{ID: 1, Desc: "待定"},
	{ID: 2, Desc: "失败"},
	{ID: 3, Desc: "成功"},
}

var statusMap = map[int32]string{
	1: "待定",
	2: "失败",
	3: "成功",
}

// Status const
const (
	_ = iota
	WAITING
	FAILURE
	SUCCESS
)

func getNextID(collection string) int {
	var generator model.Generator
	filter := bson.M{"collection": collection}
	update := bson.M{"$inc": bson.M{"id": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	if err := global.MGO.Collection(idGenerator).FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&generator); err != nil {
		global.LOG.Error("生成ID失败", zap.Any("err", err))
		panic(err)
	}
	return generator.ID
}

// GetList ...
func GetList(c *gin.Context) {
	var req request.ItemPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
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
	findOptions.SetLimit(int64(req.Limit))
	findOptions.SetSkip(int64(offset))
	count, _ := global.MGO.Collection(req.Table).CountDocuments(context.TODO(), filter)
	res, _ := global.MGO.Collection(req.Table).Find(context.TODO(), filter, findOptions)

	list := make([]bson.M, 0)
	for res.Next(context.TODO()) {
		var item bson.M
		err := res.Decode(&item)
		if err != nil {
			panic(err)
		}
		list = append(list, item)
	}
	response.OkWithData(c, gin.H{
		"list":  list,
		"count": count,
	})
}

// GetStatus ...
func GetStatus(c *gin.Context) {
	response.OkWithData(c, statusOptions)
}

// UpdateText ...
func UpdateText(c *gin.Context) { // 添加翻译
	var req request.UpdateTextReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	filter := bson.M{"id": req.ID} // 翻译
	update := bson.M{
		"$set": bson.M{"text": req.Text},
	}
	res, err := global.MGO.Collection(req.Table).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	modifiedCount := res.ModifiedCount
	response.OkWithData(c, modifiedCount)
}

// UpdateRecordText ...
func UpdateRecordText(c *gin.Context) { // 添加翻译
	var req request.UpdateTextReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	filter := bson.M{"id": req.ID} // 翻译
	update := bson.M{
		"$set": bson.M{
			req.Language + ".text":   req.Text,
			req.Language + ".status": WAITING,
		},
	}
	res, err := global.MGO.Collection(req.Table).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	modifiedCount := res.ModifiedCount
	response.OkWithData(c, modifiedCount)
}

// UpdateStatus ...
func UpdateStatus(c *gin.Context) {
	var req request.StatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, err.Error())
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
	res, err := global.MGO.Collection(req.Table).UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	modifiedCount := res.ModifiedCount
	response.OkWithData(c, modifiedCount)
}

// DeleteItem ...
func DeleteItem(c *gin.Context) {
	var req request.DeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	filter := bson.M{
		"id": bson.M{"$in": req.Ids},
	}
	res, err := global.MGO.Collection(req.Table).DeleteMany(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	deletedCount := res.DeletedCount
	response.OkWithData(c, deletedCount)
}

// ImportData ...
func ImportData(c *gin.Context) {
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
	res, err := global.MGO.Collection(table).InsertMany(context.TODO(), list)
	if err != nil {
		panic(err)
	}
	insertedIDs := res.InsertedIDs
	response.OkWithData(c, insertedIDs)
}

// ExportData ...
func ExportData(c *gin.Context) {
	table := c.Query("table")
	language := c.Query("language")
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
	res, _ := global.MGO.Collection(table).Find(context.TODO(), filter, findOptions)

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
