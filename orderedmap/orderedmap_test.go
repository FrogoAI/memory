package orderedmap

import (
	"strconv"
	"sync"
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestAdd(t *testing.T) {
	cases := []struct {
		name     string
		keys     []string
		values   []int
		wantKeys []string
		wantVals []int
	}{
		{
			name:     "single element",
			keys:     []string{"a"},
			values:   []int{1},
			wantKeys: []string{"a"},
			wantVals: []int{1},
		},
		{
			name:     "multiple elements preserve insertion order",
			keys:     []string{"c", "a", "b"},
			values:   []int{3, 1, 2},
			wantKeys: []string{"c", "a", "b"},
			wantVals: []int{3, 1, 2},
		},
		{
			name:     "duplicate key updates value without changing order",
			keys:     []string{"a", "b", "a"},
			values:   []int{1, 2, 99},
			wantKeys: []string{"a", "b"},
			wantVals: []int{99, 2},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := &OrderedMap[string, int]{}

			for i, k := range tc.keys {
				o.Add(k, tc.values[i])
			}

			keys, vals := o.GetAll()
			testutils.Equal(t, keys, tc.wantKeys)
			testutils.Equal(t, vals, tc.wantVals)
		})
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		name   string
		key    string
		want   string
		wantOK bool
	}{
		{
			name:   "existing key",
			key:    "A",
			want:   "1",
			wantOK: true,
		},
		{
			name:   "missing key returns zero value and false",
			key:    "Z",
			want:   "",
			wantOK: false,
		},
	}

	o := &OrderedMap[string, string]{}
	o.Add("A", "1")
	o.Add("B", "2")

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := o.Get(tc.key)
			testutils.Equal(t, got, tc.want)
			testutils.Equal(t, ok, tc.wantOK)
		})
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		name      string
		setup     []string
		removeKey string
		wantKeys  []string
	}{
		{
			name:      "remove middle element",
			setup:     []string{"a", "b", "c"},
			removeKey: "b",
			wantKeys:  []string{"a", "c"},
		},
		{
			name:      "remove first element",
			setup:     []string{"a", "b", "c"},
			removeKey: "a",
			wantKeys:  []string{"b", "c"},
		},
		{
			name:      "remove last element",
			setup:     []string{"a", "b", "c"},
			removeKey: "c",
			wantKeys:  []string{"a", "b"},
		},
		{
			name:      "remove non-existent key is no-op",
			setup:     []string{"a", "b"},
			removeKey: "z",
			wantKeys:  []string{"a", "b"},
		},
		{
			name:      "remove only element",
			setup:     []string{"a"},
			removeKey: "a",
			wantKeys:  []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := &OrderedMap[string, int]{}
			for i, k := range tc.setup {
				o.Add(k, i)
			}

			o.Remove(tc.removeKey)

			keys, _ := o.GetAll()
			testutils.Equal(t, keys, tc.wantKeys)
		})
	}
}

func TestExists(t *testing.T) {
	cases := []struct {
		name string
		key  string
		want bool
	}{
		{
			name: "existing key returns true",
			key:  "a",
			want: true,
		},
		{
			name: "missing key returns false",
			key:  "z",
			want: false,
		},
	}

	o := &OrderedMap[string, int]{}
	o.Add("a", 1)
	o.Add("b", 2)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, o.Exists(tc.key), tc.want)
		})
	}
}

func TestSize(t *testing.T) {
	cases := []struct {
		name  string
		count int
		want  int
	}{
		{
			name:  "empty map",
			count: 0,
			want:  0,
		},
		{
			name:  "single element",
			count: 1,
			want:  1,
		},
		{
			name:  "multiple elements",
			count: 5,
			want:  5,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := &OrderedMap[string, int]{}
			for i := range tc.count {
				o.Add(strconv.Itoa(i), i)
			}

			testutils.Equal(t, o.Size(), tc.want)
		})
	}
}

func TestCopy(t *testing.T) {
	o := &OrderedMap[string, int]{}
	o.Add("a", 1)
	o.Add("b", 2)
	o.Add("c", 3)

	cp := o.Copy()

	keys, vals := cp.GetAll()
	testutils.Equal(t, keys, []string{"a", "b", "c"})
	testutils.Equal(t, vals, []int{1, 2, 3})

	// mutating copy does not affect original
	cp.Add("d", 4)
	testutils.Equal(t, o.Size(), 3)
	testutils.Equal(t, cp.Size(), 4)
}

func TestCopyEmpty(t *testing.T) {
	o := &OrderedMap[string, int]{}

	cp := o.Copy()

	testutils.Equal(t, cp.Size(), 0)

	keys, vals := cp.GetAll()
	testutils.Equal(t, keys, []string{})
	testutils.Equal(t, vals, []int{})
}

