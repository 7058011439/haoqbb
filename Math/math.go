/*
	这里没啥好说的，提供的Factorial和Combination函数，运算量是指数增长的，也就是一不小心就要溢出，所以传入的值尽量小吧。
	如果的确有较大的值，可以调用XXXBig函数，相对可以更大一些。但是返回的值是int64，所以也别指望大太多。
*/
package Math

import (
	"math/big"
)

// Factorial 计算阶乘
func Factorial(n int) int64 {
	result := int64(1)
	for i := 1; i <= n; i++ {
		result *= int64(i)
	}
	return result
}

// Combination 计算组合数
func Combination(n, m int) int64 {
	if m > n || m < 0 {
		return 0
	}
	result := int64(1)

	count := m
	if count > n-m {
		count = n - m
	}
	// 计算 n! / (m! * (n-m)!)
	for i := 1; i <= count; i++ {
		result *= int64(n - i + 1)
		result /= int64(i)
	}

	return result
}

// CombinationBig 计算组合数大数据
func CombinationBig(n, m int) *big.Int {
	result := big.NewInt(1)

	// 计算 n! / (m! * (n-m)!)
	for i := 1; i <= m; i++ {
		result.Mul(result, big.NewInt(int64(n-i+1)))
		result.Div(result, big.NewInt(int64(i)))
	}

	return result
}
