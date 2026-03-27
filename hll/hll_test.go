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
