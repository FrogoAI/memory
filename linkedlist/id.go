package linkedlist

import "github.com/google/uuid"

// NextID generates a new unique identifier for an element.
func (l *List[V]) NextID() string {
	return uuid.New().String()
}
