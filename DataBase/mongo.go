package DataBase

import (
	"Core/Log"
	"Core/Stl"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type FunUpdateCallBack func(callbackData ...interface{})
type FunFindCallBack func(getData interface{}, callbackData ...interface{})

const (
	operateInsertOne = iota
	operateUpdateOne
	operateDeleteOne
	operateFindOne
	operateInsertMany
	operateUpdateMany
	operateDeleteMany
	operateFindMany
)

type MongoMessage struct {
	tabName      string
	condition    interface{}
	newData      interface{}
	callbackData []interface{}
	funUpdate    FunUpdateCallBack
	funFind      FunFindCallBack
	operate      int
}

type MongoDB struct {
	client     *mongo.Client
	database   *mongo.Database
	queue      []*Stl.Queue
	queueCount int
	queueIndex int
}

func NewMongoDB(ip string, port int, dbName string, username string, password string, queueCount int) *MongoDB {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", ip, port))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		Log.ErrorLog("Failed to NewMongoDB, connect error, ip = %v, port = %v, err = %v", ip, port, err)
		return nil
	}

	database := client.Database(dbName)
	if database == nil {
		Log.ErrorLog("Failed to NewMongoDB, db is nil, dbName = %v", dbName)
		return nil
	}
	if queueCount < 1 {
		queueCount = 200
	}
	ret := &MongoDB{
		client:     client,
		database:   database,
		queueCount: queueCount,
	}
	for i := 0; i < queueCount; i++ {
		ret.queue = append(ret.queue, Stl.NewQueue())
		go ret.exec(ret.queue[i])
	}

	go ret.HeartBeat()

	return ret
}

func (m *MongoDB) HeartBeat() {
	tick := time.NewTicker(time.Second * 10)
	for {
		<-tick.C
		if m.client != nil {
			err := m.client.Ping(context.TODO(), nil)
			if err != nil {
				Log.ErrorLog("Failed to HeartBeat, err = %v", err)
			}
		} else {
			Log.ErrorLog("Failed to HeartBeat, client is nil")
			break
		}
	}
}

func (m *MongoDB) exec(queue *Stl.Queue) {
	if m.client != nil {
		tick := time.NewTicker(time.Millisecond)
		for {
			<-tick.C
			for queue.Head() != nil {
				item := queue.Dequeue().(*MongoMessage)
				switch item.operate {
				case operateInsertOne:
					m.insertOne(item.tabName, item.newData, item.callbackData, item.funUpdate)
				case operateUpdateOne:
					m.updateOne(item.tabName, item.condition, item.newData, item.callbackData, item.funUpdate)
				case operateDeleteOne:
					m.deleteOne(item.tabName, item.condition, item.callbackData, item.funUpdate)
				case operateFindOne:
					m.findOne(item.tabName, item.condition, item.newData, item.callbackData, item.funFind)
				case operateInsertMany:
					m.insertMany(item.tabName, item.newData, item.callbackData, item.funUpdate)
				case operateUpdateMany:
					m.updateMany(item.tabName, item.condition, item.newData, item.callbackData, item.funUpdate)
				case operateDeleteMany:
					m.deleteMany(item.tabName, item.condition, item.callbackData, item.funUpdate)
				case operateFindMany:
					m.findMany(item.tabName, item.condition, item.newData, item.callbackData, item.funFind)
				}
			}
		}
	} else {
		Log.ErrorLog("Failed to Exec, client is nil")
	}
}

func (m *MongoDB) CloseConnect() {
	if m.client != nil {
		err := m.client.Disconnect(context.TODO())
		if err != nil {
			Log.ErrorLog("Failed to CloseConnect, err = %v", err)
		}
	}
}

func (m *MongoDB) InsertOne(tabName string, newData interface{}, index int, fun FunUpdateCallBack, callbackData ...interface{}) {
	m.putToQueue(&MongoMessage{tabName: tabName, newData: newData, callbackData: callbackData, funUpdate: fun, operate: operateInsertOne}, index)
}

func (m *MongoDB) UpdateOne(tabName string, condition interface{}, newData interface{}, index int, fun FunUpdateCallBack, callbackData ...interface{}) {
	m.putToQueue(&MongoMessage{tabName: tabName, condition: condition, newData: newData, callbackData: callbackData, funUpdate: fun, operate: operateUpdateOne}, index)
}

