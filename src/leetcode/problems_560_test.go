package leetcode

import (
	"testing"
)

func TestSubarraySum(t *testing.T) {
	var count int

	count = subarraySum2([]int{1, 1, 2, 3}, 3)
	if count != 2 {
		t.Error("ERROR! expect: ", 2, ", actual: ", count)
	}
	//
	//count = subarraySum([]int{-1, 2, 3, 1, 0, -2, 3, 1, 0}, 1)
	//if count != 3 {
	//	t.Error("ERROR! expect: ", 3, ", actual: ", count)
	//}
	//
	//count = subarraySum([]int{1, 2, 3}, 3)
	//if count != 2 {
	//	t.Error("ERROR! expect: ", 2, ", actual: ", count)
	//}
	//
	//count = subarraySum([]int{1}, 1)
	//if count != 1 {
	//	t.Error("ERROR! expect: ", 1, ", actual: ", count)
	//}

}
