package slice

// 在 a 不在 b 中的
// a := []T{"5", "2", "3", "4"}
// b := []T{"0", "1", "2", "3"}
// c => [5 4]
// DiffSliceInANotInB 返回一个包含只存在于切片a中但不在切片b中的元素的切片
func DiffSliceInANotInB[T comparable](a []T, b []T) (c []T) {
	// 创建一个map来存储切片b中的元素，以便快速查找
	temp := make(map[T]struct{})

	// 遍历b切片，将元素添加到map中
	for _, val := range b {
		temp[val] = struct{}{}
	}

	// 遍历a切片，检查哪些元素不在map中（即不在b中）
	for _, val := range a {
		if _, ok := temp[val]; !ok {
			c = append(c, val)
		}
	}

	return c
}

// MapOfElementsInANotInB 返回一个包含只存在于切片a中但不在切片b中的元素的切片map
func MapOfElementsInANotInB[T comparable](a []T, b []T) map[T]struct{} {
	// 创建一个map来存储切片b中的元素，以便快速查找
	temp := make(map[T]struct{})

	// 遍历b切片，将元素添加到map中
	for _, val := range b {
		temp[val] = struct{}{}
	}

	mp := make(map[T]struct{})

	// 遍历a切片，检查哪些元素不在map中（即不在b中）
	for _, val := range a {
		if _, ok := temp[val]; !ok {
			mp[val] = struct{}{}
		}
	}

	return mp
}
