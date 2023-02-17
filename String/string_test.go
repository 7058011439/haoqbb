package String

import (
	"testing"
)

func TestRandStr(t *testing.T) {
	times := 10000000
	strLen := 10
	ret := make(map[string]bool, times)
	for i := 0; i < times; i++ {
		ret[RandStr(strLen)] = true
	}
	if len(ret) != times {
		t.Logf("出现重复字符串, 概率 = %v%%", float64(times-len(ret))/float64(times)*100)
	}
}

func BenchmarkRandStr10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStr(10)
	}
}

func BenchmarkRandStr20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStr(20)
	}
}
