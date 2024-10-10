package Timer

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"sync"
	"sync/atomic"
)

type CountAddUp struct {
	*Timing
	name       string
	currCount  int64
	totalCount int64
	mutex      sync.Mutex
}

func NewCountAddUp(eType timingType, name string) *CountAddUp {
	return &CountAddUp{
		Timing: NewTiming(eType),
		name:   name,
	}
}

func (c *CountAddUp) Update() {
	atomic.AddInt64(&c.currCount, 1)
	atomic.AddInt64(&c.totalCount, 1)
}

func (c *CountAddUp) ReStart() {
	atomic.StoreInt64(&c.currCount, 0)
	c.Timing.ReStart()
}

func (c *CountAddUp) PrintCost(condition float64, format string, args ...interface{}) {
	atomic.AddInt64(&c.currCount, 1)
	atomic.AddInt64(&c.totalCount, 1)
	mutex.Lock()
	if c.GetCost() >= condition {
		title := fmt.Sprintf(format, args...)
		Log.Error("%v : %v cost = %v, totalCount = %v, currCount = %v, currRate = %0.3f", c.name, title, c.GetCost(), c.totalCount, c.currCount, float64(c.currCount)/c.GetCost())
		c.Timing.ReStart()
		atomic.StoreInt64(&c.currCount, 0)
	}
	mutex.Unlock()
}

func (c *CountAddUp) TotalCount() int64 {
	return atomic.LoadInt64(&c.totalCount)
}

func (c *CountAddUp) CurrCount() int64 {
	return atomic.LoadInt64(&c.currCount)
}

func (c *CountAddUp) String() string {
	return fmt.Sprintf("%s:[%v %v %v %v]", c.name, c.GetCost(), c.totalCount, c.currCount, float64(c.currCount)/c.GetCost())
}
