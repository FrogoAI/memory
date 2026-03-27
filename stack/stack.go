// Package stack provides a generic LIFO stack backed by a slice.
//
// NOT safe for concurrent use. Callers must synchronize access externally.
package stack

// Stack is a generic LIFO stack backed by a slice.
type Stack[T any] struct {
	items []T
}

// ToSlice returns a copy of the stack's elements as a slice.
func (stack *Stack[T]) ToSlice() []T {

	dst := make([]T, len(stack.items))
	copy(dst, stack.items)

	return dst
}

// Set replaces the stack's contents with the provided slice.
func (stack *Stack[T]) Set(values []T) {

	stack.items = values
}

// IsEmpty reports whether the stack contains no elements.
func (stack *Stack[T]) IsEmpty() bool {

	return len(stack.items) == 0
}

// Peek returns the top element without removing it.
// If the stack is empty, it returns the zero value.
func (stack *Stack[T]) Peek() T {

	if len(stack.items) <= 0 {
		var defaultValue T
		return defaultValue
	}

	return stack.items[len(stack.items)-1]
}

// Reverse reverses the order of elements in place and returns the stack.
func (stack *Stack[T]) Reverse() *Stack[T] {

	for i, j := 0, len(stack.items)-1; i < j; i, j = i+1, j-1 {
		stack.items[i], stack.items[j] = stack.items[j], stack.items[i]
	}

	return stack
}

// Push adds a value to the top of the stack.
func (stack *Stack[T]) Push(value T) {

	stack.items = append(stack.items, value)
}

// PushLeft adds a value to the bottom of the stack.
func (stack *Stack[T]) PushLeft(value T) {

	stack.items = append([]T{value}, stack.items...)
}

// Pop removes and returns the top element. If the stack is empty, it returns
// the zero value.
func (stack *Stack[T]) Pop() T {

	n := len(stack.items)
	if n <= 0 {
		var defaultValue T
		return defaultValue
	}

	p := stack.items[n-1]
	stack.items = stack.items[:n-1]

	return p
}

// PopLeft removes and returns the bottom element. If the stack is empty, it
// returns the zero value.
func (stack *Stack[T]) PopLeft() T {

	var item T

	n := len(stack.items)
	if n <= 0 {
		return item
	}

	item = stack.items[0]
	stack.items = stack.items[1:]

	return item
}

// Len returns the number of elements in the stack.
func (stack *Stack[T]) Len() int {

	return len(stack.items)
}

// Clear removes all items from the stack.
func (stack *Stack[T]) Clear() {

	stack.items = nil
}
