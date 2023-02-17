package service

type IService interface {
	SetName(serviceName string)                               // 安装前调用
	GetName() string                                          // 无需重写
	GetId() int                                               // 无需重写
	Regedit(serviceCfg string)                                // 无需重写
	Init() error                                              // 有需求重写
	InitMsg()                                                 // 有需求重写
	Start()                                                   // 有需求重写
	NewServiceMsg(srcServiceId int, msgType int, data []byte) // 收到其他服务消息，无需重写
	DiscoverService(string, int)                              // 发现其他(节点)服务，无需重写
	LoseService(string, int)                                  // 遗失其他(节点)服务，无需重写
}
