package Authorize

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/System"
	"io/ioutil"
	"testing"
	"time"
)

func TestNewKey(t *testing.T) {
	var key []byte
	for i := 0; i < 1; i++ {
		key = NewKey(System.GetMachineId(), "Blfx", time.Now().Add(time.Second*100))
		fmt.Println(key)
	}
	ioutil.WriteFile("F:\\Go\\Core\\Authorize\\Haoqbb.key", key, 0x644)
}

func TestRegedit(t *testing.T) {
	Regedit("Blfx", func(dateTime *time.Time, err error) {
		if err != nil {
			Log.ErrorLog(err.Error())
		} else {
			Log.Log("授权截至日期:%v", *dateTime)
		}
	})
	time.Sleep(time.Hour)
}
