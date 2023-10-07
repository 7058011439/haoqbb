package Math

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
)

var n, m = 130, 13

//[C(13, 4)]^4 / [C(52, 13) * C(39, 13) * C(26, 13) * C(13, 13)]。
func TestCombination(t *testing.T) {
	/*
		fmt.Println(Combination(52, 13))
		fmt.Println(math.Pow(float64(Combination(13, 4)), 4))
		fmt.Println(float64(Combination(52, 13)*Combination(39, 13)*Combination(26, 13)))
		fmt.Println(math.Pow(float64(Combination(13, 4)), 4) / float64(Combination(52, 13)*Combination(39, 13)*Combination(26, 13)))
	*/
	for n := 1; n < 200; n++ {
		m := rand.Intn(n)
		if Combination(n, m) != CombinationBig(n, m).Int64() {
			log.Fatalf("测试结果错误, n = %v, m = %v", n, m)
		}
	}
}

func BenchmarkCombination(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Combination(n, m)
	}
	fmt.Println(Combination(n, m))
}

func BenchmarkCombinationB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CombinationBig(n, m)
	}
	fmt.Println(CombinationBig(n, m))
}
