package xutils

// ArrayColumn 返回输入数组中指定列的值
func ArrayColumn[T any, V comparable](ss []T, fn func(v T) V, dedup bool) []V {
	if len(ss) == 0 {
		return []V{}
	}
	ids := make([]V, 0, len(ss))
	m := make(map[V]bool)
	for _, v := range ss {
		id := fn(v)
		if dedup && m[id] {
			continue
		}
		ids = append(ids, id)
		m[id] = true
	}
	return ids
}

// Cartesian 计算多个集合的笛卡尔积
func Cartesian(sets [][]int) [][]int {
	if len(sets) == 0 {
		return [][]int{}
	}
	if len(sets) == 1 {
		result := make([][]int, len(sets[0]))
		for i := range sets[0] {
			result[i] = []int{sets[0][i]}
		}
		return result
	}

	var result [][]int
	for _, set := range sets[0] {
		for _, subset := range Cartesian(sets[1:]) {
			result = append(result, append([]int{set}, subset...))
		}
	}
	return result
}

// ArrayUnique 移除数组中重复的值
func ArrayUnique[T comparable](ss []T) []T {
	if len(ss) <= 1 {
		return ss
	}
	ret := make([]T, 0, len(ss))
	m := make(map[T]struct{})
	for _, v := range ss {
		if _, ok := m[v]; !ok {
			ret = append(ret, v)
			m[v] = struct{}{}
		}
	}
	return ret
}

// InArray 检查数组中是否存在某个值
func InArray[T comparable](val T, arr []T) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

// ArrayDiff 数组取差集
// 对比 array 和其他一个或者多个数组，返回在 array 中但是不在其他 array 里的值
func ArrayDiff[T comparable](array []T, arrays ...[]T) []T {
	if len(arrays) == 0 {
		return array
	}

	var ret []T
	values := make(map[T]struct{})
	for _, arr := range arrays {
		for _, v := range arr {
			values[v] = struct{}{}
		}
	}
	for _, v := range array {
		if _, ok := values[v]; !ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// ArrayReverse 返回顺序相反的数组
func ArrayReverse[T any](s []T) []T {
	if len(s) >= 2 {
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}
	}
	return s
}

// Map 返回使用T中指定字段作为key的map
func Map[T any, V comparable](ss []T, f func(v T) V) map[V]T {
	m := make(map[V]T)
	for _, v := range ss {
		m[f(v)] = v
	}
	return m
}
