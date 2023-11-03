package protocol

// 子命令
const (
	SCmd_C2S_RT = 30001
	SCmd_S2C_RT = 30002

	SCmd_C2S_Login = 30003 // C2S_LoginWithToken
	SCmd_S2C_Login = 30004 // S2C_GameLoginResult

	SCmd_C2S_Nothing_WithReply    = 30005 // 请求一个需要回复的空消息
	SCmd_S2C_Nothing_WithReply    = 30006 // 返回一个空消息
	SCmd_C2S_Nothing_WithOutReply = 30007 // 请求一个没有回复的空消息
	SCmd_S2C_Nothing_WithOutReply = 30008 // 返回一个空消息
)
