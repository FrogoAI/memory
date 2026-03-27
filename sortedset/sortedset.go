// Package sortedset provides a sorted set backed by a skip list.
//
// NOT safe for concurrent use. Callers must synchronize access externally.
package sortedset

import (
	"fmt"
	"math/rand"

	"github.com/FrogoAI/memory/comparator"
)

const (
	skipListMaxLevel = 32
	skipListP        = 0.25
	maxLimit         = 2147483648
	probabilityMask  = 0xFFFF // 16-bit bitmask for skip list level probability
)

// SortedSet is a skip-list-backed sorted set with O(log N) insert, delete, and lookup.
type SortedSet[K comparable, V comparable] struct {
	emptyKey   K
	header     *Node[K, V]
	tail       *Node[K, V]
	dict       map[V]*Node[K, V]
	comparator comparator.Comparator
	length     uint64
	level      int
}

// compare wraps the comparator and panics on type mismatch.
func (s *SortedSet[K, V]) compare(a, b interface{}) int {
	result, err := s.comparator(a, b)
	if err != nil {
		panic(fmt.Errorf("sortedset: %w", err))
	}

	return result
}

func (s *SortedSet[K, V]) createNode(level int, key K, value V) *Node[K, V] {
	node := &Node[K, V]{
		key:   key,
		value: value,
		level: make([]Level[K, V], level),
	}

	return node
}

// RandomLevel returns a random level for the new skiplist node we are going to create.
// The return value of this function is between 1 and skipListMaxLevel
// (both inclusive), with a powerlaw-alike distribution where higher
// levels are less likely to be returned.
func (s *SortedSet[K, V]) randomLevel() int {
	level := 1
	//nolint:gosec // math/rand is intentional for skip list leveling
	for float64(rand.Int31()&probabilityMask) < float64(skipListP*probabilityMask) {
		level++
	}

	if level < skipListMaxLevel {
		return level
	}

	return skipListMaxLevel
}

func (s *SortedSet[K, V]) insertNode(key K, value V) *Node[K, V] {
	var (
		update [skipListMaxLevel]*Node[K, V]
		rank   [skipListMaxLevel]uint64
	)

	x := s.header

	for i := s.level - 1; i >= 0; i-- {
		/* store rank that is crossed to reach the insert position */
		if s.level-1 == i {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for x.level[i].forward != nil &&
			(s.compare(x.level[i].forward.key, key) < 0 ||
				(s.compare(x.level[i].forward.key, key) == 0 && // key is the same but the key is different
					x.level[i].forward.value != value)) {
			//nolint:gosec // i is bounded by s.level which never exceeds skipListMaxLevel
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}

		update[i] = x
	}

	/* we assume the key is not already inside, since we allow duplicated
	 * keys, and the re-insertion of key and redis object should never
	 * happen since the caller of Insert() should test in the hash table
	 * if the element is already inside or not. */
	level := s.randomLevel()

	if level > s.level { // add a new level
		for i := s.level; i < level; i++ {
			rank[i] = 0
			update[i] = s.header
			update[i].level[i].span = s.length
		}

		s.level = level
	}

	x = s.createNode(level, key, value)
	for i := 0; i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		/* update span covered by update[i] as x is inserted here */
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	/* increment span for untouched levels */
	for i := level; i < s.level; i++ {
		update[i].level[i].span++
	}

	if update[0] == s.header {
		x.backward = nil
	} else {
		x.backward = update[0]
	}

	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		s.tail = x
	}

	s.length++

	return x
}

func (s *SortedSet[K, V]) deleteNode(x *Node[K, V], update [skipListMaxLevel]*Node[K, V]) {
	for i := 0; i < s.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span--
		}
	}

	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		s.tail = x.backward
	}

	for s.level > 1 && s.header.level[s.level-1].forward == nil {
		s.level--
	}

	s.length--
	delete(s.dict, x.value)
}

func (s *SortedSet[K, V]) delete(key K, value V) bool {
	var update [skipListMaxLevel]*Node[K, V]

	x := s.header
	for i := s.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && s.compare(x.level[i].forward.key, key) < 0 {
			x = x.level[i].forward
		}

		update[i] = x
	}
	/* We may have multiple elements with the same key, what we need
	 * is to find the element with both the right key and object. */
	x = x.level[0].forward
	if x != nil && key == x.key && x.value == value {
		s.deleteNode(x, update)
		// free x
		return true
	}

	return false /* not found */
}

// NewSortedSet creates a new empty sorted set using the given comparator for key ordering.
func NewSortedSet[K comparable, V comparable](c comparator.Comparator) *SortedSet[K, V] {
	var (
		emptyKey   K
		emptyValue V
	)

	sortedSet := &SortedSet[K, V]{
		level:      1,
		dict:       make(map[V]*Node[K, V]),
		comparator: c,
		emptyKey:   emptyKey,
	}
	sortedSet.header = sortedSet.createNode(skipListMaxLevel, emptyKey, emptyValue)

	return sortedSet
}

