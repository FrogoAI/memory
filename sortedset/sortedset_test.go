package sortedset

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/FrogoAI/testutils"
	"github.com/spf13/cast"

	"github.com/FrogoAI/memory/comparator"
)

func TestTimeQueue(t *testing.T) {
	set := NewSortedSet[time.Time, any](comparator.TimeComparator)
	for i := -1000; i <= 10000; i++ {
		set.Upsert(time.Now().Add(time.Duration(i)*time.Minute), i)
	}

	set.Upsert(time.Now().Add(time.Minute), "test4")
	set.Upsert(time.Now().Add(time.Second), "test3")
	set.Upsert(time.Now().Add(-time.Minute), "test1")
	set.Upsert(time.Now().Add(-time.Hour), "test0")
	set.Upsert(time.Now(), "test2")

	set.GetUntilKey(time.Now().Add(-5*time.Minute), true)
	values := set.GetUntilKey(time.Now(), true)
	expectedValues := []interface{}{
		-4, -3, -2, -1, "test1", 0, "test2",
	}

	testutils.Equal(t, values, expectedValues)
}

func TestDump(t *testing.T) {
	set := NewSortedSet[time.Time, any](comparator.TimeComparator)
	set.Upsert(time.Now().Add(time.Minute), "test4")
	set.Upsert(time.Now().Add(time.Second), "test3")
	set.Upsert(time.Now().Add(-time.Minute), "test1")
	set.Upsert(time.Now().Add(-time.Hour), "test0")
	set.Upsert(time.Now(), "test2")

	dump, err := set.Dump(func(key time.Time, values any) (string, string, error) {
		return key.Format(time.RFC3339), cast.ToString(values), nil
	})
	testutils.Equal(t, err, nil)

	set = NewSortedSet[time.Time, any](comparator.TimeComparator)
	err = set.Restore(func(key string, values []string) (time.Time, []any, error) {
		rValues := make([]any, len(values))
		for i, v := range values {
			rValues[i] = v
		}

		rKey, err := time.Parse(time.RFC3339, key)

		return rKey, rValues, err
	}, dump)
	testutils.Equal(t, err, nil)

	set.GetUntilKey(time.Now().Add(-5*time.Minute), true)
	values := set.GetUntilKey(time.Now(), true)
	expectedValues := []interface{}{
		"test1", "test2",
	}

	testutils.Equal(t, values, expectedValues)
}

func newIntSet(pairs ...any) *SortedSet[int, string] {
	s := NewSortedSet[int, string](comparator.IntComparator)
	for i := 0; i < len(pairs); i += 2 {
		s.Upsert(pairs[i].(int), pairs[i+1].(string))
	}

	return s
}

func TestPeekMin(t *testing.T) {
	cases := []struct {
		name    string
		set     *SortedSet[int, string]
		wantNil bool
		wantVal string
	}{
		{"empty set", newIntSet(), true, ""},
		{"single element", newIntSet(10, "a"), false, "a"},
		{"returns smallest key", newIntSet(30, "c", 10, "a", 20, "b"), false, "a"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := tc.set.PeekMin()
			if tc.wantNil {
				testutils.Equal(t, node == nil, true)

				return
			}

			testutils.Equal(t, node.Value(), tc.wantVal)
		})
	}
}

func TestPeekMax(t *testing.T) {
	cases := []struct {
		name    string
		set     *SortedSet[int, string]
		wantNil bool
		wantVal string
	}{
		{"empty set", newIntSet(), true, ""},
		{"single element", newIntSet(10, "a"), false, "a"},
		{"returns largest key", newIntSet(10, "a", 30, "c", 20, "b"), false, "c"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := tc.set.PeekMax()
			if tc.wantNil {
				testutils.Equal(t, node == nil, true)

				return
			}

			testutils.Equal(t, node.Value(), tc.wantVal)
		})
	}
}

