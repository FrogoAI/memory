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

// Len delegates to the underlying SortedSet under a read lock.
func (s *SafeSortedSet[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.Len()
}

// PeekMin delegates to the underlying SortedSet under a read lock.
func (s *SafeSortedSet[K, V]) PeekMin() *Node[K, V] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.PeekMin()
}

// PopMin delegates to the underlying SortedSet under a write lock.
func (s *SafeSortedSet[K, V]) PopMin() *Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.PopMin()
}

// PeekMax delegates to the underlying SortedSet under a read lock.
func (s *SafeSortedSet[K, V]) PeekMax() *Node[K, V] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.PeekMax()
}

// PopMax delegates to the underlying SortedSet under a write lock.
func (s *SafeSortedSet[K, V]) PopMax() *Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.PopMax()
}

// Upsert delegates to the underlying SortedSet under a write lock.
func (s *SafeSortedSet[K, V]) Upsert(key K, value V) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.Upsert(key, value)
}

// Remove delegates to the underlying SortedSet under a write lock.
func (s *SafeSortedSet[K, V]) Remove(value V) *Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.Remove(value)
}

// GetTop delegates to the underlying SortedSet under a write lock.
func (s *SafeSortedSet[K, V]) GetTop(count int, remove bool) []*Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.GetTop(count, remove)
}

// GetRTop delegates to the underlying SortedSet under a write lock.
func (s *SafeSortedSet[K, V]) GetRTop(count int, remove bool) []*Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.GetRTop(count, remove)
}

// GetUntilKey delegates to the underlying SortedSet under a write lock.
func (s *SafeSortedSet[K, V]) GetUntilKey(untilKey K, remove bool) []any {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.GetUntilKey(untilKey, remove)
}

// GetByKeyRange delegates to the underlying SortedSet under a write lock.
func (s *SafeSortedSet[K, V]) GetByKeyRange(start K, end K, options *GetByKeyRangeOptions) []*Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.GetByKeyRange(start, end, options)
}

// GetByRankRange delegates to the underlying SortedSet under a write lock.
func (s *SafeSortedSet[K, V]) GetByRankRange(start int, end int, remove bool) []*Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.GetByRankRange(start, end, remove)
}

// GetByRank delegates to the underlying SortedSet under a write lock.
func (s *SafeSortedSet[K, V]) GetByRank(rank int, remove bool) *Node[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.GetByRank(rank, remove)
}

// GetByValue delegates to the underlying SortedSet under a read lock.
func (s *SafeSortedSet[K, V]) GetByValue(value V) *Node[K, V] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.GetByValue(value)
}

// Contains delegates to the underlying SortedSet under a read lock.
func (s *SafeSortedSet[K, V]) Contains(value V) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.Contains(value)
}

// FindRank delegates to the underlying SortedSet under a read lock.
func (s *SafeSortedSet[K, V]) FindRank(value V) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.FindRank(value)
}

// Dump delegates to the underlying SortedSet under a read lock.
func (s *SafeSortedSet[K, V]) Dump(makeDump func(key K, value V) (string, string, error)) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ss.Dump(makeDump)
}

// Restore delegates to the underlying SortedSet under a write lock.
func (s *SafeSortedSet[K, V]) Restore(keyRestore func(key string, values []string) (K, []V, error), dump string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ss.Restore(keyRestore, dump)
}
