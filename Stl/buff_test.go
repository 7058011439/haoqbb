package Stl

import (
	"fmt"
	"github.com/7058011439/haoqbb/String"
	"math/rand"
	"testing"
	"time"
)

/*
这里主要做了性能测试，测试结果，封装后的 NewBuffer 比 原生 切片要高20%性能
*/

var testData = map[int]string{}

func init() {
	for i := 10; i < 100; i++ {
		testData[i] = String.RandStr(i)
	}
}

// 这里一个奇怪的现象, 通过循环返回第一个元素的效率比下标获取要高
func getRandomData() string {
	//return testData[rand.Intn(90) + 10]
	for _, v := range testData {
		return v
	}
	return ""
}

func BenchmarkBuffer_Write1(b *testing.B) {
	buff := NewBuffer(1024)
	for i := 0; i < b.N; i++ {
		buff.WriteString(getRandomData())
		if rand.Intn(2)%5 == 0 {
			buff.OffSize(rand.Intn(buff.Len()))
		}
	}
}

func BenchmarkBuffer_Write2(b *testing.B) {
	buff := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		buff = append(buff, []byte(getRandomData())...)
		if rand.Intn(2)%5 == 0 {
			buff = buff[rand.Intn(len(buff)):]
		}
	}
}

func TestBuffer_Write1(t *testing.T) {
	buff := NewBuffer(1024)
	timeStart := time.Now()
	oldCap := buff.Cap()
	capTimes := 0
	for i := 0; i < 500000000; i++ {
		buff.WriteString(getRandomData())
		if oldCap != buff.Cap() {
			capTimes++
			oldCap = buff.Cap()
		}
		if rand.Intn(2)%5 == 0 {
			buff.OffSize(rand.Intn(buff.Len()))
		}
	}
	fmt.Println(time.Now().Sub(timeStart).Milliseconds(), capTimes)
}

func TestBuffer_Write2(t *testing.T) {
	buff := make([]byte, 1024)
	timeStart := time.Now()
	oldCap := cap(buff)
	capTimes := 0
	for i := 0; i < 500000000; i++ {
		buff = append(buff, []byte(getRandomData())...)
		if oldCap != cap(buff) {
			capTimes++
			oldCap = cap(buff)
		}
		if rand.Intn(2)%5 == 0 {
			buff = buff[rand.Intn(len(buff)):]
		}
	}
	fmt.Println(time.Now().Sub(timeStart).Milliseconds(), capTimes)
}
