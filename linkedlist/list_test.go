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

func TestEmptyListOps(t *testing.T) {
	cases := []struct {
		name string
		fn   func(*List[int]) any
		want any
	}{
		{"Front returns nil", func(l *List[int]) any { return l.Front() }, (*Element[int])(nil)},
		{"Back returns nil", func(l *List[int]) any { return l.Back() }, (*Element[int])(nil)},
		{"Len returns 0", func(l *List[int]) any { return l.Len() }, 0},
		{"List returns nil slice", func(l *List[int]) any { return l.List() }, ([]int)(nil)},
		{"ByID returns nil", func(l *List[int]) any { return l.ByID("nonexistent") }, (*Element[int])(nil)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := New[int]()
			testutils.Equal(t, tc.fn(l), tc.want)
		})
	}
}

func TestPushFrontPushBack(t *testing.T) {
	cases := []struct {
		name     string
		ops      func(*List[int])
		wantList []int
		wantLen  int
	}{
		{
			"PushFront single",
			func(l *List[int]) { l.PushFront(1) },
			[]int{1},
			1,
		},
		{
			"PushBack single",
			func(l *List[int]) { l.PushBack(1) },
			[]int{1},
			1,
		},
		{
			"PushFront ordering",
			func(l *List[int]) {
				l.PushFront(1)
				l.PushFront(2)
				l.PushFront(3)
			},
			[]int{1, 2, 3}, // List() is back-to-front; front is 3, back is 1
			3,
		},
		{
			"PushBack ordering",
			func(l *List[int]) {
				l.PushBack(1)
				l.PushBack(2)
				l.PushBack(3)
			},
			[]int{3, 2, 1}, // List() is back-to-front; front is 1, back is 3
			3,
		},
		{
			"mixed PushFront and PushBack",
			func(l *List[int]) {
				l.PushBack(2)
				l.PushFront(1)
				l.PushBack(3)
			},
			[]int{3, 2, 1}, // front: 1->2->3 :back; List() reverses
			3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := New[int]()
			tc.ops(l)
			testutils.Equal(t, l.Len(), tc.wantLen)
			testutils.Equal(t, l.List(), tc.wantList)
		})
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		name      string
		setup     func(*List[int]) []*Element[int]
		removeIdx []int
		wantList  []int
		wantLen   int
	}{
		{
			"remove middle element",
			func(l *List[int]) []*Element[int] {
				return []*Element[int]{l.PushBack(1), l.PushBack(2), l.PushBack(3)}
			},
			[]int{1},
			[]int{3, 1}, // back-to-front: 3, 1
			2,
		},
		{
			"remove front element",
			func(l *List[int]) []*Element[int] {
				return []*Element[int]{l.PushBack(1), l.PushBack(2), l.PushBack(3)}
			},
			[]int{0},
			[]int{3, 2}, // back-to-front: 3, 2
			2,
		},
		{
			"remove back element",
			func(l *List[int]) []*Element[int] {
				return []*Element[int]{l.PushBack(1), l.PushBack(2), l.PushBack(3)}
			},
			[]int{2},
			[]int{2, 1}, // back-to-front: 2, 1
			2,
		},
		{
			"remove all elements",
			func(l *List[int]) []*Element[int] {
				return []*Element[int]{l.PushBack(1), l.PushBack(2)}
			},
			[]int{0, 1},
			([]int)(nil),
			0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := New[int]()
			elems := tc.setup(l)

			for _, idx := range tc.removeIdx {
				l.Remove(elems[idx])
			}

			testutils.Equal(t, l.Len(), tc.wantLen)
			testutils.Equal(t, l.List(), tc.wantList)
		})
	}
}

func TestRemoveElementFromWrongList(t *testing.T) {
	l1 := New[int]()
	l2 := New[int]()
	e := l1.PushBack(1)
	l2.PushBack(2)

	// Remove element from wrong list — should be a no-op
	val := l2.Remove(e)
	testutils.Equal(t, val, 1) // returns the value regardless
	testutils.Equal(t, l1.Len(), 1)
	testutils.Equal(t, l2.Len(), 1)
}

func TestByIDNotFound(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)

	testutils.Equal(t, l.ByID("nonexistent"), (*Element[int])(nil))
}

func TestAppend(t *testing.T) {
	cases := []struct {
		name     string
		values   []int
		wantList []int
	}{
		{"single value", []int{1}, []int{1}},
		{"multiple values", []int{1, 2, 3}, []int{3, 2, 1}},
		{"no values", []int{}, ([]int)(nil)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := New[int]()
			l.Append(tc.values...)
			testutils.Equal(t, l.List(), tc.wantList)
		})
	}
}