func TestClear(t *testing.T) {
	o := &OrderedMap[string, int]{}
	o.Add("a", 1)
	o.Add("b", 2)

	o.Clear()

	testutils.Equal(t, o.Size(), 0)
	testutils.Equal(t, o.Exists("a"), false)

	keys, vals := o.GetAll()
	testutils.Equal(t, keys, []string{})
	testutils.Equal(t, vals, []int{})
}

func TestClearEmpty(t *testing.T) {
	o := &OrderedMap[string, int]{}

	o.Clear()

	testutils.Equal(t, o.Size(), 0)
}

func TestSetKeys(t *testing.T) {
	o := &OrderedMap[string, int]{}

	o.SetKeys([]string{"x", "y", "z"}, 0)

	testutils.Equal(t, o.Size(), 3)

	x, xOK := o.Get("x")
	testutils.Equal(t, x, 0)
	testutils.Equal(t, xOK, true)

	y, yOK := o.Get("y")
	testutils.Equal(t, y, 0)
	testutils.Equal(t, yOK, true)

	z, zOK := o.Get("z")
	testutils.Equal(t, z, 0)
	testutils.Equal(t, zOK, true)

	keys, _ := o.GetAll()
	testutils.Equal(t, keys, []string{"x", "y", "z"})
}

func TestKeys(t *testing.T) {
	cases := []struct {
		name     string
		keys     []string
		values   []int
		wantKeys []string
	}{
		{
			name:     "empty map",
			keys:     nil,
			values:   nil,
			wantKeys: []string{},
		},
		{
			name:     "preserves insertion order",
			keys:     []string{"c", "a", "b"},
			values:   []int{3, 1, 2},
			wantKeys: []string{"c", "a", "b"},
		},
		{
			name:     "single element",
			keys:     []string{"x"},
			values:   []int{42},
			wantKeys: []string{"x"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := &OrderedMap[string, int]{}
			for i, k := range tc.keys {
				o.Add(k, tc.values[i])
			}

			got := o.Keys()

			if len(tc.wantKeys) == 0 {
				testutils.Equal(t, len(got), 0)
			} else {
				testutils.Equal(t, got, tc.wantKeys)
			}
		})
	}
}

func TestKeysReturnsCopy(t *testing.T) {
	o := &OrderedMap[string, int]{}
	o.Add("a", 1)
	o.Add("b", 2)

	keys := o.Keys()
	keys[0] = "mutated"

	testutils.Equal(t, o.Keys()[0], "a")
}

