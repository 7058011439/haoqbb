package test

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Probability"
	"github.com/7058011439/haoqbb/Timer"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/player"
)

const (
	randModuleId = 0
)

type testModule struct {
	id         int         // id
	fun        testFun     // 具体执行内容
	nextModule map[int]int // map[模块id]权重
}

type testFun func(player.IPlayer) bool

var mapTestModule = make(map[int]*testModule)
var listEntranceModule []int
var weight = Probability.NewWeights()
var mapTestFun = make(map[int]testFun)
var qps = 0

func InitFunc(id int, fun testFun) {
	mapTestFun[id] = fun
}

func InitOver(q int) {
	qps = q
	if qps < 1 {
		qps = 1
	}
	for id, module := range mapTestModule {
		if len(module.nextModule) == 0 {
			weight.AddWeight(id, 0, 1)
		}
		for v, w := range module.nextModule {
			if _, ok := mapTestModule[v]; !ok && v != randModuleId {
				Log.WarningLog("Failed to init the test module, next module not exist, id = %v", id)
			} else {
				weight.AddWeight(id, v, w)
			}
		}
	}
}

func InsertTestModule(id int, entrance int, nextModule map[int]int) {
	if module := mapTestModule[id]; module == nil {
		if fun, ok := mapTestFun[id]; ok {
			mapTestModule[id] = &testModule{
				id:         id,
				fun:        fun,
				nextModule: nextModule,
			}
			if entrance > 0 {
				weight.AddWeight(randModuleId, id, entrance)
			}
		} else {
			Log.ErrorLog("Failed to InsertTestModule, no exec func", id)
		}
	} else {
		Log.ErrorLog("Failed to InsertTestModule, id repeated = %v", id)
	}
}

func GetRandomModule() int {
	return weight.Value(randModuleId)
}

func Run(_ Timer.TimerID, args ...interface{}) {
	clientId := args[0].(uint64)
	player := player.GetPlayerByClientId(clientId)
	if player != nil && player.IsLogin() {
		if module := mapTestModule[player.TestModule()]; module != nil {
			for i := 0; i < qps; i++ {
				module.fun(player)
			}
			nextModuleId := weight.Value(module.id)
			if nextModuleId == 0 {
				nextModuleId = weight.Value(randModuleId)
			}
			player.SetTestModule(nextModuleId)
		} else {
			Log.ErrorLog("Failed to Run, test module is nil, id = %v", player.TestModule())
		}
	}
}
