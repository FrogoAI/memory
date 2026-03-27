// Package orderedmap provides an insertion-ordered map.
//
// NOT safe for concurrent use. Callers must synchronize access externally.
package orderedmap

type OrderedMap[K comparable, V any] struct {
	values map[K]V // nolint:structcheck
	keys   []K     // nolint:structcheck
}

func (o *OrderedMap[K, V]) Copy() *OrderedMap[K, V] {

	s := &OrderedMap[K, V]{}
	for _, key := range o.keys {
		s.Add(key, o.values[key])
	}

	return s
}

func (o *OrderedMap[K, V]) Clear() {

	o.values = map[K]V{}
	o.keys = []K{}
}

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

func (o *OrderedMap[K, V]) Remove(key K) {

	for i, k := range o.keys {
		if k == key {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}

	delete(o.values, key)
}

func (o *OrderedMap[K, V]) Get(key K) V {

	return o.values[key]
}

func (o *OrderedMap[K, V]) Exists(key K) bool {

	_, exists := o.values[key]

	return exists
}

func (o *OrderedMap[K, V]) Size() int {

	return len(o.keys)
}

func (o *OrderedMap[K, V]) SetKeys(keys []K, def V) {
	for _, key := range keys {
		o.Add(key, def)
	}
}

func (o *OrderedMap[K, V]) GetAll() ([]K, []V) {

	keys := make([]K, len(o.keys))
	copy(keys, o.keys)

	values := make([]V, len(keys))
	for i, key := range keys {
		values[i] = o.values[key]
	}

	return keys, values
}

func (o *OrderedMap[K, V]) GetMap() map[K]V {

	res := map[K]V{}
	for _, key := range o.keys {
		res[key] = o.values[key]
	}

	return res
}

func (o *OrderedMap[K, V]) SetAll(values []V) {

	if len(o.keys) > len(values) {
		return
	}

	for i := range o.keys {
		key := o.keys[i]
		o.values[key] = values[i]
	}
}

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
