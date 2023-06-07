package service

import (
	"fmt"
	"github.com/7058011439/haoqbb/DataBase"
	"github.com/7058011439/haoqbb/Http"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/Timer"
	"time"
)

const (
	msgTypeTimer           = "timer"
	msgTypeServiceMsg      = "serviceMsg"
	msgTypeTcpMsg          = "tcpMsg"
	msgTypeMongo           = "mongo"
	msgTypeHttp            = "http"
	msgTypeDiscoverService = "discoverService"
	msgTypeLoseService     = "loseService"

	warningTime = 16.0 // 单次执行报警时间
)

type timerData struct {
	timeId   Timer.TimerID
	backData []interface{}
	callFun  Timer.TimerFun
}

type serviceMsgData struct {
	srcServiceId int
	msgType      int
	data         []byte
}

type mongoData struct {
	backData      []interface{}
	data          interface{}
	callFunGet    DataBase.FunFindCallBack
	callFunUpdate DataBase.FunUpdateCallBack
}

type tcpMsgData struct {
	clientId uint64
	data     []byte
}

type httpData struct {
	backData []interface{}
	data     map[string]interface{}
	callFun  func(map[string]interface{}, ...interface{})
}

type serviceStatus struct {
	serviceName string // 服务名
	serviceId   int    // 服务Id
}

type queue struct {
	chanTimer             chan *timerData
	chanServiceMsg        chan *serviceMsgData
	chanMongo             chan *mongoData
	chanHttp              chan *httpData
	chanTcpMsg            chan *tcpMsgData
	chanDiscoverService   chan *serviceStatus
	chanLoseService       chan *serviceStatus
	MongoDB               *DataBase.MongoDB
	RedisDB               *DataBase.RedisDB
	tcpMsgHandler         func(uint64, []byte)
	name                  string
	lastTime              time.Time
	performance           map[string]perform
	serviceMsgHandle      map[int]func(srcServiceId int, data []byte)
	discoverServiceHandle map[string]func(int)
	loseServiceHandle     map[string]func(int)
}

func NewQueue(name string) *queue {
	return &queue{
		chanTimer:           make(chan *timerData, 65535*64),      // 定时器消息队列
		chanServiceMsg:      make(chan *serviceMsgData, 65535*64), // 服务间消息队列
		chanMongo:           make(chan *mongoData, 65535*64),      // mongo回调消息队列
		chanHttp:            make(chan *httpData, 65535*64),       // http回调消息队列
		chanTcpMsg:          make(chan *tcpMsgData, 65535*64),     // tcp消息队列
		chanDiscoverService: make(chan *serviceStatus, 64),        // 发现服务消息队列
		chanLoseService:     make(chan *serviceStatus, 64),        // 丢失服务消息队列

		name:                  name,
		lastTime:              time.Now(),
		performance:           map[string]perform{},
		serviceMsgHandle:      map[int]func(srcServiceId int, data []byte){},
		discoverServiceHandle: map[string]func(int){},
		loseServiceHandle:     map[string]func(int){},
	}
}

type perform struct {
	msgCount int   // 处理消息数量
	costTime int64 // 耗时
}

func (t *perform) update(costTime int64) {
	t.costTime += costTime
	t.msgCount++
}

func (q *queue) run() {
	cost := Timer.NewTiming(Timer.Millisecond)
	msgType := ""
	select {
	case msg := <-q.chanTcpMsg:
		cost.ReStart()
		q.tcpMsgHandler(msg.clientId, msg.data)
		msgType = msgTypeTcpMsg
	case data := <-q.chanServiceMsg:
		cost.ReStart()
		q.serviceMsgHandle[data.msgType](data.srcServiceId, data.data)
		msgType = msgTypeServiceMsg
	case db := <-q.chanMongo:
		cost.ReStart()
		if db.callFunGet != nil {
			db.callFunGet(db.data, db.backData...)
		}
		if db.callFunUpdate != nil {
			db.callFunUpdate(db.backData)
		}
		msgType = msgTypeMongo
	case http := <-q.chanHttp:
		cost.ReStart()
		http.callFun(http.data, http.backData...)
		msgType = msgTypeHttp
	case timer := <-q.chanTimer:
		cost.ReStart()
		timer.callFun(timer.timeId, timer.backData...)
		msgType = msgTypeTimer
	case status := <-q.chanDiscoverService:
		cost.ReStart()
		q.discoverServiceHandle[status.serviceName](status.serviceId)
		msgType = msgTypeDiscoverService
	case status := <-q.chanLoseService:
		cost.ReStart()
		q.loseServiceHandle[status.serviceName](status.serviceId)
		msgType = msgTypeLoseService
	}
	//cost.PrintCost(warningTime, false, "%v %s callFun timeout", q.name, msgType)
	per := q.performance[msgType]
	per.update(int64(cost.GetCost()))
	q.performance[msgType] = per

	if gaps := time.Now().Sub(q.lastTime).Seconds(); gaps > 1 {
		//Log.Debug("queue run cost 1 second, name = %v, gaps = %vs, info = %v, NumGoroutine = %v", q.name, gaps, q.performance, runtime.NumGoroutine())
		q.lastTime = time.Now()
		q.performance = map[string]perform{}
	}
}

