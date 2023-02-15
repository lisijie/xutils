package xutils

// IF 如果条件为true返回a，否则返回b
func IF[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}
