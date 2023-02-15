package xutils

// MapFlip 交换map中的键和值
func MapFlip[K comparable, V comparable](m map[K]V) map[V]K {
	ret := make(map[V]K)
	for k, v := range m {
		ret[v] = k
	}
	return ret
}
