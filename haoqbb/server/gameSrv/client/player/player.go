package player

import (
	"github.com/7058011439/haoqbb/Timer"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/interface"
	"github.com/golang/protobuf/proto"
)

type User struct {
	UserId int
}

type IPlayer interface {
	IsLogin() bool
	SetTestModule(id int)
	TestModule() int
	UpdateData(data interface{})
	Data() interface{}
	SendMsgToServer(cmdId int16, msg proto.Message)
	SetTimerId(id Timer.TimerID)
	TimerId() Timer.TimerID
}

type Player struct {
	data       interface{}
	clientId   uint64 // 客户端id
	testModule int    // 测试模块
	timerId    Timer.TimerID
}

func (p *Player) IsLogin() bool {
	return p.data != nil
}

func (p *Player) SetTestModule(id int) {
	if id == 3 {
		id = 3
	}
	p.testModule = id
}

func (p *Player) TestModule() int {
	return p.testModule
}

func (p *Player) UpdateData(data interface{}) {
	p.data = data
}

func (p *Player) Data() interface{} {
	return p.data
}

func (p *Player) SendMsgToServer(cmdId int16, msg proto.Message) {
	Interface.SendMsgToServer(p.clientId, cmdId, msg)
}

func (p *Player) SetTimerId(id Timer.TimerID) {
	p.timerId = id
}

func (p *Player) TimerId() Timer.TimerID {
	return p.timerId
}
