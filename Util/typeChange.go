package Util

import (
	"fmt"
	"strconv"
)

func StrToInt16(value string) int16 {
	if value == "" {
		return 0
	}
	if data, err := strconv.ParseInt(value, 10, 16); err != nil {
		fmt.Printf("Failed to StrToInt16, err = %v\r\n", err)
		return 0
	} else {
		return int16(data)
	}
}

func StrToInt(value string) int {
	if value == "" {
		return 0
	}
	if data, err := strconv.Atoi(value); err != nil {
		fmt.Printf("Failed to StrToInt, err = %v\r\n", err)
		return 0
	} else {
		return data
	}
}

func StrToInt64(value string) int64 {
	if value == "" {
		return 0
	}
	if data, err := strconv.ParseInt(value, 10, 64); err != nil {
		fmt.Printf("Failed to StrToInt64, err = %v\r\n", err)
		return 0
	} else {
		return data
	}
}

func StrToFloat64(value string) float64 {
	if value == "" {
		return 0
	}
	if data, err := strconv.ParseFloat(value, 64); err != nil {
		fmt.Printf("Failed to StrToFloat64, err = %v\r\n", err)
		return 0
	} else {
		return data
	}
}

func StrToFloat32(value string) float32 {
	if value == "" {
		return 0
	}
	if data, err := strconv.ParseFloat(value, 32); err != nil {
		fmt.Printf("Failed to StrToFloat32, err = %v\r\n", err)
		return 0
	} else {
		return float32(data)
	}
}
