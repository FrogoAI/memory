//nolint:revive // package name is intentional
package utils

import (
	"sync"
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestSafeMapNewWithData(t *testing.T) {
	cases := []struct {
		name string
		data map[string]int
		want map[string]int
	}{
		{
			name: "nil_map",
			data: nil,
			want: map[string]int{},
		},
		{
			name: "empty_map",
			data: map[string]int{},
			want: map[string]int{},
		},
		{
			name: "populated_map",
			data: map[string]int{"a": 1, "b": 2},
			want: map[string]int{"a": 1, "b": 2},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sm := NewSafeMap(tc.data)
			testutils.Equal(t, sm.GetMap(), tc.want)
		})
	}
}

func TestSafeMapSetAndGet(t *testing.T) {
	cases := []struct {
		name   string
		key    string
		value  string
		exists bool
	}{
		{
			name:   "set_and_get_existing",
			key:    "hello",
			value:  "world",
			exists: true,
		},
		{
			name:   "get_nonexistent",
			key:    "missing",
			value:  "",
			exists: false,
		},
	}

	sm := NewSafeMap[string, string](nil)
	sm.Set("hello", "world")

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := sm.Get(tc.key)
			testutils.Equal(t, ok, tc.exists)
			testutils.Equal(t, val, tc.value)
		})
	}
}

func TestSafeMapExists(t *testing.T) {
	sm := NewSafeMap[string, int](nil)
	sm.Set("key1", 42)

	cases := []struct {
		name string
		key  string
		want bool
	}{
		{name: "exists", key: "key1", want: true},
		{name: "not_exists", key: "key2", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, sm.Exists(tc.key), tc.want)
		})
	}
}

func TestSafeMapRemove(t *testing.T) {
	sm := NewSafeMap[string, int](map[string]int{"a": 1, "b": 2})

	sm.Remove("a")
	testutils.Equal(t, sm.Exists("a"), false)
	testutils.Equal(t, sm.Exists("b"), true)

	// removing nonexistent key should not panic
	sm.Remove("nonexistent")
}

func TestSafeMapOverwrite(t *testing.T) {
	sm := NewSafeMap[string, int](nil)
	sm.Set("key", 1)

	val, ok := sm.Get("key")
	testutils.Equal(t, ok, true)
	testutils.Equal(t, val, 1)

	sm.Set("key", 2)

	val, ok = sm.Get("key")
	testutils.Equal(t, ok, true)
	testutils.Equal(t, val, 2)
}

func TestSafeMapGetMapReturnsCopy(t *testing.T) {
	sm := NewSafeMap[string, int](map[string]int{"a": 1})

	m := sm.GetMap()
	m["b"] = 2

	// original should be unchanged
	testutils.Equal(t, sm.Exists("b"), false)
}

func TestSafeMapConcurrentOps(t *testing.T) {
	sm := NewSafeMap[string, int](nil)

	const goroutines = 50

	const opsPerGoroutine = 100

	var wg sync.WaitGroup

	wg.Add(goroutines)

	for i := range goroutines {
		go func(id int) {
			defer wg.Done()

			for j := range opsPerGoroutine {
				key := "key"
				sm.Set(key, id*opsPerGoroutine+j)
				sm.Get(key)
				sm.Exists(key)
				sm.GetMap()
			}
		}(i)
	}

	wg.Wait()

	// map should still be operational after concurrent access
	sm.Set("final", 999)

	val, ok := sm.Get("final")
	testutils.Equal(t, ok, true)
	testutils.Equal(t, val, 999)
}

func TestSafeMapConcurrentReadWrite(t *testing.T) {
	sm := NewSafeMap[string, int](nil)

	const writers = 10

	const readers = 20

	const ops = 200

	var wg sync.WaitGroup

	wg.Add(writers + readers)

	// writers
	for i := range writers {
		go func(id int) {
			defer wg.Done()

			for j := range ops {
				sm.Set("key", id*ops+j)
			}
		}(i)
	}

	// readers
	for range readers {
		go func() {
			defer wg.Done()

			for range ops {
				sm.Get("key")
				sm.Exists("key")
			}
		}()
	}

	wg.Wait()
}

func TestSafeMapConcurrentRemove(t *testing.T) {
	sm := NewSafeMap[string, int](nil)

	const goroutines = 20

	var wg sync.WaitGroup

	wg.Add(goroutines * 2)

	for i := range goroutines {
		go func(id int) {
			defer wg.Done()

			sm.Set("key", id)
		}(i)

		go func() {
			defer wg.Done()

			sm.Remove("key")
		}()
	}

	wg.Wait()
}
