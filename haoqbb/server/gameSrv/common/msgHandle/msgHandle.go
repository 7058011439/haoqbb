package msgHandle

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Timer"
	"github.com/golang/protobuf/proto"
)

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		msgRoute: make(map[int32]*msgHandle, 1024),
	}
}

type HandleFun func(msg *ClientMsg)

type msgHandle struct {
	msg proto.Message // 消息结构
	fun HandleFun     // 对应函数
}

type ClientMsg struct {
	ClientId uint64      // 客户端id
	UserId   int         // userId
	Data     interface{} // 数据
}

type Dispatcher struct {
	msgRoute map[int32]*msgHandle // map[子命令]消息处理器
}

func (d *Dispatcher) DispatchMsg(clientId uint64, userId int, cmdId int32, data []byte) {
	if info, ok := d.msgRoute[cmdId]; !ok {
		Log.Error("Failed to DispatchMsg, unknown cmdId, cmdId = %v", cmdId)
		return
	} else {
		if err := proto.Unmarshal(data, info.msg); err != nil {
			Log.Error("Failed to DispatchMsg, proto.Unmarshal error, cmdId = %v, error = %v", cmdId, err.Error())
			return
		} else {
			cost := Timer.NewTiming(Timer.Millisecond)
			info.fun(&ClientMsg{
				ClientId: clientId,
				UserId:   userId,
				Data:     info.msg,
			})
			cost.PrintCost(10, true, "DispatchMsg info.fun, clientId = %v, userId = %v, cmdId = %v, info.Msg = %v", clientId, userId, cmdId, info.msg)
		}
	}
}

func (d *Dispatcher) RegeditMsgHandle(cmdId int32, msg proto.Message, fun HandleFun) {
	if d.msgRoute == nil {
		d.msgRoute = make(map[int32]*msgHandle)
	}
	if d.msgRoute[cmdId] != nil {
		Log.Warn("Failed to RegeditMsgHandle, cmd repeat regedit, cmdId = %v", cmdId)
	}
	d.msgRoute[cmdId] = &msgHandle{msg: msg, fun: fun}
}
