package Stl

import (
	"fmt"
	"sync"
)

type DoubleMap struct {
	mutex       sync.RWMutex
	mapKeyValue map[interface{}]interface{}
	mapValueKey map[interface{}]interface{}
}

func NewDoubleMap() *DoubleMap {
	return &DoubleMap{
		mapKeyValue: map[interface{}]interface{}{},
		mapValueKey: map[interface{}]interface{}{},
	}
}

func (d *DoubleMap) Add(key interface{}, value interface{}) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if key == nil || value == nil {
		fmt.Printf("Failed to Add, key == nil || value == nil, key = %v, value = %v", key, value)
		return
	}
	d.mapKeyValue[key] = value
	d.mapValueKey[value] = key
}

func (d *DoubleMap) RemoveByKey(key interface{}) {
	d.mutex.Lock()
	d.mutex.Unlock()
	value := d.mapKeyValue[key]
	delete(d.mapKeyValue, key)
	delete(d.mapValueKey, value)
}

func (d *DoubleMap) RemoveByValue(value interface{}) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	key := d.mapValueKey[value]
	delete(d.mapKeyValue, key)
	delete(d.mapValueKey, value)
}

func (d *DoubleMap) GetKey(value interface{}) interface{} {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.mapValueKey[value]
}

func (d *DoubleMap) GetValue(key interface{}) interface{} {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.mapKeyValue[key]
}

func (d *DoubleMap) Len() int {
	d.mutex.RLock()
	d.mutex.RUnlock()
	return len(d.mapKeyValue)
}
