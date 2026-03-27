package hll

import (
	"fmt"
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestRestore(t *testing.T) {
	h1, err := New()
	testutils.Equal(t, err, nil)

	for i := 0; i < 10000; i++ {
		err = h1.AddAny(i)
		testutils.Equal(t, err, nil)
	}

	testutils.Equal(t, h1.Count(), uint64(10000))

	dump := h1.ToBytes()

	h1, err = FromBytes(dump)
	testutils.Equal(t, err, nil)

	testutils.Equal(t, h1.Count(), uint64(10000))
}

func TestAdd(t *testing.T) {
	cases := []struct {
		name  string
		count int
	}{
		{name: "single item", count: 1},
		{name: "hundred items", count: 100},
		{name: "ten thousand items", count: 10000},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := New()
			testutils.Equal(t, err, nil)

			for i := range tc.count {
				h.Add([]byte(fmt.Sprintf("key-%d", i)))
			}

			testutils.Equal(t, h.Count(), uint64(tc.count))
		})
	}
}

func TestLen(t *testing.T) {
	cases := []struct {
		name string
		n    int
	}{
		{name: "empty", n: 0},
		{name: "single", n: 1},
		{name: "thousand", n: 1000},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := New()
			testutils.Equal(t, err, nil)

			for i := range tc.n {
				h.Add([]byte(fmt.Sprintf("item-%d", i)))
			}

			testutils.Equal(t, h.Len(), tc.n)
		})
	}
}

func TestClear(t *testing.T) {
	h, err := New()
	testutils.Equal(t, err, nil)

	for i := range 1000 {
		h.Add([]byte(fmt.Sprintf("item-%d", i)))
	}

	testutils.Equal(t, h.Count(), uint64(1000))

	err = h.Clear()
	testutils.Equal(t, err, nil)
	testutils.Equal(t, h.Count(), uint64(0))
	testutils.Equal(t, h.Len(), 0)
}

func TestSerializationRoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		count int
	}{
		{name: "empty", count: 0},
		{name: "single item", count: 1},
		{name: "large set", count: 100000},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := New()
			testutils.Equal(t, err, nil)

			for i := range tc.count {
				h.Add([]byte(fmt.Sprintf("ser-%d", i)))
			}

			raw := h.ToBytes()
			restored, err := FromBytes(raw)
			testutils.Equal(t, err, nil)
			testutils.Equal(t, restored.Count(), h.Count())
		})
	}
}

func TestMarshalBinaryRoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		count int
	}{
		{name: "empty", count: 0},
		{name: "single item", count: 1},
		{name: "large set", count: 100000},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := New()
			testutils.Equal(t, err, nil)

			for i := range tc.count {
				h.Add([]byte(fmt.Sprintf("marshal-%d", i)))
			}

			data, err := h.MarshalBinary()
			testutils.Equal(t, err, nil)

			restored := &HyperLogLog{}
			err = restored.UnmarshalBinary(data)
			testutils.Equal(t, err, nil)
			testutils.Equal(t, restored.Count(), h.Count())
		})
	}
}

func TestUnmarshalBinary_InvalidData(t *testing.T) {
	h := &HyperLogLog{}

	err := h.UnmarshalBinary([]byte{0xFF, 0xFE, 0xFD})
	if err == nil {
		t.Error("expected error from invalid bytes, got nil")
	}
}

func TestUnionEdgeCases(t *testing.T) {
	t.Run("union with empty", func(t *testing.T) {
		h1, err := New()
		testutils.Equal(t, err, nil)

		for i := range 1000 {
			h1.Add([]byte(fmt.Sprintf("u-%d", i)))
		}

		empty, err := New()
		testutils.Equal(t, err, nil)

		u, err := h1.Union(empty)
		testutils.Equal(t, err, nil)
		testutils.Equal(t, u.Count(), uint64(1000))
	})

	t.Run("union disjoint sets", func(t *testing.T) {
		h1, err := New()
		testutils.Equal(t, err, nil)

		h2, err := New()
		testutils.Equal(t, err, nil)

		for i := range 5000 {
			h1.Add([]byte(fmt.Sprintf("a-%d", i)))
			h2.Add([]byte(fmt.Sprintf("b-%d", i)))
		}

		count, err := h1.UnionCount(h2)
		testutils.Equal(t, err, nil)
		testutils.Equal(t, count, uint64(10000))
	})
}

