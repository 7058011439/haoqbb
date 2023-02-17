package Timer

import (
	"Core/Log"
	"fmt"
	"time"
)

type timingType = int

const (
	Second      timingType = 0
	Millisecond timingType = 1
	Microsecond timingType = 2
	Nanosecond  timingType = 3
)

type Timing struct {
	startTime time.Time
	eType     timingType
}

var mapDesc = map[timingType]string{
	Second:      "s",
	Millisecond: "ms",
	Microsecond: "us",
	Nanosecond:  "ns",
}

func NewTiming(eType timingType) *Timing {
	ret := Timing{
		eType:     eType,
		startTime: time.Now(),
	}

	return &ret
}

func (t *Timing) getDesc() string {
	return mapDesc[t.eType]
}

func (t *Timing) ReStart() float64 {
	ret := t.GetCost()
	t.startTime = time.Now()
	return ret
}

func (t *Timing) GetCost() float64 {
	gaps := time.Now().Sub(t.startTime)
	switch t.eType {
	case Second:
		return gaps.Seconds()
	case Millisecond:
		return float64(gaps.Milliseconds())
	case Microsecond:
		return float64(gaps.Microseconds())
	case Nanosecond:
		return float64(gaps.Nanoseconds())
	}
	return 0
}

func (t *Timing) String() string {
	return fmt.Sprintf("%v %v", t.GetCost(), mapDesc[t.eType])
}

func (t *Timing) PrintCost(condition float64, restart bool, format string, args ...interface{}) {
	if t.GetCost() >= condition {
		title := fmt.Sprintf(format, args...)
		Log.ErrorLog("%v : cost = %v", title, t.GetCost())
	}
	if restart {
		t.ReStart()
	}
}
