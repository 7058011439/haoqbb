package common

import (
	"encoding/json"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/Util"
)

/* 简写描述
Gs: 游戏服务器
Hs: 大厅服务器
Ds: 调度器
Gw: 网关
Cl: 客户端
Srv: 所有自定义服(比如:游戏服,大厅服,聊天服等)
Forward: 转发
*/

const (
	GwToDsStatus = iota + 1 // 网关(发送自己状态)到适配器服, 对应结构 GsInfoTag

	GwForwardClToSrv // 网关转发客户端消息到其他服, 对应结构 GwForwardClToSrvTag
	GwForwardSrvToCl // 网关转发其他服消息到客户端, 对应结构 GwForwardSrvToClTag

	SrvPlayerOnLine  // 其他服(大厅、游戏等)玩家在线
	SrvPlayerOffLine // 其他服(大厅、游戏等)玩家离线

	GwClConnect    // 网关有客户端连接，全域广播, 对应结构 Uint64
	GwClDisconnect // 网关客户端断开连接，全域广播, 对应结构 Uint64

	EventLoginSrvLogin // 登录(发送验证结果)到指定游戏服，对应结构 LoginSrvToGameSrv
	MsgMax
)

type ServiceMsg interface {
	Marshal() []byte
	Unmarshal(data []byte)
}

// GsInfoTag 网关状态, 对应 GwToDsStatus 消息
type GsInfoTag struct {
	Addr         string
	MemRate      float64 // 内存占用率
	CpuRate      float64 // Cpu使用率
	NetRate      float64 // 带宽
	ConnectCount int
}

func (g *GsInfoTag) Marshal() []byte {
	// 28 是float64 * 3 + int = 8 * 3 + 4 = 24 + 4 = 28
	buff := Stl.NewBuffer(28 + len(g.Addr))
	buff.WriteFloat64(g.MemRate)
	buff.WriteFloat64(g.CpuRate)
	buff.WriteFloat64(g.NetRate)
	buff.WriteInt(g.ConnectCount)
	buff.WriteString(g.Addr)
	return buff.Bytes()
}

func (g *GsInfoTag) Unmarshal(data []byte) {
	g.MemRate = Util.Float64(data[0:8])
	g.CpuRate = Util.Float64(data[8:16])
	g.NetRate = Util.Float64(data[16:24])
	g.ConnectCount = Util.Int(data[24:28])
	if len(data) > 28 {
		g.Addr = string(data[28:])
	}
}

// GwForwardClToSrvTag 网关转发客户端消息到服务端
type GwForwardClToSrvTag struct {
	ClientId uint64
	CmdId    int
	Data     []byte
}

func (g *GwForwardClToSrvTag) Marshal() []byte {
	buff := Stl.NewBuffer(12 + len(g.Data))
	buff.WriteUInt64(g.ClientId)
	buff.WriteInt(g.CmdId)
	buff.Write(g.Data)
	return buff.Bytes()
}

func (g *GwForwardClToSrvTag) Unmarshal(data []byte) {
	g.ClientId = Util.Uint64(data[0:8])
	g.CmdId = Util.Int(data[8:12])
	if len(data) > 12 {
		g.Data = make([]byte, len(data)-12)
		copy(g.Data, data[12:])
	}
}

// GwForwardSrvToClTag 网关转发服务端消息到客户端
type GwForwardSrvToClTag struct {
	ClientId []uint64 // 这里支持发送给多个用户相同的消息, 如果该字段为空, 则对该服所有连接广播
	CmdId    int
	Data     []byte
}

func (g *GwForwardSrvToClTag) Marshal() []byte {
	// 其中cmdId长度4, ClientId 包头2(列表大小)
	buff := Stl.NewBuffer(len(g.ClientId)*8 + len(g.Data) + 4 + 2)
	buff.WriteInt(g.CmdId)
	buff.WriteInt16(int16(len(g.ClientId)))
	for i := 0; i < len(g.ClientId); i++ {
		buff.WriteUInt64(g.ClientId[i])
	}
	buff.Write(g.Data)
	return buff.Bytes()
}

func (g *GwForwardSrvToClTag) Unmarshal(data []byte) {
	g.CmdId = Util.Int(data[0:4])
	g.ClientId = nil
	l := Util.Int16(data[4:6])
	index := 6
	for i := int16(0); i < l; i++ {
		g.ClientId = append(g.ClientId, Util.Uint64(data[index:index+8]))
		index += 8
	}
	if len(data) > index {
		g.Data = make([]byte, len(data)-index)
		copy(g.Data, data[index:])
	}
}

// LoginSrvToGameSrv 玩家登录结果, 对应 EventLoginSrvLogin
type LoginSrvToGameSrv struct {
	ClientId uint64 `json:"c"`
	OpenId   string `json:"o"`
	Msg      string `json:"m"`
}

func (l *LoginSrvToGameSrv) Marshal() []byte {
	// 这里偷个懒，直接用Json库
	data, _ := json.Marshal(l)
	return data
}

func (l *LoginSrvToGameSrv) Unmarshal(data []byte) {
	// 这里偷个懒，直接用Json库
	json.Unmarshal(data, l)
}

type Uint64 struct {
	Data uint64
}

func (u *Uint64) Marshal() []byte {
	return Util.Uint64ToBytes(u.Data)
}

func (u *Uint64) Unmarshal(data []byte) {
	u.Data = Util.Uint64(data)
}
