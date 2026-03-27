package linkedlist

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

type MockEntity struct {
	id string
}

func (m *MockEntity) ID() string {
	return m.id
}

func TestCustomID(t *testing.T) {
	l := New[*MockEntity]()
	test0 := l.PushFront(&MockEntity{id: "test1"})
	l.PushFront(&MockEntity{id: "test2"})
	test1 := l.PushFront(&MockEntity{id: "test2"}) // Replace existing
	test2 := l.PushFront(&MockEntity{id: "test3"})

	testutils.Equal(t, l.ByID(test0.ID()).Root().Value, &MockEntity{id: "test3"})
	testutils.Equal(t, l.ByID(test1.ID()).Root().Value, &MockEntity{id: "test3"})
	testutils.Equal(t, l.ByID(test2.ID()).Root().Value, &MockEntity{id: "test3"})
	testutils.Equal(t, l.ByID(test2.ID()).Root().Next().Value, &MockEntity{id: "test2"})

	testutils.Equal(t, l.ByID(test2.ID()).Value, &MockEntity{id: "test3"})
	testutils.Equal(t, l.ByID(test2.ID()).Next().Value, &MockEntity{id: "test2"})
	testutils.Equal(t, l.ByID(test2.ID()).Next().Next().Value, &MockEntity{id: "test1"})
}

// Prevent compiler optimization of benchmark results.
var benchResult any

func benchmarkSizes() []struct {
	name string
	size int
} {
	return []struct {
		name string
		size int
	}{
		{"n=100", 100},
		{"n=1000", 1000},
		{"n=10000", 10000},
	}
}

func BenchmarkPushBack(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				l := New[int]()
				for i := range tc.size {
					l.PushBack(i)
				}
			}
		})
	}
}

func BenchmarkPushFront(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				l := New[int]()
				for i := range tc.size {
					l.PushFront(i)
				}
			}
		})
	}
}

func BenchmarkByID(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			l := New[int]()
			ids := make([]string, tc.size)
			for i := range tc.size {
				e := l.PushBack(i)
				ids[i] = e.ID()
			}

			b.ResetTimer()

			for i := range b.N {
				benchResult = l.ByID(ids[i%tc.size])
			}
		})
	}
}

func BenchmarkRemove(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			for range b.N {
				b.StopTimer()

				l := New[int]()
				elements := make([]*Element[int], tc.size)
				for i := range tc.size {
					elements[i] = l.PushBack(i)
				}

				b.StartTimer()

				for _, e := range elements {
					l.Remove(e)
				}
			}
		})
	}
}

func BenchmarkList(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			l := New[int]()
			for i := range tc.size {
				l.PushBack(i)
			}

			b.ResetTimer()

			for range b.N {
				benchResult = l.List()
			}
		})
	}
}

func BenchmarkMoveToFront(b *testing.B) {
	l := New[int]()
	elements := make([]*Element[int], 1000)
	for i := range 1000 {
		elements[i] = l.PushBack(i)
	}

	b.ResetTimer()

	for i := range b.N {
		l.MoveToFront(elements[i%1000])
	}
}

func TestListCopy(t *testing.T) {
	l := New[string]()
	l.PushFront("test")
	test1 := l.PushFront("test1")
	test2 := l.PushFront("test2")
	test3 := l.PushBack("test3")
	l.PushBack("test4")

	testutils.Equal(t, l.List(), []string{"test4", "test3", "test", "test1", "test2"})
	testutils.Equal(t, l.Front().Next().Value, "test1")
	testutils.Equal(t, l.Back().Prev().Value, "test3")

	l2 := New[string]()
	l2.Append(l.List()...) // copy
	testutils.Equal(t, l2.List(), []string{"test2", "test1", "test", "test3", "test4"})
	testutils.Equal(t, l2.Front().Next().Value, "test3")
	testutils.Equal(t, l2.Back().Prev().Value, "test1")

	testutils.Equal(t, l.ByID(test2.ID()).Value, "test2")
	testutils.Equal(t, l.ByID(test2.ID()).Next().Value, "test1")

	testutils.Equal(t, l.ByID(test1.ID()).Value, "test1")
	testutils.Equal(t, l.ByID(test1.ID()).Next().Value, "test")
	testutils.Equal(t, l.ByID(test1.ID()).Prev().Value, "test2")

	testutils.Equal(t, l.ByID(test3.ID()).Value, "test3")
	testutils.Equal(t, l.ByID(test3.ID()).Next().Value, "test4")
	testutils.Equal(t, l.ByID(test3.ID()).Prev().Value, "test")
}