func TestPopMin(t *testing.T) {
	cases := []struct {
		name    string
		set     *SortedSet[int, string]
		wantNil bool
		wantVal string
		wantLen int
	}{
		{"empty set", newIntSet(), true, "", 0},
		{"single element", newIntSet(10, "a"), false, "a", 0},
		{"removes min", newIntSet(30, "c", 10, "a", 20, "b"), false, "a", 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := tc.set.PopMin()
			if tc.wantNil {
				testutils.Equal(t, node == nil, true)

				return
			}

			testutils.Equal(t, node.Value(), tc.wantVal)
			testutils.Equal(t, tc.set.Len(), tc.wantLen)
		})
	}
}

func TestPopMax(t *testing.T) {
	cases := []struct {
		name    string
		set     *SortedSet[int, string]
		wantNil bool
		wantVal string
		wantLen int
	}{
		{"empty set", newIntSet(), true, "", 0},
		{"single element", newIntSet(10, "a"), false, "a", 0},
		{"removes max", newIntSet(10, "a", 30, "c", 20, "b"), false, "c", 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := tc.set.PopMax()
			if tc.wantNil {
				testutils.Equal(t, node == nil, true)

				return
			}

			testutils.Equal(t, node.Value(), tc.wantVal)
			testutils.Equal(t, tc.set.Len(), tc.wantLen)
		})
	}
}

func TestUpsertScoreUpdate(t *testing.T) {
	cases := []struct {
		name      string
		firstKey  int
		secondKey int
		value     string
		wantAdded bool
		wantKey   int
	}{
		{"insert new", 10, 0, "a", true, 10},
		{"update same key", 10, 10, "a", false, 10},
		{"update different key", 10, 20, "a", false, 20},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSortedSet[int, string](comparator.IntComparator)
			s.Upsert(tc.firstKey, tc.value)

			if tc.secondKey != 0 {
				added := s.Upsert(tc.secondKey, tc.value)
				testutils.Equal(t, added, tc.wantAdded)
			}

			node := s.GetByValue(tc.value)
			testutils.Equal(t, node.Key(), tc.wantKey)
		})
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		name    string
		value   string
		wantNil bool
	}{
		{"existing element", "a", false},
		{"non-existing element", "z", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := newIntSet(10, "a", 20, "b")
			node := s.Remove(tc.value)
			testutils.Equal(t, node == nil, tc.wantNil)
		})
	}
}

func TestGetByRank(t *testing.T) {
	s := newIntSet(10, "a", 20, "b", 30, "c")

	cases := []struct {
		name    string
		rank    int
		wantNil bool
		wantVal string
	}{
		{"rank 1", 1, false, "a"},
		{"rank 2", 2, false, "b"},
		{"rank 3", 3, false, "c"},
		{"rank -1 is last", -1, false, "c"},
		{"rank -2 is second last", -2, false, "b"},
		{"rank -3 is first", -3, false, "a"},
		{"rank beyond length", 4, true, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := s.GetByRank(tc.rank, false)
			if tc.wantNil {
				testutils.Equal(t, node == nil, true)

				return
			}

			testutils.Equal(t, node.Value(), tc.wantVal)
		})
	}
}

func TestGetByRankWithRemove(t *testing.T) {
	s := newIntSet(10, "a", 20, "b", 30, "c")
	node := s.GetByRank(1, true)
	testutils.Equal(t, node.Value(), "a")
	testutils.Equal(t, s.Len(), 2)
	testutils.Equal(t, s.Contains("a"), false)
}

func TestGetByRankRange(t *testing.T) {
	s := newIntSet(10, "a", 20, "b", 30, "c", 40, "d", 50, "e")

	cases := []struct {
		name     string
		start    int
		end      int
		wantVals []string
	}{
		{"forward range", 2, 4, []string{"b", "c", "d"}},
		{"reverse range", 4, 2, []string{"d", "c", "b"}},
		{"single element", 3, 3, []string{"c"}},
		{"negative indices", -1, -1, []string{"e"}},
		{"full range", 1, -1, []string{"a", "b", "c", "d", "e"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nodes := s.GetByRankRange(tc.start, tc.end, false)
			testutils.Equal(t, len(nodes), len(tc.wantVals))

			for i, n := range nodes {
				testutils.Equal(t, n.Value(), tc.wantVals[i])
			}
		})
	}
}

