package linkedlist

import "github.com/google/uuid"

func (l *List[V]) NextID() string {
	return uuid.New().String()
}