func TestListOrdering(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	testutils.Equal(t, l.Front().Value, 1)
	testutils.Equal(t, l.Back().Value, 3)
	testutils.Equal(t, l.List(), []int{3, 2, 1}) // List() returns back-to-front

	// Traverse front to back via Next()
	e := l.Front()
	testutils.Equal(t, e.Value, 1)
	e = e.Next()
	testutils.Equal(t, e.Value, 2)
	e = e.Next()
	testutils.Equal(t, e.Value, 3)
	testutils.Equal(t, e.Next().IsEmpty(), true) // past the end
}

func TestClear(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	testutils.Equal(t, l.Len(), 3)

	l.Clear()

	testutils.Equal(t, l.Len(), 0)
	testutils.Equal(t, l.Front(), (*Element[int])(nil))
	testutils.Equal(t, l.Back(), (*Element[int])(nil))
	testutils.Equal(t, l.List(), ([]int)(nil))

	// Can reuse after clear
	l.PushBack(10)
	testutils.Equal(t, l.Len(), 1)
	testutils.Equal(t, l.Front().Value, 10)
}

func TestInsertBeforeAfter(t *testing.T) {
	cases := []struct {
		name     string
		ops      func(*List[int])
		wantList []int
	}{
		{
			"InsertBefore middle",
			func(l *List[int]) {
				l.PushBack(1)
				mark := l.PushBack(3)
				l.InsertBefore(2, mark)
			},
			[]int{3, 2, 1}, // front: 1->2->3 :back; List() back-to-front
		},
		{
			"InsertAfter middle",
			func(l *List[int]) {
				mark := l.PushBack(1)
				l.PushBack(3)
				l.InsertAfter(2, mark)
			},
			[]int{3, 2, 1}, // front: 1->2->3 :back; List() back-to-front
		},
		{
			"InsertBefore front",
			func(l *List[int]) {
				mark := l.PushBack(2)
				l.InsertBefore(1, mark)
			},
			[]int{2, 1}, // front: 1->2 :back; List() back-to-front
		},
		{
			"InsertAfter back",
			func(l *List[int]) {
				mark := l.PushBack(1)
				l.InsertAfter(2, mark)
			},
			[]int{2, 1}, // front: 1->2 :back; List() back-to-front
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := New[int]()
			tc.ops(l)
			testutils.Equal(t, l.List(), tc.wantList)
		})
	}
}

func TestInsertBeforeAfterWrongList(t *testing.T) {
	l1 := New[int]()
	l2 := New[int]()

	mark := l1.PushBack(1)

	testutils.Equal(t, l2.InsertBefore(0, mark), (*Element[int])(nil))
	testutils.Equal(t, l2.InsertAfter(0, mark), (*Element[int])(nil))
}

func TestMoveToFrontBack(t *testing.T) {
	cases := []struct {
		name     string
		ops      func(*List[int])
		wantList []int
	}{
		{
			"MoveToFront from back",
			func(l *List[int]) {
				l.PushBack(1)
				l.PushBack(2)
				e := l.PushBack(3)
				l.MoveToFront(e)
			},
			[]int{2, 1, 3}, // front: 3->1->2 :back; List() back-to-front
		},
		{
			"MoveToBack from front",
			func(l *List[int]) {
				e := l.PushBack(1)
				l.PushBack(2)
				l.PushBack(3)
				l.MoveToBack(e)
			},
			[]int{1, 3, 2}, // front: 2->3->1 :back; List() back-to-front
		},
		{
			"MoveToFront already at front is no-op",
			func(l *List[int]) {
				e := l.PushBack(1)
				l.PushBack(2)
				l.MoveToFront(e)
			},
			[]int{2, 1}, // front: 1->2 :back; List() back-to-front
		},
		{
			"MoveToBack already at back is no-op",
			func(l *List[int]) {
				l.PushBack(1)
				e := l.PushBack(2)
				l.MoveToBack(e)
			},
			[]int{2, 1}, // front: 1->2 :back; List() back-to-front
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := New[int]()
			tc.ops(l)
			testutils.Equal(t, l.List(), tc.wantList)
		})
	}
}

func TestMoveToFrontBackWrongList(t *testing.T) {
	l1 := New[int]()
	l2 := New[int]()

	e := l1.PushBack(1)
	l2.PushBack(2)

	// Should be no-ops
	l2.MoveToFront(e)
	l2.MoveToBack(e)
	testutils.Equal(t, l1.Len(), 1)
	testutils.Equal(t, l2.Len(), 1)
}

