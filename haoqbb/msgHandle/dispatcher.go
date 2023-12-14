package msgHandle

type IDispatcher interface {
	DispatchMsg(clientId uint64, userId int64, cmdId int32, data []byte)
	RegeditMsgHandle(cmdId int32, msgType interface{}, fun HandleFun)
}

type HandleFun func(msg *ClientMsg)

type ClientMsg struct {
	ClientId uint64      // 客户端id
	UserId   int64       // userId
	Data     interface{} // 数据
}

type msgHandle struct {
	msg interface{} // 消息结构
	fun HandleFun   // 对应函数
}
