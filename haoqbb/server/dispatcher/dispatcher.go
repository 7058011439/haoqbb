package dispatcher

import (
	"encoding/json"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	"github.com/7058011439/haoqbb/haoqbb/server/dispatcher/interface"
	"github.com/7058011439/haoqbb/haoqbb/service"
	"github.com/mitchellh/mapstructure"
)

type hitType = int

const (
	random       hitType = iota // 随机
	cpuRate                     // cpu使用率
	memRate                     // 内存使用率
	netRate                     // 网络带宽
	connectCount                // 连接数量
	max
)

type dispatcherConfig struct {
	HitType hitType
	Port    int
}

type Dispatcher struct {
	service.Service
	Net.INetPool
	mapOptimalGate map[hitType]*common.GateInfo // 最优网关
	mapAllGate     map[int]*common.GateInfo     // 所有网关
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
	d.refreshOptimalGate(newGate)
}

func (d *Dispatcher) refreshOptimalGate(gate *common.GateInfo) {
	for t := random; t < max; t++ {
		oldGate := d.mapOptimalGate[t]
		if oldGate == nil {
			d.mapOptimalGate[t] = gate
		} else {
			switch t {
			case random:
				d.mapOptimalGate[t] = gate
			case cpuRate:
				if oldGate.CpuRate > gate.CpuRate {
					d.mapOptimalGate[t] = gate
				}
			case memRate:
				if oldGate.MemRate > gate.MemRate {
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
			}
		}
	}
}

func (d *Dispatcher) loseGateWay(gateWayId int) {
	delete(d.mapAllGate, gateWayId)
	d.mapOptimalGate = map[int]*common.GateInfo{}
	Log.Log("有网关丢失, gateWayId = %v, 剩余网关数量 = %v", gateWayId, len(d.mapAllGate))

	for _, gate := range d.mapAllGate {
		d.refreshOptimalGate(gate)
	}
}
