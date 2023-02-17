package node

import (
	"Core/Log"
	"Core/Stl"
	"Core/System"
	"Core/Util"
	"Core/haoqbb/config"
	"Core/haoqbb/service/interface/service"
	"fmt"
	"math"
	"os"
	"reflect"
	"time"
)

var preSetupService = map[string]service.IService{} // 预安装服务(all)
var localNodeService = map[int]service.IService{}   // 本节点服务
var localNodeServiceName = Stl.NewDoubleMap()

func initLog(nodeId int) {
	Log.Init(Log.LevelLog, false, 0, fmt.Sprintf("Logs/Node_%v", nodeId))
}

func startNodeService() {
	for _, s := range localNodeService {
		s.Init()
	}
	for _, s := range localNodeService {
		s.InitMsg()
	}
	for _, s := range localNodeService {
		s.Start()
	}
	// 发现本地服务
	for _, s := range localNodeService {
		for _, o := range localNodeService {
			if s != o {
				s.DiscoverService(o.GetName(), o.GetId())
			}
		}
	}
	for _, s := range localNodeService {
		Log.Log("start %v finished", s.GetName())
	}
}

func Start() {
	args := os.Args
	if len(args) < 2 {
		Log.ErrorLog("Failed to start, nodeId is nil")
		return
	}
	nodeID := Util.StrToInt(args[1])
	nodeConfig := config.GetAllNodeConfig()
	for _, node := range nodeConfig {
		if node.NodeId == nodeID {
			initLog(node.NodeId)
			config.SetCurrNodeID(nodeID)
			System.SetTitle(node.NodeName)
			for _, startServerName := range node.ServiceList {
				if s, ok := preSetupService[startServerName]; ok {
					s.Regedit(config.GetServiceConfig(startServerName))
					localNodeService[s.GetId()] = s
					localNodeServiceName.Add(s.GetId(), s.GetName())
				} else {
					Log.ErrorLog("Failed to start service, service not setup, service name = %v", startServerName)
				}
			}
			break
		}
	}
	startNodeService()
	StartServer()
	StartClient()
	for {
		time.Sleep(time.Second)
	}
}

func Setup(s service.IService) {
	if s.GetName() == "" {
		s.SetName(reflect.Indirect(reflect.ValueOf(s)).Type().Name())
	}
	preSetupService[s.GetName()] = s
}

func getLocalServiceId(serviceName string) int {
	if ret := localNodeServiceName.GetKey(serviceName); ret != nil {
		return ret.(int)
	}
	return math.MaxInt32
}
