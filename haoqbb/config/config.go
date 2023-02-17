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
	NodeId      int
	NodeName    string
	ListenAddr  string
	ServiceList []string
}

var nodeConfig = make(map[int]*nodeCfg, 20)
var serviceConfig = make(map[string]string, 20)
var nodeID int

func init() {
	fileData, _ := ioutil.ReadFile("cfg/Cluster.json")
	var data []*nodeCfg
	err := json.Unmarshal(fileData, &data)
	if err != nil {
		Log.ErrorLog("Failed to init nodeConfig, err = %v, fileDate = %v", err, fileData)
		return
	}
	for _, node := range data {
		nodeConfig[node.NodeId] = node
	}
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

func GetAllNodeConfig() map[int]*nodeCfg {
	return nodeConfig
}

func GetNodeConfig(nodeId int) *nodeCfg {
	return nodeConfig[nodeId]
}

func GetServiceConfig(serviceName string) string {
	cfg := serviceConfig[fmt.Sprintf("%v_%v", serviceName, nodeID)]
	if cfg == "" {
		cfg = serviceConfig[serviceName]
	}
	return cfg
}

func SetCurrNodeID(nodeId int) {
	nodeID = nodeId
}

func GetCurrNodeId() int {
	return nodeID
}

func GetCurrNodeListenAddr() string {
	if nodeConfig[nodeID] != nil {
		return nodeConfig[nodeID].ListenAddr
	}
	return ""
}
