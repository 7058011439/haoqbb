package Net

import (
	"encoding/json"
	"github.com/7058011439/haoqbb/Http"
	"github.com/7058011439/haoqbb/Log"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// GetOutBoundIP 获取外网ip地址
func GetOutBoundIP() string {
	/*
		conn, err := net.Dial("udp", "8.8.8.8:53")
		if err != nil {
			Log.ErrorLog("Failed to GetOutBoundIP, err = %v", err)
			return ""
		}
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return strings.Split(localAddr.String(), ":")[0]
	*/
	// 使用一个公共的IP查询API（例如：httpbin.org）来获取公网IP
	resp, err := http.Get("https://httpbin.org/ip")
	if err != nil {
		Log.ErrorLog("无法获取公网IP, err = %v", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应内容
	ipData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.ErrorLog("无法读取IP数据, err = %v", err)
		return ""
	}

	dataTemp := struct {
		Origin string `json:"origin"`
	}{}

	json.Unmarshal(ipData, &dataTemp)

	return dataTemp.Origin
}

// GetInputBoundIP 获取内网ip地址
func GetInputBoundIP() string {
	/*
		conn, err := net.Dial("udp", "8.8.8.8:53")
		if err != nil {
			Log.ErrorLog("Failed to GetOutBoundIP, err = %v", err)
			return ""
		}
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return strings.Split(localAddr.String(), ":")[0]
	*/
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		Log.ErrorLog("获取网络接口出错:%v", err)
		return ""
	}

	// 遍历每个网络接口
	for _, iface := range interfaces {
		if strings.Index(iface.Name, "VPN") != -1 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			Log.ErrorLog("获取地址出错:%v", err)
			continue
		}

		// 遍历每个地址
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}
	return ""
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
