package sortedset

// Node is an element in the sorted set, holding a key for ordering and an associated value.
type Node[K comparable, V comparable] struct {
	value    V // associated data
	key      K // key to determine the order of this node in the set
	backward *Node[K, V]
	level    []Level[K, V]
}

// Value returns the node's associated data.
func (s *Node[K, V]) Value() V {
	return s.value
}

// Key returns the node's ordering key.
func (s *Node[K, V]) Key() K {
	return s.key
}

// Level represents a single skip list level, holding a forward pointer and span.
type Level[K comparable, V comparable] struct {
	forward *Node[K, V]
	span    uint64
}
