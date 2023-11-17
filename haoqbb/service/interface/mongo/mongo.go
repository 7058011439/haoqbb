package IMongo

import (
	"github.com/7058011439/haoqbb/DataBase"
)

type IMongo interface {
	GetMongoAsync(tabName string, condition interface{}, getData interface{}, index int, fun DataBase.FunFindCallBack, callbackData ...interface{})
	GetMongoSync(tabName string, condition interface{}, getData interface{})
	InsertMongoAsync(tabName string, data interface{}, index int, fun DataBase.FunUpdateCallBack, callBackData ...interface{})
	UpdateMongoAsync(tabName string, condition interface{}, data interface{}, index int, fun DataBase.FunUpdateCallBack, callBackData ...interface{})
	GetName() string
}

type mongoDB struct {
	i map[string]IMongo
}

var db = mongoDB{i: make(map[string]IMongo)}

func SetMongoAgent(d IMongo) {
	db.i[d.GetName()] = d
}

func FindOneSync(serviceName string, tabName string, condition interface{}, getData interface{}) {
	db.i[serviceName].GetMongoSync(tabName, condition, getData)
}

func FindOne(serviceName string, tabName string, condition interface{}, getData interface{}, index int, fun DataBase.FunFindCallBack, callbackData ...interface{}) {
	db.i[serviceName].GetMongoAsync(tabName, condition, getData, index, fun, callbackData...)
}

func InsertOne(serviceName string, tabName string, data interface{}, index int) {
	db.i[serviceName].InsertMongoAsync(tabName, data, index, nil)
}

func UpdateOne(serviceName string, tabName string, condition interface{}, data interface{}, index int, fun DataBase.FunUpdateCallBack, callbackData ...interface{}) {
	db.i[serviceName].UpdateMongoAsync(tabName, condition, data, index, fun, callbackData...)
}
