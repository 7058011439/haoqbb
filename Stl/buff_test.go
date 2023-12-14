package Stl

import (
	"fmt"
	"github.com/7058011439/haoqbb/String"
	"math/rand"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

/*
这里主要做了性能测试，测试结果，封装后的 NewBuffer 比 原生 切片要高20%性能
*/

var testData = make([]string, 100)
var index = 0

func init() {
	for i := 0; i < 100; i++ {
		testData[i] = String.RandStr(i)
	}
}

// 这里一个奇怪的现象, 通过循环返回第一个元素的效率比下标获取要高
var strData = "abcdefghijklnmoqprstuvwxyzabcdefghijklnmoqprstuvwxyzabcdefghijklnmoqprstuvwxyz"

func getRandomData() string {
	return strData
	/*
		defer func(){
			index++
			if index >= 100 {
				index = 0
			}
		}()
		return testData[index]
	*/
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

func print(desc string, data []byte) {
	fmt.Println(desc+":", len(data), cap(data), (*reflect.SliceHeader)(unsafe.Pointer(&data)).Data, string(data))
}

func getPtr(data []byte) uintptr {
	return (*reflect.SliceHeader)(unsafe.Pointer(&data)).Data
}

func TestBuffer_Bytes(t *testing.T) {
	ids1 := []byte("abcdefg")
	print("ids1", ids1)
	copy(ids1[0:4], ids1[3:])
	print("ids1", ids1)

	/*
		buff := NewBuffer(2)
		buff.WriteString("ab")
		buff.WriteString("cd")
		buff.OffSize(3)
		/*
		buff := NewBuffer(6)
		buff.WriteString("abcde")
		buff.OffSize(3)
	*/
}

func TestBuffer_Write1(t *testing.T) {
	buff := NewBuffer(1024)
	timeStart := time.Now()
	oldPtr := getPtr(buff.cs)
	capTimes := 0
	for i := 0; i < 500000000; i++ {
		buff.WriteString(getRandomData())
		if oldPtr != getPtr(buff.cs) {
			capTimes++
			oldPtr = getPtr(buff.cs)
		}
		if rand.Intn(2) == 0 {
			buff.OffSize(rand.Intn(buff.Len()))
		}
	}
	fmt.Println(time.Now().Sub(timeStart).Milliseconds(), capTimes)
}

func TestBuffer_Write2(t *testing.T) {
	buff := make([]byte, 0, 1024)
	timeStart := time.Now()
	oldPtr := getPtr(buff)
	capTimes := 0
	for i := 0; i < 500000000; i++ {
		buff = append(buff, []byte(getRandomData())...)
		if oldPtr != getPtr(buff) {
			capTimes++
			oldPtr = getPtr(buff)
		}
		if rand.Intn(2) == 0 {
			buff = buff[rand.Intn(len(buff)):]
		}
	}
	fmt.Println(time.Now().Sub(timeStart).Milliseconds(), capTimes)
}
