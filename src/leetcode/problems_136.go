package leetcode

//给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
//
//说明：
//
//你的算法应该具有线性时间复杂度。 你可以不使用额外空间来实现吗？
//
//示例 1:
//
//输入: [2,2,1]
//输出: 1
//示例 2:
//
//输入: [4,1,2,1,2]
//输出: 4
//
//来源：力扣（LeetCode）
//链接：https://leetcode-cn.com/problems/single-number
//著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。

// by me
func singleNumber(nums []int) int {
find:
	//fmt.Println("nums:", nums)
	for j := 1; j < len(nums); j++ {
		//fmt.Println("j:", j)
		if nums[0] == nums[j] {
			if j < len(nums)-1 {
				nums = append(nums[1:j], nums[j+1:]...)
			} else {
				nums = nums[1:j]
			}
			goto find
		}
	}
	//fmt.Println("num: ", nums[0])
	return nums[0]
}

// by leetcode
func singleNumber2(nums []int) int {
	single := 0
	for _, num := range nums {
		single ^= num
	}
	return single
}
