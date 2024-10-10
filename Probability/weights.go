package Probability

import (
	"github.com/7058011439/haoqbb/Log"
	"math/rand"
)

type intervalValue struct {
	startWeight int
	endWeight   int
	value       int
}

type weights struct {
	maxWeight int
	allValue  []*intervalValue
}

type Weights struct {
	allWeight map[int]*weights
}

func NewWeights() *Weights {
	return &Weights{
		allWeight: make(map[int]*weights),
	}
}

func (w *Weights) Reset(eType int) {
	w.allWeight = make(map[int]*weights)
}

func (w *Weights) AddWeight(eType int, value int, weight int) {
	wi := w.allWeight[eType]
	if wi == nil {
		wi = &weights{}
		w.allWeight[eType] = wi
	}
	wi.allValue = append(wi.allValue, &intervalValue{
		startWeight: wi.maxWeight,
		endWeight:   wi.maxWeight + weight,
		value:       value,
	})
	wi.maxWeight += weight
}

func (w *Weights) AddWeights(eType int, value []int, weight []int) {
	if len(value) != len(weight) {
		Log.Error("Failed to AddWeights, value len not equal weight len")
		return
	}
	for i, v := range value {
		w.AddWeight(eType, v, weight[i])
	}
}

func (w *Weights) Value(eType int) int {
	if pro, ok := w.allWeight[eType]; ok {
		rate := rand.Intn(pro.maxWeight)
		for _, value := range pro.allValue {
			if value.startWeight <= rate && rate < value.endWeight {
				return value.value
			}
		}
	} else {
		Log.Error("Failed to get value type error, type = %v", eType)
	}
	return 0
}
