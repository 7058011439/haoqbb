package Stl

import (
	"sync"
	"testing"
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
	sign.Wait()
}
