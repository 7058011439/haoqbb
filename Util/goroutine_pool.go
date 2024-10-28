package Util

import (
	"container/list"
	"math/rand"
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
			// 如果执行协程空闲，则直接通知执行协程执行任务，否则放入队列中。
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
	pool           []*Coroutine
	idxCoroutineId map[int64]*Coroutine // idx 对应的队列, 保证有序执行(RunOrder)是, 同一个idx由同一个协程执行, 保证执行顺序
	sync.Mutex
}

func NewCoroutinePool() *CoroutinePool {
	count := runtime.NumCPU()
	ret := &CoroutinePool{}
	for i := 0; i < count; i++ {
		ret.pool = append(ret.pool, NewCoroutine())
	}
	ret.idxCoroutineId = make(map[int64]*Coroutine, 1024)
	return ret
}

// RunOrder 有序执行任务, 根据idx确定任务队列，然后队列有序
func (c *CoroutinePool) RunOrder(idx int64, fun func(...interface{}), args ...interface{}) {
	c.getOrderPool(idx).chNewTask <- &task{fun: fun, args: args}
}

// Run 有序执行任务, 根据idx确定任务队列，然后队列有序
func (c *CoroutinePool) Run(fun func(...interface{}), args ...interface{}) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.pool[rand.Intn(len(c.pool))].chNewTask <- &task{fun: fun, args: args}
}

func (c *CoroutinePool) Empty() bool {
	totalStatus := int32(0)
	for _, v := range c.pool {
		totalStatus += atomic.LoadInt32(&v.status)
	}
	return totalStatus == 0
}

func (c *CoroutinePool) getOrderPool(idx int64) *Coroutine {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if data, ok := c.idxCoroutineId[idx]; ok {
		return data
	} else {
		minCount := 9999999
		var ret *Coroutine
		for _, v := range c.pool {
			if v.idxCount == 0 {
				ret = v
				break
			}

			if v.idxCount < minCount {
				minCount = v.idxCount
				ret = v
			}
		}
		if ret != nil {
			ret.idxCount += 1
			c.idxCoroutineId[idx] = ret
		}
		return ret
	}
}
