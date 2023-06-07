package config

import (
	"encoding/json"
	"fmt"
	"github.com/7058011439/haoqbb/File"
	"github.com/7058011439/haoqbb/Log"
	"io/ioutil"
	"strings"
)

type nodeCfg struct {
	NodeId      int `json:"-"`
	NodeName    string
	ListenAddr  string
	ServiceList []string
}

type clusterCfg struct {
	CenterAddr string
	NodeId     int
	NodeList   map[int]*nodeCfg
}

func (c *clusterCfg) Init() {
	for nodeId, node := range c.NodeList {
		node.NodeId = nodeId
	}
}

var clusterConfig clusterCfg

var serviceConfig = make(map[string]string, 20)

func init() {
	fileData, _ := ioutil.ReadFile("cfg/Cluster.json")
	err := json.Unmarshal(fileData, &clusterConfig)
	if err != nil {
		Log.ErrorLog("Failed to init nodeConfig, err = %v, fileDate = %v", err, fileData)
		return
	}
	clusterConfig.Init()
	fileList, _ := File.WalkFile("cfg", ".json", "")
	for _, file := range fileList {
		pos := strings.LastIndex(file, "\\")
		if strings.LastIndex(file, "/") > pos {
			pos = strings.LastIndex(file, "/")
		}
		name := file[pos+1 : len(file)-5]
		data, _ := ioutil.ReadFile(file)
		serviceConfig[name] = string(data)
	}
}

func GetServiceConfig(serviceName string) string {
	cfg := serviceConfig[fmt.Sprintf("%v_%v", serviceName, clusterConfig.NodeId)]
	if cfg == "" {
		cfg = serviceConfig[serviceName]
	}
	return cfg
}

func GetNodeID() int {
	return clusterConfig.NodeId
}

func GetNodeConfig() *nodeCfg {
	return clusterConfig.NodeList[clusterConfig.NodeId]
}

// GetListenAddr 获取当前节点监听地址
func GetListenAddr() string {
	if clusterConfig.NodeList[clusterConfig.NodeId] != nil {
		return clusterConfig.NodeList[clusterConfig.NodeId].ListenAddr
	}
	return ""
}

// GetCenterAddr 获取中心节点监听地址
func GetCenterAddr() string {
	return clusterConfig.CenterAddr
}

// IsCenterNode 是否中心节点
func IsCenterNode() bool {
	return clusterConfig.NodeId == 0
}
