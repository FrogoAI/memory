package bloom

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/FrogoAI/testutils"
)

var (
	foo = []byte("foo")
	bar = []byte("bar")
	baz = []byte("baz")
)

func TestCountingFilter(t *testing.T) {
	f, err := NewCounting(3000, 0.01)
	testutils.Equal(t, err, nil)

	f.Add(foo)
	f.Add(foo)
	f.Remove(foo)

	if !f.Test(foo) {
		t.Error("foo not in bloom filter")
	}

	f.Remove(foo)

	if f.Test(foo) {
		t.Error("foo still in bloom filter")
	}
}

func TestNewCounting_Errors(t *testing.T) {
	cases := []struct {
		name    string
		n       int
		p       float64
		wantErr error
	}{
		{
			name:    "bitset overflow from huge capacity",
			n:       math.MaxInt32,
			p:       1e-300,
			wantErr: ErrBitsetTooBig,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewCounting(tc.n, tc.p)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("got err %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestRemoveNonExistent(t *testing.T) {
	f, err := NewCounting(100, 0.01)
	testutils.Equal(t, err, nil)

	cases := []struct {
		name string
		data []byte
	}{
		{name: "never added element", data: []byte("ghost")},
		{name: "empty key", data: []byte{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f.Remove(tc.data)

			if f.Test(tc.data) {
				t.Errorf("element %q should not be in filter after removing non-existent", tc.data)
			}
		})
	}
}

func TestCounterOverflow(t *testing.T) {
	f, err := NewCounting(100, 0.01)
	testutils.Equal(t, err, nil)

	for range 300 {
		f.Add(foo)
	}

	if !f.Test(foo) {
		t.Error("element should still test positive after overflow-protected adds")
	}

	for range 300 {
		f.Remove(foo)
	}

	if f.Test(foo) {
		t.Error("element should test negative after all removes")
	}
}

func TestNewCountingFromBytes_Errors(t *testing.T) {
	cases := []struct {
		name    string
		data    []byte
		wantErr error
	}{
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: ErrEmptyDump,
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: ErrEmptyDump,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewCountingFromBytes(tc.data)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("got err %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestSerializationRoundTrip(t *testing.T) {
	cases := []struct {
		name   string
		add    [][]byte
		remove [][]byte
	}{
		{
			name: "empty filter",
		},
		{
			name: "single element",
			add:  [][]byte{foo},
		},
		{
			name:   "add then remove",
			add:    [][]byte{foo, bar},
			remove: [][]byte{foo},
		},
		{
			name: "many elements",
			add:  [][]byte{foo, bar, baz, []byte("qux"), []byte("quux")},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewCounting(100, 0.01)
			testutils.Equal(t, err, nil)

			for _, d := range tc.add {
				f.Add(d)
			}

			for _, d := range tc.remove {
				f.Remove(d)
			}

			dump, err := f.ToBytes()
			testutils.Equal(t, err, nil)

			restored, err := NewCountingFromBytes(dump)
			testutils.Equal(t, err, nil)

			for _, d := range tc.add {
				expected := true

				for _, r := range tc.remove {
					if string(d) == string(r) {
						expected = false

						break
					}
				}

				testutils.Equal(t, restored.Test(d), expected)
			}

			testutils.Equal(t, restored.Test([]byte("definitely-not-in-filter")), false)
		})
	}
}

func benchFilter(b *testing.B, n int) *CountingFilter {
	b.Helper()

	f, err := NewCounting(n, 0.01)
	if err != nil {
		b.Fatal(err)
	}

	return f
}

func benchData(n int) [][]byte {
	data := make([][]byte, n)
	for i := range n {
		data[i] = []byte(fmt.Sprintf("key-%d", i))
	}

	return data
}

func BenchmarkAdd(b *testing.B) {
	sizes := []int{100, 1_000, 10_000}
	for _, n := range sizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			f := benchFilter(b, n)
			data := benchData(n)

			b.ResetTimer()

			for range b.N {
				f.Add(data[b.N%n])
			}
		})
	}
}

func BenchmarkTest(b *testing.B) {
	sizes := []int{100, 1_000, 10_000}
	for _, n := range sizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			f := benchFilter(b, n)
			data := benchData(n)

			for _, d := range data {
				f.Add(d)
			}

			b.ResetTimer()

			for range b.N {
				f.Test(data[b.N%n])
			}
		})
	}
}

func BenchmarkRemove(b *testing.B) {
	sizes := []int{100, 1_000, 10_000}
	for _, n := range sizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			f := benchFilter(b, n)
			data := benchData(n)

			for _, d := range data {
				f.Add(d)
			}

			b.ResetTimer()

			for range b.N {
				f.Remove(data[b.N%n])
			}
		})
	}
}

func TestDump(t *testing.T) {
	f, err := NewCounting(3000, 0.01)
	testutils.Equal(t, err, nil)

	f.Add(foo)
	f.Add(bar)
	f.Remove(foo)

	testutils.Equal(t, f.Test(foo), false)
	testutils.Equal(t, f.Test(bar), true)

	dump, err := f.ToBytes()
	testutils.Equal(t, err, nil)

	f2, err := NewCountingFromBytes(dump)
	testutils.Equal(t, err, nil)

	testutils.Equal(t, f2.Test(foo), false)
	testutils.Equal(t, f2.Test(bar), true)
	testutils.Equal(t, f2.Test(baz), false)

	f2.Add(baz)

	testutils.Equal(t, f2.Test(baz), true)
}
