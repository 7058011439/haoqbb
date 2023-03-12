package Stl

import (
	"math/rand"
	"testing"
)

func TestNewSet(t *testing.T) {
	set := NewSet()
	count := 20
	for i := 0; i < count; i++ {
		set.Add(i)
	}
	if set.Len() != count {
		t.Errorf("len错误, 希望值 = %v, 实际值 = %v", count, set.Len())
	}
	for i := 0; i < count; i++ {
		if !set.Exist(i) {
			t.Errorf("exist错误, value = %v", i)
		}
	}

	set.Add(rand.Intn(count))
	set.Add(rand.Intn(count))
	set.Add(rand.Intn(count))
	set.Add(rand.Intn(count))
	if set.Len() != count {
		t.Errorf("len错误, 希望值 = %v, 实际值 = %v", count, set.Len())
	}

	randomData := rand.Intn(count)
	if set.Del(randomData) == false {
		t.Errorf("del错误, 值不存在, 值 = %v", randomData)
	}
	if set.Exist(randomData) {
		t.Errorf("del错误, 删除无效， 值 = %v", randomData)
	}

	if set.Len() != count-1 {
		t.Errorf("len错误, 希望值 = %v, 实际值 = %v", count, set.Len())
	}
}
