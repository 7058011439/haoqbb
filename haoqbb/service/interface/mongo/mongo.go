package IMongo

import (
	"github.com/7058011439/haoqbb/DataBase"
)

type IMongo interface {
	GetMongoAsync(tabName string, condition interface{}, getData interface{}, index int, fun DataBase.FunFindCallBack, callbackData ...interface{})
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

func FindOne(serviceName string, tabName string, condition interface{}, getData interface{}, index int, fun DataBase.FunFindCallBack, callbackData ...interface{}) {
	db.i[serviceName].GetMongoAsync(tabName, condition, getData, index, fun, callbackData...)
}

func InsertOne(serviceName string, tabName string, data interface{}, index int) {
	db.i[serviceName].InsertMongoAsync(tabName, data, index, nil)
}

func UpdateOne(serviceName string, tabName string, condition interface{}, data interface{}, index int) {
	//var data1 reflect.Value
	//typ := reflect.ValueOf(data)
	//if typ.Kind() == reflect.Ptr {
	//	data1 = typ.Elem()
	//} else {
	//	data1 = typ
	//}
	//fmt.Println(data, data1)
	db.i[serviceName].UpdateMongoAsync(tabName, condition, data, index, nil)
}
