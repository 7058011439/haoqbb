package main

import (
	"encoding/json"
	"fmt"
	"github.com/7058011439/haoqbb/Http"
	"github.com/7058011439/haoqbb/Timer"
	"github.com/7058011439/haoqbb/haoqbb/node"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	_ "github.com/7058011439/haoqbb/haoqbb/server/dispatcher"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/bujidao/client"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/bujidao/server"
	_ "github.com/7058011439/haoqbb/haoqbb/server/gateWay"
	_ "github.com/7058011439/haoqbb/haoqbb/server/loginSrv"
	"strings"
)

func init() {
	gameServer := new(server.BuJiDaoSrv)
	gameServer.SetName(common.GameServerBuJiDao)

	gameClient := new(client.GameClient)
	gameClient.SetName(common.GameClientBuJiDao)

	node.Setup(gameServer)
	node.Setup(gameClient)
}

func main() {
	byteData, _ := Http.GetHttpSync("http://192.168.13.220:8086/", Http.NewHead(nil))
	data := string(byteData)
	//fmt.Println(data)
	sessionKey := "var _SessionID = "
	posBegin := strings.Index(data, sessionKey)
	posEnd := strings.Index(data[posBegin:], "\n")
	session := data[posBegin+len(sessionKey)+1 : posBegin+posEnd-2]
	fmt.Println(session)

	serverKey := "var serverLabels = "
	posBegin = strings.Index(data, serverKey)
	posEnd = strings.Index(data[posBegin:], "\n")
	server := data[posBegin+len(serverKey) : posBegin+posEnd-1]
	//fmt.Println(server)

	var serverList []string
	json.Unmarshal([]byte(server), &serverList)

	fmt.Println(serverList)
	t := Timer.NewTiming(Timer.Second)
	query, _ := Http.PostHttpSync("http://192.168.13.220:8086/API/CreateQuery.html",
		Http.NewHead(nil),
		Http.NewBody(nil).
			Add("servernid", "76Rx光辉岁月2021.11.25").
			Add("session", session).
			Add("account", "").
			Add("name", "").
			Add("mac", "").
			Add("uuid", "").
			Add("ip", "").
			Add("map", "").
			Add("eventchrname", "").
			Add("itemmakeindex", "").
			Add("desc", "").
			Add("event", 300).
			Add("starttime", "2023-01-05T00:00:00").
			Add("stoptime", "2023-01-05T18:33:06").
			//Add("itemname", "太阳水").
			Add("Numberof", 100))
	fmt.Println(len(query["data"].([]interface{})))
	fmt.Println(t.GetCost())
	node.Start()
}
