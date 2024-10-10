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
	ServiceList []string
	NeedService []string
}

type clusterCfg struct {
	CenterAddr string
	Sign       string
	PerformLog bool
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
		Log.Error("Failed to init nodeConfig, err = %v, fileDate = %v", err, fileData)
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

func SetNodeId(nodeId int) {
	clusterConfig.NodeId = nodeId
}

func GetNodeConfig() *nodeCfg {
	return clusterConfig.NodeList[clusterConfig.NodeId]
}

// GetNeedService 获取该节点服务所需其他服务列表
func GetNeedService() []string {
	if clusterConfig.NodeList[clusterConfig.NodeId] != nil {
		return clusterConfig.NodeList[clusterConfig.NodeId].NeedService
	}
	return nil
}

// GetCenterAddr 获取中心节点监听地址
func GetCenterAddr() string {
	return clusterConfig.CenterAddr
}

// GetSign 获取签名
func GetSign() string {
	return clusterConfig.Sign
}

// IsCenterNode 是否中心节点
func IsCenterNode() bool {
	return clusterConfig.NodeId == 0
}
