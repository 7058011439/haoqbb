package common

const (
	GateToGameSrvClientMsg = iota + 1 // 网关(转发客户端消息)到游戏服, 对应结构 GateWayToGameSrv
	GameSrvToGateClientMsg            // 游戏服到网关(转发消息到客户端), 对应结构 GameSrvToGateWay

	GateToLoginSrvClientMsg // 网关(转发客户端消息-目前就登录命令)到登录服务器, 对应结构 GateWayToLoginSrv

	GateToDispatcherStatus // 网关(发送自己状态)到适配器服, 对应结构 GateInfo

	GameSrvPlayerOnLine  // 游戏端有玩家上线，全域广播
	GameSrvPlayerOffLine // 游戏端有玩家离线，全域广播

	GateWayClientConnect    // 网关有客户端连接，全域广播
	GateWayClientDisconnect // 网关客户端断开连接，全域广播

	EventLoginSrvLogin // 登录(发送验证结果)到指定游戏服，对应结构 LoginSrvToGameSrv
)

// GateInfo 网关状态，对应EventG2DGateWayRegedit消息
type GateInfo struct {
	Addr         string
	MemRate      float64
	NetRate      float64
	ConnectCount int
}

// LoginSrvToGameSrv 玩家登录结果，对应 EventLoginSrvLogin，从登录服到游戏服
type LoginSrvToGameSrv struct {
	ClientId uint64
	Ret      int
	OpenId   int
	Msg      string
}

// GateWayToLoginSrv 网关转发登录消息 GateToLoginSrvClientMsg
type GateWayToLoginSrv struct {
	GameSrvId int
	ClientId  uint64
	CmdId     int
	Data      []byte
}

type GateWayToGameSrv struct {
	ClientId uint64
	CmdId    int
	Data     []byte
}

type GameSrvToGateWay struct {
	ClientId []uint64
	CmdId    int
	Data     []byte
}
