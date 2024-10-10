package Util

import "sync"

var slicePool sync.Map

func formatNumber(value int) (ret int) {
	if value < 1 {
		return 1
	}

	// 从 1 开始，找出大于等于 value 的最近 2 的幂次值
	ret = 1
	for ret < value {
		ret <<= 1 // 左移一位，相当于乘以 2
	}

	return ret
}

// createOrLoadSlicePool 封装池创建或加载的逻辑
func createOrLoadSlicePool(size int) *sync.Pool {
	size = formatNumber(size)
	if pool, ok := slicePool.Load(size); ok {
		return pool.(*sync.Pool)
	}

	pool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, size)
		},
	}
	slicePool.Store(size, pool)
	return pool
}

func GetNewSlice(size int) []byte {
	return createOrLoadSlicePool(size).Get().([]byte)
}

func GiveBackSlice(data []byte) {
	// 使用[:0]复用切片，确保切片长度为0，避免后续访问时出现意外
	createOrLoadSlicePool(cap(data)).Put(data[:0])
}

var objectPool sync.Map

type IObject interface {
	Reset()
	Key() string
	New() interface{}
}

func createOrLoadObjectPool(object IObject) *sync.Pool {
	key := object.Key()
	if pool, ok := objectPool.Load(key); ok {
		return pool.(*sync.Pool)
	}

	pool := &sync.Pool{
		New: object.New,
	}
	objectPool.Store(key, pool)
	return pool
}

func GetNewObj(object IObject) interface{} {
	return createOrLoadObjectPool(object).Get()
}

func GiveBackObj(object IObject) {
	object.Reset()
	createOrLoadObjectPool(object).Put(object)
}
