//nolint:revive // package name is intentional
package utils

import "sync"

// SafeMap is a thread-safe string-keyed map protected by a read-write mutex.
type SafeMap[K string, V any] struct {
	data map[K]V
	mu   *sync.RWMutex
}

// NewSafeMap creates a new SafeMap initialized with a copy of the given data.
func NewSafeMap[K string, V any](data map[K]V) *SafeMap[K, V] {
	s := &SafeMap[K, V]{
		data: map[K]V{},
		mu:   &sync.RWMutex{},
	}
	for k, v := range data {
		s.data[k] = v
	}

	return s
}

// Set inserts or updates the value for the given key.
func (s *SafeMap[K, V]) Set(name K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = map[K]V{}
	}

	s.data[name] = value
}

// Get returns the value for the given key and whether it exists.
func (s *SafeMap[K, V]) Get(name K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.data[name]

	return v, ok
}

// Exists reports whether the given key is present.
func (s *SafeMap[K, V]) Exists(name K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[name]

	return ok
}

// Remove deletes the entry for the given key.
func (s *SafeMap[K, V]) Remove(name K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, name)
}

// GetMap returns a copy of the underlying map.
func (s *SafeMap[K, V]) GetMap() map[K]V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := map[K]V{}
	for k, v := range s.data {
		result[k] = v
	}

	return result
}
