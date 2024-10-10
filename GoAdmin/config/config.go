package config

import (
	"encoding/json"
	"github.com/7058011439/haoqbb/Log"
	"io/ioutil"
)

type mysql struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"userName"`
	PassWord string `json:"password"`
	DBName   string `json:"dbName"`
	Log      bool   `json:"log"`
	Init     bool   `json:"init"`
	Cache    bool   `json:"cache"`
}

type http struct {
	Port    int    `json:"port"`
	Version string `json:"version"`
	JWTKey  []byte `json:"jwtKey"`
}

type config struct {
	Mysql mysql `json:"mysql"`
	Http  http  `json:"http"`
}

var stConfig config

func init() {
	Init()
}

func Init() {
	if fileData, err := ioutil.ReadFile("config_core.json"); err != nil {
		Log.Error("加载配置文件失败, err = %v", err)
	} else {
		json.Unmarshal(fileData, &stConfig)
	}
}

func MysqlHost() string {
	return stConfig.Mysql.Host
}

func MysqlPort() int {
	return stConfig.Mysql.Port
}

func MysqlUserName() string {
	return stConfig.Mysql.UserName
}

func MysqlPassWord() string {
	return stConfig.Mysql.PassWord
}

func MysqlDBName() string {
	return stConfig.Mysql.DBName
}

func MysqlLog() bool {
	return stConfig.Mysql.Log
}

func MysqlInit() bool {
	return stConfig.Mysql.Init
}

func HttpPort() int {
	return stConfig.Http.Port
}

func HttpVersion() string {
	return stConfig.Http.Version
}

func HttpJWTKey() []byte {
	return stConfig.Http.JWTKey
}
