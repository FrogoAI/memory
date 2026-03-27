// Package lru provides an LRU cache that evicts the least recently used entry
// when capacity is reached.
//
// NOT safe for concurrent use. Callers must synchronize access externally.
package lru

import (
	"container/list"
	"fmt"
)

// Value wraps a cached item along with its position in the eviction list.
type Value[K any] struct {
	Value   K
	Element *list.Element
}

// Cache is an LRU cache that evicts the least recently used entry
// when capacity is reached.
type Cache[K any] struct {
	list     *list.List
	index    map[string]*Value[K]
	capacity int
}

// NewLRUCache creates a new Cache with the given maximum capacity.
func NewLRUCache[K any](capacity int) *Cache[K] {
	return &Cache[K]{
		list:     list.New(),
		index:    map[string]*Value[K]{},
		capacity: capacity,
	}
}

// Len returns the number of items in the cache.
func (c *Cache[K]) Len() int {
	return len(c.index)
}

// Copy returns a deep copy of the cache. The new cache has independent storage
// but shares the same value references. Eviction order is preserved.
func (c *Cache[K]) Copy() *Cache[K] {
	clone := NewLRUCache[K](c.capacity)

	// Walk from back (LRU) to front (MRU) so that the most recently used
	// entry ends up at the front of the clone's eviction list.
	for e := c.list.Back(); e != nil; e = e.Prev() {
		key := fmt.Sprint(e.Value)
		clone.Put(key, c.index[key].Value)
	}

	return clone
}

// Clear removes all items from the cache.
func (c *Cache[K]) Clear() {
	c.list.Init()
	c.index = map[string]*Value[K]{}
}

// Get returns the value associated with key and promotes it to the front of
// the eviction list. The second return value reports whether the key was found.
func (c *Cache[K]) Get(key string) (K, bool) {
	var val K

	v, ok := c.index[key]
	if !ok {
		return val, false
	}

	c.list.MoveToFront(v.Element)

	return v.Value, true
}

// Iterator returns an iterator over cache entries from most recently used
// to least recently used. It can be used directly in a for-range loop:
//
//	for key, value := range cache.Iterator() {
//	    // ...
//	}
func (c *Cache[K]) Iterator() func(yield func(string, K) bool) {
	return func(yield func(string, K) bool) {
		for e := c.list.Front(); e != nil; e = e.Next() {
			key := fmt.Sprint(e.Value)
			if !yield(key, c.index[key].Value) {
				return
			}
		}
	}
}

// Put adds or updates an entry in the cache and evicts the least recently used
// entry if the cache is at capacity.
func (c *Cache[K]) Put(key string, value K) {
	v, ok := c.index[key]
	if !ok {
		c.list.PushFront(key)
	} else {
		c.list.MoveToFront(v.Element)
	}

	c.index[key] = &Value[K]{
		Value:   value,
		Element: c.list.Front(),
	}

	listSize := c.list.Len()

	if listSize > c.capacity {
		prevKey := c.list.Back()
		c.list.Remove(prevKey)

		delete(c.index, fmt.Sprint(prevKey.Value))
	}
}