// Len returns the number of elements in the sorted set.
func (s *SortedSet[K, V]) Len() int {
	return int(s.length) //nolint:gosec // length is bounded by available memory; overflow beyond MaxInt is not practical
}

// PeekMin get the element with minimum key, nil if the set is empty
// Time complexity of this method is : O(log(N))
func (s *SortedSet[K, V]) PeekMin() *Node[K, V] {
	f := s.header.level[0].forward

	return f
}

// PopMin get and remove the element with minimal key, nil if the set is empty
// Time complexity of this method is : O(log(N))
func (s *SortedSet[K, V]) PopMin() *Node[K, V] {
	x := s.header.level[0].forward
	if x != nil {
		s.Remove(x.value)
	}

	return x
}

// PeekMax get the element with maximum key, nil if the set is empty
// Time Complexity : O(1)
func (s *SortedSet[K, V]) PeekMax() *Node[K, V] {
	t := s.tail

	return t
}

// PopMax get and remove the element with maximum key, nil if the set is empty
// Time complexity of this method is : O(log(N))
func (s *SortedSet[K, V]) PopMax() *Node[K, V] {
	x := s.tail
	if x != nil {
		s.Remove(x.value)
	}

	return x
}

// Upsert add an element into the sorted set with specific key / value / key.
// if the element is added, this method returns true; otherwise false means updated
// Time complexity of this method is : O(log(N))
func (s *SortedSet[K, V]) Upsert(key K, value V) bool {
	var newNode *Node[K, V]

	found := s.dict[value]
	if found != nil {
		// key does not change, only update value
		if s.compare(found.key, key) == 0 {
			found.value = value
		} else { // key changes, delete and re-insert
			s.delete(found.key, found.value)
			newNode = s.insertNode(key, value)
		}
	} else {
		newNode = s.insertNode(key, value)
	}

	if newNode != nil {
		s.dict[value] = newNode
	}

	return found == nil
}

// Remove delete element specified by key
// Time complexity of this method is : O(log(N))
func (s *SortedSet[K, V]) Remove(value V) *Node[K, V] {
	found := s.dict[value]
	if found != nil {
		s.delete(found.key, found.value)
		return found
	}

	return nil
}

// GetTop return top data
func (s *SortedSet[K, V]) GetTop(count int, remove bool) (result []*Node[K, V]) {
	return s.GetByRankRange(-1, -count, remove)
}

// GetRTop return top from end data
func (s *SortedSet[K, V]) GetRTop(count int, remove bool) (result []*Node[K, V]) {
	return s.GetByRankRange(1, count, remove)
}

// GetUntilKey get all values until given key
func (s *SortedSet[K, V]) GetUntilKey(untilKey K, remove bool) []any {
	nodes := s.GetByKeyRange(s.emptyKey, untilKey, &GetByKeyRangeOptions{
		Remove: remove,
	})

	data := make([]any, len(nodes))
	for i, nd := range nodes {
		data[i] = nd.Value()
	}

	return data
}

// seekPosition traverses skip list levels to find the rightmost node
// whose key is less than (or equal to, if inclusive) the given boundary.
func (s *SortedSet[K, V]) seekPosition(boundary K, inclusive bool) *Node[K, V] {
	x := s.header

	for i := s.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil {
			cmp := s.compare(x.level[i].forward.key, boundary)
			if cmp > 0 || (cmp == 0 && !inclusive) {
				break
			}

			x = x.level[i].forward
		}
	}

	return x
}

// collectForward gathers nodes walking forward from x while keys remain
// within the upper bound. excludeEnd controls whether the end boundary
// is exclusive.
func (s *SortedSet[K, V]) collectForward(x *Node[K, V], end K, excludeEnd bool, limit int, remove bool) []*Node[K, V] {
	var nodes []*Node[K, V]

	for x != nil && limit > 0 {
		cmp := s.compare(x.key, end)
		if cmp > 0 || (cmp == 0 && excludeEnd) {
			break
		}

		next := x.level[0].forward
		nodes = append(nodes, x)

		if remove {
			s.delete(x.Key(), x.Value())
		}

		limit--
		x = next
	}

	return nodes
}

// collectReverse gathers nodes walking backward from x while keys remain
// within the lower bound. excludeStart controls whether the start boundary
// is exclusive.
func (s *SortedSet[K, V]) collectReverse(
	x *Node[K, V], start K, excludeStart bool, limit int, remove bool,
) []*Node[K, V] {
	var nodes []*Node[K, V]

	for x != nil && limit > 0 {
		cmp := s.compare(x.key, start)
		if cmp < 0 || (cmp == 0 && excludeStart) {
			break
		}

		next := x.backward
		nodes = append(nodes, x)

		if remove {
			s.delete(x.Key(), x.Value())
		}

		limit--
		x = next
	}

	return nodes
}

