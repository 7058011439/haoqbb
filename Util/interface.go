package Util

import (
	"fmt"
	"reflect"
	"strings"
)

func CompareLess(dataA interface{}, dataB interface{}) bool {
	typeA := reflect.TypeOf(dataA).Kind()
	typeB := reflect.TypeOf(dataB).Kind()
	if typeA != typeB {
		return false
	}
	switch typeA {
	case reflect.Int:
		return dataA.(int) < dataB.(int)
	case reflect.Int8:
		return dataA.(int8) < dataB.(int8)
	case reflect.Int16:
		return dataA.(int16) < dataB.(int16)
	case reflect.Int32:
		return dataA.(int32) < dataB.(int32)
	case reflect.Int64:
		return dataA.(int64) < dataB.(int64)
	case reflect.Uint:
		return dataA.(uint) < dataB.(uint)
	case reflect.Uint8:
		return dataA.(uint8) < dataB.(uint8)
	case reflect.Uint16:
		return dataA.(uint16) < dataB.(uint16)
	case reflect.Uint32:
		return dataA.(uint32) < dataB.(uint32)
	case reflect.Uint64:
		return dataA.(uint64) < dataB.(uint64)
	case reflect.Float32:
		return dataA.(float32) < dataB.(float32)
	case reflect.Float64:
		return dataA.(float64) < dataB.(float64)
	case reflect.String:
		return strings.Compare(dataA.(string), dataB.(string)) < 0
	default:
		fmt.Printf("Failed to CompareLess, unknown type = %v", typeA)
		return false
	}
}

func CompareEqual(dataA interface{}, dataB interface{}) bool {
	typeA := reflect.TypeOf(dataA).Kind()
	typeB := reflect.TypeOf(dataB).Kind()
	if typeA != typeB {
		return false
	}
	switch typeA {
	case reflect.Int:
		return dataA.(int) == dataB.(int)
	case reflect.Int8:
		return dataA.(int8) == dataB.(int8)
	case reflect.Int16:
		return dataA.(int16) == dataB.(int16)
	case reflect.Int32:
		return dataA.(int32) == dataB.(int32)
	case reflect.Int64:
		return dataA.(int64) == dataB.(int64)
	case reflect.Uint:
		return dataA.(uint) == dataB.(uint)
	case reflect.Uint8:
		return dataA.(uint8) == dataB.(uint8)
	case reflect.Uint16:
		return dataA.(uint16) == dataB.(uint16)
	case reflect.Uint32:
		return dataA.(uint32) == dataB.(uint32)
	case reflect.Uint64:
		return dataA.(uint64) == dataB.(uint64)
	case reflect.Float32:
		return dataA.(float32) == dataB.(float32)
	case reflect.Float64:
		return dataA.(float64) == dataB.(float64)
	case reflect.Bool:
		return dataA.(bool) == dataB.(bool)
	case reflect.String:
		return strings.Compare(dataA.(string), dataB.(string)) == 0
	default:
		fmt.Printf("Failed to CompareEqual, unknown type = %v", typeA)
		return false
	}
}
