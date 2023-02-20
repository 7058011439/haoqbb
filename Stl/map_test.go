package Stl

import "testing"

func TestDoubleMap(t *testing.T) {
	data := NewDoubleMap()
	count := 10000
	for i := 0; i < count; i++ {
		data.Add(i, i*2)
	}

	for i := 0; i < count; i++ {
		if value := data.GetValue(i); value == nil {
			t.Errorf("获取value数据失败, key = %v", i)
		}
	}

	for i := 0; i < count; i++ {
		if value := data.GetKey(i * 2); value == nil {
			t.Errorf("获取key数据失败, key = %v", i)
		}
	}
}