func TestValues(t *testing.T) {
	cases := []struct {
		name     string
		keys     []string
		values   []int
		wantVals []int
	}{
		{
			name:     "empty map",
			keys:     nil,
			values:   nil,
			wantVals: []int{},
		},
		{
			name:     "preserves insertion order",
			keys:     []string{"c", "a", "b"},
			values:   []int{3, 1, 2},
			wantVals: []int{3, 1, 2},
		},
		{
			name:     "single element",
			keys:     []string{"x"},
			values:   []int{42},
			wantVals: []int{42},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := &OrderedMap[string, int]{}
			for i, k := range tc.keys {
				o.Add(k, tc.values[i])
			}

			got := o.Values()

			if len(tc.wantVals) == 0 {
				testutils.Equal(t, len(got), 0)
			} else {
				testutils.Equal(t, got, tc.wantVals)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	cases := []struct {
		name     string
		keys     []string
		values   []string
		wantKeys []string
		wantVals []string
	}{
		{
			name:     "empty map",
			keys:     nil,
			values:   nil,
			wantKeys: []string{},
			wantVals: []string{},
		},
		{
			name:     "preserves insertion order",
			keys:     []string{"c", "a", "b"},
			values:   []string{"3", "1", "2"},
			wantKeys: []string{"c", "a", "b"},
			wantVals: []string{"3", "1", "2"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := &OrderedMap[string, string]{}
			for i, k := range tc.keys {
				o.Add(k, tc.values[i])
			}

			keys, vals := o.GetAll()
			testutils.Equal(t, keys, tc.wantKeys)
			testutils.Equal(t, vals, tc.wantVals)
		})
	}
}

func TestGetAllReturnsCopy(t *testing.T) {
	o := &OrderedMap[string, int]{}
	o.Add("a", 1)

	keys, _ := o.GetAll()
	keys[0] = "mutated"

	origKeys, _ := o.GetAll()
	testutils.Equal(t, origKeys[0], "a")
}

func TestGetMap(t *testing.T) {
	o := &OrderedMap[string, int]{}
	o.Add("a", 1)
	o.Add("b", 2)

	m := o.GetMap()

	testutils.Equal(t, len(m), 2)
	testutils.Equal(t, m["a"], 1)
	testutils.Equal(t, m["b"], 2)
}

func TestGetMapEmpty(t *testing.T) {
	o := &OrderedMap[string, int]{}

	m := o.GetMap()

	testutils.Equal(t, len(m), 0)
}

func TestSetAll(t *testing.T) {
	cases := []struct {
		name     string
		keys     []string
		initVals []int
		newVals  []int
		wantVals []int
	}{
		{
			name:     "update all values",
			keys:     []string{"a", "b", "c"},
			initVals: []int{1, 2, 3},
			newVals:  []int{10, 20, 30},
			wantVals: []int{10, 20, 30},
		},
		{
			name:     "more values than keys is idempotent",
			keys:     []string{"a", "b"},
			initVals: []int{1, 2},
			newVals:  []int{10, 20, 30, 40},
			wantVals: []int{10, 20},
		},
		{
			name:     "fewer values than keys is no-op",
			keys:     []string{"a", "b", "c"},
			initVals: []int{1, 2, 3},
			newVals:  []int{10},
			wantVals: []int{1, 2, 3},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := &OrderedMap[string, int]{}
			for i, k := range tc.keys {
				o.Add(k, tc.initVals[i])
			}

			o.SetAll(tc.newVals)

			_, vals := o.GetAll()
			testutils.Equal(t, vals, tc.wantVals)
		})
	}
}

func TestIterator(t *testing.T) {
	cases := []struct {
		name     string
		keys     []string
		values   []int
		bufSize  int
		wantVals []int
	}{
		{
			name:     "iterates in insertion order",
			keys:     []string{"a", "b", "c"},
			values:   []int{1, 2, 3},
			bufSize:  3,
			wantVals: []int{1, 2, 3},
		},
		{
			name:     "empty map yields no values",
			keys:     nil,
			values:   nil,
			bufSize:  1,
			wantVals: []int{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := &OrderedMap[string, int]{}
			for i, k := range tc.keys {
				o.Add(k, tc.values[i])
			}

			ch := o.Iterator(tc.bufSize)

			var got []int
			for v := range ch {
				got = append(got, v)
			}

			if len(tc.wantVals) == 0 {
				testutils.Equal(t, len(got), 0)
			} else {
				testutils.Equal(t, got, tc.wantVals)
			}
		})
	}
}

func TestOrderAfterMixedOperations(t *testing.T) {
	o := &OrderedMap[string, string]{}
	o.Add("A", "1")
	o.Add("B", "2")
	o.Add("C", "3")
	o.Add("D", "4")
	o.Remove("B")
	o.Add("E", "5")
	o.Remove("A")

	keys, vals := o.GetAll()
	testutils.Equal(t, keys, []string{"C", "D", "E"})
	testutils.Equal(t, vals, []string{"3", "4", "5"})
}

func TestAddAfterClear(t *testing.T) {
	o := &OrderedMap[string, int]{}
	o.Add("a", 1)
	o.Clear()
	o.Add("b", 2)

	testutils.Equal(t, o.Size(), 1)

	b, bOK := o.Get("b")
	testutils.Equal(t, b, 2)
	testutils.Equal(t, bOK, true)
	testutils.Equal(t, o.Exists("a"), false)
}

// --- SafeOrderedMap tests ---

func TestSafeAdd(t *testing.T) {
	var s SafeOrderedMap[string, int]

	s.Add("a", 1)
	s.Add("b", 2)

	testutils.Equal(t, s.Size(), 2)

	a, aOK := s.Get("a")
	testutils.Equal(t, a, 1)
	testutils.Equal(t, aOK, true)

	b, bOK := s.Get("b")
	testutils.Equal(t, b, 2)
	testutils.Equal(t, bOK, true)
}

func TestSafeRemove(t *testing.T) {
	var s SafeOrderedMap[string, int]

	s.Add("a", 1)
	s.Add("b", 2)
	s.Remove("a")

	testutils.Equal(t, s.Size(), 1)
	testutils.Equal(t, s.Exists("a"), false)
	testutils.Equal(t, s.Exists("b"), true)
}

func TestSafeCopy(t *testing.T) {
	var s SafeOrderedMap[string, int]

	s.Add("a", 1)
	s.Add("b", 2)

	cp := s.Copy()

	testutils.Equal(t, cp.Size(), 2)

	a, aOK := cp.Get("a")
	testutils.Equal(t, a, 1)
	testutils.Equal(t, aOK, true)

	cp.Add("c", 3)
	testutils.Equal(t, s.Size(), 2)
}

func TestSafeClear(t *testing.T) {
	var s SafeOrderedMap[string, int]

	s.Add("a", 1)
	s.Clear()

	testutils.Equal(t, s.Size(), 0)
}

func TestSafeSetKeys(t *testing.T) {
	var s SafeOrderedMap[string, int]

	s.SetKeys([]string{"x", "y"}, 0)

	testutils.Equal(t, s.Size(), 2)

	x, xOK := s.Get("x")
	testutils.Equal(t, x, 0)
	testutils.Equal(t, xOK, true)
}

func TestSafeKeys(t *testing.T) {
	var s SafeOrderedMap[string, int]

	s.Add("a", 1)
	s.Add("b", 2)

	testutils.Equal(t, s.Keys(), []string{"a", "b"})
}

func TestSafeValues(t *testing.T) {
	var s SafeOrderedMap[string, int]

	s.Add("a", 1)
	s.Add("b", 2)

	testutils.Equal(t, s.Values(), []int{1, 2})
}

func TestSafeGetAll(t *testing.T) {
	var s SafeOrderedMap[string, int]

	s.Add("a", 1)
	s.Add("b", 2)

	keys, vals := s.GetAll()
	testutils.Equal(t, keys, []string{"a", "b"})
	testutils.Equal(t, vals, []int{1, 2})
}

func TestSafeGetMap(t *testing.T) {
	var s SafeOrderedMap[string, int]

	s.Add("a", 1)
	s.Add("b", 2)

	m := s.GetMap()
	testutils.Equal(t, len(m), 2)
	testutils.Equal(t, m["a"], 1)
}

func TestSafeSetAll(t *testing.T) {
	var s SafeOrderedMap[string, int]

	s.Add("a", 1)
	s.Add("b", 2)
	s.SetAll([]int{10, 20})

	a, _ := s.Get("a")
	testutils.Equal(t, a, 10)

	b, _ := s.Get("b")
	testutils.Equal(t, b, 20)
}

func TestSafeIterator(t *testing.T) {
	var s SafeOrderedMap[string, int]

	s.Add("a", 1)
	s.Add("b", 2)

	ch := s.Iterator(2)

	var got []int
	for v := range ch {
		got = append(got, v)
	}

	testutils.Equal(t, got, []int{1, 2})
}

func TestSafeConcurrentAddGet(t *testing.T) {
	var s SafeOrderedMap[int, int]

	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			s.Add(n, n*10)
		}(i)
	}

	wg.Wait()

	testutils.Equal(t, s.Size(), 100)

	for i := range 100 {
		v, ok := s.Get(i)
		testutils.Equal(t, v, i*10)
		testutils.Equal(t, ok, true)
	}
}

