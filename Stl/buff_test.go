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

var sliceData = []byte("abcdefghijklnmoqprstuvwxyzabcdefghijklnmoqprstuvwxyzabcdefghijklnmoqprstuvwxyz")

func BenchmarkBuffer_Write1(b *testing.B) {
	buff := NewBuffer(128)
	for i := 0; i < b.N; i++ {
		buff.Write(sliceData)
		if i%5 == 0 {
			buff.OffSize(rand.Intn(buff.Len()))
		}
	}
}

func BenchmarkBuffer_Write2(b *testing.B) {
	buff := make([]byte, 0, 128)
	for i := 0; i < b.N; i++ {
		buff = append(buff, sliceData...)
		if i%5 == 0 {
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

func parse(data []byte) ([]byte, int) {
	os := rand.Intn(len(data))
	return data[0:os], os
}

func TestBuffer_Bytes(t *testing.T) {
	/*
		ids1 := []byte("abcdefg")
		print("ids1", ids1)
		copy(ids1[0:4], ids1[3:])
		print("ids1", ids1)
	*/

	buff := NewBuffer(11)
	print("初始值", buff.cs)
	buff.WriteString("ab")
	print("插入 ab 后", buff.cs)
	buff.WriteString("cde")
	print("插入 cde 后", buff.cs)
	buff.WriteString("fghijk")
	print("插入 fghijk 后", buff.cs)
	buff.OffSize(5)
	print("偏移 5 后", buff.cs)
	buff.Reset()
	print("重置后", buff.cs)
	buff.WriteString("abcdefghij")
	print("插入 abcdefghij 后", buff.cs)
}

func TestBuffer_Write1(t *testing.T) {
	buff := NewBuffer(1024)
	timeStart := time.Now()
	oldPtr := getPtr(buff.cs)
	capTimes := 0
	for i := 0; i < 500000000; i++ {
		buff.Write(sliceData)
		if oldPtr != getPtr(buff.cs) {
			capTimes++
			oldPtr = getPtr(buff.cs)
		}
		if i%2 == 0 {
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
		buff = append(buff, sliceData...)
		if oldPtr != getPtr(buff) {
			capTimes++
			oldPtr = getPtr(buff)
		}
		if i%2 == 0 {
			buff = buff[rand.Intn(len(buff)):]
		}
	}
	fmt.Println(time.Now().Sub(timeStart).Milliseconds(), capTimes)
}

func TestSyncPool(t *testing.T) {
	fmt.Println(FormatNumber(3))    // 输出 4
	fmt.Println(FormatNumber(5))    // 输出 8
	fmt.Println(FormatNumber(129))  // 输出 256
	fmt.Println(FormatNumber(255))  // 输出 256
	fmt.Println(FormatNumber(513))  // 输出 1024
	fmt.Println(FormatNumber(1023)) // 输出 1024
}

func FormatNumber(value int) (ret int) {
	if value < 1 {
		return 1
	}

	// 从 1 开始，找出大于等于 value 的最近 2 的幂次值
	ret = 1
	for ret < value {
		ret <<= 1 // 左移一位，相当于乘以 2
	}

	return ret
}
