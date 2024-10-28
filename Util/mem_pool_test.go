package Util

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func print(desc string, data []byte) {
	fmt.Println(desc+":", len(data), cap(data), (*reflect.SliceHeader)(unsafe.Pointer(&data)).Data, string(data))
}

func sliceString(data []byte) string {
	return fmt.Sprintf("%v, %v, %v, %v", len(data), cap(data), (*reflect.SliceHeader)(unsafe.Pointer(&data)).Data, string(data))
}

func TestGetSlice(t *testing.T) {
	dataLen := 2
	var data = make([][]byte, 0, dataLen*2)
	for i := 0; i < dataLen; i++ {
		data = append(data, GetNewSlice(5))
	}
	for i := 0; i < dataLen; i++ {
		data = append(data, GetNewSlice(6))
	}
	for i, d := range data {
		t.Logf("第一轮获取[%2d]次数据: %v", i+i, sliceString(d))
		GiveBackSlice(d)
	}
	for i := 0; i < dataLen*2; i++ {
		d := GetNewSlice(8)
		t.Logf("第二轮获取[%2d]次数据: %v", i+i, sliceString(d))
	}
}

var size = 1024 * 64
var buff []byte

func BenchmarkSliceNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buff = make([]byte, 0, i%size+1)
	}
}

func BenchmarkSlicePool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buff = GetNewSlice(i%size + 1)
		GiveBackSlice(buff)
	}
}

// 结构中有大内存结构
type player struct {
	Name   string
	Age    int
	Money  int64
	Addr   string
	Remark [1024]byte
}

func (p *player) Reset() {
	p.Name = ""
	p.Age = 0
	p.Age = 0
	p.Addr = ""
	p.Remark = [1024]byte{}
}

func (p *player) Key() string {
	return "player"
}

func (p *player) New() interface{} {
	return &player{}
}

// 结构中小数据
type school struct {
	Name string
	Addr []byte
}

func (s *school) Reset() {
	s.Name = ""
	s.Addr = s.Addr[:0]
}

func (s *school) Key() string {
	return "school"
}

func (s *school) New() interface{} {
	return &school{}
}

var p *player

func BenchmarkObjectPlayerNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p = &player{}
	}
}

func BenchmarkObjectPlayerPool(b *testing.B) {
	temp := &player{}
	for i := 0; i < b.N; i++ {
		p = GetNewObj(temp).(*player)
		GiveBackObj(p)
	}
}

var s *school

func BenchmarkObjectSchoolNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s = &school{}
	}
}

func BenchmarkObjectSchoolPool(b *testing.B) {
	temp := &school{}
	for i := 0; i < b.N; i++ {
		s = GetNewObj(temp).(*school)
		GiveBackObj(s)
	}
}