func (m *MongoDB) DeleteOne(tabName string, condition interface{}, index int, fun FunUpdateCallBack, callbackData ...interface{}) {
	m.putToQueue(&MongoMessage{tabName: tabName, condition: condition, callbackData: callbackData, funUpdate: fun, operate: operateDeleteOne}, index)
}

func (m *MongoDB) FindOne(tabName string, condition interface{}, getData interface{}, index int, fun FunFindCallBack, callbackData ...interface{}) {
	m.putToQueue(&MongoMessage{tabName: tabName, condition: condition, newData: getData, callbackData: callbackData, funFind: fun, operate: operateFindOne}, index)
}

func (m *MongoDB) InsertMany(tabName string, newData []interface{}, index int, fun FunUpdateCallBack, callbackData ...interface{}) {
	m.putToQueue(&MongoMessage{tabName: tabName, newData: newData, callbackData: callbackData, funUpdate: fun, operate: operateInsertMany}, index)
}

func (m *MongoDB) UpdateMany(tabName string, condition interface{}, newData interface{}, index int, fun FunUpdateCallBack, callbackData ...interface{}) {
	m.putToQueue(&MongoMessage{tabName: tabName, condition: condition, newData: newData, callbackData: callbackData, funUpdate: fun, operate: operateUpdateMany}, index)
}

func (m *MongoDB) DeleteMany(tabName string, condition interface{}, index int, fun FunUpdateCallBack, callbackData ...interface{}) {
	m.putToQueue(&MongoMessage{tabName: tabName, condition: condition, callbackData: callbackData, funUpdate: fun, operate: operateDeleteMany}, index)
}

func (m *MongoDB) FindMany(tabName string, condition interface{}, getData interface{}, index int, fun FunFindCallBack, callbackData ...interface{}) {
	m.putToQueue(&MongoMessage{tabName: tabName, condition: condition, newData: getData, callbackData: callbackData, funFind: fun, operate: operateFindMany}, index)
}

func (m *MongoDB) putToQueue(message *MongoMessage, index int) {
	m.queue[index%m.queueCount].Enqueue(message)
}

func (m *MongoDB) insertOne(tabName string, newData interface{}, callBackData []interface{}, fun FunUpdateCallBack) {
	collection := m.database.Collection(tabName)

	_, err := collection.InsertOne(context.TODO(), newData)
	if err != nil {
		Log.ErrorLog("insertOne error, collection.InsertOne, err = %v", err)
		return
	}

	if fun != nil {
		fun(callBackData...)
	}
}

func (m *MongoDB) updateOne(tabName string, condition interface{}, newData interface{}, callbackData []interface{}, fun FunUpdateCallBack) {
	collection := m.database.Collection(tabName)

	filter, err := bson.Marshal(&condition)
	if err != nil {
		Log.ErrorLog("updateOne error, bson.Marshal, err = %v", err)
		return
	}

	maps := bson.M{}
	if err := bson.Unmarshal(filter, maps); err != nil {
		Log.ErrorLog("updateOne error, bson.Unmarshal, err = ", err)
		return
	}

	update := bson.M{
		"$set": newData,
	}

	if _, err := collection.UpdateOne(context.TODO(), maps, update); err != nil {
		Log.ErrorLog("updateOne error, collection.UpdateOne, err = %v", err)
		return
	}

	if fun != nil {
		fun(callbackData...)
	}
}

func (m *MongoDB) deleteOne(tabName string, condition interface{}, callbackData []interface{}, fun FunUpdateCallBack) {
	collection := m.database.Collection(tabName)

	filter, err := bson.Marshal(&condition)
	if err != nil {
		Log.ErrorLog("deleteOne error, bson.Marshal, err = %v", err)
		return
	}

	maps := bson.M{}
	if err := bson.Unmarshal(filter, maps); err != nil {
		Log.ErrorLog("deleteOne error, bson.Unmarshal, err = %v", err)
		return
	}

	if _, err := collection.DeleteOne(context.TODO(), maps); err != nil {
		Log.ErrorLog("deleteOne error, collection.DeleteOne, err = %v", err)
		return
	}

	if fun != nil {
		fun(callbackData...)
	}
}

