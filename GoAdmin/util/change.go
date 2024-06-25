package util

import (
	"encoding/json"
)

func StructToMap(s interface{}) map[string]interface{} {
	ret := map[string]interface{}{}
	if arr, err := json.Marshal(s); err != nil {
		return ret
	} else {
		err = json.Unmarshal(arr, &ret)
		if err != nil {
			return ret
		}
	}
	return ret
}

func RemoveDuplicates[T comparable](slice []T, key func(interface{}) interface{}) []T {
	uniqueMap := make(map[interface{}]bool)
	uniqueSlice := []T{}

	for _, item := range slice {
		if _, exists := uniqueMap[key(item)]; !exists {
			uniqueMap[key(item)] = true
			uniqueSlice = append(uniqueSlice, item)
		}
	}

	return uniqueSlice
}
