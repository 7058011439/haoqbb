package client

import (
	"Core/Log"
	"Core/Net"
	"Core/Timer"
	"Core/Util"
	"Core/haoqbb/server/gameSrv/client/capability"
	"Core/haoqbb/server/gameSrv/client/interface"
	"Core/haoqbb/server/gameSrv/client/login"
	"Core/haoqbb/server/gameSrv/client/player"
	"Core/haoqbb/server/gameSrv/client/test"
	cProtocol "Core/haoqbb/server/gameSrv/common/protocol"
	"Core/haoqbb/service"
	"Core/haoqbb/service/interface/timer"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"net"
)

type clientConfig struct {
	DispatcherIp   string
	DispatcherPort int
	GateWayIp      string
	GateWayPort    int
	StartID        int
	MaxConn        int
}

type GameClient struct {
	service.Service
	config *clientConfig
}

func (g *GameClient) Init() error {
	if err := mapstructure.Decode(g.ServiceCfg.Other, &g.config); err != nil {
		Log.ErrorLog("Failed to parse client Config, err = %v", err)
		return err
	}
	login.SetStartID(g.config.StartID)
	Interface.SetServiceAgent(g)
	return nil
}

func (g *GameClient) Start() {
	Interface.NewConnManager(g.newConnect, g.disConnect, g.parseProtocol, g.NewTcpMsg)
	ITimer.SetRepeatTimer(Interface.GetServiceName(), 20, g.NewClient)
}

func (g *GameClient) InitMsg() {
	g.RegeditHandleTcpMsg(g.msgHandle)
	g.RegeditMsgHandle(cProtocol.SCmd_S2C_Login, &cProtocol.S2C_GameLoginResult{}, login.S2CLogin)
	g.RegeditMsgHandle(cProtocol.SCmd_S2C_RT, &cProtocol.S2C_Test_RT{}, capability.S2CRT)
}

func (g *GameClient) msgHandle(clientId uint64, data []byte) {
	cmdId := Util.Int16(data[2:4])
	g.DispatchMsg(clientId, 0, int32(cmdId), data[6:])
}

func (g *GameClient) NewClient(Timer.TimerID, ...interface{}) {
	if Interface.GetPlayerCount() >= g.config.MaxConn {
		return
	}
	if g.config.DispatcherIp != "" {
		conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", g.config.DispatcherIp, g.config.DispatcherPort))
		if err != nil {
			Log.ErrorLog("Failed to allPlayer dispatcher, err = %v", err)
			return
		}
		defer conn.Close()
		buff := make([]byte, 1024)
		if n, err := conn.Read(buff); err != nil || n < 10 {
			Log.ErrorLog("Failed to get data from dispatcher, err = %v, n = %v", err, n)
			return
		} else {
			if newConn, err := net.Dial("tcp", string(buff[0:n])); err != nil {
				Log.ErrorLog("Failed to connect gateway, err = %v", err)
				return
			} else {
				Interface.NewClient(newConn)
			}
		}
	} else {
		if newConn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", g.config.GateWayIp, g.config.GateWayPort)); err != nil {
			Log.ErrorLog("Failed to connect gateway, err = %v", err)
			return
		} else {
			Interface.NewClient(newConn)
		}
	}
}

func (g *GameClient) newConnect(client Net.IClient) {
	clientId := client.GetId()
	player := player.NewPlayer(clientId)
	if Interface.GetPlayerCount() == 1 {
		player.SetTimerId(ITimer.SetRepeatTimer(Interface.GetServiceName(), 1000, capability.Main, clientId))
	} else {
		player.SetTimerId(ITimer.SetRepeatTimer(Interface.GetServiceName(), 500, test.Run, clientId))
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