// NewServiceMsg 收到其他服务消息
func (q *queue) NewServiceMsg(srcServiceId int, msgType int, data []byte) {
	if fun := q.serviceMsgHandle[msgType]; fun == nil {
		return
	}
	if len(q.chanServiceMsg) == cap(q.chanServiceMsg) {
		Log.ErrorLog("%v Failed to insert serviceMsg to chan, chan full", q.name)
		return
	}
	q.chanServiceMsg <- &serviceMsgData{
		srcServiceId: srcServiceId,
		msgType:      msgType,
		data:         data,
	}
}

// DiscoverService 发现其他(节点)服务
func (q *queue) DiscoverService(serviceName string, serviceId int) {
	if handle := q.discoverServiceHandle[serviceName]; handle == nil {
		return
	}
	if len(q.chanDiscoverService) == cap(q.chanDiscoverService) {
		Log.ErrorLog("%v Failed to insert discoverService to chan, chan full", q.name)
		return
	}
	q.chanDiscoverService <- &serviceStatus{
		serviceName: serviceName,
		serviceId:   serviceId,
	}
}

// LoseService 遗失其他(节点)服务
func (q *queue) LoseService(serviceName string, serviceId int) {
	if handle := q.loseServiceHandle[serviceName]; handle == nil {
		return
	}
	if len(q.chanLoseService) == cap(q.chanLoseService) {
		Log.ErrorLog("%v Failed to insert loseService to chan, chan full", q.name)
		return
	}
	q.chanLoseService <- &serviceStatus{
		serviceName: serviceName,
		serviceId:   serviceId,
	}
}

// RegeditHandleTcpMsg 注册tcp消息处理
func (q *queue) RegeditHandleTcpMsg(fun func(clientId uint64, data []byte)) {
	q.tcpMsgHandler = fun
}

// NewTcpMsg 收到tcp包处理(目前就网关和客户端)
func (q *queue) NewTcpMsg(client Net.IClient, data []byte) {
	if len(q.chanTcpMsg) == cap(q.chanTcpMsg) {
		Log.ErrorLog("%v Failed to insert tcpMsg to chan, chan full", q.name)
		return
	}
	q.chanTcpMsg <- &tcpMsgData{
		clientId: client.GetId(),
		data:     data,
	}
}

// SetRepeatTimer 设置循环定时器
func (q *queue) SetRepeatTimer(duration Timer.TimeWheel, funcName Timer.TimerFun, args ...interface{}) Timer.TimerID {
	return Timer.AddRepeatTimer(duration, func(timerId Timer.TimerID, args ...interface{}) {
		if len(q.chanTimer) == cap(q.chanTimer) {
			Log.ErrorLog("%v Failed to insert timer to chan, chan full", q.name)
			return
		}
		data := &timerData{
			timeId:   timerId,
			backData: args,
			callFun:  funcName,
		}
		q.chanTimer <- data
	}, args...)
}

// SetOnceTimer 设置一次性定时器
func (q *queue) SetOnceTimer(duration Timer.TimeWheel, funcName Timer.TimerFun, args ...interface{}) Timer.TimerID {
	return Timer.AddOnceTimer(duration, func(timerId Timer.TimerID, args ...interface{}) {
		if len(q.chanTimer) == cap(q.chanTimer) {
			Log.ErrorLog("%v Failed to insert timer to chan, chan full", q.name)
			return
		}
		data := &timerData{
			timeId:   timerId,
			backData: args,
			callFun:  funcName,
		}
		q.chanTimer <- data
	}, args...)
}

