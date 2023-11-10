package DataBase

import (
	"fmt"
	"github.com/7058011439/haoqbb/String"
	"github.com/7058011439/haoqbb/Timer"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

var mongoDB = NewMongoDB("127.0.0.1", 27017, "Test", "", "", 20)

type home struct {
	Name  string
	Level int
}

type bag struct {
	Data map[int]int
}

type mission struct {
	Id         int `bson:"id,omitempty"`
	CurrentPro int
	TotalPro   int
}

type missions struct {
	MissionList []*mission `bson:"mission_list,omitempty"`
	IntData     int
	StringData  string
}

type tagMongo struct {
	Id      uint64    `bson:"id,omitempty"`
	Name    string    `bson:"name,omitempty"`
	Addr    string    `bson:"addr,omitempty"`
	Home    *home     `bson:"home,omitempty"`
	Bag     *bag      `bson:"bag,omitempty"`
	Mission *missions `bson:"mission,omitempty"`
	Friends []uint64  `bson:"friends,omitempty"`
}

func newBag() *bag {
	data := map[int]int{}
	count := rand.Intn(21) + 30
	for i := 0; i < count; i++ {
		data[rand.Intn(100)] = rand.Intn(10)
	}
	return &bag{
		Data: data,
	}
}

func newHome() *home {
	return &home{
		Name:  String.RandStr(20),
		Level: rand.Intn(100),
	}
}

func newFriends() []uint64 {
	var ret []uint64
	count := rand.Intn(16) + 5
	for i := 0; i < count; i++ {
		ret = append(ret, uint64(rand.Int63()))
	}
	return ret
}

func newMissions() *missions {
	ret := &missions{
		IntData:    rand.Intn(1000),
		StringData: String.RandStr(10),
	}
	count := rand.Intn(6) + 10
	for i := 0; i < count; i++ {
		ret.MissionList = append(ret.MissionList, &mission{
			Id:         rand.Intn(10000),
			CurrentPro: rand.Intn(10),
			TotalPro:   rand.Intn(20),
		})
	}
	return ret
}

func newData() *tagMongo {
	return &tagMongo{
		Addr:    String.RandStr(20),
		Home:    newHome(),
		Bag:     newBag(),
		Friends: newFriends(),
		Mission: newMissions(),
	}
}

func newPartData() *tagMongo {
	data := &tagMongo{}
	switch rand.Intn(5) {
	case 0:
		data.Addr = String.RandStr(20)
	case 1:
		data.Home = newHome()
	case 2:
		data.Bag = newBag()
	case 3:
		data.Friends = newFriends()
	case 4:
		data.Mission = newMissions()
	}
	return data
}

func init() {
	rand.Seed(time.Now().UnixMilli())
}

var count = int32(0)
var times = uint64(1000000)
var ch = make(chan int, 1)
var isQueue = true

func TestMongoDB_InsertOne(t *testing.T) {
	cost := Timer.NewTiming(Timer.Millisecond)

	for i := uint64(1); i <= times; i++ {
		data := newData()
		data.Addr = fmt.Sprintf("昵称_%v", i)
		data.Id = i
		if isQueue {
			mongoDB.InsertOne("test", data, int(data.Id), func(callbackData ...interface{}) {
				atomic.AddInt32(&count, 1)
				if uint64(count) == times {
					ch <- 1
				}
			}, nil)
		} else {
			mongoDB.insertOne("test", data, nil, nil)
		}
	}
	if isQueue {
		<-ch
	}
	fmt.Println(cost)
}

func TestMongoDB_FindOne(t *testing.T) {
	cost := Timer.NewTiming(Timer.Millisecond)

	ret := &tagMongo{}
	for i := uint64(1); i <= times; i++ {
		condition := tagMongo{Id: i}
		if isQueue {
			mongoDB.FindOne("test", condition, ret, int(condition.Id), func(getData interface{}, callbackData ...interface{}) {
				atomic.AddInt32(&count, 1)
				if uint64(count) == times {
					ch <- 1
				}
			}, nil)
		} else {
			mongoDB.findOne("test", condition, ret, nil, nil)
		}
	}
	fmt.Println(ret)
	if isQueue {
		<-ch
	}
	fmt.Println(cost)
}

func TestMongoDB_UpdateOne(t *testing.T) {
	cost := Timer.NewTiming(Timer.Millisecond)

	for i := uint64(1); i <= times; i++ {
		condition := &tagMongo{Id: i}
		if isQueue {
			mongoDB.UpdateOne("test", condition, *newData(), int(condition.Id), func(callbackData ...interface{}) {
				atomic.AddInt32(&count, 1)
				if uint64(count) == times {
					ch <- 1
				}
			}, nil)
		} else {
			mongoDB.updateOne("test", condition, newPartData(), nil, nil)
		}
	}
	if isQueue {
		<-ch
	}
	fmt.Println(cost)
}
