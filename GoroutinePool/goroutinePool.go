package GoroutinePool

import (
	"container/list"
	"sync/atomic"
)

type task struct {
	fun  func(...interface{})
	args []interface{}
}

type GoRoutinePool struct {
	chNewTask          chan *task    // 新任务
	chWorkGoroutine    chan *task    // 通知工作协程有任务
	chIdleGoroutine    chan struct{} // 通知有空闲协程
	workGoroutineCount int32         // 开辟协程数量
	idleGoroutineCount int32         // 空闲协程数量
	taskList           *list.List    // 任务缓冲列表
}

func NewPool(count int) *GoRoutinePool {
	ret := &GoRoutinePool{
		chNewTask:          make(chan *task, count),
		chWorkGoroutine:    make(chan *task, count),
		chIdleGoroutine:    make(chan struct{}, count),
		workGoroutineCount: int32(count),
		idleGoroutineCount: int32(count),
		taskList:           list.New(),
	}
	for i := 0; i < count; i++ {
		go ret.doTask()
	}
	go ret.run()
	return ret
}

func (g *GoRoutinePool) Run(fun func(...interface{}), args ...interface{}) {
	g.chNewTask <- &task{fun: fun, args: args}
}

// Empty 协程池是否完全空闲
func (g *GoRoutinePool) Empty() bool {
	return atomic.LoadInt32(&g.idleGoroutineCount) == g.workGoroutineCount
}

// HaveIdle 是否空闲协程池
func (g *GoRoutinePool) HaveIdle() bool {
	return atomic.LoadInt32(&g.idleGoroutineCount) > 0
}

func (g *GoRoutinePool) run(args ...interface{}) {
	for {
		select {
		// 新任务需要执行
		case newTask := <-g.chNewTask:
			// 有空闲协程，将任务放入工作协程池中(协程池自己抢占)，空闲西城数量-1；否则 将任务放入缓冲队列中
			if g.HaveIdle() {
				atomic.AddInt32(&g.idleGoroutineCount, -1)
				g.chWorkGoroutine <- newTask
			} else {
				g.taskList.PushBack(newTask)
			}
		// 工作协程有空闲
		case <-g.chIdleGoroutine:
			// 如果缓存队列有任务，则取出任务，放入工作协程池中(协程池自己抢占)；否则 空闲协程数量+1
			front := g.taskList.Front()
			if front != nil {
				g.taskList.Remove(front)
				g.chWorkGoroutine <- front.Value.(*task)
			} else {
				atomic.AddInt32(&g.idleGoroutineCount, 1)
			}
		}
	}
}

func (g *GoRoutinePool) doTask() {
	for {
		select {
		case task := <-g.chWorkGoroutine:
			task.fun(task.args...)
			g.chIdleGoroutine <- struct{}{}
		}
	}
}
