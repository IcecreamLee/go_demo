package leetcode

//287. 寻找重复数
//给定一个包含 n + 1 个整数的数组 nums，其数字都在 1 到 n 之间（包括 1 和 n），可知至少存在一个重复的整数。假设只有一个重复的整数，找出这个重复的数。
//
//示例 1:
//
//输入: [1,3,4,2,2]
//输出: 2
//示例 2:
//
//输入: [3,1,3,4,2]
//输出: 3
//说明：
//
//1. 不能更改原数组（假设数组是只读的）。
//2. 只能使用额外的 O(1) 的空间。
//3. 时间复杂度小于 O(n2) 。
//4. 数组中只有一个重复的数字，但它可能不止重复出现一次。

//解题思路：
//假设 nums = [1,2,2,3,4,5]
//最小的数为1，最大数为5
//使用二分查找的思路，假设重复的数字为中间值：(1 + 5) / 2 = 3
//查找数组中小于等于3的数量count，若count <= 3，则代表1 2 3中没有重复的，如count > 3, 则代表123中有重复的

func findDuplicate(nums []int) int {
	n := len(nums)
	left := 1
	right := n - 1
	duplicate := -1
	for left <= right {
		mid := (left + right) >> 1
		count := 0
		for i := 0; i < n; i++ {
			if nums[i] <= mid {
				count++
			}
		}
		if count <= mid {
			left = mid + 1
		} else {
			right = mid - 1
			duplicate = mid
		}
	}
	return duplicate
}