func TestIntersectionEdgeCases(t *testing.T) {
	t.Run("disjoint sets", func(t *testing.T) {
		h1, err := New()
		testutils.Equal(t, err, nil)

		h2, err := New()
		testutils.Equal(t, err, nil)

		for i := range 5000 {
			h1.Add([]byte(fmt.Sprintf("x-%d", i)))
			h2.Add([]byte(fmt.Sprintf("y-%d", i)))
		}

		count, err := h1.IntersectionCount(h2)
		testutils.Equal(t, err, nil)
		testutils.Equal(t, count, uint64(0))
	})

	t.Run("identical sets", func(t *testing.T) {
		h1, err := New()
		testutils.Equal(t, err, nil)

		h2, err := New()
		testutils.Equal(t, err, nil)

		for i := range 5000 {
			data := []byte(fmt.Sprintf("z-%d", i))
			h1.Add(data)
			h2.Add(data)
		}

		count, err := h1.IntersectionCount(h2)
		testutils.Equal(t, err, nil)
		testutils.Equal(t, count, uint64(5000))
	})
}

func TestLargeCardinality(t *testing.T) {
	h, err := New()
	testutils.Equal(t, err, nil)

	n := 500000

	for i := range n {
		h.Add([]byte(fmt.Sprintf("big-%d", i)))
	}

	count := h.Count()
	tolerance := uint64(n) / 100 // 1% error tolerance

	if count < uint64(n)-tolerance || count > uint64(n)+tolerance {
		t.Errorf("large cardinality out of tolerance: got %d, want ~%d (±%d)", count, n, tolerance)
	}
}

func TestFromBytesInvalid(t *testing.T) {
	_, err := FromBytes([]byte{0xFF, 0xFE, 0xFD})
	if err == nil {
		t.Error("expected error from invalid bytes, got nil")
	}
}

func TestHLL(t *testing.T) {
	h1, err := New()
	testutils.Equal(t, err, nil)

	h2, err := New()
	testutils.Equal(t, err, nil)

	for i := 0; i < 10000; i++ {
		err = h1.AddAny(i)
		testutils.Equal(t, err, nil)
	}

	for i := 5000; i < 15000; i++ {
		err = h2.AddAny(i)
		testutils.Equal(t, err, nil)
	}

	c, err := h1.IntersectionCount(h2)
	testutils.Equal(t, err, nil)

	u, err := h1.UnionCount(h2)
	testutils.Equal(t, err, nil)

	testutils.Equal(t, c, uint64(5000))
	testutils.Equal(t, h1.Count(), uint64(10000))
	testutils.Equal(t, h2.Count(), uint64(10000))
	testutils.Equal(t, u, uint64(15000))
}

// Prevent compiler optimization of benchmark results.
var benchCount uint64

func benchHLL(b *testing.B, n int) *HyperLogLog {
	b.Helper()

	h, err := New()
	if err != nil {
		b.Fatal(err)
	}

	for i := range n {
		h.Add([]byte(fmt.Sprintf("key-%d", i)))
	}

	return h
}

func BenchmarkAdd(b *testing.B) {
	h, err := New()
	if err != nil {
		b.Fatal(err)
	}

	data := make([][]byte, 10000)
	for i := range data {
		data[i] = []byte(fmt.Sprintf("key-%d", i))
	}

	b.ResetTimer()

	for i := range b.N {
		h.Add(data[i%len(data)])
	}
}

func BenchmarkAddAny(b *testing.B) {
	h, err := New()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := range b.N {
		_ = h.AddAny(i)
	}
}

func BenchmarkCount(b *testing.B) {
	h := benchHLL(b, 10000)

	b.ResetTimer()

	for range b.N {
		benchCount = h.Count()
	}
}

func BenchmarkUnion(b *testing.B) {
	h1 := benchHLL(b, 5000)
	h2 := benchHLL(b, 5000)

	b.ResetTimer()

	for range b.N {
		_, _ = h1.Union(h2)
	}
}

func BenchmarkToBytes(b *testing.B) {
	h := benchHLL(b, 10000)

	b.ResetTimer()

	for range b.N {
		_ = h.ToBytes()
	}
}

func BenchmarkFromBytes(b *testing.B) {
	h := benchHLL(b, 10000)
	raw := h.ToBytes()

	b.ResetTimer()

	for range b.N {
		_, _ = FromBytes(raw)
	}
}
