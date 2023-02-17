package common

import (
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/Util"
	"math"
)

func EncodeMsgType(serverId int16, msgType int16) int {
	return int(serverId)<<16 + int(msgType)
}

func DecodeMsgType(msgType int) (int16, int16) {
	return int16(msgType >> 16), int16(msgType % math.MaxInt16)
}

func EncodeSendMsg(serverId int16, mainCmdId int16, subCmdId int16, data []byte) []byte {
	buff := Stl.NewBuffer(len(data) + 12)
	buff.WriteByte(0xFE)
	buff.Write(Util.Int16ToBytes(0))
	buff.Write(Util.Int16ToBytes(int16(len(data))))
	buff.Write(Util.Int16ToBytes(mainCmdId))
	buff.Write(Util.Int16ToBytes(subCmdId))
	buff.Write(Util.Int16ToBytes(serverId))
	buff.Write(data)
	buff.WriteByte(0xEE)
	return buff.Bytes()
}
