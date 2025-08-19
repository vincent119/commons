package slicex

// Contains 檢查 slice 是否包含指定元素。
func Contains[T comparable](s []T, v T) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}
	return false
}

// IndexOf 回傳元素第一次出現的索引，若不存在回傳 -1。
func IndexOf[T comparable](s []T, v T) int {
	for i, e := range s {
		if e == v {
			return i
		}
	}
	return -1
}

// Filter 回傳符合條件的子 slice（不修改原 slice）。
func Filter[T any](s []T, f func(T) bool) []T {
	res := make([]T, 0, len(s))
	for _, e := range s {
		if f(e) {
			res = append(res, e)
		}
	}
	return res
}

// Map 轉換 slice 中的每個元素，並產生新 slice。
func Map[T any, R any](s []T, f func(T) R) []R {
	res := make([]R, 0, len(s))
	for _, e := range s {
		res = append(res, f(e))
	}
	return res
}