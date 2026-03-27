// Package orderedmap provides an insertion-ordered map.
//
// NOT safe for concurrent use. Callers must synchronize access externally.
package orderedmap

// OrderedMap is a map that preserves insertion order of keys.
type OrderedMap[K comparable, V any] struct {
	values map[K]V
	keys   []K
}

// Copy returns a shallow copy of the ordered map.
func (o *OrderedMap[K, V]) Copy() *OrderedMap[K, V] {
	s := &OrderedMap[K, V]{}
	for _, key := range o.keys {
		s.Add(key, o.values[key])
	}

	return s
}

// Clear removes all entries from the map.
func (o *OrderedMap[K, V]) Clear() {
	o.values = map[K]V{}
	o.keys = []K{}
}

// Add inserts or updates the value for the given key, preserving insertion order.
func (o *OrderedMap[K, V]) Add(key K, val V) {
	if o.values == nil {
		o.values = map[K]V{}
	}

	_, ok := o.values[key]
	if !ok {
		o.keys = append(o.keys, key)
	}

	o.values[key] = val
}

// Remove deletes the entry for the given key.
func (o *OrderedMap[K, V]) Remove(key K) {
	for i, k := range o.keys {
		if k == key {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}

	delete(o.values, key)
}

// Get returns the value for the given key and whether it exists.
func (o *OrderedMap[K, V]) Get(key K) (V, bool) {
	val, ok := o.values[key]

	return val, ok
}

// Exists reports whether the given key is present in the map.
func (o *OrderedMap[K, V]) Exists(key K) bool {
	_, exists := o.values[key]

	return exists
}

// Size returns the number of entries in the map.
func (o *OrderedMap[K, V]) Size() int {
	return len(o.keys)
}

// SetKeys adds multiple keys with the same default value.
func (o *OrderedMap[K, V]) SetKeys(keys []K, def V) {
	for _, key := range keys {
		o.Add(key, def)
	}
}

// Keys returns a copy of all keys in insertion order.
func (o *OrderedMap[K, V]) Keys() []K {
	keys := make([]K, len(o.keys))
	copy(keys, o.keys)

	return keys
}

// Values returns all values in insertion order.
func (o *OrderedMap[K, V]) Values() []V {
	values := make([]V, len(o.keys))
	for i, key := range o.keys {
		values[i] = o.values[key]
	}

	return values
}

// GetAll returns copies of all keys and values in insertion order.
func (o *OrderedMap[K, V]) GetAll() ([]K, []V) {
	keys := make([]K, len(o.keys))
	copy(keys, o.keys)

	values := make([]V, len(keys))
	for i, key := range keys {
		values[i] = o.values[key]
	}

	return keys, values
}

// GetMap returns a plain map copy of all key-value pairs.
func (o *OrderedMap[K, V]) GetMap() map[K]V {
	res := map[K]V{}
	for _, key := range o.keys {
		res[key] = o.values[key]
	}

	return res
}

// SetAll replaces all values in key order; the slice must be at least as long as the number of keys.
func (o *OrderedMap[K, V]) SetAll(values []V) {
	if len(o.keys) > len(values) {
		return
	}

	for i := range o.keys {
		key := o.keys[i]
		o.values[key] = values[i]
	}
}

// Iterator returns a buffered channel that yields all values in insertion order.
func (o *OrderedMap[K, V]) Iterator(size int) chan V {
	ch := make(chan V, size)

	go func() {
		for _, key := range o.keys {
			ch <- o.values[key]
		}

		close(ch)
	}()

	return ch
}
