package leetcode

import (
	"testing"
)

func TestSingleNumber(t *testing.T) {
	var num int

	num = singleNumber([]int{2, 2, 1})
	if num != 1 {
		t.Error("ERROR! expect: ", 1, ", actual: ", num)
	}

	num = singleNumber([]int{4, 1, 2, 1, 2})
	if num != 4 {
		t.Error("ERROR! expect: ", 4, ", actual: ", num)
	}

	num = singleNumber([]int{1, 1, 2, 3, 3})
	if num != 2 {
		t.Error("ERROR! expect: ", 2, ", actual: ", num)
	}
}
