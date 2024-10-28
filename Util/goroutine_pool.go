package Util

import (
	"container/list"
	"runtime"
	"sync"
	"sync/atomic"
)

type task struct {
	fun  func(...interface{})
	args []interface{}
}

type Coroutine struct {
	chNewTask    chan *task
	chDoTask     chan *task
	chTaskFinish chan struct{}
	taskList     *list.List
	status       int32 // 状态: 0-空闲, 1-正在执行任务
	idxCount     int   // 有序执行的时候, 表示有多少标志的任务由他执行
}

func NewCoroutine() *Coroutine {
	ret := &Coroutine{
		chNewTask:    make(chan *task),
		chDoTask:     make(chan *task),
		chTaskFinish: make(chan struct{}),
		taskList:     list.New(),
	}
	go ret.dispatch()
	go ret.doTask()
	return ret
}

func (c *Coroutine) doTask() {
	for {
		select {
		case task := <-c.chDoTask:
			task.fun(task.args...)
			c.chTaskFinish <- struct{}{}
		}
	}
}

func (c *Coroutine) dispatch() bool {
	for {
		select {
		// 新任务需要执行
		case newTask := <-c.chNewTask:
			// 如果执行协程空闲，则直接通知执行协程执行任务；否则放入队列中。
			if atomic.LoadInt32(&c.status) == 0 {
				atomic.StoreInt32(&c.status, 1)
				c.chDoTask <- newTask
			} else {
				c.taskList.PushBack(newTask)
			}
		// 工作协程有空闲
		case <-c.chTaskFinish:
			// 如果缓存队列有任务，则取出任务，通知执行协程执行任务；否则设置状态为空闲
			if front := c.taskList.Front(); front != nil {
				c.taskList.Remove(front)
				c.chDoTask <- front.Value.(*task)
			} else {
				atomic.StoreInt32(&c.status, 0)
			}
		}
	}
}

type CoroutinePool struct {
	pool           []*Coroutine // 协程列表
	idxCoroutineId sync.Map     // map[idx]*Coroutine idx 对应的协程, 保证有序执行(RunOrder)是, 同一个idx由同一个协程执行, 保证执行顺序
	next           uint64       // 用于实现轮询的 atomic 自增索引
}

func NewCoroutinePool() *CoroutinePool {
	count := runtime.NumCPU()
	ret := &CoroutinePool{}
	for i := 0; i < count; i++ {
		ret.pool = append(ret.pool, NewCoroutine())
	}
	return ret
}

// RunOrder 有序执行任务, 根据idx确定任务队列，然后队列有序
/*
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
可以保证hello world A 比 hello world B 先执行; hello world C 比 hello world D 先执行。但是hello world A 和 hello world C 会随机执行
*/
func (c *CoroutinePool) RunOrder(idx int64, fun func(...interface{}), args ...interface{}) {
	c.getOrderPool(idx).chNewTask <- &task{fun: fun, args: args}
}

// Run 无序执行
/*
pool := NewCoroutinePool()
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
这里虽然顺序写的E,F,G,H, 但是实际执行顺序不一定, 因为4个任务会被分配到不同的协程,各协程间独立, 不保证顺序
*/
func (c *CoroutinePool) Run(fun func(...interface{}), args ...interface{}) {
	idx := atomic.AddUint64(&c.next, 1) % uint64(len(c.pool))
	c.pool[idx].chNewTask <- &task{fun: fun, args: args}
}

func (c *CoroutinePool) Empty() bool {
	totalStatus := int32(0)
	for _, v := range c.pool {
		totalStatus += atomic.LoadInt32(&v.status)
	}
	return totalStatus == 0
}

func (c *CoroutinePool) getOrderPool(idx int64) *Coroutine {
	if v, ok := c.idxCoroutineId.Load(idx); ok {
		return v.(*Coroutine)
	}
	// 新建协程并更新映射
	minCount := int(^uint(0) >> 1) // 最大整数
	var ret *Coroutine
	for _, v := range c.pool {
		if v.idxCount < minCount {
			minCount = v.idxCount
			ret = v
		}
		if minCount == 0 {
			break
		}
	}
	if ret != nil {
		ret.idxCount += 1
		c.idxCoroutineId.Store(idx, ret)
	}
	return ret
}