func TestMoveBeforeAfter(t *testing.T) {
	cases := []struct {
		name     string
		ops      func(*List[int])
		wantList []int
	}{
		{
			"MoveBefore",
			func(l *List[int]) {
				l.PushBack(1)
				mark := l.PushBack(2)
				e := l.PushBack(3)
				l.MoveBefore(e, mark)
			},
			[]int{2, 3, 1}, // front: 1->3->2 :back; List() back-to-front
		},
		{
			"MoveAfter",
			func(l *List[int]) {
				e := l.PushBack(1)
				mark := l.PushBack(2)
				l.PushBack(3)
				l.MoveAfter(e, mark)
			},
			[]int{3, 1, 2}, // front: 2->1->3 :back; List() back-to-front
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := New[int]()
			tc.ops(l)
			testutils.Equal(t, l.List(), tc.wantList)
		})
	}
}

func TestMoveBeforeAfterEdgeCases(t *testing.T) {
	l1 := New[int]()
	l2 := New[int]()

	e := l1.PushBack(1)
	mark := l1.PushBack(2)
	l2.PushBack(3)
	wrongE := l2.PushBack(4)

	// Same element — no-op
	l1.MoveBefore(e, e)
	l1.MoveAfter(e, e)
	testutils.Equal(t, l1.List(), []int{2, 1})

	// Element from wrong list — no-op
	l1.MoveBefore(wrongE, mark)
	l1.MoveAfter(wrongE, mark)
	testutils.Equal(t, l1.Len(), 2)
	testutils.Equal(t, l2.Len(), 2)

	// Mark from wrong list — no-op
	l1.MoveBefore(e, wrongE)
	l1.MoveAfter(e, wrongE)
	testutils.Equal(t, l1.Len(), 2)
}

func TestPushBackList(t *testing.T) {
	cases := []struct {
		name     string
		l1Vals   []int
		l2Vals   []int
		wantList []int
	}{
		{"both non-empty", []int{1, 2}, []int{3, 4}, []int{4, 3, 2, 1}},
		{"source empty", []int{1, 2}, []int{}, []int{2, 1}},
		{"dest empty", []int{}, []int{3, 4}, []int{4, 3}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l1 := New[int]()
			for _, v := range tc.l1Vals {
				l1.PushBack(v)
			}

			l2 := New[int]()
			for _, v := range tc.l2Vals {
				l2.PushBack(v)
			}

			l1.PushBackList(l2)
			testutils.Equal(t, l1.List(), tc.wantList)
		})
	}
}

func TestPushFrontList(t *testing.T) {
	cases := []struct {
		name     string
		l1Vals   []int
		l2Vals   []int
		wantList []int
	}{
		{"both non-empty", []int{3, 4}, []int{1, 2}, []int{4, 3, 2, 1}},
		{"source empty", []int{1, 2}, []int{}, []int{2, 1}},
		{"dest empty", []int{}, []int{1, 2}, []int{2, 1}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l1 := New[int]()
			for _, v := range tc.l1Vals {
				l1.PushBack(v)
			}

			l2 := New[int]()
			for _, v := range tc.l2Vals {
				l2.PushBack(v)
			}

			l1.PushFrontList(l2)
			testutils.Equal(t, l1.List(), tc.wantList)
		})
	}
}

func TestElementIsEmpty(t *testing.T) {
	cases := []struct {
		name string
		elem *Element[int]
		want bool
	}{
		{"empty element", &Element[int]{}, true},
		{"element from list", nil, false}, // set below
	}

	l := New[int]()
	e := l.PushBack(42)
	cases[1].elem = e

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, tc.elem.IsEmpty(), tc.want)
		})
	}
}

func TestElementPrevAtFront(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)

	front := l.Front()
	prev := front.Prev()
	testutils.Equal(t, prev.IsEmpty(), true)
}

func TestLazyInit(t *testing.T) {
	// Use zero-value List (not via New)
	var l List[int]
	l.PushBack(1)
	l.PushBack(2)

	testutils.Equal(t, l.Len(), 2)
	testutils.Equal(t, l.List(), []int{2, 1})
}

func TestInitClearsList(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)
	testutils.Equal(t, l.Len(), 2)

	l.Init()
	testutils.Equal(t, l.Len(), 0)
	testutils.Equal(t, l.Front(), (*Element[int])(nil))
}

func TestCopy(t *testing.T) {
	cases := []struct {
		name   string
		values []int
	}{
		{
			name:   "empty list",
			values: nil,
		},
		{
			name:   "single element",
			values: []int{1},
		},
		{
			name:   "multiple elements",
			values: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := New[int]()
			for _, v := range tc.values {
				l.PushBack(v)
			}

			clone := l.Copy()

			testutils.Equal(t, clone.Len(), l.Len())
			testutils.Equal(t, clone.List(), l.List())

			// Verify independence: modify original, clone unaffected
			l.PushBack(99)
			testutils.NotEqual(t, clone.Len(), l.Len())
		})
	}
}
