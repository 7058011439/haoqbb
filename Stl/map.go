package Stl

import (
	"fmt"
)

type DoubleMap struct {
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
	if key == nil || value == nil {
		fmt.Printf("Failed to Add, key == nil || value == nil, key = %v, value = %v", key, value)
		return
	}
	d.mapKeyValue[key] = value
	d.mapValueKey[value] = key
}

func (d *DoubleMap) RemoveByKey(key interface{}) {
	value := d.mapKeyValue[key]
	delete(d.mapKeyValue, key)
	delete(d.mapValueKey, value)
}

func (d *DoubleMap) RemoveByValue(value interface{}) {
	key := d.mapValueKey[value]
	delete(d.mapKeyValue, key)
	delete(d.mapValueKey, value)
}

func (d *DoubleMap) GetKey(value interface{}) interface{} {
	return d.mapValueKey[value]
}

func (d *DoubleMap) GetValue(key interface{}) interface{} {
	return d.mapKeyValue[key]
}

func (d *DoubleMap) Len() int {
	return len(d.mapKeyValue)
}
