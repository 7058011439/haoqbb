package Net

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/Timer"
	"net"
	"strings"
	"sync"
	"sync/atomic"
)

type Client struct {
	conn       net.Conn      // 网络连接
	id         uint64        // 客户端id
	customData interface{}   // 自定义数据
	recvBuff   *Stl.Buffer   // 接受缓存池
	sendBuff   *Stl.Buffer   // 缓存池
	chClose    chan struct{} // 关闭客户端(先通知接受协程结束，接受协程通知发送协程)
	sendMutex  sync.RWMutex  // 发送数据锁
	timerId    Timer.TimerID // 延时发送定时器ID
	INetPool
}

func (c *Client) GetId() uint64 {
	return c.id
}

func (c *Client) CustomData() interface{} {
	return c.customData
}

func (c *Client) SetCustomData(data interface{}) {
	c.customData = data
}

func (c *Client) GetIp() string {
	if c.conn != nil {
		return strings.Split(c.GetAddr(), ":")[0]
	} else {
		Log.ErrorLog("Failed to GetIp, tcpConn is nil")
		return "0.0.0.0"
	}
}

func (c *Client) GetAddr() string {
	if c.conn != nil {
		return c.conn.RemoteAddr().String()
	} else {
		Log.ErrorLog("Failed to GetIp, tcpConn is nil")
		return "0.0.0.0:0"
	}
}

// Close 关闭连接，先通知接受协程退出，接受协程退出后通知发送协程处理(将待发送数据发送，然后关闭端口)
func (c *Client) Close() {
	c.chClose <- struct{}{}
}

func (c *Client) SendMsg(data []byte) {
	//if _, err := c.conn.Write(data); err != nil {
	//	Log.ErrorLog("Failed to conn.write, err = %v, data = %v, clientId = %v", err, data, c.GetId())
	//}
	c.sendMutex.Lock()
	c.sendBuff.Write(data)
	c.sendMutex.Unlock()
	if c.sendBuff.Len() >= c.getSendPackageSize() {
		c.send(0)
	} else {
		if atomic.LoadInt64(&c.timerId) == 0 {
			atomic.StoreInt64(&c.timerId, Timer.AddOnceTimer(1, c.send))
		}
	}
}

func (c *Client) send(timerId Timer.TimerID, _ ...interface{}) {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()
	if c.sendBuff.Len() < 1 {
		return
	}
	if _, err := c.conn.Write(c.sendBuff.Bytes()); err != nil {
		Log.ErrorLog("Failed to conn.write, err = %v, data = %v, clientId = %v", err, c.sendBuff.Bytes(), c.GetId())
	}
	c.sendBuff.Reset()
	if timerId == 0 && atomic.LoadInt64(&c.timerId) != 0 {
		Timer.CloseTimer(c.timerId)
	}
	atomic.StoreInt64(&c.timerId, 0)
}

func (c *Client) revMsg() {
	defer func() {
		c.send(0)
		c.conn.Close()
		c.onDisconnect(c)
	}()
	buf := make([]byte, revCacheSize)
	for {
		select {
		case <-c.chClose:
			return
		default:
			n, err := c.conn.Read(buf)
			if err == nil && n > 0 {
				c.recvBuff.Write(buf[:n])
				buff := c.recvBuff.Bytes()
				i := 0
				for i = 0; i < len(buff); {
					if data, offSize := c.onParseProtocol(buff[i:]); offSize > 0 {
						msg := make([]byte, len(data))
						copy(msg, data)
						c.onHandleMsg(c, msg)
						i += offSize
					} else {
						break
					}
				}
				if i > 0 {
					c.recvBuff.OffSize(i)
				}
				if c.recvBuff.Len() > c.getRecvPackageLimit() {
					Log.ErrorLog("rev buff to long, size = %v", c.recvBuff.Len())
					return
				}
			} else {
				//Log.ErrorLog("Failed to read from client, err = %v, clientId = %v", err, c.GetId())
				return
			}
		}
	}
}
