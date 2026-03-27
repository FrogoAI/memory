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
