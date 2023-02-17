package Net

import (
	"github.com/7058011439/haoqbb/Http"
	"github.com/7058011439/haoqbb/Log"
	"net"
	"strings"
)

// GetOutBoundIP 获取外网ip地址
func GetOutBoundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		Log.ErrorLog("Failed to GetOutBoundIP, err = %v", err)
		return ""
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.String(), ":")[0]
}

// GetProvince 获取ip所在省份
func GetProvince(ip string) string {
	data, _ := Http.PostHttpSync("http://api.tianapi.com/ipquery/index",
		Http.NewHead(nil),
		Http.NewBody(nil).Add("key", "9cb1a3fb5cb7616edb135741fb6b81ef").Add("ip", ip))

	if data["code"].(float64) == 200 {
		return data["newslist"].([]interface{})[0].(map[string]interface{})["province"].(string)
	}
	return ""
}
