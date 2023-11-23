package node

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/System"
	"github.com/7058011439/haoqbb/Util"
	"github.com/7058011439/haoqbb/haoqbb/config"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/service"
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
	if len(os.Args) > 1 {
		config.SetNodeId(Util.StrToInt(os.Args[1]))
	}
	initLog(config.GetNodeID())
	if config.IsCenterNode() {
		System.SetTitle("中心节点")
		StartCenterServer()
		Log.Log("启动中心节点完成")
	} else {
		if node := config.GetNodeConfig(); node != nil {
			System.SetTitle(fmt.Sprintf("%v_%v", node.NodeName, node.NodeId))
			for _, startServerName := range node.ServiceList {
				if s, ok := preSetupService[startServerName]; ok {
					s.Regedit(config.GetServiceConfig(startServerName))
					localNodeService[s.GetId()] = s
					localNodeServiceName.Add(s.GetId(), s.GetName())
				} else {
					Log.ErrorLog("Failed to start service, service not setup, service name = %v", startServerName)
				}
			}
			startNodeService()
			StartServer()
			StartClient()
		} else {
			Log.ErrorLog("未找到节点id配置, nodeId = %v", config.GetNodeID())
		}
	}
	for {
		time.Sleep(time.Second)
	}

	/*
		for _, node := range nodeConfig {
			if node.NodeId == nodeID {
				initLog(node.NodeId)
				config.SetNodeID(nodeID)
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
	*/
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
