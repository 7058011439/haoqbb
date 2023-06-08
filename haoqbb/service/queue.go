package service

import (
	"github.com/7058011439/haoqbb/DataBase"
	"github.com/7058011439/haoqbb/Http"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/Timer"
	"time"
)

const (
	warningTime = 16.0 // 单次执行报警时间
)

type msgType = int

const (
	typeTimer msgType = iota // 定时器
	typeServiceMsg
	typeTcpMsg
	typeMongoMsg
	typeHttpMsg
	typeDiscoverService
	typeLoseService
)

type queueData struct {
	eType msgType
	data  interface{}
}

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

type serviceData struct {
	serviceName string // 服务名
	serviceId   int    // 服务Id
}

type queue struct {
	name                  string
	chanAll               chan *queueData
	MongoDB               *DataBase.MongoDB
	RedisDB               *DataBase.RedisDB
	tcpMsgHandler         func(uint64, []byte)
	serviceMsgHandle      map[int]func(srcServiceId int, data []byte)
	discoverServiceHandle map[string]func(int)
	loseServiceHandle     map[string]func(int)
}

func NewQueue(name string) *queue {
	return &queue{
		name:                  name,
		chanAll:               make(chan *queueData, 65535*64),
		serviceMsgHandle:      map[int]func(srcServiceId int, data []byte){},
		discoverServiceHandle: map[string]func(int){},
		loseServiceHandle:     map[string]func(int){},
	}
}

type perform struct {
	msgCount int     // 处理消息数量
	costTime float64 // 耗时
}

func (t *perform) update(costTime float64) {
	t.costTime += costTime
	t.msgCount++
}

func (q *queue) run() {
	cost := Timer.NewTiming(Timer.Millisecond)
	lastTime := time.Now()
	performance := map[int]*perform{}
	select {
	case msg := <-q.chanAll:
		cost.ReStart()
		switch msg.eType {
		case typeTcpMsg:
			data := msg.data.(*tcpMsgData)
			if q.tcpMsgHandler != nil {
				q.tcpMsgHandler(data.clientId, data.data)
			}
		case typeServiceMsg:
			data := msg.data.(*serviceMsgData)
			if q.serviceMsgHandle[data.msgType] != nil {
				q.serviceMsgHandle[data.msgType](data.srcServiceId, data.data)
			}
		case typeMongoMsg:
			data := msg.data.(*mongoData)
			if data.callFunGet != nil {
				data.callFunGet(data.data, data.backData...)
			}
			if data.callFunUpdate != nil {
				data.callFunUpdate(data.backData...)
			}
		case typeHttpMsg:
			data := msg.data.(*httpData)
			if data.callFun != nil {
				data.callFun(data.data, data.backData...)
			}
		case typeTimer:
			data := msg.data.(*timerData)
			if data.callFun != nil {
				data.callFun(data.timeId, data.backData...)
			}
		case typeDiscoverService:
			data := msg.data.(*serviceData)
			if q.discoverServiceHandle[data.serviceName] != nil {
				q.discoverServiceHandle[data.serviceName](data.serviceId)
			}
		case typeLoseService:
			data := msg.data.(*serviceData)
			if q.loseServiceHandle[data.serviceName] != nil {
				q.loseServiceHandle[data.serviceName](data.serviceId)
			}
		}
		cost.PrintCost(warningTime, false, "%v(%v) callFun timeout", q.name, msg.eType)
		if performance[msg.eType] != nil {
			performance[msg.eType].update(cost.GetCost())
		} else {
			performance[msg.eType] = &perform{
				msgCount: 1,
				costTime: cost.GetCost(),
			}
		}
		if gaps := time.Now().Sub(lastTime).Seconds(); gaps > 1 {
			//Log.Debug("queue run cost 1 second, name = %v, gaps = %vs, info = %v, NumGoroutine = %v", q.name, gaps, q.performance, runtime.NumGoroutine())
			lastTime = time.Now()
			performance = map[int]*perform{}
		}
	}
}

// NewServiceMsg 收到其他服务消息
func (q *queue) NewServiceMsg(srcServiceId int, msgType int, data []byte) {
	if fun := q.serviceMsgHandle[msgType]; fun == nil {
		return
	}
	if len(q.chanAll) == cap(q.chanAll) {
		Log.ErrorLog("%v Failed to insert chan, chan full", q.name)
		return
	}
	q.chanAll <- &queueData{
		eType: typeServiceMsg,
		data: &serviceMsgData{
			srcServiceId: srcServiceId,
			msgType:      msgType,
			data:         data,
		},
	}
}

// DiscoverService 发现其他(节点)服务
func (q *queue) DiscoverService(serviceName string, serviceId int) {
	if handle := q.discoverServiceHandle[serviceName]; handle == nil {
		return
	}
	if len(q.chanAll) == cap(q.chanAll) {
		Log.ErrorLog("%v Failed to insert chan, chan full", q.name)
		return
	}
	q.chanAll <- &queueData{
		eType: typeDiscoverService,
		data: &serviceData{
			serviceName: serviceName,
			serviceId:   serviceId,
		},
	}
}

// LoseService 遗失其他(节点)服务
func (q *queue) LoseService(serviceName string, serviceId int) {
	if handle := q.loseServiceHandle[serviceName]; handle == nil {
		return
	}
	if len(q.chanAll) == cap(q.chanAll) {
		Log.ErrorLog("%v Failed to insert chan, chan full", q.name)
		return
	}
	q.chanAll <- &queueData{
		eType: typeLoseService,
		data: &serviceData{
			serviceName: serviceName,
			serviceId:   serviceId,
		},
	}
}

