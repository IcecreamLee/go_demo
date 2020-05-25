package leetcode

import (
	"testing"
)

func TestLRU(t *testing.T) {
	cache := Constructor(2 /* 缓存容量 */)

	var result int
	cache.Put(1, 1)
	cache.Put(2, 2)
	result = cache.Get(1) // 返回  1
	if result != 1 {
		t.Error("ERROR! expect: ", 1, ", actual: ", result)
	}
	cache.Put(3, 3)       // 该操作会使得密钥 2 作废
	result = cache.Get(2) // 返回 -1 (未找到)
	if result != -1 {
		t.Error("ERROR! expect: ", -1, ", actual: ", result)
	}
	cache.Put(4, 4)       // 该操作会使得密钥 1 作废
	result = cache.Get(1) // 返回 -1 (未找到)
	if result != -1 {
		t.Error("ERROR! expect: ", -1, ", actual: ", result)
	}
	result = cache.Get(3) // 返回  3
	if result != 3 {
		t.Error("ERROR! expect: ", 3, ", actual: ", result)
	}
	result = cache.Get(4) // 返回  4
	if result != 4 {
		t.Error("ERROR! expect: ", 4, ", actual: ", result)
	}
}
