//nolint:revive // package name is intentional
package utils

import "sync"

// SafeList is a thread-safe append-only list protected by a read-write mutex.
type SafeList[V any] struct {
	list []V
	mu   sync.RWMutex
}

// NewSafeList creates a new SafeList optionally initialized with the given elements.
func NewSafeList[K any](data ...K) *SafeList[K] {
	s := &SafeList[K]{
		list: data,
	}

	return s
}

// Add appends a value to the list.
func (s *SafeList[V]) Add(value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.list = append(s.list, value)
}

// List returns a copy of all elements.
func (s *SafeList[V]) List() []V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dst := make([]V, len(s.list))
	copy(dst, s.list)

	return dst
}

// Clear removes all elements and returns them.
func (s *SafeList[V]) Clear() []V {
	s.mu.Lock()
	defer s.mu.Unlock()

	dst := make([]V, len(s.list))
	copy(dst, s.list)
	s.list = []V{}

	return dst
}

// Count returns the number of elements in the list.
func (s *SafeList[V]) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.list)
}