func TestGetByRankRangeWithRemove(t *testing.T) {
	s := newIntSet(10, "a", 20, "b", 30, "c", 40, "d")
	nodes := s.GetByRankRange(2, 3, true)
	testutils.Equal(t, len(nodes), 2)
	testutils.Equal(t, nodes[0].Value(), "b")
	testutils.Equal(t, nodes[1].Value(), "c")
	testutils.Equal(t, s.Len(), 2)
}

func TestGetTop(t *testing.T) {
	s := newIntSet(10, "a", 20, "b", 30, "c", 40, "d")

	cases := []struct {
		name     string
		count    int
		wantVals []string
	}{
		{"top 1", 1, []string{"d"}},
		{"top 2", 2, []string{"d", "c"}},
		{"top 3", 3, []string{"d", "c", "b"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nodes := s.GetTop(tc.count, false)
			testutils.Equal(t, len(nodes), len(tc.wantVals))

			for i, n := range nodes {
				testutils.Equal(t, n.Value(), tc.wantVals[i])
			}
		})
	}
}

func TestGetRTop(t *testing.T) {
	s := newIntSet(10, "a", 20, "b", 30, "c", 40, "d")

	cases := []struct {
		name     string
		count    int
		wantVals []string
	}{
		{"bottom 1", 1, []string{"a"}},
		{"bottom 2", 2, []string{"a", "b"}},
		{"bottom 3", 3, []string{"a", "b", "c"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nodes := s.GetRTop(tc.count, false)
			testutils.Equal(t, len(nodes), len(tc.wantVals))

			for i, n := range nodes {
				testutils.Equal(t, n.Value(), tc.wantVals[i])
			}
		})
	}
}

func TestFindRank(t *testing.T) {
	s := newIntSet(10, "a", 20, "b", 30, "c")

	cases := []struct {
		name     string
		value    string
		wantRank int
	}{
		{"first element", "a", 1},
		{"middle element", "b", 2},
		{"last element", "c", 3},
		{"not found", "z", 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, s.FindRank(tc.value), tc.wantRank)
		})
	}
}

func TestGetByKeyRangeForward(t *testing.T) {
	s := newIntSet(10, "a", 20, "b", 30, "c", 40, "d", 50, "e")

	cases := []struct {
		name     string
		start    int
		end      int
		options  *GetByKeyRangeOptions
		wantVals []string
	}{
		{"inclusive range", 20, 40, nil, []string{"b", "c", "d"}},
		{"exclude start", 20, 40, &GetByKeyRangeOptions{ExcludeStart: true}, []string{"c", "d"}},
		{"exclude end", 20, 40, &GetByKeyRangeOptions{ExcludeEnd: true}, []string{"b", "c"}},
		{"exclude both", 20, 40, &GetByKeyRangeOptions{ExcludeStart: true, ExcludeEnd: true}, []string{"c"}},
		{"with limit", 10, 50, &GetByKeyRangeOptions{Limit: 2}, []string{"a", "b"}},
		{"empty result", 60, 70, nil, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nodes := s.GetByKeyRange(tc.start, tc.end, tc.options)
			testutils.Equal(t, len(nodes), len(tc.wantVals))

			for i, n := range nodes {
				testutils.Equal(t, n.Value(), tc.wantVals[i])
			}
		})
	}
}