// GetByKeyRange get the nodes whose key within the specific range
// If options is nil, it `searches` in interval [start, end] without any limit by default
// Time complexity of this method is : O(log(N))
func (s *SortedSet[K, V]) GetByKeyRange(start K, end K, options *GetByKeyRangeOptions) []*Node[K, V] {
	limit := maxLimit
	if options != nil && options.Limit > 0 {
		limit = options.Limit
	}

	var remove bool
	if options != nil {
		remove = options.Remove
	}

	excludeStart := options != nil && options.ExcludeStart
	excludeEnd := options != nil && options.ExcludeEnd

	reverse := s.compare(start, end) > 0
	if reverse {
		start, end = end, start
		excludeStart, excludeEnd = excludeEnd, excludeStart
	}

	if s.length == 0 {
		return nil
	}

	if reverse {
		x := s.seekPosition(end, !excludeEnd)

		return s.collectReverse(x, start, excludeStart, limit, remove)
	}

	x := s.seekPosition(start, excludeStart)

	return s.collectForward(x.level[0].forward, end, excludeEnd, limit, remove)
}

// GetByRankRange get nodes within specific rank range [start, end]
// Note that the rank is 1-based integer. Rank 1 means the first node; Rank -1 means the last node;
// If start is greater than end, the returned array is in reserved order
// If remove is true, the returned nodes are removed
// Time complexity of this method is : O(log(N))
func (s *SortedSet[K, V]) GetByRankRange(start int, end int, remove bool) []*Node[K, V] {
	/* Sanitize indexes. */
	if start < 0 {
		start = int(s.length) + start + 1 //nolint:gosec // length is bounded by available memory
	}

	if end < 0 {
		end = int(s.length) + end + 1 //nolint:gosec // length is bounded by available memory
	}

	if start <= 0 {
		start = 1
	}

	if end <= 0 {
		end = 1
	}

	reverse := start > end
	if reverse { // swap start and end
		start, end = end, start
	}

	var (
		update [skipListMaxLevel]*Node[K, V]
		nodes  []*Node[K, V]
	)

	traversed := 0

	x := s.header

	for i := s.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			traversed+int(x.level[i].span) < start { //nolint:gosec // span is bounded by set length
			traversed += int(x.level[i].span) //nolint:gosec // span is bounded by set length
			x = x.level[i].forward
		}

		if remove {
			update[i] = x
		} else if traversed+1 == start {
			break
		}
	}

	traversed++
	x = x.level[0].forward

	for x != nil && traversed <= end {
		next := x.level[0].forward
		nodes = append(nodes, x)

		if remove {
			s.deleteNode(x, update)
		}

		traversed++
		x = next
	}

	if reverse {
		for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
			nodes[i], nodes[j] = nodes[j], nodes[i]
		}
	}

	return nodes
}

// GetByRank get  node by rank.
// Note that the rank is 1-based integer. Rank 1 means the first node; Rank -1 means the last node;
// If remove is true, the returned nodes are removed
// If node is not found at specific rank, nil is returned
// Time complexity of this method is : O(log(N))
func (s *SortedSet[K, V]) GetByRank(rank int, remove bool) *Node[K, V] {
	nodes := s.GetByRankRange(rank, rank, remove)
	if len(nodes) == 1 {
		return nodes[0]
	}

	return nil
}

// GetByValue get node by value
// If node is not found, nil is returned
// Time complexity : O(1)
func (s *SortedSet[K, V]) GetByValue(value V) *Node[K, V] {
	n := s.dict[value]

	return n
}

// Contains reports whether the given value exists in the sorted set.
func (s *SortedSet[K, V]) Contains(value V) bool {
	return s.GetByValue(value) != nil
}

// Copy returns a deep copy of the sorted set. The new set has independent
// storage but shares the same key and value references. Element order and
// scores are preserved.
func (s *SortedSet[K, V]) Copy() *SortedSet[K, V] {
	clone := NewSortedSet[K, V](s.comparator)

	for x := s.header.level[0].forward; x != nil; x = x.level[0].forward {
		clone.Upsert(x.key, x.value)
	}

	return clone
}

// FindRank find the rank of the node specified by key
// Note that the rank is 1-based integer. Rank 1 means the first node
// If the node is not found, 0 is returned. Otherwise rank(> 0) is returned
// Time complexity of this method is : O(log(N))
func (s *SortedSet[K, V]) FindRank(value V) int {
	rank := 0
	node := s.dict[value]

	if node != nil {
		x := s.header
		for i := s.level - 1; i >= 0; i-- {
			for x.level[i].forward != nil &&
				(s.compare(x.level[i].forward.key, node.key) < 0 ||
					(s.compare(x.level[i].forward.key, node.key) == 0 &&
						x.level[i].forward.value != node.value)) {
				rank += int(x.level[i].span) //nolint:gosec // span is bounded by set length
				x = x.level[i].forward
			}

			if x.level[i].forward == node {
				rank += int(x.level[i].span) //nolint:gosec // span is bounded by set length
				return rank
			}
		}
	}

	return 0
}