// SetOffTimer 关闭定时器
func (q *queue) SetOffTimer(id Timer.TimerID) {
	Timer.CloseTimer(id)
}

func (q *queue) GetTimerCount() int {
	return Timer.Count()
}

func (q *queue) GetMongoAsync(tabName string, condition interface{}, getData interface{}, index int, fun DataBase.FunFindCallBack, callbackData ...interface{}) {
	cost := Timer.NewTiming(Timer.Millisecond)
	q.MongoDB.FindOne(tabName, condition, getData, index, func(getData interface{}, callbackData ...interface{}) {
		if len(q.chanMongo) == cap(q.chanMongo) {
			Log.ErrorLog("%v Failed to insert mongo to chan, chan full", q.name)
			return
		}
		q.chanMongo <- &mongoData{
			data:       getData,
			backData:   callbackData,
			callFunGet: fun,
		}
		cost.PrintCost(1000, false, "Get mongo data time out, tabName = %v, condition = %v, index = %v", tabName, condition, index)
	}, callbackData...)
}

func (q *queue) InsertMongoAsync(tabName string, data interface{}, index int, fun DataBase.FunUpdateCallBack, callBackData ...interface{}) {
	if fun != nil {
		q.MongoDB.InsertOne(tabName, data, index, func(data ...interface{}) {
			if len(q.chanMongo) == cap(q.chanMongo) {
				Log.ErrorLog("%v Failed to insert mongo to chan, chan full", q.name)
				return
			}
			q.chanMongo <- &mongoData{
				backData:      data,
				callFunUpdate: fun,
			}
		}, callBackData...)
	} else {
		q.MongoDB.InsertOne(tabName, data, index, fun, callBackData...)
	}
}

func (q *queue) UpdateMongoAsync(tabName string, condition interface{}, data interface{}, index int, fun DataBase.FunUpdateCallBack, callBackData ...interface{}) {
	if fun != nil {
		q.MongoDB.UpdateOne(tabName, condition, data, index, func(data ...interface{}) {
			if len(q.chanMongo) == cap(q.chanMongo) {
				Log.ErrorLog("%v Failed to insert mongo to chan, chan full", q.name)
				return
			}
			q.chanMongo <- &mongoData{
				backData:      data,
				callFunUpdate: fun,
			}
		}, callBackData...)
	} else {
		q.MongoDB.UpdateOne(tabName, condition, data, index, fun, callBackData...)
	}
}

func (q *queue) GetHttpAsync(url string, header map[string]string, fun func(getData map[string]interface{}, backData ...interface{}), backData ...interface{}) {
	Http.GetHttpAsync(url, Http.NewHead(header), func(getData map[string]interface{}, _ error, backData ...interface{}) {
		if len(q.chanHttp) == cap(q.chanHttp) {
			Log.ErrorLog("%v Failed to insert http to chan, chan full", q.name)
			return
		}
		q.chanHttp <- &httpData{
			data:     getData,
			backData: backData,
			callFun:  fun,
		}
	}, backData...)
}

func (q *queue) PostHttpAsync(url string, header map[string]string, body map[string]interface{}, callback func(map[string]interface{}, ...interface{}), backData ...interface{}) {
	Http.PostHttpAsync(url, Http.NewHead(header), Http.NewBody(body), func(getData map[string]interface{}, _ error, backData ...interface{}) {
		if len(q.chanHttp) == cap(q.chanHttp) {
			Log.ErrorLog("%v Failed to insert http to chan, chan full", q.name)
			return
		}
		q.chanHttp <- &httpData{
			data:     getData,
			backData: backData,
			callFun:  callback,
		}
	}, backData...)
}

func (q *queue) GetRedisSync(key string, field string) string {
	return q.RedisDB.HGetString(key, field)
}

func (q *queue) SetRedisSync(key string, field string, value interface{}) {
	q.RedisDB.HSetValue(key, field, value)
}

func (q *queue) IncRedisSyn(key string, field string, number int64) int64 {
	return q.RedisDB.HIncBy(key, field, number)
}

func (q *queue) GetQueueLen() string {
	return fmt.Sprintf("%v_%v", len(q.chanServiceMsg), len(q.chanTcpMsg))
}
