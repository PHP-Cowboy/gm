package slice

// Contains 检查切片 s 中是否包含元素 value
func Contains[T comparable](s []T, value T) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}
