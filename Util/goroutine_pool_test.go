package Util

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

/*
测试效果:
1.通过协程池方式运算耗时33.61s, 运算结果=357913941332992000
2.通过普通方式运算耗时265.51s, 运算结果=357913941332992000
3.本测试环境为8核16线程CPU, 采用协程池可以将cpu利用率拉满
*/
func calcSum(value int) (ret int64) {
	for index := 1; index <= value; index++ {
		ret += int64(index)
	}
	return ret
}

func TestCalcSumWithPool(t *testing.T) {
	times := 1024000
	coroutineCount := 16
	pool := NewCoroutinePool()
	total := int64(0)
	data := make([]int64, 0, times/coroutineCount)
	for i := 0; i < times; i++ {
		// 有序执行, 通过取余算法, 确保平均间隔N个单位是同一协程执行, 然后通过isSortSlice 验证是否顺序执行
		pool.RunOrder(int64(i%coroutineCount), func(i ...interface{}) {
			value := i[0].(int)
			ret := calcSum(value)
			atomic.AddInt64(&total, ret)
			if value%coroutineCount == 0 {
				data = append(data, ret)
			}
		}, i)

		// 随机(协程)执行,不保证顺序
		pool.Run(func(i ...interface{}) {
			value := i[0].(int)
			ret := calcSum(value)
			atomic.AddInt64(&total, ret)
		}, i)
	}
	tStart := time.Now()
	for {
		if pool.Empty() {
			t.Logf("total = %v, wait cost = %v ms", total, time.Since(tStart).Milliseconds())
			break
		}
		time.Sleep(time.Millisecond)
	}
	if !isSortSlice(data) {
		t.Errorf("非顺序执行 data = %v", data)
	}
}

func TestCalcSumWithoutPool(t *testing.T) {
	total := int64(0)
	for i := 0; i < 1024000; i++ {
		total += calcSum(i)
		total += calcSum(i)
	}

	t.Logf("total = %v", total)
}

func isSortSlice(arr []int64) bool {
	for i := 0; i < len(arr)-2; i++ {
		if arr[i] > arr[i+1] {
			return false
		}
	}
	return true
}

func ExampleNewCoroutine() {
	pool := NewCoroutinePool()
	pool.RunOrder(10000, func(i ...interface{}) {
		fmt.Println("hello world A")
	})
	pool.RunOrder(10000, func(i ...interface{}) {
		fmt.Println("hello world B")
	})
	pool.RunOrder(10086, func(i ...interface{}) {
		fmt.Println("hello world C")
	})
	pool.RunOrder(10086, func(i ...interface{}) {
		fmt.Println("hello world D")
	})
	// 可以保证hello world A 比 hello world B 先执行; hello world C 比 hello world D 先执行。但是hello world A 和 hello world C 会随机执行

	pool.Run(func(i ...interface{}) {
		fmt.Println("hello world E")
	})
	pool.Run(func(i ...interface{}) {
		fmt.Println("hello world F")
	})
	pool.Run(func(i ...interface{}) {
		fmt.Println("hello world G")
	})
	pool.Run(func(i ...interface{}) {
		fmt.Println("hello world H")
	})
	// 这里虽然顺序写的E,F,G,H, 但是实际执行顺序不一定, 因为4个任务会被分配到不同的协程,各协程间独立, 不保证顺序
}
