// Package linkedlist provides a generic doubly linked list with ID-based indexing.
//
// NOT safe for concurrent use. Callers must synchronize access externally.
package linkedlist

// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type List[V any] struct {
	index map[string]*Element[V] // hash-map to search Element by ID
	root Element[V] // sentinel list element, only &root, root.prev, and root.next are used
	len  int // current list length excluding (this) sentinel element
}

func (l *List[V]) init() {
	l.index = map[string]*Element[V]{}

	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
}

// Init initializes or clears the list and returns it.
func (l *List[V]) Init() *List[V] {

	l.init()

	return l
}

// New creates and initializes a new empty List.
func New[V any]() *List[V] { return new(List[V]).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List[V]) Len() int {

	return l.len
}

func (l *List[V]) front() *Element[V] {
	if l.len == 0 {
		return nil
	}

	return l.root.next
}

// Front returns the first element of the list, or nil if the list is empty.
func (l *List[V]) Front() *Element[V] {

	return l.front()
}

func (l *List[V]) back() *Element[V] {
	if l.len == 0 {
		return nil
	}

	return l.root.prev
}

// Back returns the last element of the list, or nil if the list is empty.
func (l *List[V]) Back() *Element[V] {

	return l.back()
}

func (l *List[V]) lazyInit() {
	if l.root.next == nil {
		l.init()
	}
}

func (l *List[V]) insert(e, at *Element[V]) *Element[V] {
	var id string
	if entity, ok := l.ValueToAny(e.Value).(Entity); ok {
		id = entity.ID()
	} else {
		id = l.NextID()
	}

	e.id = id

	// Remove previous value
	if prev, exists := l.index[e.id]; exists {
		l.remove(prev)
	}

	l.index[e.id] = e

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++

	return e
}

func (l *List[V]) insertValue(v V, at *Element[V]) *Element[V] {
	return l.insert(&Element[V]{Value: v}, at)
}

func (l *List[V]) remove(e *Element[V]) {
	delete(l.index, e.id)

	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
}

func (l *List[V]) move(e, at *Element[V]) {
	if e == at {
		return
	}

	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// Remove removes e from l if e belongs to the list, and returns its value.
func (l *List[V]) Remove(e *Element[V]) V {

	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}

	return e.Value
}

// PushFront inserts a new element with value v at the front of the list and returns it.
func (l *List[V]) PushFront(v V) *Element[V] {

	l.lazyInit()

	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element with value v at the back of the list and returns it.
func (l *List[V]) PushBack(v V) *Element[V] {

	l.lazyInit()

	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element with value v immediately before mark and returns it.
// If mark does not belong to l, the list is not modified and nil is returned.
func (l *List[V]) InsertBefore(v V, mark *Element[V]) *Element[V] {

	if mark.list != l {
		return nil
	}

	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element with value v immediately after mark and returns it.
// If mark does not belong to l, the list is not modified and nil is returned.
func (l *List[V]) InsertAfter(v V, mark *Element[V]) *Element[V] {

	if mark.list != l {
		return nil
	}

	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
func (l *List[V]) MoveToFront(e *Element[V]) {

	if e.list != l || l.root.next == e {
		return
	}

	l.move(e, &l.root)
}

// MoveToBack moves element e to the back of list l.
func (l *List[V]) MoveToBack(e *Element[V]) {

	if e.list != l || l.root.prev == e {
		return
	}

	l.move(e, l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
func (l *List[V]) MoveBefore(e, mark *Element[V]) {

	if e.list != l || e == mark || mark.list != l {
		return
	}

	l.move(e, mark.prev)
}

// MoveAfter moves element e to its new position after mark.
func (l *List[V]) MoveAfter(e, mark *Element[V]) {

	if e.list != l || e == mark || mark.list != l {
		return
	}

	l.move(e, mark)
}

// PushBackList inserts a copy of another list at the back of list l.
// The other list may be the same as l.
func (l *List[V]) PushBackList(other *List[V]) {
	l.lazyInit()

	for i, e := other.len, other.root.next; i > 0; i, e = i-1, e.next {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of another list at the front of list l.
// The other list may be the same as l.
func (l *List[V]) PushFrontList(other *List[V]) {
	l.lazyInit()

	for i, e := other.len, other.root.prev; i > 0; i, e = i-1, e.prev {
		l.insertValue(e.Value, &l.root)
	}
}

// ByID returns the element with the given ID, or nil if not found.
func (l *List[V]) ByID(id string) *Element[V] {

	return l.index[id]
}

// List returns all element values in front-to-back order as a slice.
func (l *List[V]) List() (result []V) {

	e := l.back()
	for e != nil && e != &l.root {
		result = append(result, e.Value)
		e = e.prev
	}

	return result
}

// Append pushes one or more values to the back of the list.
func (l *List[V]) Append(elements ...V) {

	l.lazyInit()

	for _, v := range elements {
		l.insertValue(v, l.root.prev)
	}
}

// Clear removes all elements from the list.
func (l *List[V]) Clear() {

	l.init()
}

// ValueToAny converts a value of type V to the any interface for type assertion.
func (l *List[V]) ValueToAny(v V) any {
	return v
}
