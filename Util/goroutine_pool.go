package Util

import (
	"container/list"
	"math"
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
	idle         bool         // 当前调度链是否空闲 true-空闲 false-已有任务在执行或排队
	idxCount     int          // 有序执行的时候, 表示有多少标志的任务由他执行
	pending      atomic.Int64 // 当前待处理任务数(执行中 + 排队中)
}

func newCoroutine() *Coroutine {
	ret := &Coroutine{
		chNewTask:    make(chan *task),
		chDoTask:     make(chan *task),
		chTaskFinish: make(chan struct{}),
		taskList:     list.New(),
		idle:         true,
	}
	go ret.dispatch()
	go ret.doTask()
	return ret
}

func (c *Coroutine) addTask(fun func(...interface{}), args ...interface{}) {
	c.pending.Add(1)
	c.chNewTask <- &task{fun: fun, args: args}
}

func (c *Coroutine) doTask() {
	for {
		task := <-c.chDoTask
		task.fun(task.args...)
		c.pending.Add(-1)
		c.chTaskFinish <- struct{}{}
	}
}

func (c *Coroutine) dispatch() bool {
	for {
		select {
		// 新任务需要执行
		case newTask := <-c.chNewTask:
			// 如果执行协程空闲，则直接通知执行协程执行任务；否则放入队列中。
			if c.idle {
				c.idle = false
				c.chDoTask <- newTask
			} else {
				c.taskList.PushBack(newTask)
			}
		// 工作协程有空闲
		case <-c.chTaskFinish:
			// 如果缓存队列有任务，则取出任务，通知执行协程执行任务；否则说明当前协程已空闲
			if front := c.taskList.Front(); front != nil {
				c.taskList.Remove(front)
				c.chDoTask <- front.Value.(*task)
			} else {
				c.idle = true
			}
		}
	}
}

type CoroutinePool struct {
	pool           []*Coroutine  // 协程列表
	idxCoroutineId sync.Map      // map[idx]*Coroutine idx 对应的协程, 保证有序执行(RunOrder)时, 同一个idx由同一个协程执行, 保证执行顺序
	next           atomic.Uint64 // 用于实现轮询的 atomic 自增索引
	mutex          sync.Mutex    // 保护 getOrderPool 中 idxCoroutineId / idxCount 的并发访问
}

func NewCoroutinePool() *CoroutinePool {
	count := runtime.NumCPU()
	ret := &CoroutinePool{}
	for i := 0; i < count; i++ {
		ret.pool = append(ret.pool, newCoroutine())
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
	c.getOrderPool(idx).addTask(fun, args...)
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
	pool := c.pool[c.next.Add(1)%uint64(len(c.pool))]
	if pool.pending.Load() >= 5 {
		c.getBestPool().addTask(fun, args...)
	} else {
		pool.addTask(fun, args...)
	}
}

func (c *CoroutinePool) getOrderPool(idx int64) *Coroutine {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 锁内再检查一次
	if v, ok := c.idxCoroutineId.Load(idx); ok {
		return v.(*Coroutine)
	}

	minCount := math.MaxInt
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
		ret.idxCount++
		c.idxCoroutineId.Store(idx, ret)
	}
	return ret
}

func (c *CoroutinePool) getBestPool() (ret *Coroutine) {
	minTaskCount := int64(math.MaxInt64)
	for _, p := range c.pool {
		if taskCount := p.pending.Load(); taskCount < minTaskCount {
			minTaskCount = taskCount
			ret = p
		}
	}
	return ret
}