func TestGetByKeyRangeReverse(t *testing.T) {
	s := newIntSet(10, "a", 20, "b", 30, "c", 40, "d", 50, "e")

	cases := []struct {
		name     string
		start    int
		end      int
		options  *GetByKeyRangeOptions
		wantVals []string
	}{
		{"reverse inclusive", 40, 20, nil, []string{"d", "c", "b"}},
		{"reverse exclude start", 40, 20, &GetByKeyRangeOptions{ExcludeStart: true}, []string{"c", "b"}},
		{"reverse exclude end", 40, 20, &GetByKeyRangeOptions{ExcludeEnd: true}, []string{"d", "c"}},
		{"reverse with limit", 50, 10, &GetByKeyRangeOptions{Limit: 2}, []string{"e", "d"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nodes := s.GetByKeyRange(tc.start, tc.end, tc.options)
			testutils.Equal(t, len(nodes), len(tc.wantVals))

			for i, n := range nodes {
				testutils.Equal(t, n.Value(), tc.wantVals[i])
			}
		})
	}
}

func TestGetByKeyRangeWithRemove(t *testing.T) {
	s := newIntSet(10, "a", 20, "b", 30, "c", 40, "d")
	nodes := s.GetByKeyRange(20, 30, &GetByKeyRangeOptions{Remove: true})
	testutils.Equal(t, len(nodes), 2)
	testutils.Equal(t, s.Len(), 2)
	testutils.Equal(t, s.Contains("b"), false)
	testutils.Equal(t, s.Contains("c"), false)
}

func TestGetByKeyRangeEmpty(t *testing.T) {
	s := newIntSet()
	nodes := s.GetByKeyRange(10, 20, nil)
	testutils.Equal(t, len(nodes), 0)
}

func TestDumpError(t *testing.T) {
	s := newIntSet(10, "a")
	errDump := errors.New("dump failed")

	_, err := s.Dump(func(_ int, _ string) (string, string, error) {
		return "", "", errDump
	})

	testutils.Equal(t, err, errDump)
}

func TestRestoreErrors(t *testing.T) {
	cases := []struct {
		name    string
		dump    string
		wantErr bool
	}{
		{"invalid json", "not json", true},
		{"valid json", `{"10":["a"]}`, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSortedSet[int, string](comparator.IntComparator)
			err := s.Restore(func(key string, values []string) (int, []string, error) {
				k, err := strconv.Atoi(key)

				return k, values, err
			}, tc.dump)

			testutils.Equal(t, err != nil, tc.wantErr)
		})
	}
}

func TestRestoreCallbackError(t *testing.T) {
	s := NewSortedSet[int, string](comparator.IntComparator)
	errRestore := errors.New("restore failed")

	err := s.Restore(func(_ string, _ []string) (int, []string, error) {
		return 0, nil, errRestore
	}, `{"10":["a"]}`)

	testutils.Equal(t, err, errRestore)
}

func TestDumpRestore(t *testing.T) {
	s := newIntSet(10, "a", 20, "b", 30, "c")

	dump, err := s.Dump(func(key int, value string) (string, string, error) {
		return strconv.Itoa(key), value, nil
	})
	testutils.Equal(t, err, nil)

	s2 := NewSortedSet[int, string](comparator.IntComparator)

	err = s2.Restore(func(key string, values []string) (int, []string, error) {
		k, err := strconv.Atoi(key)

		return k, values, err
	}, dump)
	testutils.Equal(t, err, nil)
	testutils.Equal(t, s2.Len(), 3)
	testutils.Equal(t, s2.Contains("a"), true)
	testutils.Equal(t, s2.Contains("b"), true)
	testutils.Equal(t, s2.Contains("c"), true)
}

