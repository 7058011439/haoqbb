package db

import (
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/Timer"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/service"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/mongo"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/timer"
)

const (
	interval = 7
)

//某玩家的农场或者家园有哪些访问者
type visitRecord struct {
	hostVisitors  map[int64]*Stl.Set // 被访问记录  我的家园有哪些浏览者
	guestVisitIds map[int64]*Stl.Set // 浏览痕迹 我浏览了哪些人的家园
}

func newVisitRecord() visitRecord {
	return visitRecord{
		hostVisitors:  make(map[int64]*Stl.Set),
		guestVisitIds: make(map[int64]*Stl.Set),
	}
}

func (a *visitRecord) add(hostId, guestId int64) {
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
func (a *visitRecord) clearHistoryVisit(guestId int64) []int64 {
	if a.guestVisitIds[guestId] == nil {
		return nil
	}
	var ret []int64
	a.guestVisitIds[guestId].Range(func(i interface{}) {
		hostId := i.(int64)
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

// ShareDataMgrSync 某个玩家的农场或者庄园的数据 多人访问 管理类
type ShareDataMgrSync struct {
	dbIndex      int64
	operate      IDBTool
	record       visitRecord
	collectName  string
	shareDataMap map[int64]IDBData
	updateData   map[int64]bool // 有哪些玩家数据发生变化
}

func NewShareDataMgrSync(collectName string, operate IDBTool) *ShareDataMgrSync {
	ret := &ShareDataMgrSync{
		operate:      operate,
		record:       newVisitRecord(),
		collectName:  collectName,
		shareDataMap: make(map[int64]IDBData),
		updateData:   make(map[int64]bool),
	}
	ITimer.SetRepeatTimer(service.GetServiceName(), 1000, ret.updateRun)
	return ret
}

// LogOut 玩家退出时，清理本模块的数据
func (s *ShareDataMgrSync) LogOut(guestId int64) {
	emptyPlaceIds := s.record.clearHistoryVisit(guestId)
	for _, id := range emptyPlaceIds {
		if _, ok := s.updateData[id]; ok {
			IMongo.UpdateOne(service.GetServiceName(), s.collectName, s.operate.NewCondition(id), s.shareDataMap[id].Data(), int(id), func(callbackData ...interface{}) {
				delete(s.shareDataMap, id)
			})
			delete(s.updateData, id)
		} else {
			delete(s.shareDataMap, id)
		}
	}
}

func (s *ShareDataMgrSync) GetData(guestId, hostId int64) IDBData {
	s.record.add(hostId, guestId)
	if data, ok := s.shareDataMap[hostId]; ok {
		return data
	} else {
		data = s.operate.NewObj()
		IMongo.FindOneSync(service.GetServiceName(), s.collectName, s.operate.NewCondition(hostId), data)
		s.shareDataMap[hostId] = data
		return data
	}
}

// UpdateData 有数据变更
func (s *ShareDataMgrSync) UpdateData(hostId int64) {
	s.updateData[hostId] = true
}

func (s *ShareDataMgrSync) InsertData(data IDBData) {
	userId := data.GetUserId()
	IMongo.InsertOne(service.GetServiceName(), s.collectName, data.Data(), int(userId))
	s.shareDataMap[userId] = data
}

func (s *ShareDataMgrSync) updateRun(Timer.TimerID, ...interface{}) {
	for hostId := range s.updateData {
		if hostId%interval == s.dbIndex {
			obj := s.shareDataMap[hostId]
			if obj != nil {
				IMongo.UpdateOne(service.GetServiceName(), s.collectName, s.operate.NewCondition(hostId), s.shareDataMap[hostId].Data(), int(hostId), nil)
			}
			delete(s.updateData, hostId)
		}
	}
	s.dbIndex++
	if s.dbIndex >= interval {
		s.dbIndex = 0
	}
}
