package helper

// InSliceInt64 是否在陣列中
func InSliceInt64(nums []int64, target int64) bool {
	for _, n := range nums {
		if n == target {
			return true
		}
	}
	return false
}
