package Stl

import (
	"fmt"
	"github.com/7058011439/haoqbb/String"
	"sync"
	"testing"
	"time"
)

func TestDoubleMap(t *testing.T) {
	data := NewDoubleMap()
	count := 100000
	for i := 0; i < count; i++ {
		data.Add(i, i*2)
	}

	var sign sync.WaitGroup
	sign.Add(3)

	go func() {
		for i := count; i < count*2; i++ {
			data.Add(i, i*2)
		}
		sign.Done()
	}()

	go func() {
		for i := 0; i < count; i++ {
			if value := data.GetValue(i); value == nil {
				t.Errorf("获取value数据失败, key = %v", i)
			}
		}
		sign.Done()
	}()

	go func() {
		for i := 0; i < count; i++ {
			if value := data.GetKey(i * 2); value == nil {
				t.Errorf("获取key数据失败, key = %v", i)
			}
		}
		sign.Done()
	}()
	// 这是一句很操蛋的话，永远不要管他
	sign.Wait()
}

var mutex sync.RWMutex
var testData1 = map[int64]string{}
var testData2 sync.Map
var testTimes = int64(10000000)

func writeData1(key int64, value string) {
	mutex.Lock()
	defer mutex.Unlock()
	testData1[key] = value
}

func writeData2(key int64, value string) {
	testData2.Store(key, strData)
}

func accessData1(key int64) string {
	mutex.RLock()
	defer mutex.RUnlock()
	if data, ok := testData1[key]; ok {
		return data
	} else {
		return ""
	}
}

func accessData2(key int64) string {
	if data, _ := testData2.Load(key); data != nil {
		return data.(string)
	}
	return ""
}

func TestMapAccess1(t *testing.T) {
	data := String.RandStr(20)
	var wait sync.WaitGroup
	startTime := time.Now()
	wait.Add(2)
	go func() {
		for i := int64(0); i < testTimes; i++ {
			writeData1(i, data)
		}
		wait.Done()
	}()
	for j := 0; j < 1; j++ {
		go func() {
			for i := int64(0); i < testTimes*10; i++ {
				accessData1(i)
			}
			wait.Done()
		}()
	}
	wait.Wait()
	fmt.Println(time.Now().Sub(startTime).Seconds())

	startTime = time.Now()
	wait.Add(2)
	go func() {
		for i := int64(0); i < testTimes; i++ {
			writeData2(i, data)
		}
		wait.Done()
	}()
	for j := 0; j < 1; j++ {
		go func() {
			for i := int64(0); i < testTimes*10; i++ {
				accessData2(i)
			}
			wait.Done()
		}()
	}
	wait.Wait()
	fmt.Println(time.Now().Sub(startTime).Seconds())
}
