package DataBase

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/String"
	"math/rand"
	"testing"
	"time"
)

type OrderInfo struct {
	SDataA string  `json:"SA"`
	SDataB string  `json:"SB"`
	SDataC string  `json:"SC"`
	FDataA float64 `json:"FA"`
	FDataB float64 `json:"FB"`
	FDataC float64 `json:"BC"`
}

func TestRedis(t *testing.T) {
	var redisLocal = NewRedisDB("1.2.3.4", 56, "foobared10086", 3)
	for i := 0; i < 10; i++ {
		if err := redisLocal.Set(String.RandStr(44), &OrderInfo{
			SDataA: "dataA" + String.RandStr(5),
			SDataB: "dataB" + String.RandStr(5),
			SDataC: "dataC" + String.RandStr(5),
			FDataA: rand.Float64(),
			FDataB: rand.Float64(),
			FDataC: rand.Float64(),
		}, time.Second*300); err != nil {
			Log.Error("写入表错误, err = %v", err)
		}
	}
}
