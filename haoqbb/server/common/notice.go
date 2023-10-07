package common

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

	GwClConnect    // 网关有客户端连接，全域广播
	GwClDisconnect // 网关客户端断开连接，全域广播

	EventLoginSrvLogin // 登录(发送验证结果)到指定游戏服，对应结构 LoginSrvToGameSrv
	MsgMax
)

// GsInfoTag 网关状态, 对应 GwToDsStatus 消息
type GsInfoTag struct {
	Addr         string
	MemRate      float64 // 内存占用率
	CpuRate      float64 // Cpu使用率
	NetRate      float64 // 带框
	ConnectCount int
}

// GwForwardClToSrvTag 网关转发客户端消息到服务端
type GwForwardClToSrvTag struct {
	ClientId uint64
	CmdId    int
	Data     []byte
}

// GwForwardSrvToClTag 网关转发服务端消息到客户端
type GwForwardSrvToClTag struct {
	ClientId []uint64 // 这里支持发送给多个用户相同的消息, 如果该字段为空, 则对该服所有连接广播
	CmdId    int
	Data     []byte
}

// LoginSrvToGameSrv 玩家登录结果, 对应 EventLoginSrvLogin
type LoginSrvToGameSrv struct {
	ClientId uint64
	UserId   int
	Msg      string
}
