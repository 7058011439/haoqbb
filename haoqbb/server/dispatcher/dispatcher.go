package dispatcher

import (
	"Core/Log"
	"Core/Net"
	"Core/haoqbb/server/common"
	"Core/haoqbb/server/dispatcher/interface"
	"Core/haoqbb/service"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
)

type hitType = int

const (
	cpuRate      hitType = iota // cpu使用率
	netRate                     // 网络带宽
	connectCount                // 连接数量
	random                      // 随机
)

type dispatcherConfig struct {
	HitType hitType
	Port    int
}

type Dispatcher struct {
	service.Service
	Net.INetPool
	mapOptimalGate map[hitType]*common.GateInfo
	mapAllGate     map[int]*common.GateInfo
	config         *dispatcherConfig
}

func (d *Dispatcher) Init() error {
	if err := mapstructure.Decode(d.ServiceCfg.Other, &d.config); err != nil {
		Log.ErrorLog("Failed to parse dispatcher Config, err = %v", err)
	}
	d.INetPool = Net.NewTcpServer(d.config.Port, d.connect, d.disConnect, nil, nil, Net.WithPoolId(d.GetId()))
	d.mapOptimalGate = map[hitType]*common.GateInfo{}
	d.mapAllGate = map[int]*common.GateInfo{}
	Interface.SetServiceAgent(d)
	return nil
}

func (d *Dispatcher) Start() {
	d.StartServer()
	d.RegeditLoseService(common.GateWay, d.loseGateWay)
}

func (d *Dispatcher) InitMsg() {
	d.RegeditServiceMsg(common.GateToDispatcherStatus, d.gateWayRegedit)
}

func (d *Dispatcher) connect(client Net.IClient) {
	data := "error"
	if gate, ok := d.mapOptimalGate[d.config.HitType]; ok {
		data = gate.Addr
	} else {
		Log.ErrorLog("Failed to Dispatcher connect, d.mapOptimalGate = %v", d.mapOptimalGate)
	}
	client.SendMsg([]byte(data))
	client.Close()
	//Log.Debug("Client request gateway, ip = %v, return = %v", client.GetAddr(), data)
}

func (d *Dispatcher) disConnect(_ Net.IClient) {

}

func (d *Dispatcher) gateWayRegedit(srcServiceId int, data []byte) {
	var newGate = &common.GateInfo{}
	if err := json.Unmarshal(data, newGate); err != nil {
		Log.ErrorLog("Failed to json.Unmarshal on gateWayRegedit, err = %v, data = %v", err, data)
		return
	}
	d.mapAllGate[srcServiceId] = newGate
	for t := cpuRate; t <= random; t++ {
		oldGate := d.mapOptimalGate[t]
		if oldGate == nil {
			d.mapOptimalGate[t] = newGate
		} else {
			switch t {
			case cpuRate:
				if oldGate.CpuRate > newGate.CpuRate {
					d.mapOptimalGate[t] = newGate
				}
			case netRate:
				if oldGate.NetRate > newGate.NetRate {
					d.mapOptimalGate[t] = newGate
				}
			case connectCount:
				if oldGate.ConnectCount > newGate.ConnectCount {
					d.mapOptimalGate[t] = newGate
				}
			case random:
				d.mapOptimalGate[t] = newGate
			}
		}
	}
}

func (d *Dispatcher) loseGateWay(gateWayId int) {
	delete(d.mapAllGate, gateWayId)
	d.mapOptimalGate = map[int]*common.GateInfo{}
	Log.WarningLog("有网关丢失, gateWayId = %v, 剩余网关数量 = %v", gateWayId, len(d.mapOptimalGate))

	for _, gate := range d.mapAllGate {
		for t := cpuRate; t <= random; t++ {
			oldGate := d.mapOptimalGate[t]
			if oldGate == nil {
				d.mapOptimalGate[t] = gate
			} else {
				switch t {
				case cpuRate:
					if oldGate.CpuRate > gate.CpuRate {
						d.mapOptimalGate[t] = gate
					}
				case netRate:
					if oldGate.NetRate > gate.NetRate {
						d.mapOptimalGate[t] = gate
					}
				case connectCount:
					if oldGate.ConnectCount > gate.ConnectCount {
						d.mapOptimalGate[t] = gate
					}
				case random:
					d.mapOptimalGate[t] = gate
				}
			}
		}
	}
}
