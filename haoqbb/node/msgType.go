package node

import (
	"encoding/json"
	"fmt"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/Util"
)

type INodeMsg interface {
	Marshal() []byte
	Unmarshal(data []byte)
}

type N2NMsg struct {
	DestServiceId int
	SrcServerId   int
	MsgType       int
	Data          []byte
}

func (n *N2NMsg) Marshal() []byte {
	newBuff := Stl.NewBuffer(12 + len(n.Data))
	newBuff.WriteInt(n.DestServiceId)
	newBuff.WriteInt(n.SrcServerId)
	newBuff.WriteInt(n.MsgType)
	newBuff.Write(n.Data)
	return newBuff.Bytes()
}

func (n *N2NMsg) Unmarshal(data []byte) {
	n.DestServiceId = Util.Int(data[0:4])
	n.SrcServerId = Util.Int(data[4:8])
	n.MsgType = Util.Int(data[8:12])
	n.Data = make([]byte, len(data)-12)
	copy(n.Data, data[12:])
}

func (n N2NMsg) String() string {
	return fmt.Sprintf("destServiceId:%v srcServerId:%v msgType:%v data:\"%v\"", n.DestServiceId, n.SrcServerId, n.MsgType, string(n.Data))
}

// NodeInfo 节点信息
type NodeInfo struct {
	NodeId      int32
	NodeName    string
	Addr        string
	NeedService []string
	ServiceList []string
}

func (n *NodeInfo) Marshal() []byte {
	ret, _ := json.Marshal(n)
	return ret
}

func (n *NodeInfo) Unmarshal(data []byte) {
	json.Unmarshal(data, n)
}

// NodeList 节点列表
type NodeList struct {
	NodeList []*NodeInfo
}

func (n *NodeList) Marshal() []byte {
	ret, _ := json.Marshal(n)
	return ret
}

func (n *NodeList) Unmarshal(data []byte) {
	json.Unmarshal(data, n)
}

type ServiceInfo struct {
	ServiceName string
	ServiceId   int
}

func (s *ServiceInfo) Marshal() []byte {
	ret, _ := json.Marshal(s)
	return ret
}

func (s *ServiceInfo) Unmarshal(data []byte) {
	json.Unmarshal(data, s)
}

type N2NRegedit struct {
	ServiceList []*ServiceInfo
}

func (n *N2NRegedit) Marshal() []byte {
	ret, _ := json.Marshal(n)
	return ret
}

func (n *N2NRegedit) Unmarshal(data []byte) {
	json.Unmarshal(data, n)
}
