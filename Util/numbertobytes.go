package Util

import (
	"encoding/binary"
	"math"
)

var littleEndian = true

func Int16ToBytes(data int16) []byte {
	return Uint16ToBytes(uint16(data))
}

func Uint16ToBytes(data uint16) []byte {
	ret := make([]byte, 2)
	if littleEndian {
		binary.LittleEndian.PutUint16(ret, data)
	} else {
		binary.BigEndian.PutUint16(ret, data)
	}
	return ret
}

func IntToBytes(data int) []byte {
	return Uint32ToBytes(uint32(data))
}

func UintToBytes(data uint) []byte {
	return Uint32ToBytes(uint32(data))
}

func Int32ToBytes(data int32) []byte {
	return Uint32ToBytes(uint32(data))
}

func Uint32ToBytes(data uint32) []byte {
	ret := make([]byte, 4)
	if littleEndian {
		binary.LittleEndian.PutUint32(ret, data)
	} else {
		binary.BigEndian.PutUint32(ret, data)
	}
	return ret
}

func Int64ToBytes(data int64) []byte {
	return Uint64ToBytes(uint64(data))
}

func Uint64ToBytes(data uint64) []byte {
	ret := make([]byte, 8)
	if littleEndian {
		binary.LittleEndian.PutUint64(ret, data)
	} else {
		binary.BigEndian.PutUint64(ret, data)
	}
	return ret
}

func Float32ToBytes(data float32) []byte {
	ret := make([]byte, 4)
	if littleEndian {
		binary.LittleEndian.PutUint32(ret, math.Float32bits(data))
	} else {
		binary.BigEndian.PutUint32(ret, math.Float32bits(data))
	}
	return ret
}

func Float64ToBytes(data float64) []byte {
	ret := make([]byte, 8)
	if littleEndian {
		binary.LittleEndian.PutUint64(ret, math.Float64bits(data))
	} else {
		binary.BigEndian.PutUint64(ret, math.Float64bits(data))
	}
	return ret
}

func Int16(data []byte) int16 {
	return int16(Uint16(data))
}

func Uint16(data []byte) uint16 {
	if littleEndian {
		return binary.LittleEndian.Uint16(data)
	} else {
		return binary.BigEndian.Uint16(data)
	}
}

func Int(data []byte) int {
	return int(Uint32(data))
}

func Uint(data []byte) uint {
	return uint(Uint32(data))
}

func Int32(data []byte) int32 {
	return int32(Uint32(data))
}

func Uint32(data []byte) uint32 {
	if littleEndian {
		return binary.LittleEndian.Uint32(data)
	} else {
		return binary.BigEndian.Uint32(data)
	}
}

func Int64(data []byte) int64 {
	return int64(Uint64(data))
}

func Uint64(data []byte) uint64 {
	if littleEndian {
		return binary.LittleEndian.Uint64(data)
	} else {
		return binary.BigEndian.Uint64(data)
	}
}

func Float32(data []byte) float32 {
	if littleEndian {
		return math.Float32frombits(binary.LittleEndian.Uint32(data))
	} else {
		return math.Float32frombits(binary.BigEndian.Uint32(data))
	}
}

func Float64(data []byte) float64 {
	if littleEndian {
		return math.Float64frombits(binary.LittleEndian.Uint64(data))
	} else {
		return math.Float64frombits(binary.BigEndian.Uint64(data))
	}
}
