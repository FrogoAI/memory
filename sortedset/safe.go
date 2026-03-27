package sortedset

import (
	"sync"

	"github.com/FrogoAI/memory/comparator"
)

// SafeSortedSet is a thread-safe wrapper around SortedSet.
// All methods are protected by a read-write mutex.
type SafeSortedSet[K comparable, V comparable] struct {
	mu sync.RWMutex
	ss *SortedSet[K, V]
}

// NewSafeSortedSet creates a thread-safe sorted set.
func NewSafeSortedSet[K comparable, V comparable](c comparator.Comparator) *SafeSortedSet[K, V] {
	return &SafeSortedSet[K, V]{
		ss: NewSortedSet[K, V](c),
	}
}

func (s *SafeSortedSet[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.Len()
}

func (s *SafeSortedSet[K, V]) PeekMin() *Node[K, V] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.PeekMin()
}

func (s *SafeSortedSet[K, V]) PopMin() *Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.PopMin()
}

func (s *SafeSortedSet[K, V]) PeekMax() *Node[K, V] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.PeekMax()
}

func (s *SafeSortedSet[K, V]) PopMax() *Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.PopMax()
}

func (s *SafeSortedSet[K, V]) Upsert(key K, value V) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.Upsert(key, value)
}

func (s *SafeSortedSet[K, V]) Remove(value V) *Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.Remove(value)
}

func (s *SafeSortedSet[K, V]) GetTop(count int, remove bool) []*Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.GetTop(count, remove)
}

func (s *SafeSortedSet[K, V]) GetRTop(count int, remove bool) []*Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.GetRTop(count, remove)
}

func (s *SafeSortedSet[K, V]) GetUntilKey(untilKey K, remove bool) []any {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.GetUntilKey(untilKey, remove)
}

func (s *SafeSortedSet[K, V]) GetByKeyRange(start K, end K, options *GetByKeyRangeOptions) []*Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.GetByKeyRange(start, end, options)
}

func (s *SafeSortedSet[K, V]) GetByRankRange(start int, end int, remove bool) []*Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.GetByRankRange(start, end, remove)
}

func (s *SafeSortedSet[K, V]) GetByRank(rank int, remove bool) *Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.GetByRank(rank, remove)
}

func (s *SafeSortedSet[K, V]) GetByValue(value V) *Node[K, V] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.GetByValue(value)
}

func (s *SafeSortedSet[K, V]) Contains(value V) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.Contains(value)
}

func (s *SafeSortedSet[K, V]) FindRank(value V) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.FindRank(value)
}

func (s *SafeSortedSet[K, V]) Dump(makeDump func(key K, value V) (string, string, error)) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.Dump(makeDump)
}

func (s *SafeSortedSet[K, V]) Restore(keyRestore func(key string, values []string) (K, []V, error), dump string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.Restore(keyRestore, dump)
}
