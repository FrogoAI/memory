package orderedmap

import "sync"

// SafeOrderedMap is a thread-safe wrapper around OrderedMap.
// All methods are protected by a read-write mutex.
type SafeOrderedMap[K comparable, V any] struct {
	mu sync.RWMutex
	om OrderedMap[K, V]
}

func (s *SafeOrderedMap[K, V]) Copy() *SafeOrderedMap[K, V] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	inner := s.om.Copy()

	return &SafeOrderedMap[K, V]{om: *inner}
}

func (s *SafeOrderedMap[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.om.Clear()
}

func (s *SafeOrderedMap[K, V]) Add(key K, val V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.om.Add(key, val)
}

func (s *SafeOrderedMap[K, V]) Remove(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.om.Remove(key)
}

func (s *SafeOrderedMap[K, V]) Get(key K) V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.om.Get(key)
}

func (s *SafeOrderedMap[K, V]) Exists(key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.om.Exists(key)
}

func (s *SafeOrderedMap[K, V]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.om.Size()
}

func (s *SafeOrderedMap[K, V]) SetKeys(keys []K, def V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.om.SetKeys(keys, def)
}

func (s *SafeOrderedMap[K, V]) GetAll() ([]K, []V) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.om.GetAll()
}

func (s *SafeOrderedMap[K, V]) GetMap() map[K]V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.om.GetMap()
}

func (s *SafeOrderedMap[K, V]) SetAll(values []V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.om.SetAll(values)
}

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