func TestSafeSortedSet(t *testing.T) {
	s := NewSafeSortedSet[int, string](comparator.IntComparator)

	t.Run("upsert and len", func(t *testing.T) {
		testutils.Equal(t, s.Upsert(10, "a"), true)
		testutils.Equal(t, s.Upsert(20, "b"), true)
		testutils.Equal(t, s.Upsert(30, "c"), true)
		testutils.Equal(t, s.Upsert(40, "d"), true)
		testutils.Equal(t, s.Upsert(50, "e"), true)
		testutils.Equal(t, s.Len(), 5)
	})

	t.Run("contains and get by value", func(t *testing.T) {
		testutils.Equal(t, s.Contains("a"), true)
		testutils.Equal(t, s.Contains("z"), false)
		node := s.GetByValue("b")
		testutils.Equal(t, node.Key(), 20)
	})

	t.Run("find rank", func(t *testing.T) {
		testutils.Equal(t, s.FindRank("a"), 1)
		testutils.Equal(t, s.FindRank("z"), 0)
	})

	t.Run("peek min and max", func(t *testing.T) {
		testutils.Equal(t, s.PeekMin().Value(), "a")
		testutils.Equal(t, s.PeekMax().Value(), "e")
	})

	t.Run("get by rank", func(t *testing.T) {
		testutils.Equal(t, s.GetByRank(1, false).Value(), "a")
		testutils.Equal(t, s.GetByRank(-1, false).Value(), "e")
	})

	t.Run("get by rank range", func(t *testing.T) {
		nodes := s.GetByRankRange(1, 3, false)
		testutils.Equal(t, len(nodes), 3)
	})

	t.Run("get by key range", func(t *testing.T) {
		nodes := s.GetByKeyRange(20, 40, nil)
		testutils.Equal(t, len(nodes), 3)
	})

	t.Run("get top and rtop", func(t *testing.T) {
		top := s.GetTop(2, false)
		testutils.Equal(t, len(top), 2)
		testutils.Equal(t, top[0].Value(), "e")

		rtop := s.GetRTop(2, false)
		testutils.Equal(t, len(rtop), 2)
		testutils.Equal(t, rtop[0].Value(), "a")
	})

	t.Run("get until key", func(t *testing.T) {
		vals := s.GetUntilKey(25, false)
		testutils.Equal(t, len(vals), 2)
	})

	t.Run("dump and restore", func(t *testing.T) {
		dump, err := s.Dump(func(key int, value string) (string, string, error) {
			return strconv.Itoa(key), value, nil
		})
		testutils.Equal(t, err, nil)

		s2 := NewSafeSortedSet[int, string](comparator.IntComparator)

		err = s2.Restore(func(key string, values []string) (int, []string, error) {
			k, err := strconv.Atoi(key)

			return k, values, err
		}, dump)
		testutils.Equal(t, err, nil)
		testutils.Equal(t, s2.Len(), 5)
	})

	t.Run("remove", func(t *testing.T) {
		removed := s.Remove("c")
		testutils.Equal(t, removed.Value(), "c")
		testutils.Equal(t, s.Len(), 4)
	})

	t.Run("pop min and max", func(t *testing.T) {
		minNode := s.PopMin()
		testutils.Equal(t, minNode.Value(), "a")

		maxNode := s.PopMax()
		testutils.Equal(t, maxNode.Value(), "e")
		testutils.Equal(t, s.Len(), 2)
	})
}

func TestSafeSortedSetConcurrent(t *testing.T) {
	s := NewSafeSortedSet[int, string](comparator.IntComparator)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			s.Upsert(n, strconv.Itoa(n))
		}(i)
	}

	wg.Wait()
	testutils.Equal(t, s.Len(), 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			s.Contains(strconv.Itoa(n))
			s.FindRank(strconv.Itoa(n))
		}(i)
	}

	wg.Wait()
}

/*
BenchmarkUpsert-12               1000000              1065 ns/op             241 B/op          4 allocs/op
BenchmarkUpsertAlt-12            3403590               351.5 ns/op           335 B/op          3 allocs/op
BenchmarkGetUntilKey-12            44418             26979 ns/op           32784 B/op         13 allocs/op
BenchmarkGetUntilKeyAlt-12         70933             17045 ns/op           32832 B/op         13 allocs/op
BenchmarkTop-12                  4050680               314.8 ns/op           120 B/op          4 allocs/op
BenchmarkTopWithRemove-12       32379015                34.64 ns/op            0 B/op          0 allocs/op
*/

