package Util

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"testing"
)

/*
测试效果本测试环境为8核16线程CPU:
一、测试calcSum(密集计算型), 协程池优势明显
1.通过协程池方式耗时33.61s, 运算结果=357913941332992000
2.通过普通方式耗时265.51s, 运算结果=357913941332992000

二、测试json.UmMarshal(耗时操作型), 协程池有一定优势
1.通过协程池方式耗时5.62s
2.通过普通方式耗时17.58s

三、测试json.Marshal(简单操作), 协程池毫无优势
1.通过协程池方式耗时4.71s
2.通过普通方式耗时3.99s
*/
func calcSum(value int) (ret int64) {
	for index := 1; index <= value; index++ {
		ret += int64(index)
	}
	return ret
	// return int64(value) * int64(value)
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
	pool.Wait()
	t.Logf("total = %v", total)
	if ok, index := isSortSlice(data); !ok {
		t.Errorf("非顺序执行 index = %v, data = %v", index, data)
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

type testPlayer struct {
	Name   string
	Age    int
	Score  int64
	Money  int64
	Addr   string
	School []byte
}

func initTestPlayer() (*testPlayer, []byte) {
	p := &testPlayer{
		Name:   "hello",
		Age:    18,
		Score:  100,
		Money:  200,
		Addr:   "127.0.0.1",
		School: []byte("bei jing da xue & qing hua da da xue & zhe jiang da xue & fu dan da xue"),
	}
	data, _ := json.Marshal(p)
	return p, data
}

func TestUnMarshalWithPool(t *testing.T) {
	times := 10240000
	pool := NewCoroutinePool()
	p, d := initTestPlayer()
	for i := 0; i < times; i++ {
		pool.Run(func(i ...interface{}) {
			if err := json.Unmarshal(i[0].([]byte), p); err != nil {
				t.Errorf("反序列化错误, err = %v", err)
			}
		}, d)
	}
	pool.Wait()
}

func TestUnMarshalWithoutPool(t *testing.T) {
	times := 10240000
	p, d := initTestPlayer()
	for i := 0; i < times; i++ {
		if err := json.Unmarshal(d, p); err != nil {
			t.Errorf("反序列化错误, err = %v", err)
		}
	}
}

func TestMarshalWithPool(t *testing.T) {
	times := 10240000
	pool := NewCoroutinePool()
	p, _ := initTestPlayer()
	for i := 0; i < times; i++ {
		pool.Run(func(i ...interface{}) {
			if _, err := json.Marshal(i[0].(*testPlayer)); err != nil {
				t.Errorf("反序列化错误, err = %v", err)
			}
		}, p)
	}
	pool.Wait()
}

func TestMarshalWithoutPool(t *testing.T) {
	times := 10240000
	p, _ := initTestPlayer()
	for i := 0; i < times; i++ {
		if _, err := json.Marshal(p); err != nil {
			t.Errorf("序列化错误, err = %v", err)
		}
	}
}

func isSortSlice(arr []int64) (bool, int) {
	for i := 0; i < len(arr)-2; i++ {
		if arr[i] > arr[i+1] {
			return false, i
		}
	}
	return true, -1
}

func TestNewCoroutinePool(t *testing.T) {
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
	pool.Wait()
}
