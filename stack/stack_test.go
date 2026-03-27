package stack

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestPush(t *testing.T) {
	cases := []struct {
		name   string
		values []int
		want   []int
	}{
		{
			name:   "push single element",
			values: []int{1},
			want:   []int{1},
		},
		{
			name:   "push multiple elements",
			values: []int{1, 2, 3},
			want:   []int{1, 2, 3},
		},
		{
			name:   "push to empty stack",
			values: []int{},
			want:   []int{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s Stack[int]

			for _, v := range tc.values {
				s.Push(v)
			}

			testutils.Equal(t, s.Len(), len(tc.want))
			testutils.Equal(t, s.ToSlice(), tc.want)
		})
	}
}

func TestPushLeft(t *testing.T) {
	cases := []struct {
		name   string
		values []int
		want   []int
	}{
		{
			name:   "push left single element",
			values: []int{1},
			want:   []int{1},
		},
		{
			name:   "push left multiple prepends in order",
			values: []int{1, 2, 3},
			want:   []int{3, 2, 1},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s Stack[int]

			for _, v := range tc.values {
				s.PushLeft(v)
			}

			testutils.Equal(t, s.ToSlice(), tc.want)
		})
	}
}

func TestPop(t *testing.T) {
	cases := []struct {
		name     string
		setup    []int
		popCount int
		wantPops []int
		wantLen  int
	}{
		{
			name:     "pop from empty stack returns zero value",
			setup:    nil,
			popCount: 1,
			wantPops: []int{0},
			wantLen:  0,
		},
		{
			name:     "pop single element",
			setup:    []int{42},
			popCount: 1,
			wantPops: []int{42},
			wantLen:  0,
		},
		{
			name:     "pop returns LIFO order",
			setup:    []int{1, 2, 3},
			popCount: 3,
			wantPops: []int{3, 2, 1},
			wantLen:  0,
		},
		{
			name:     "pop partial",
			setup:    []int{10, 20, 30},
			popCount: 1,
			wantPops: []int{30},
			wantLen:  2,
		},
		{
			name:     "pop beyond empty returns zero values",
			setup:    []int{1},
			popCount: 3,
			wantPops: []int{1, 0, 0},
			wantLen:  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s Stack[int]
			s.Set(tc.setup)

			for i := range tc.popCount {
				got := s.Pop()
				testutils.Equal(t, got, tc.wantPops[i])
			}

			testutils.Equal(t, s.Len(), tc.wantLen)
		})
	}
}

func TestPopLeft(t *testing.T) {
	cases := []struct {
		name     string
		setup    []string
		popCount int
		wantPops []string
		wantLen  int
	}{
		{
			name:     "pop left from empty stack returns zero value",
			setup:    nil,
			popCount: 1,
			wantPops: []string{""},
			wantLen:  0,
		},
		{
			name:     "pop left returns FIFO order",
			setup:    []string{"a", "b", "c"},
			popCount: 3,
			wantPops: []string{"a", "b", "c"},
			wantLen:  0,
		},
		{
			name:     "pop left partial",
			setup:    []string{"x", "y", "z"},
			popCount: 1,
			wantPops: []string{"x"},
			wantLen:  2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s Stack[string]
			s.Set(tc.setup)

			for i := range tc.popCount {
				got := s.PopLeft()
				testutils.Equal(t, got, tc.wantPops[i])
			}

			testutils.Equal(t, s.Len(), tc.wantLen)
		})
	}
}

func TestPeek(t *testing.T) {
	cases := []struct {
		name  string
		setup []int
		want  int
	}{
		{
			name:  "peek empty stack returns zero value",
			setup: nil,
			want:  0,
		},
		{
			name:  "peek single element",
			setup: []int{5},
			want:  5,
		},
		{
			name:  "peek returns top element",
			setup: []int{1, 2, 3},
			want:  3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s Stack[int]
			s.Set(tc.setup)

			got := s.Peek()
			testutils.Equal(t, got, tc.want)
			testutils.Equal(t, s.Len(), len(tc.setup))
		})
	}
}

func TestPeekDoesNotRemove(t *testing.T) {
	var s Stack[int]
	s.Push(99)

	for range 3 {
		testutils.Equal(t, s.Peek(), 99)
	}

	testutils.Equal(t, s.Len(), 1)
}

