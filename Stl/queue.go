package Stl

import (
	"container/list"
	"sync"
)

func NewQueue() *Queue {
	return &Queue{
		listData: list.New(),
	}
}

type Queue struct {
	listData *list.List
	mutex    sync.RWMutex
}

func (q *Queue) Enqueue(data interface{}) {
	defer q.mutex.Unlock()
	q.mutex.Lock()
	q.listData.PushBack(data)
}

func (q *Queue) Dequeue() interface{} {
	defer q.mutex.Unlock()
	q.mutex.Lock()
	ret := q.listData.Front()
	return q.listData.Remove(ret)
}

func (q *Queue) Head() interface{} {
	defer q.mutex.RUnlock()
	q.mutex.RLock()
	if ret := q.listData.Front(); ret == nil {
		return nil
	} else {
		return ret.Value
	}
}

func (q *Queue) Len() int {
	defer q.mutex.RUnlock()
	q.mutex.RLock()
	return q.listData.Len()
}
