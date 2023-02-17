package main

import (
	"Core/Http"
	"Core/Timer"
	"Core/haoqbb/node"
	"Core/haoqbb/server/common"
	_ "Core/haoqbb/server/dispatcher"
	"Core/haoqbb/server/gameSrv/bujidao/client"
	"Core/haoqbb/server/gameSrv/bujidao/server"
	_ "Core/haoqbb/server/gateWay"
	_ "Core/haoqbb/server/loginSrv"
	"encoding/json"
	"fmt"
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
	data := string(Http.GetHttpSyncNew("http://192.168.13.220:8086/", Http.NewHead(nil)))
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
	query := Http.PostHttpSyncNew("http://192.168.13.220:8086/API/CreateQuery.html",
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