// RegeditHandleTcpMsg 注册tcp消息处理
func (q *queue) RegeditHandleTcpMsg(fun func(clientId uint64, data []byte)) {
	q.tcpMsgHandler = fun
}

// NewTcpMsg 收到tcp包处理(目前就网关和客户端)
func (q *queue) NewTcpMsg(client Net.IClient, data []byte) {
	if len(q.chanAll) == cap(q.chanAll) {
		Log.ErrorLog("%v Failed to insert chan, chan full", q.name)
		return
	}
	q.chanAll <- &queueData{
		eType: typeHttpMsg,
		data: &tcpMsgData{
			clientId: client.GetId(),
			data:     data,
		},
	}
}

// SetRepeatTimer 设置循环定时器
func (q *queue) SetRepeatTimer(duration Timer.TimeWheel, funcName Timer.TimerFun, args ...interface{}) Timer.TimerID {
	return Timer.AddRepeatTimer(duration, func(timerId Timer.TimerID, args ...interface{}) {
		if len(q.chanAll) == cap(q.chanAll) {
			Log.ErrorLog("%v Failed to insert chan, chan full", q.name)
			return
		}
		q.chanAll <- &queueData{
			eType: typeTimer,
			data: &timerData{
				timeId:   timerId,
				backData: args,
				callFun:  funcName,
			},
		}
	}, args...)
}

// SetOnceTimer 设置一次性定时器
func (q *queue) SetOnceTimer(duration Timer.TimeWheel, funcName Timer.TimerFun, args ...interface{}) Timer.TimerID {
	return Timer.AddOnceTimer(duration, func(timerId Timer.TimerID, args ...interface{}) {
		if len(q.chanAll) == cap(q.chanAll) {
			Log.ErrorLog("%v Failed to insert chan, chan full", q.name)
			return
		}
		q.chanAll <- &queueData{
			eType: typeTimer,
			data: &timerData{
				timeId:   timerId,
				backData: args,
				callFun:  funcName,
			},
		}
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
		if len(q.chanAll) == cap(q.chanAll) {
			Log.ErrorLog("%v Failed to insert chan, chan full", q.name)
			return
		}
		q.chanAll <- &queueData{
			eType: typeMongoMsg,
			data: &mongoData{
				data:       getData,
				backData:   callbackData,
				callFunGet: fun,
			},
		}
		cost.PrintCost(1000, false, "Get mongo data time out, tabName = %v, condition = %v, index = %v", tabName, condition, index)
	}, callbackData...)
}

func (q *queue) InsertMongoAsync(tabName string, data interface{}, index int, fun DataBase.FunUpdateCallBack, callBackData ...interface{}) {
	if fun != nil {
		q.MongoDB.InsertOne(tabName, data, index, func(data ...interface{}) {
			if len(q.chanAll) == cap(q.chanAll) {
				Log.ErrorLog("%v Failed to insert chan, chan full", q.name)
				return
			}
			q.chanAll <- &queueData{
				eType: typeMongoMsg,
				data: &mongoData{
					backData:      data,
					callFunUpdate: fun,
				},
			}
		}, callBackData...)
	} else {
		q.MongoDB.InsertOne(tabName, data, index, fun, callBackData...)
	}
}

func (q *queue) UpdateMongoAsync(tabName string, condition interface{}, data interface{}, index int, callBack DataBase.FunUpdateCallBack, callBackData ...interface{}) {
	if callBack != nil {
		q.MongoDB.UpdateOne(tabName, condition, data, index, func(data ...interface{}) {
			if len(q.chanAll) == cap(q.chanAll) {
				Log.ErrorLog("%v Failed to insert chan, chan full", q.name)
				return
			}
			q.chanAll <- &queueData{
				eType: typeMongoMsg,
				data: &mongoData{
					backData:      data,
					callFunUpdate: callBack,
				},
			}
		}, callBackData...)
	} else {
		q.MongoDB.UpdateOne(tabName, condition, data, index, callBack, callBackData...)
	}
}

func (q *queue) GetHttpAsync(url string, header map[string]string, callback func(getData map[string]interface{}, backData ...interface{}), backData ...interface{}) {
	Http.GetHttpAsync(url, Http.NewHead(header), func(getData map[string]interface{}, _ error, backData ...interface{}) {
		if len(q.chanAll) == cap(q.chanAll) {
			Log.ErrorLog("%v Failed to insert chan, chan full", q.name)
			return
		}
		q.chanAll <- &queueData{
			eType: typeHttpMsg,
			data: &httpData{
				data:     getData,
				backData: backData,
				callFun:  callback,
			},
		}
	}, backData...)
}

func (q *queue) PostHttpAsync(url string, header map[string]string, body map[string]interface{}, callback func(map[string]interface{}, ...interface{}), backData ...interface{}) {
	Http.PostHttpAsync(url, Http.NewHead(header), Http.NewBody(body), func(getData map[string]interface{}, _ error, backData ...interface{}) {
		if len(q.chanAll) == cap(q.chanAll) {
			Log.ErrorLog("%v Failed to insert chan, chan full", q.name)
			return
		}
		q.chanAll <- &queueData{
			eType: typeHttpMsg,
			data: &httpData{
				data:     getData,
				backData: backData,
				callFun:  callback,
			},
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