func TestSafeConcurrentReadWrite(t *testing.T) {
	var s SafeOrderedMap[string, int]

	s.Add("key", 0)

	var wg sync.WaitGroup

	for i := range 50 {
		wg.Add(2)

		go func(n int) {
			defer wg.Done()

			s.Add("key", n)
		}(i)

		go func() {
			defer wg.Done()

			s.Get("key")
			s.Exists("key")
			s.Size()
		}()
	}

	wg.Wait()

	testutils.Equal(t, s.Size(), 1)
	testutils.Equal(t, s.Exists("key"), true)
}

// --- Benchmarks ---

func BenchmarkAdd(b *testing.B) {
	b.Run("new_key", func(b *testing.B) {
		o := &OrderedMap[string, int]{}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			o.Add(strconv.Itoa(i), i)
		}
	})

	b.Run("update_existing", func(b *testing.B) {
		o := &OrderedMap[string, int]{}
		for i := 0; i < 1000; i++ {
			o.Add(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			o.Add(strconv.Itoa(i%1000), i)
		}
	})
}

func BenchmarkGet(b *testing.B) {
	b.Run("hit", func(b *testing.B) {
		o := &OrderedMap[string, int]{}
		for i := 0; i < 1000; i++ {
			o.Add(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			o.Get(strconv.Itoa(i % 1000))
		}
	})

	b.Run("miss", func(b *testing.B) {
		o := &OrderedMap[string, int]{}
		for i := 0; i < 1000; i++ {
			o.Add(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			o.Get("miss" + strconv.Itoa(i))
		}
	})
}

func BenchmarkIterate(b *testing.B) {
	b.Run("get_all", func(b *testing.B) {
		o := &OrderedMap[string, int]{}
		for i := 0; i < 1000; i++ {
			o.Add(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			o.GetAll()
		}
	})

	b.Run("iterator", func(b *testing.B) {
		o := &OrderedMap[string, int]{}
		for i := 0; i < 1000; i++ {
			o.Add(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			ch := o.Iterator(1000)
			for v := range ch {
				_ = v
			}
		}
	})
}
