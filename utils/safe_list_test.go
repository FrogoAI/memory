//nolint:revive // package name is intentional
package utils

import (
	"sync"
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestSafeListNewEmpty(t *testing.T) {
	sl := NewSafeList[int]()
	testutils.Equal(t, sl.Count(), 0)
	testutils.Equal(t, sl.List(), []int{})
}

func TestSafeListNewWithData(t *testing.T) {
	cases := []struct {
		name  string
		input []int
		count int
	}{
		{name: "single_element", input: []int{42}, count: 1},
		{name: "multiple_elements", input: []int{1, 2, 3}, count: 3},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSafeList(tc.input...)
			testutils.Equal(t, sl.Count(), tc.count)
			testutils.Equal(t, sl.List(), tc.input)
		})
	}
}

func TestSafeListAdd(t *testing.T) {
	sl := NewSafeList[string]()

	sl.Add("a")
	testutils.Equal(t, sl.Count(), 1)
	testutils.Equal(t, sl.List(), []string{"a"})

	sl.Add("b")
	sl.Add("c")
	testutils.Equal(t, sl.Count(), 3)
	testutils.Equal(t, sl.List(), []string{"a", "b", "c"})
}

func TestSafeListListReturnsCopy(t *testing.T) {
	sl := NewSafeList(1, 2, 3)

	result := sl.List()
	result[0] = 999

	// original should be unchanged
	testutils.Equal(t, sl.List(), []int{1, 2, 3})
}

func TestSafeListClear(t *testing.T) {
	sl := NewSafeList(1, 2, 3)

	cleared := sl.Clear()
	testutils.Equal(t, cleared, []int{1, 2, 3})
	testutils.Equal(t, sl.Count(), 0)
	testutils.Equal(t, sl.List(), []int{})
}

func TestSafeListClearEmpty(t *testing.T) {
	sl := NewSafeList[int]()

	cleared := sl.Clear()
	testutils.Equal(t, cleared, []int{})
	testutils.Equal(t, sl.Count(), 0)
}

func TestSafeListClearThenAdd(t *testing.T) {
	sl := NewSafeList(1, 2, 3)
	sl.Clear()

	sl.Add(10)
	testutils.Equal(t, sl.Count(), 1)
	testutils.Equal(t, sl.List(), []int{10})
}

func TestSafeListCount(t *testing.T) {
	cases := []struct {
		name  string
		items []int
		want  int
	}{
		{name: "empty", items: nil, want: 0},
		{name: "one", items: []int{1}, want: 1},
		{name: "many", items: []int{1, 2, 3, 4, 5}, want: 5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sl := NewSafeList(tc.items...)
			testutils.Equal(t, sl.Count(), tc.want)
		})
	}
}

func TestSafeListConcurrentAdd(t *testing.T) {
	sl := NewSafeList[int]()

	const goroutines = 50

	const opsPerGoroutine = 100

	var wg sync.WaitGroup

	wg.Add(goroutines)

	for i := range goroutines {
		go func(id int) {
			defer wg.Done()

			for j := range opsPerGoroutine {
				sl.Add(id*opsPerGoroutine + j)
			}
		}(i)
	}

	wg.Wait()

	testutils.Equal(t, sl.Count(), goroutines*opsPerGoroutine)
}

func TestSafeListConcurrentReadWrite(t *testing.T) {
	sl := NewSafeList[int]()

	const writers = 10

	const readers = 20

	const ops = 100

	var wg sync.WaitGroup

	wg.Add(writers + readers)

	for i := range writers {
		go func(id int) {
			defer wg.Done()

			for j := range ops {
				sl.Add(id*ops + j)
			}
		}(i)
	}

	for range readers {
		go func() {
			defer wg.Done()

			for range ops {
				sl.List()
				sl.Count()
			}
		}()
	}

	wg.Wait()

	testutils.Equal(t, sl.Count(), writers*ops)
}

func TestSafeListConcurrentClear(t *testing.T) {
	sl := NewSafeList[int]()

	const goroutines = 20

	var wg sync.WaitGroup

	wg.Add(goroutines * 2)

	for i := range goroutines {
		go func(id int) {
			defer wg.Done()

			sl.Add(id)
		}(i)

		go func() {
			defer wg.Done()

			sl.Clear()
		}()
	}

	wg.Wait()
}
