package linkedlist

// Entity is implemented by values that provide their own identity. When a value
// stored in the list satisfies Entity, its ID is used as the element key
// instead of generating a random one.
type Entity interface {
	ID() string
}

// Element is a node in the doubly linked list, holding a value and pointers
// to adjacent elements.
type Element[V any] struct {
	// The value stored with this element.
	Value V

	next *Element[V]
	prev *Element[V]

	// The list to which this element belongs.
	list *List[V]

	id string // auto-assigned ID
}

// Next returns the next list element or an empty Element if e is at the back.
func (e *Element[V]) Next() *Element[V] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}

	return &Element[V]{}
}

// Prev returns the previous list element or an empty Element if e is at the front.
func (e *Element[V]) Prev() *Element[V] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}

	return &Element[V]{}
}

// Root returns the first element of the list that e belongs to.
func (e *Element[V]) Root() *Element[V] {
	return e.list.root.next
}

// IsEmpty reports whether e is an empty (sentinel) element with no assigned ID.
func (e *Element[V]) IsEmpty() bool {
	return e.id == ""
}

// ID returns the unique identifier assigned to this element.
func (e *Element[V]) ID() string {
	return e.id
}