var resultList []*Node[time.Time, any]

func BenchmarkUpsert(b *testing.B) {
	s := NewSortedSet[time.Time, any](comparator.TimeComparator)

	for i := 0; i < b.N; i++ {
		j := i
		if j%2 == 0 {
			j *= -1
		}

		s.Upsert(time.Now().Add(time.Duration(j)*time.Minute), i)
	}
}

func BenchmarkGetUntilKey(b *testing.B) {
	b.StopTimer()

	s := NewSortedSet[time.Time, any](comparator.TimeComparator)
	for i := -1000; i <= 10000; i++ {
		s.Upsert(time.Now().Add(time.Duration(i)*time.Minute), i)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = s.GetUntilKey(time.Now(), false)
	}
}

func BenchmarkTop(b *testing.B) {
	b.StopTimer()

	s := NewSortedSet[time.Time, any](comparator.TimeComparator)
	for i := -1000; i <= 10000; i++ {
		s.Upsert(time.Now().Add(time.Duration(i)*time.Minute), i)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		resultList = s.GetTop(5, false)
	}
}

func BenchmarkTopWithRemove(b *testing.B) {
	b.StopTimer()

	s := NewSortedSet[time.Time, any](comparator.TimeComparator)
	for i := -1000; i <= 10000; i++ {
		s.Upsert(time.Now().Add(time.Duration(i)*time.Minute), i)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		resultList = s.GetTop(5, true)
	}
}

var (
	benchNode  *Node[int, int]
	benchNodes []*Node[int, int]
)

func BenchmarkAdd(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("n_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s := NewSortedSet[int, int](comparator.IntComparator)
				for j := 0; j < size; j++ {
					s.Upsert(j, j)
				}
			}
		})
	}
}

func BenchmarkGetByRank(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("n_%d", size), func(b *testing.B) {
			s := NewSortedSet[int, int](comparator.IntComparator)
			for j := 0; j < size; j++ {
				s.Upsert(j, j)
			}

			mid := size / 2

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				benchNode = s.GetByRank(mid, false)
			}
		})
	}
}

func BenchmarkRange(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("n_%d", size), func(b *testing.B) {
			s := NewSortedSet[int, int](comparator.IntComparator)
			for j := 0; j < size; j++ {
				s.Upsert(j, j)
			}

			start := size / 4
			end := size * 3 / 4

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				benchNodes = s.GetByKeyRange(start, end, nil)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	cases := []struct {
		name   string
		keys   []int
		values []string
	}{
		{
			name: "empty set",
		},
		{
			name:   "single element",
			keys:   []int{10},
			values: []string{"a"},
		},
		{
			name:   "multiple elements",
			keys:   []int{30, 10, 20},
			values: []string{"c", "a", "b"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSortedSet[int, string](comparator.IntComparator)
			for i, v := range tc.values {
				s.Upsert(tc.keys[i], v)
			}

			clone := s.Copy()

			testutils.Equal(t, clone.Len(), s.Len())

			// Verify all elements present in clone
			for _, v := range tc.values {
				node := clone.GetByValue(v)
				if node == nil {
					t.Fatalf("value %q not found in clone", v)
				}
			}

			// Verify order is preserved
			if s.Len() > 0 {
				origMin := s.PeekMin()
				cloneMin := clone.PeekMin()

				testutils.Equal(t, cloneMin.Key(), origMin.Key())
				testutils.Equal(t, cloneMin.Value(), origMin.Value())
			}

			// Verify independence: remove from original, clone unaffected
			for _, v := range tc.values {
				s.Remove(v)
			}

			testutils.Equal(t, s.Len(), 0)
			testutils.Equal(t, clone.Len(), len(tc.values))
		})
	}
}
