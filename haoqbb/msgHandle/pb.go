package msgHandle

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/golang/protobuf/proto"
)

func NewPBDispatcher() *PBDispatcher {
	return &PBDispatcher{
		msgRoute: make(map[int32]*msgHandle, 1024),
	}
}

type PBDispatcher struct {
	msgRoute map[int32]*msgHandle // map[子命令]消息处理器
}

func (d *PBDispatcher) DispatchMsg(clientId uint64, userId int64, cmdId int32, data []byte) {
	if info, ok := d.msgRoute[cmdId]; !ok {
		Log.ErrorLog("Failed to DispatchMsg, unknown cmdId, cmdId = %v", cmdId)
		return
	} else {
		if err := proto.Unmarshal(data, info.msg.(proto.Message)); err != nil {
			Log.ErrorLog("Failed to DispatchMsg, proto.Unmarshal error, cmdId = %v, error = %v", cmdId, err.Error())
			return
		} else {
			info.fun(&ClientMsg{
				ClientId: clientId,
				UserId:   userId,
				Data:     info.msg,
			})
		}
	}
}

func (d *PBDispatcher) RegeditMsgHandle(cmdId int32, msgType interface{}, fun HandleFun) {
	if d.msgRoute == nil {
		d.msgRoute = make(map[int32]*msgHandle)
	}
	if d.msgRoute[cmdId] != nil {
		Log.WarningLog("Failed to RegeditMsgHandle, cmd repeat regedit, cmdId = %v", cmdId)
	}
	d.msgRoute[cmdId] = &msgHandle{msg: msgType, fun: fun}
}
