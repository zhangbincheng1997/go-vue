package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"main/constant"
	"main/global"
	"main/model"
	"main/model/request"
	"mime/multipart"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const idGenerator = "id_generator"

func getNextID(collection string) int {
	var generator model.Generator
	filter := bson.M{"collection": collection}
	update := bson.M{"$inc": bson.M{"id": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	if err := global.MGO.Collection(idGenerator).FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&generator); err != nil {
		global.LOG.Error("生成ID失败！！！", zap.Any("err", err))
		panic(err)
	}
	return generator.ID
}

// GetItemList ...
func GetItemList(req request.ItemPageReq) (interface{}, int64, error) {
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

	var list []interface{}
	for res.Next(context.TODO()) {
		var item interface{}
		err := res.Decode(&item)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, item)
	}
	return list, count, nil
}

// UpdateText ...
func UpdateText(req request.UpdateTextReq) error {
	filter := bson.M{"id": req.ID}
	update := bson.M{
		"$set": bson.M{"text": req.Text},
	}
	_, err := global.MGO.Collection(req.Table).UpdateOne(context.TODO(), filter, update)
	return err
}

// UpdateRecordText ...
func UpdateRecordText(req request.UpdateTextReq) error {
	filter := bson.M{"id": req.ID} // 翻译
	update := bson.M{
		"$set": bson.M{
			req.Language + ".text":   req.Text,
			req.Language + ".status": constant.WAITING,
		},
	}
	_, err := global.MGO.Collection(req.Table).UpdateOne(context.TODO(), filter, update)
	return err
}

// UpdateStatus ...
func UpdateStatus(req request.StatusReq) error {
	filter := bson.M{
		"id": bson.M{"$in": req.Ids},
	}
	update := bson.M{
		"$set": bson.M{
			req.Language + ".status": req.Status,
		},
	}
	_, err := global.MGO.Collection(req.Table).UpdateMany(context.TODO(), filter, update)
	return err
}

// DeleteItem ...
func DeleteItem(req request.DeleteReq) error {
	filter := bson.M{
		"id": bson.M{"$in": req.Ids},
	}
	_, err := global.MGO.Collection(req.Table).DeleteMany(context.TODO(), filter)
	return err
}

// ImportData ...
func ImportData(file *multipart.FileHeader, table string) error {
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
			return err
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
	_, err := global.MGO.Collection(table).InsertMany(context.TODO(), list)
	return err
}

// ExportData ...
func ExportData(filename string, table string, language string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
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
			return err
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
			content = append(content, constant.StatusMap[l["status"].(int32)])
		} else {
			content = append(content, "")
		}
		data = append(data, content)
	}
	w.WriteAll(data)
	w.Flush()
	return nil
}
