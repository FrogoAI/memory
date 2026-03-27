package orderedmap

import "sync"

// SafeOrderedMap is a thread-safe wrapper around OrderedMap.
// All methods are protected by a read-write mutex.
type SafeOrderedMap[K comparable, V any] struct {
	mu sync.RWMutex
	om OrderedMap[K, V]
}

// Copy delegates to the underlying OrderedMap under a read lock.
func (s *SafeOrderedMap[K, V]) Copy() *SafeOrderedMap[K, V] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	inner := s.om.Copy()

	return &SafeOrderedMap[K, V]{om: *inner}
}

// Clear delegates to the underlying OrderedMap under a write lock.
func (s *SafeOrderedMap[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.om.Clear()
}

// Add delegates to the underlying OrderedMap under a write lock.
func (s *SafeOrderedMap[K, V]) Add(key K, val V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.om.Add(key, val)
}

// Remove delegates to the underlying OrderedMap under a write lock.
func (s *SafeOrderedMap[K, V]) Remove(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.om.Remove(key)
}

// Get delegates to the underlying OrderedMap under a read lock.
func (s *SafeOrderedMap[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.om.Get(key)
}

// Exists delegates to the underlying OrderedMap under a read lock.
func (s *SafeOrderedMap[K, V]) Exists(key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.om.Exists(key)
}

// Size delegates to the underlying OrderedMap under a read lock.
func (s *SafeOrderedMap[K, V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.om.Size()
}

// SetKeys delegates to the underlying OrderedMap under a write lock.
func (s *SafeOrderedMap[K, V]) SetKeys(keys []K, def V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.om.SetKeys(keys, def)
}

// Keys returns a copy of all keys in insertion order.
func (s *SafeOrderedMap[K, V]) Keys() []K {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.om.Keys()
}

// Values returns all values in insertion order.
func (s *SafeOrderedMap[K, V]) Values() []V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.om.Values()
}

// GetAll delegates to the underlying OrderedMap under a read lock.
func (s *SafeOrderedMap[K, V]) GetAll() ([]K, []V) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.om.GetAll()
}

// GetMap delegates to the underlying OrderedMap under a read lock.
func (s *SafeOrderedMap[K, V]) GetMap() map[K]V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.om.GetMap()
}

// SetAll delegates to the underlying OrderedMap under a write lock.
func (s *SafeOrderedMap[K, V]) SetAll(values []V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.om.SetAll(values)
}

// Iterator delegates to the underlying OrderedMap under a read lock.
func (s *SafeOrderedMap[K, V]) Iterator(size int) chan V {
	ch := make(chan V, size)

	go func() {
		s.mu.RLock()
		defer s.mu.RUnlock()

		for _, key := range s.om.keys {
			ch <- s.om.values[key]
		}

		close(ch)
	}()

	return ch
}
