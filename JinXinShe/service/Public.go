package service

import (
	"sync/atomic"
)

type WriteCloser interface {
	Close()
}

// AutomicInt64 原子计数器
type AutomicInt64 int64

// NewAutomicInt64 创建一个AutomicInt64
func NewAutomicInt64(iniVal int64) *AutomicInt64 {
	a := AutomicInt64(iniVal)
	return &a
}

// Get 返回int64的值
func (a *AutomicInt64) Get() int64 {
	return int64(*a)
}

// CompareAndSet 原子操作，比较*a和expect，相等则*a=update
func (a *AutomicInt64) CompareAndSet(expect, update int64) bool {
	return atomic.CompareAndSwapInt64((*int64)(a), expect, update)
}

// GetAndIncrement 获取值并加1（原子操作）
func (a *AutomicInt64) GetAndIncrement() int64 {
	for {
		curCount := a.Get()
		next := curCount + 1
		if a.CompareAndSet(curCount, next) {
			return curCount
		}
	}
}

// GetAndDecrement 获取值并减1（原子操作）
func (a *AutomicInt64) GetAndDecrement() int64 {
	for {
		curCount := a.Get()
		next := curCount - 1
		if a.CompareAndSet(curCount, next) {
			return curCount
		}
	}
}
