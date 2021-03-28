package main

import "sort"

type ArgSortIntSlice struct {
	sort.IntSlice
	indices []int
}

func (s ArgSortIntSlice) Swap(i, j int) {
	s.IntSlice.Swap(i, j)
	s.indices[i], s.indices[j] = s.indices[j], s.indices[i]
}

func ArgSort(n []int) []int {
	array := make([]int, len(n))
	indices := make([]int, len(n))
	for idx := range n {
		array[idx] = n[idx]
		indices[idx] = idx
	}

	s := ArgSortIntSlice{IntSlice: sort.IntSlice(array), indices: indices}
	sort.Stable(s)
	return s.indices
}
