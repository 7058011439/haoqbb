package Net

import (
	"sync"
	"time"
)

type speed struct {
	allMsgLen   int64     // 累计消息长度
	msgLenPreS  int64     // 当前秒消息长度
	bandwidth   float64   // 带宽
	lastMsgTime time.Time // 上一次消息时间
	mutex       sync.Mutex
}

func (s *speed) RefreshMsgLen(msgLen int) {
	s.mutex.Lock()
	s.allMsgLen += int64(msgLen)
	s.msgLenPreS += int64(msgLen)
	gas := time.Now().Sub(s.lastMsgTime).Seconds()
	if gas > 1 {
		s.bandwidth = float64(s.msgLenPreS) / gas / 1024
		s.lastMsgTime = time.Now()
		s.msgLenPreS = 0
	}
	s.mutex.Unlock()
}

func (s *speed) GetBandwidth() float64 {
	return s.bandwidth
}

func (s *speed) GetMsgLen() int64 {
	return s.allMsgLen
}
