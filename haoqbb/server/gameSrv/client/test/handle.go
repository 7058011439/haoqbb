package test

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Probability"
	"github.com/7058011439/haoqbb/Timer"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/player"
	"math/rand"
)

type testModule struct {
	id         int         // id
	fun        testFun     // 具体执行内容
	randRate   int         // 完全随机概率 N%
	nextModule map[int]int // map[模块id]权重
}

type testFun func(player.IPlayer) bool

var mapTestModule = make(map[int]*testModule)
var listEntranceModule []int
var weight = Probability.NewWeights()

func OnInitOver() {
	for id, module := range mapTestModule {
		for v, w := range module.nextModule {
			weight.AddWeight(id, v, w)
			if _, ok := mapTestModule[v]; !ok {
				Log.ErrorLog("Failed to init, next module not exist, id = %v", id)
				return
			}
		}
	}
}

func InsertTestModule(id int, entrance bool, fun testFun, randRate int, nextModule map[int]int) {
	if module := mapTestModule[id]; module == nil {
		mapTestModule[id] = &testModule{
			id:         id,
			fun:        fun,
			randRate:   randRate,
			nextModule: nextModule,
		}
		if entrance {
			listEntranceModule = append(listEntranceModule, id)
		}
	} else {
		Log.ErrorLog("Failed to InsertTestModule, id repeated = %v", id)
	}
}

func Run(_ Timer.TimerID, args ...interface{}) {
	clientId := args[0].(uint64)
	player := player.GetPlayerByClientId(clientId)
	if player != nil && player.IsLogin() {
		if module := mapTestModule[player.TestModule()]; module != nil {
			if module.fun(player) {
				nextModuleId := 0
				if rand.Intn(100) < module.randRate {
					if len(listEntranceModule) > 0 {
						nextModuleId = listEntranceModule[rand.Intn(len(listEntranceModule))]
					}
				} else {
					nextModuleId = weight.Value(module.id)
				}
				player.SetTestModule(nextModuleId)
			}
		} else {
			Log.ErrorLog("Failed to Run, test module is nil, id = %v", player.TestModule())
		}
	}
}
