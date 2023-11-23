package client

import (
	"encoding/json"
	"github.com/7058011439/haoqbb/Http"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/Timer"
	"github.com/7058011439/haoqbb/Util"
	"github.com/7058011439/haoqbb/haoqbb/msgHandle"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/capability"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/interface"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/login"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/other"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/player"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/test"
	cProtocol "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
	"github.com/7058011439/haoqbb/haoqbb/service"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/timer"
	_ "github.com/mitchellh/mapstructure"
	"net"
)

type testConfig struct {
	Open       bool        // 测试用例是否生效
	Entrance   int         // 直接进入该模块权重
	NextModule map[int]int // 指定下一模块权重
}

type clientConfig struct {
	DispatcherAddr string
	GateWayAddr    string
	StartID        int
	MaxConn        int
	Test           map[int]*testConfig
}

type GameClient struct {
	service.Service
	config *clientConfig
}

func interfaceToStruct(data interface{}, data1 interface{}) {
	d, _ := json.Marshal(data)
	json.Unmarshal(d, data1)
}

func (g *GameClient) Init() error {
	d, _ := json.Marshal(g.ServiceCfg.Other)
	json.Unmarshal(d, &g.config)
	for id, t := range g.config.Test {
		if t.Open {
			test.InsertTestModule(id, t.Entrance, t.NextModule)
		}
	}
	test.InitOver()
	login.SetStartID(g.config.StartID)
	Interface.SetServiceAgent(g)
	return nil
}

func (g *GameClient) Start() {
	Interface.NewConnManager(g.newConnect, g.disConnect, g.parseProtocol, g.NewTcpMsg)
	ITimer.SetRepeatTimer(Interface.GetServiceName(), 50, g.NewClient)
}

func (g *GameClient) InitMsg() {
	g.RegeditHandleTcpMsg(g.msgHandle)
	g.IDispatcher = msgHandle.NewPBDispatcher()
	g.RegeditMsgHandle(cProtocol.SCmd_S2C_Login, &cProtocol.S2C_GameLoginResult{}, login.S2CLogin)
	g.RegeditMsgHandle(cProtocol.SCmd_S2C_RT, &cProtocol.S2C_Test_RT{}, capability.S2CRT)
	g.RegeditMsgHandle(cProtocol.SCmd_S2C_Nothing_WithReply, &cProtocol.S2C_Test_RT{}, other.S2CNothingWithReply)
}

func (g *GameClient) msgHandle(clientId uint64, data []byte) {
	cmdId := Util.Int16(data[2:4])
	g.DispatchMsg(clientId, 0, int32(cmdId), data[6:])
}

func (g *GameClient) NewClient(Timer.TimerID, ...interface{}) {
	if Interface.GetPlayerCount() >= g.config.MaxConn {
		return
	}
	gateWayAddr := g.config.GateWayAddr
	if g.config.DispatcherAddr != "" {
		if data, err := Http.GetHttpSync(g.config.DispatcherAddr, nil, nil); err == nil {
			mapData := map[string]interface{}{}
			json.Unmarshal(data, &mapData)
			if mapData["code"].(float64) == 200 {
				gateWayAddr = mapData["data"].(string)
			} else {
				Log.ErrorLog("获取网关信息失败, err = %v", mapData["msg"].(string))
			}
		} else {
			Log.ErrorLog("获取网关信息失败, err = %v", err)
			return
		}
	}

	if newConn, err := net.Dial("tcp", gateWayAddr); err != nil {
		Log.ErrorLog("连接网关失败, err = %v", err)
		return
	} else {
		Interface.NewClient(newConn)
	}
}

func (g *GameClient) newConnect(client Net.IClient) {
	clientId := client.GetId()
	player := player.NewPlayer(clientId)
	if Interface.GetPlayerCount() == 1 {
		player.SetTimerId(ITimer.SetRepeatTimer(Interface.GetServiceName(), 1000, capability.Main, clientId))
	} else {
		player.SetTimerId(ITimer.SetRepeatTimer(Interface.GetServiceName(), 1000, test.Run, clientId))
	}
	login.C2SLogin(clientId)
}

func (g *GameClient) disConnect(client Net.IClient) {
	player.RemovePlayer(client.GetId())
}

func (g *GameClient) parseProtocol(data []byte) ([]byte, int) {
	if len(data) < 12 {
		return nil, 0
	}
	endPos := int(Util.Uint16(data[3:5])) + 12
	if len(data) >= endPos {
		return data[5 : endPos-1], endPos
	}
	return nil, 0
}
