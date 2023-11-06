package dispatcher

import (
	"github.com/7058011439/haoqbb/Http"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	"github.com/7058011439/haoqbb/haoqbb/server/dispatcher/Interface"
	"github.com/7058011439/haoqbb/haoqbb/service"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"sync"
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
	*Http.Server
	service.Service
	mapOptimalGate map[hitType]*common.GsInfoTag // 最优网关
	mapAllGate     map[int]*common.GsInfoTag     // 所有网关
	config         *dispatcherConfig
	mutex          sync.Mutex
}

func (d *Dispatcher) Init() error {
	if err := mapstructure.Decode(d.ServiceCfg.Other, &d.config); err != nil {
		Log.ErrorLog("Failed to parse dispatcher Config, err = %v", err)
	}
	//d.INetPool = Net.NewTcpServer(d.config.Port, d.connect, nil, nil, nil, Net.WithPoolId(d.GetId()))
	d.Server = Http.NewHttpServer("release")
	d.mapOptimalGate = map[hitType]*common.GsInfoTag{}
	d.mapAllGate = map[int]*common.GsInfoTag{}
	Interface.SetServiceAgent(d)
	return nil
}

func (d *Dispatcher) Start() {
	d.Server.Start(d.config.Port)
	d.RegeditLoseService(common.GateWay, d.loseGateWay)
}

func (d *Dispatcher) InitMsg() {
	d.RegeditServiceMsg(common.GwToDsStatus, d.gateWayRegedit)
	d.Server.RegeditApi(Http.TypeGet, "getgw", d.getGw)
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

func (d *Dispatcher) getGw(c *gin.Context) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ret := Http.NewResult(c)
	if gate, ok := d.mapOptimalGate[d.config.HitType]; ok {
		ret.Success("ok", gate.Addr)
	} else {
		ret.Fail("", nil)
	}
}

func (d *Dispatcher) gateWayRegedit(srcServiceId int, data []byte) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	var newGate = &common.GsInfoTag{}
	newGate.Unmarshal(data)
	d.mapAllGate[srcServiceId] = newGate
	d.refreshOptimalGate()
}

func (d *Dispatcher) loseGateWay(gateWayId int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.mapAllGate, gateWayId)
	Log.Log("有网关丢失, gateWayId = %v, 剩余网关数量 = %v", gateWayId, len(d.mapAllGate))
	d.refreshOptimalGate()
}

func (d *Dispatcher) refreshOptimalGate() {
	d.mapOptimalGate = map[int]*common.GsInfoTag{}
	if len(d.mapAllGate) < 1 {
		return
	}
	for _, gate := range d.mapAllGate {
		d.mapOptimalGate[random] = gate
		break
	}
	for _, gate := range d.mapAllGate {
		if oldGate := d.mapOptimalGate[cpuRate]; oldGate == nil || oldGate.CpuRate > gate.CpuRate {
			d.mapOptimalGate[cpuRate] = gate
		}
	}
	for _, gate := range d.mapAllGate {
		if oldGate := d.mapOptimalGate[memRate]; oldGate == nil || oldGate.MemRate > gate.MemRate {
			d.mapOptimalGate[memRate] = gate
		}
	}
	for _, gate := range d.mapAllGate {
		if oldGate := d.mapOptimalGate[netRate]; oldGate == nil || oldGate.NetRate > gate.NetRate {
			d.mapOptimalGate[netRate] = gate
		}
	}
	for _, gate := range d.mapAllGate {
		if oldGate := d.mapOptimalGate[connectCount]; oldGate == nil || oldGate.ConnectCount > gate.ConnectCount {
			d.mapOptimalGate[connectCount] = gate
		}
	}
}