func TestLen(t *testing.T) {
	cases := []struct {
		name  string
		setup []int
		want  int
	}{
		{
			name:  "empty stack",
			setup: nil,
			want:  0,
		},
		{
			name:  "single element",
			setup: []int{1},
			want:  1,
		},
		{
			name:  "multiple elements",
			setup: []int{1, 2, 3, 4, 5},
			want:  5,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s Stack[int]
			s.Set(tc.setup)

			testutils.Equal(t, s.Len(), tc.want)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	cases := []struct {
		name  string
		setup []int
		want  bool
	}{
		{
			name:  "zero value stack is empty",
			setup: nil,
			want:  true,
		},
		{
			name:  "non-empty stack",
			setup: []int{1},
			want:  false,
		},
		{
			name:  "empty after pop all",
			setup: nil,
			want:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s Stack[int]
			s.Set(tc.setup)

			testutils.Equal(t, s.IsEmpty(), tc.want)
		})
	}
}

func TestIsEmptyAfterPopAll(t *testing.T) {
	var s Stack[int]
	s.Push(1)
	s.Push(2)
	s.Pop()
	s.Pop()

	testutils.Equal(t, s.IsEmpty(), true)
}

func TestToSlice(t *testing.T) {
	cases := []struct {
		name  string
		setup []int
		want  []int
	}{
		{
			name:  "empty stack returns empty slice",
			setup: nil,
			want:  []int{},
		},
		{
			name:  "returns copy of elements",
			setup: []int{1, 2, 3},
			want:  []int{1, 2, 3},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s Stack[int]
			s.Set(tc.setup)

			got := s.ToSlice()
			testutils.Equal(t, got, tc.want)
		})
	}
}

func TestToSliceReturnsCopy(t *testing.T) {
	var s Stack[int]
	s.Set([]int{1, 2, 3})

	got := s.ToSlice()
	got[0] = 999

	testutils.Equal(t, s.Peek(), 3)

	original := s.ToSlice()
	testutils.Equal(t, original[0], 1)
}

func TestSet(t *testing.T) {
	var s Stack[int]

	s.Set([]int{10, 20, 30})
	testutils.Equal(t, s.Len(), 3)
	testutils.Equal(t, s.Peek(), 30)

	s.Set([]int{99})
	testutils.Equal(t, s.Len(), 1)
	testutils.Equal(t, s.Peek(), 99)

	s.Set(nil)
	testutils.Equal(t, s.IsEmpty(), true)
}

func TestReverse(t *testing.T) {
	cases := []struct {
		name  string
		setup []int
		want  []int
	}{
		{
			name:  "empty stack",
			setup: nil,
			want:  []int{},
		},
		{
			name:  "single element",
			setup: []int{1},
			want:  []int{1},
		},
		{
			name:  "odd count",
			setup: []int{1, 2, 3},
			want:  []int{3, 2, 1},
		},
		{
			name:  "even count",
			setup: []int{1, 2, 3, 4},
			want:  []int{4, 3, 2, 1},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s Stack[int]
			s.Set(tc.setup)

			result := s.Reverse()
			testutils.Equal(t, s.ToSlice(), tc.want)
			testutils.Equal(t, result, &s)
		})
	}
}

func TestClear(t *testing.T) {
	var s Stack[int]
	s.Set([]int{1, 2, 3})

	s.Clear()
	testutils.Equal(t, s.IsEmpty(), true)
	testutils.Equal(t, s.Len(), 0)
}

func TestClearEmpty(t *testing.T) {
	var s Stack[int]

	s.Clear()
	testutils.Equal(t, s.IsEmpty(), true)
}

func TestLargeStack(t *testing.T) {
	const size = 10000

	var s Stack[int]

	for i := range size {
		s.Push(i)
	}

	testutils.Equal(t, s.Len(), size)
	testutils.Equal(t, s.Peek(), size-1)

	for i := size - 1; i >= 0; i-- {
		got := s.Pop()
		testutils.Equal(t, got, i)
	}

	testutils.Equal(t, s.IsEmpty(), true)
}

func TestZeroValueStack(t *testing.T) {
	var s Stack[int]

	testutils.Equal(t, s.IsEmpty(), true)
	testutils.Equal(t, s.Len(), 0)
	testutils.Equal(t, s.Peek(), 0)
	testutils.Equal(t, s.Pop(), 0)
	testutils.Equal(t, s.PopLeft(), 0)
}

func TestStringStack(t *testing.T) {
	var s Stack[string]

	s.Push("hello")
	s.Push("world")

	testutils.Equal(t, s.Pop(), "world")
	testutils.Equal(t, s.Pop(), "hello")
	testutils.Equal(t, s.Pop(), "")
}
