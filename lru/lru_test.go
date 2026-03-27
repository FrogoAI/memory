package lru

import (
	"strconv"
	"testing"

	"github.com/FrogoAI/testutils"
)

func BenchmarkPut(b *testing.B) {
	b.Run("under_capacity", func(b *testing.B) {
		c := NewLRUCache[int](b.N + 1)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			c.Put(strconv.Itoa(i), i)
		}
	})

	b.Run("update_existing", func(b *testing.B) {
		c := NewLRUCache[int](1000)
		for i := 0; i < 1000; i++ {
			c.Put(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			c.Put(strconv.Itoa(i%1000), i)
		}
	})
}

func BenchmarkGet(b *testing.B) {
	b.Run("hit", func(b *testing.B) {
		c := NewLRUCache[int](1000)
		for i := 0; i < 1000; i++ {
			c.Put(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			c.Get(strconv.Itoa(i % 1000))
		}
	})

	b.Run("miss", func(b *testing.B) {
		c := NewLRUCache[int](1000)
		for i := 0; i < 1000; i++ {
			c.Put(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			c.Get("missing" + strconv.Itoa(i))
		}
	})
}

func BenchmarkEviction(b *testing.B) {
	b.Run("at_capacity", func(b *testing.B) {
		c := NewLRUCache[int](1000)
		for i := 0; i < 1000; i++ {
			c.Put(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			c.Put("evict"+strconv.Itoa(i), i)
		}
	})

	b.Run("small_capacity", func(b *testing.B) {
		c := NewLRUCache[int](10)
		for i := 0; i < 10; i++ {
			c.Put(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			c.Put("evict"+strconv.Itoa(i), i)
		}
	})
}

func TestLRU(t *testing.T) {
	l := NewLRUCache[any](10)

	for i := 0; i < 20; i++ {
		k := strconv.Itoa(i)
		l.Put("key"+k, "val"+k)
	}

	v, ok := l.Get("test")
	testutils.Equal(t, ok, false)
	testutils.Equal(t, v, nil)

	v, ok = l.Get("key12")
	testutils.Equal(t, ok, true)
	testutils.Equal(t, v, "val12")

	for i := 20; i < 28; i++ {
		k := strconv.Itoa(i)
		l.Put("key"+k, "val"+k)
	}

	v, ok = l.Get("key12")
	testutils.Equal(t, ok, true)
	testutils.Equal(t, v, "val12")

	v, ok = l.Get("key22")
	testutils.Equal(t, ok, true)
	testutils.Equal(t, v, "val22")

	v, ok = l.Get("key13")
	testutils.Equal(t, ok, false)
	testutils.Equal(t, v, nil)
}

func TestCapacityOne(t *testing.T) {
	cases := []struct {
		name      string
		puts      []struct{ key, val string }
		getKey    string
		wantVal   string
		wantFound bool
	}{
		{
			name:      "single item stays",
			puts:      []struct{ key, val string }{{"a", "alpha"}},
			getKey:    "a",
			wantVal:   "alpha",
			wantFound: true,
		},
		{
			name: "second put evicts first",
			puts: []struct{ key, val string }{
				{"a", "alpha"},
				{"b", "beta"},
			},
			getKey:    "a",
			wantVal:   "",
			wantFound: false,
		},
		{
			name: "second put keeps second",
			puts: []struct{ key, val string }{
				{"a", "alpha"},
				{"b", "beta"},
			},
			getKey:    "b",
			wantVal:   "beta",
			wantFound: true,
		},
		{
			name: "update same key keeps it",
			puts: []struct{ key, val string }{
				{"a", "alpha"},
				{"a", "updated"},
			},
			getKey:    "a",
			wantVal:   "updated",
			wantFound: true,
		},
		{
			name: "third put evicts second",
			puts: []struct{ key, val string }{
				{"a", "alpha"},
				{"b", "beta"},
				{"c", "gamma"},
			},
			getKey:    "b",
			wantVal:   "",
			wantFound: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewLRUCache[string](1)

			for _, p := range tc.puts {
				c.Put(p.key, p.val)
			}

			got, ok := c.Get(tc.getKey)
			testutils.Equal(t, ok, tc.wantFound)
			testutils.Equal(t, got, tc.wantVal)
		})
	}
}

func TestEvictionOrder(t *testing.T) {
	cases := []struct {
		name         string
		capacity     int
		puts         []string
		accessBefore []string // keys to Get before the eviction-triggering put
		evictPut     string   // the put that triggers eviction
		wantEvicted  []string // keys that should be gone
		wantPresent  []string // keys that should still be present
	}{
		{
			name:        "LRU item evicted first",
			capacity:    3,
			puts:        []string{"a", "b", "c"},
			evictPut:    "d",
			wantEvicted: []string{"a"},
			wantPresent: []string{"b", "c", "d"},
		},
		{
			name:         "access promotes item, next-oldest evicted",
			capacity:     3,
			puts:         []string{"a", "b", "c"},
			accessBefore: []string{"a"},
			evictPut:     "d",
			wantEvicted:  []string{"b"},
			wantPresent:  []string{"a", "c", "d"},
		},
		{
			name:         "multiple accesses change eviction order",
			capacity:     3,
			puts:         []string{"a", "b", "c"},
			accessBefore: []string{"a", "b"},
			evictPut:     "d",
			wantEvicted:  []string{"c"},
			wantPresent:  []string{"a", "b", "d"},
		},
		{
			name:         "update existing does not evict",
			capacity:     3,
			puts:         []string{"a", "b", "c"},
			accessBefore: []string{},
			evictPut:     "a",
			wantEvicted:  []string{},
			wantPresent:  []string{"a", "b", "c"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewLRUCache[string](tc.capacity)

			for _, k := range tc.puts {
				c.Put(k, "val-"+k)
			}

			for _, k := range tc.accessBefore {
				c.Get(k)
			}

			c.Put(tc.evictPut, "val-"+tc.evictPut)

			for _, k := range tc.wantEvicted {
				_, ok := c.Get(k)
				testutils.Equal(t, ok, false)
			}

			for _, k := range tc.wantPresent {
				v, ok := c.Get(k)
				testutils.Equal(t, ok, true)
				testutils.Equal(t, v, "val-"+k)
			}
		})
	}
}

func TestGetAfterEvict(t *testing.T) {
	cases := []struct {
		name    string
		wantVal int
		wantOK  bool
	}{
		{
			name:    "evicted key returns zero value and false",
			wantVal: 0,
			wantOK:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewLRUCache[int](2)
			c.Put("first", 1)
			c.Put("second", 2)
			c.Put("third", 3) // evicts "first"

			got, ok := c.Get("first")
			testutils.Equal(t, ok, tc.wantOK)
			testutils.Equal(t, got, tc.wantVal)
		})
	}
}

func TestGetAfterEvictAndReinsert(t *testing.T) {
	cases := []struct {
		name    string
		wantVal int
		wantOK  bool
	}{
		{
			name:    "re-added key is accessible after eviction",
			wantVal: 100,
			wantOK:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewLRUCache[int](2)
			c.Put("a", 1)
			c.Put("b", 2)
			c.Put("c", 3) // evicts "a"

			_, ok := c.Get("a")
			testutils.Equal(t, ok, false)

			c.Put("a", tc.wantVal) // re-add "a", evicts "b"

			got, ok := c.Get("a")
			testutils.Equal(t, ok, tc.wantOK)
			testutils.Equal(t, got, tc.wantVal)

			_, ok = c.Get("b")
			testutils.Equal(t, ok, false)
		})
	}
}

func TestLenAfterEviction(t *testing.T) {
	cases := []struct {
		name     string
		capacity int
		numPuts  int
		wantLen  int
	}{
		{
			name:     "at capacity after overflow",
			capacity: 3,
			numPuts:  10,
			wantLen:  3,
		},
		{
			name:     "capacity 1 after many puts",
			capacity: 1,
			numPuts:  50,
			wantLen:  1,
		},
		{
			name:     "under capacity",
			capacity: 10,
			numPuts:  5,
			wantLen:  5,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewLRUCache[int](tc.capacity)

			for i := range tc.numPuts {
				c.Put(strconv.Itoa(i), i)
			}

			testutils.Equal(t, c.Len(), tc.wantLen)
		})
	}
}

