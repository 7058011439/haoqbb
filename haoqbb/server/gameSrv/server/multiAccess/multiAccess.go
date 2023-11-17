package multiAccess

import (
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/Timer"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/service"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/mongo"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/timer"
)

//某玩家的农场或者家园有哪些访问者
type visitRecord struct {
	hostVisitors  map[int]*Stl.Set // 被访问记录  我的家园有哪些浏览者
	guestVisitIds map[int]*Stl.Set // 浏览痕迹 我浏览了哪些人的家园
}

func newVisitRecord() *visitRecord {
	return &visitRecord{
		hostVisitors:  make(map[int]*Stl.Set),
		guestVisitIds: make(map[int]*Stl.Set),
	}
}

func (a *visitRecord) add(hostId, guestId int) {
	if a.hostVisitors[hostId] == nil {
		a.hostVisitors[hostId] = Stl.NewSet()
	}
	if a.guestVisitIds[guestId] == nil {
		a.guestVisitIds[guestId] = Stl.NewSet()
	}

	a.hostVisitors[hostId].Add(guestId)
	a.guestVisitIds[guestId].Add(hostId)
}

//返回值：所有 家园没人 的 主人ID集合
func (a *visitRecord) clearHistoryVisit(guestId int) []int {
	if a.guestVisitIds[guestId] == nil {
		return nil
	}
	ret := make([]int, 16)
	a.guestVisitIds[guestId].Range(func(i interface{}) {
		hostId := i.(int)
		if a.hostVisitors[hostId] != nil {
			a.hostVisitors[hostId].Del(guestId)
			if a.hostVisitors[hostId].Empty() {
				ret = append(ret, hostId)
			}
		}
	})

	delete(a.guestVisitIds, guestId)
	return ret
}

type IDBData interface {
	DataOK() bool
	Condition() map[string]interface{}
	OnLoadEnd()
}

var mgrCount Timer.TimeWheel

// ShareDataMgr 某个玩家的农场或者庄园的数据 多人访问 管理类
type ShareDataMgr struct {
	newObj       FunNewObj
	record       *visitRecord
	shareDataMap map[int]interface{}
	collectName  string
	updateData   map[int]map[string]interface{} // 有哪些玩家数据发生变化
	waitCall     map[int][]*dbCallBackMsg       //多人访问某个玩家的农场，用一个队列来缓存所有的访问队列
}

func NewShareDataMgr(fun FunNewObj, collectName string) *ShareDataMgr {
	ret := &ShareDataMgr{
		newObj:       fun,
		record:       newVisitRecord(),
		collectName:  collectName,
		shareDataMap: make(map[int]interface{}),
		updateData:   make(map[int]map[string]interface{}),
		waitCall:     make(map[int][]*dbCallBackMsg),
	}
	mgrCount++
	ITimer.SetOnceTimer(service.GetServiceName(), mgrCount*1000, func(timerId Timer.TimerID, args ...interface{}) {
		ITimer.SetRepeatTimer(service.GetServiceName(), 7000, ret.updateRun)
	})
	return ret
}

// IsExist 某人的家园数据是否加载好了
func (s *ShareDataMgr) IsExist(hostId int) bool {
	_, ok := s.shareDataMap[hostId]
	return ok
}

// GetData A访问B的农场或者家园
func (s *ShareDataMgr) GetData(hostId, guestId int) interface{} {
	if s.shareDataMap[hostId] == nil {
		if len(s.waitCall[hostId]) == 0 {
			s.readDB(hostId)
		}
		return nil
	}

	s.record.add(hostId, guestId)
	return s.shareDataMap[hostId]
}

// LogOut 玩家退出时，清理本模块的数据
func (s *ShareDataMgr) LogOut(guestId int) {
	emptyPlaceIds := s.record.clearHistoryVisit(guestId)
	for _, id := range emptyPlaceIds {
		delete(s.shareDataMap, id)
	}
}

func (s *ShareDataMgr) GetDataAndDo(guestId, hostId int, fun func(data interface{}, args ...interface{}), args ...interface{}) {
	s.ExecWhenDataReady(func(guestId int, hostId int, args ...interface{}) {
		data := s.GetData(hostId, guestId)
		fun(data, args...)
	}, guestId, hostId, args...)
}

// ExecWhenDataReady 数据存在时，调用方法。可能延迟执行，因为还没加载完数据
func (s *ShareDataMgr) ExecWhenDataReady(callBack DBCallBackFun, guestId int, hostId int, args ...interface{}) {
	//如果数据已经在了，直接执行方法并返回
	if s.GetData(hostId, guestId) != nil {
		callBack(guestId, hostId, args...)
		return
	}

	//如果没有数据，就加入到缓存队列里面
	s.addPendingCall(callBack, guestId, hostId, args...)
}

// UpdateData 有数据变更
func (s *ShareDataMgr) UpdateData(hostId int, key string, value interface{}) {
	if s.updateData[hostId] == nil {
		s.updateData[hostId] = make(map[string]interface{})
	}
	if value != nil {
		s.updateData[hostId][key] = value
	}
}

//添加A的家园数据到本管理类
func (s *ShareDataMgr) addData(hostId int, data interface{}) {
	s.shareDataMap[hostId] = data
}

func (s *ShareDataMgr) updateRun(Timer.TimerID, ...interface{}) {
	for hostId, value := range s.updateData {
		data := s.GetData(hostId, hostId).(IDBData)
		if len(value) == 0 {
			IMongo.UpdateOne(service.GetServiceName(), s.collectName, data.Condition(), data, hostId, nil)
		} else {
			IMongo.UpdateOne(service.GetServiceName(), s.collectName, value, nil, hostId, nil)
		}
	}
	s.updateData = make(map[int]map[string]interface{})
}

// 发起数据库读取
func (s *ShareDataMgr) readDB(targetId int) {
	condition := s.newObj(targetId).Condition()
	data := s.newObj(0)
	IMongo.FindOne(service.GetServiceName(), s.collectName, condition, data, targetId, func(getData interface{}, callbackData ...interface{}) {
		data = getData.(IDBData)
		if !data.DataOK() {
			data = s.newObj(targetId)
			IMongo.InsertOne(service.GetServiceName(), s.collectName, data, targetId)
			//Log.Debug("new data, collectName = %v, data = %v", s.collectName, data)
		}
		s.addData(targetId, data)
		data.OnLoadEnd()
		s.doPendingCall(targetId)
	}, targetId, s)
}

//添加访问缓存
func (s *ShareDataMgr) addPendingCall(callBack DBCallBackFun, guestId int, hostId int, args ...interface{}) {
	s.waitCall[hostId] = append(s.waitCall[hostId], &dbCallBackMsg{hostId: hostId, guestId: guestId, params: args, callFun: callBack})
}

//执行访问缓存
func (s *ShareDataMgr) doPendingCall(userId int) {
	for _, info := range s.waitCall[userId] {
		if info.callFun != nil {
			info.callFun(info.guestId, info.hostId, info.params...)
		}
	}
	delete(s.waitCall, userId)
}

type DBCallBackFun func(guestId int, hostId int, args ...interface{})
type dbCallBackMsg struct {
	hostId  int           // 访问者id
	guestId int           // 被访问者id
	params  []interface{} // 额外参数
	callFun DBCallBackFun // 回调函数
}

type FunNewObj func(id int) IDBData