func (m *MongoDB) findOneBack(tabName string, condition map[string]interface{}, getData interface{}, callbackData []interface{}, fun FunFindCallBack) {
	collection := m.database.Collection(tabName)

	var filter bson.D
	for key, value := range condition {
		filter = append(filter, primitive.E{Key: key, Value: value})
	}

	if err := collection.FindOne(context.TODO(), filter).Decode(getData); err != nil {
		Log.ErrorLog("findOneBack error, collection.find, err = %v", err)
		return
	}

	if fun != nil {
		fun(getData, callbackData...)
	}
}

func (m *MongoDB) findOne(tabName string, condition interface{}, getData interface{}, callbackData []interface{}, fun FunFindCallBack) {
	collection := m.database.Collection(tabName)

	filter, err := bson.Marshal(&condition)
	if err != nil {
		Log.ErrorLog("findOne error, bson.Marshal, err = %v", err)
		return
	}

	maps := bson.M{}
	if err := bson.Unmarshal(filter, maps); err != nil {
		Log.ErrorLog("findOne error, bson.Unmarshal, err = %v", err)
		return
	}

	if err := collection.FindOne(context.TODO(), maps).Decode(getData); err != nil {
		if !strings.Contains(err.Error(), "no documents in result") {
			Log.ErrorLog("findOne error, collection.find, err = %v", err)
			return
		}
	}

	if fun != nil {
		fun(getData, callbackData...)
	}
}

func (m *MongoDB) insertMany(tabName string, newData interface{}, callBackData []interface{}, fun FunUpdateCallBack) {
	collection := m.database.Collection(tabName)

	_, err := collection.InsertMany(context.TODO(), newData.([]interface{}))
	if err != nil {
		Log.ErrorLog("insertMany error, collection.InsertMany, err = %v", err)
		return
	}

	if fun != nil {
		fun(callBackData...)
	}
}

func (m *MongoDB) updateMany(tabName string, condition interface{}, newData interface{}, callbackData []interface{}, fun FunUpdateCallBack) {
	collection := m.database.Collection(tabName)

	filter, err := bson.Marshal(&condition)
	if err != nil {
		Log.ErrorLog("updateMany error, bson.Marshal, err = %v", err)
		return
	}

	maps := bson.M{}
	if err := bson.Unmarshal(filter, maps); err != nil {
		Log.ErrorLog("updateMany error, bson.Unmarshal, err = %v", err)
		return
	}

	update := bson.M{
		"$set": newData,
	}

	if _, err := collection.UpdateMany(context.TODO(), maps, update); err != nil {
		Log.ErrorLog("updateMany error, collection.updateOne, err = %v", err)
		return
	}

	if fun != nil {
		fun(callbackData...)
	}
}

func (m *MongoDB) deleteMany(tabName string, condition interface{}, callbackData []interface{}, fun FunUpdateCallBack) {
	collection := m.database.Collection(tabName)

	filter, err := bson.Marshal(&condition)
	if err != nil {
		Log.ErrorLog("deleteMany error, bson.Marshal, err = %v", err)
		return
	}

	maps := bson.M{}
	if err := bson.Unmarshal(filter, maps); err != nil {
		Log.ErrorLog("deleteMany error, bson.Unmarshal, err = %v", err)
		return
	}

	if _, err := collection.DeleteMany(context.TODO(), maps); err != nil {
		Log.ErrorLog("deleteMany error, collection.DeleteMany, err = %v", err)
		return
	}

	if fun != nil {
		fun(callbackData...)
	}
}

func (m *MongoDB) findMany(tabName string, condition interface{}, getData interface{}, callBackData []interface{}, fun FunFindCallBack) {
	collection := m.database.Collection(tabName)

	filter, err := bson.Marshal(&condition)
	if err != nil {
		Log.ErrorLog("findMany error, bson.Marshal, err = %v", err)
		return
	}

	maps := bson.M{}
	if err := bson.Unmarshal(filter, maps); err != nil {
		Log.ErrorLog("findMany error, bson.Unmarshal, err = %v", err)
		return
	}

	if cursor, err := collection.Find(context.TODO(), maps); err != nil {
		Log.ErrorLog("findMany error, collection.FindOne, err = %v", err)
	} else {
		//延迟关闭游标
		defer func() {
			if err := cursor.Close(context.TODO()); err != nil {
				Log.ErrorLog("findMany error, cursor.Close, err = %v", err)
			}
		}()

		if err := cursor.All(context.TODO(), &getData); err != nil {
			Log.ErrorLog("findMany error, cursor.All, err = %v", err)
			return
		}

		if fun != nil {
			fun(getData, callBackData...)
		}
	}
}