func TestCopy(t *testing.T) {
	cases := []struct {
		name     string
		capacity int
		keys     []string
		values   []int
	}{
		{
			name:     "empty cache",
			capacity: 5,
		},
		{
			name:     "single entry",
			capacity: 5,
			keys:     []string{"a"},
			values:   []int{1},
		},
		{
			name:     "multiple entries",
			capacity: 5,
			keys:     []string{"a", "b", "c"},
			values:   []int{1, 2, 3},
		},
		{
			name:     "at capacity",
			capacity: 2,
			keys:     []string{"a", "b"},
			values:   []int{1, 2},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewLRUCache[int](tc.capacity)

			for i, k := range tc.keys {
				c.Put(k, tc.values[i])
			}

			clone := c.Copy()

			testutils.Equal(t, clone.Len(), c.Len())

			// Verify all entries present in clone
			for i, k := range tc.keys {
				val, ok := clone.Get(k)
				testutils.Equal(t, ok, true)
				testutils.Equal(t, val, tc.values[i])
			}

			// Verify independence: add to original, clone unaffected
			c.Put("new", 99)

			_, ok := clone.Get("new")
			testutils.Equal(t, ok, false)
		})
	}
}

func TestIterator(t *testing.T) {
	cases := []struct {
		name     string
		capacity int
		puts     []struct{ key, val string }
		access   []string // keys to Get (promote) before iterating
		wantKeys []string // expected keys in MRU-to-LRU order
	}{
		{
			name:     "empty cache",
			capacity: 5,
			wantKeys: nil,
		},
		{
			name:     "single entry",
			capacity: 5,
			puts:     []struct{ key, val string }{{"a", "alpha"}},
			wantKeys: []string{"a"},
		},
		{
			name:     "multiple entries MRU to LRU",
			capacity: 5,
			puts: []struct{ key, val string }{
				{"a", "alpha"},
				{"b", "beta"},
				{"c", "gamma"},
			},
			wantKeys: []string{"c", "b", "a"},
		},
		{
			name:     "order reflects Get promotion",
			capacity: 5,
			puts: []struct{ key, val string }{
				{"a", "alpha"},
				{"b", "beta"},
				{"c", "gamma"},
			},
			access:   []string{"a"},
			wantKeys: []string{"a", "c", "b"},
		},
		{
			name:     "after eviction",
			capacity: 2,
			puts: []struct{ key, val string }{
				{"a", "alpha"},
				{"b", "beta"},
				{"c", "gamma"}, // evicts "a"
			},
			wantKeys: []string{"c", "b"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewLRUCache[string](tc.capacity)

			for _, p := range tc.puts {
				c.Put(p.key, p.val)
			}

			for _, k := range tc.access {
				c.Get(k)
			}

			var gotKeys []string

			var gotVals []string

			for key, val := range c.Iterator() {
				gotKeys = append(gotKeys, key)
				gotVals = append(gotVals, val)
			}

			testutils.Equal(t, len(gotKeys), len(tc.wantKeys))

			for i, wantKey := range tc.wantKeys {
				testutils.Equal(t, gotKeys[i], wantKey)
			}

			// Verify values match their keys
			for i, key := range gotKeys {
				expected, ok := c.Get(key)
				testutils.Equal(t, ok, true)
				testutils.Equal(t, gotVals[i], expected)
			}
		})
	}
}

func TestIteratorEarlyBreak(t *testing.T) {
	c := NewLRUCache[int](5)

	for i := range 5 {
		c.Put(strconv.Itoa(i), i)
	}

	count := 0

	for range c.Iterator() {
		count++
		if count == 2 {
			break
		}
	}

	testutils.Equal(t, count, 2)
}

func TestCopyEvictionOrder(t *testing.T) {
	c := NewLRUCache[int](3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	c.Get("a") // promote "a" to MRU

	clone := c.Copy()

	// Adding a new entry should evict "b" (LRU) in both original and clone
	clone.Put("d", 4)

	_, ok := clone.Get("b")
	testutils.Equal(t, ok, false) // "b" was LRU and should be evicted

	_, ok = clone.Get("a")
	testutils.Equal(t, ok, true) // "a" was promoted, should survive
}
