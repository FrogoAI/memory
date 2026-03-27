package utils

import (
	"sort"

	"github.com/FrogoAI/memory/comparator"
)

// Sort sorts values (in-place) with respect to the given comparator.
// Returns the first comparator error encountered during sorting, if any.
//
// Uses Go's sort (hybrid of quicksort for large and then insertion sort for smaller slices).
func Sort(values []interface{}, c comparator.Comparator) error {
	s := &sortable{values: values, comparator: c}
	sort.Sort(s)

	return s.err
}

type sortable struct {
	comparator comparator.Comparator
	values     []interface{}
	err        error
}

// Len return len of slice
func (s *sortable) Len() int {
	return len(s.values)
}

// Swap implement swap items
func (s *sortable) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

// Less return true if i-th less of j-th item
func (s *sortable) Less(i, j int) bool {
	result, err := s.comparator(s.values[i], s.values[j])
	if err != nil && s.err == nil {
		s.err = err
	}

	return result < 0
}
