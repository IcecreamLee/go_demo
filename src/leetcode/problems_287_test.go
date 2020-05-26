package leetcode

import (
	"testing"
)

func TestFindDuplicate(t *testing.T) {
	result := findDuplicate([]int{1, 2, 2, 3, 4, 5})
	if result != 2 {
		t.Error("ERROR! expect: ", 2, ", actual: ", result)
	}
}
