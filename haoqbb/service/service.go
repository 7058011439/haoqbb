package service

import (
	"encoding/json"
	"github.com/7058011439/haoqbb/DataBase"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/haoqbb/msgHandle"
	"github.com/7058011439/haoqbb/haoqbb/node"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/http"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/mongo"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/redis"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/timer"
)

type configRedis struct {
	Ip       string
	Port     int
	PassWord string
	Index    int
}

type configMongo struct {
	Ip     string
	Port   int
	DBName string
}

type configMysql struct {
	Ip       string
	Port     int
	DBName   string
	UserName string
	PassWord string
}

type serviceCfg struct {
	Redis *configRedis
	Mongo *configMongo
	Mysql *configMysql
	Id    int
	Other interface{}
}

type Service struct {
	*queue
	Net.INetPool
	ServiceCfg serviceCfg
	msgHandle.IDispatcher
	name string
}

func (s *Service) run() {
	s.queue.run()
}

func (s *Service) SetName(name string) {
	s.name = name
}

func (s *Service) GetName() string {
	return s.name
}

func (s *Service) GetId() int {
	return s.ServiceCfg.Id
}

func (s *Service) Regedit(serviceCfg string) {
	if err := json.Unmarshal([]byte(serviceCfg), &s.ServiceCfg); err != nil {
		Log.ErrorLog("Failed to json.Unmarshal on RegeditApi, err = %v", err)
	}
	s.queue = NewQueue(s.name)
	s.setAgent()
	go s.run()
}

func (s *Service) Init() error {
	return nil
}

func (s *Service) setAgent() {
	if s.ServiceCfg.Mongo != nil {
		cfg := s.ServiceCfg.Mongo
		s.MongoDB = DataBase.NewMongoDB(cfg.Ip, cfg.Port, cfg.DBName, "", "", 0)
		IMongo.SetMongoAgent(s)
	}
	if s.ServiceCfg.Redis != nil {
		cfg := s.ServiceCfg.Redis
		s.RedisDB = DataBase.NewRedisDB(cfg.Ip, cfg.Port, cfg.PassWord, cfg.Index)
		IRedis.SetRedisAgent(s)
	}

	ITimer.SetTimerAgent(s)
	IHttp.SetHttpAgent(s)
}

func (s *Service) Start() {

}

func (s *Service) InitMsg() {

}

func (s *Service) InitTcpServer(port int, connect Net.ConnectHandle, disconnect Net.ConnectHandle, parse Net.ParseProtocol, fun func(clientId uint64, data []byte), options ...Net.Options) {
	options = append(options, Net.WithPoolId(s.GetId()))
	s.INetPool = Net.NewTcpServer(port, connect, disconnect, parse, s.NewTcpMsg, options...)
	s.RegeditHandleTcpMsg(fun)
}

func (s *Service) RegeditServiceMsg(msgType int, fun func(srcServiceId int, data []byte)) {
	s.serviceMsgHandle[msgType] = fun
}

func (s *Service) RegeditDiscoverService(serviceName string, fun func(int)) {
	s.discoverServiceHandle[serviceName] = fun
}

func (s *Service) RegeditLoseService(serviceName string, fun func(int)) {
	s.loseServiceHandle[serviceName] = fun
}

func (s *Service) PublicEventByName(serviceName string, eventType int, data interface{}) {
	sendData, _ := json.Marshal(data)
	s.SendMsgToServiceByName(serviceName, eventType, sendData)
}

func (s *Service) PublicEventById(serviceId int, eventType int, data interface{}) {
	sendData, _ := json.Marshal(data)
	s.SendMsgToServiceById(serviceId, eventType, sendData)
}

func (s *Service) SendMsgToServiceByName(serviceName string, msgType int, data []byte) {
	node.SendMsgByName(s.GetId(), serviceName, msgType, data)
}

func (s *Service) SendMsgToServiceById(serviceId int, msgType int, data []byte) {
	node.SendMsgById(s.GetId(), serviceId, msgType, data)
}
