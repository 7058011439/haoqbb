package protocol

// 子命令
const (
	SCmd_C2S_Anything_Add = 10013 // C2S_Anything_Add
	SCmd_S2C_Anything_Add = 20013 // S2C_Anything_Add
	SCmd_C2S_Anything_Sub = 10014 // C2S_Anything_Sub
	SCmd_S2C_Anything_Sub = 20014 // S2C_Anything_Sub

	SCmd_C2S_HomeUp = 10055 // 请求家园升级(扩容)     C2S_HomeUp
	SCmd_S2C_HomeUp = 20055 // 家园升级(扩容)结果     S2C_OperationResult
)
