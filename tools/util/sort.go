package util

// Int64SliceSort int64 slice sort util
type Int64SliceSort struct {
	Values []int64
}

// Len Get length
func (i64sort *Int64SliceSort) Len() int {
	return len(i64sort.Values)
}

// Less compare i,  j
func (i64sort *Int64SliceSort) Less(i, j int) bool {
	return i64sort.Values[i] < i64sort.Values[j]
}

// Swap swap i <=> j
func (i64sort *Int64SliceSort) Swap(i, j int) {
	i64sort.Values[i], i64sort.Values[j] = i64sort.Values[j], i64sort.Values[i]
}
