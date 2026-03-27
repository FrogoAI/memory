// Package bloom provides a counting Bloom filter for probabilistic
// membership testing with element removal support.
//
// NOT safe for concurrent use. Callers must synchronize access externally.
package bloom

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"hash"
	"hash/fnv"
	"math"

	"github.com/FrogoAI/packer"
)

// Compile-time interface checks.
var (
	_ encoding.BinaryMarshaler   = (*CountingFilter)(nil)
	_ encoding.BinaryUnmarshaler = (*CountingFilter)(nil)
)

const (
	wordAlignPadding = 31 // bitsPerWord(32) - 1, for rounding up to 32-bit boundary
	wordAlignShift   = 5  // log2(bitsPerWord), for dividing by 32
)

var ln2 = math.Log(2) //nolint:gochecknoglobals,mnd // precomputed ln(2) mathematical constant

type filter struct {
	m uint32
	k uint32
	h hash.Hash64
}

func (f *filter) bits(data []byte) []uint32 {
	f.h.Reset()

	if _, err := f.h.Write(data); err != nil {
		return nil
	}

	d := f.h.Sum(nil)
	a := binary.BigEndian.Uint32(d[4:8])
	b := binary.BigEndian.Uint32(d[0:4])
	is := make([]uint32, f.k)

	for i := uint32(0); i < f.k; i++ {
		is[i] = (a + b*i) % f.m
	}

	return is
}

func newFilter(m, k uint32) *filter {
	return &filter{
		m: m,
		k: k,
		h: fnv.New64(),
	}
}

func estimates(n uint32, p float64) (uint32, uint32, error) {
	nf := float64(n)
	m := -1 * nf * math.Log(p) / (ln2 * ln2)
	k := math.Ceil(ln2 * m / nf)

	words := m + wordAlignPadding>>wordAlignShift
	if words >= math.MaxInt32 || m > math.MaxUint32 {
		return 0, 0, fmt.Errorf("%w: n=%d p=%f requires %.0f bits", ErrBitsetTooBig, n, p, m)
	}

	return uint32(m), uint32(k), nil
}

// CountingFilter is a probabilistic data structure that supports
// both membership testing and element removal.
type CountingFilter struct {
	*filter
	counters []byte
}

// NewCounting creates an optimized counting bloom filter.
// It returns ErrBitsetTooBig when n and p require more bits than a 32-bit filter supports.
func NewCounting(n int, p float64) (*CountingFilter, error) {
	//nolint:gosec // n is validated by estimates(); negative n wraps to large uint32 and is caught
	m, k, err := estimates(uint32(n), p)
	if err != nil {
		return nil, err
	}

	return &CountingFilter{
		filter:   newFilter(m, k),
		counters: make([]byte, m),
	}, nil
}

// Test checks if an item is likely in the set.
// It returns true if the item might be in the filter, false if it is definitely not.
func (f *CountingFilter) Test(data []byte) bool {
	for _, i := range f.bits(data) {
		if f.counters[i] == 0 {
			return false
		}
	}

	return true
}

// Add inserts data into the filter by incrementing its counters.
// It protects against counter overflow (stopping at 255).
func (f *CountingFilter) Add(data []byte) {
	for _, i := range f.bits(data) {
		// Prevent overflow
		if f.counters[i] < math.MaxUint8 {
			f.counters[i]++
		}
	}
}

// Remove deletes data from the filter by decrementing its counters.
// It protects against counter underflow (stopping at 0).
func (f *CountingFilter) Remove(data []byte) {
	for _, i := range f.bits(data) {
		// Prevent underflow
		if f.counters[i] > 0 {
			f.counters[i]--
		}
	}
}

// Copy returns a deep copy of the CountingFilter.
func (f *CountingFilter) Copy() *CountingFilter {
	counters := make([]byte, len(f.counters))
	copy(counters, f.counters)

	return &CountingFilter{
		filter:   newFilter(f.m, f.k),
		counters: counters,
	}
}

// Clear resets all counters in the filter.
func (f *CountingFilter) Clear() {
	for i := range f.counters {
		f.counters[i] = 0
	}
}

// ToBytes serializes the CountingFilter to a byte slice.
func (f *CountingFilter) ToBytes() ([]byte, error) {
	enc := packer.NewBinaryEncoder()

	err := enc.Encode(f.m)
	if err != nil {
		return nil, fmt.Errorf("encode m param: %w", err)
	}

	err = enc.Encode(f.k)
	if err != nil {
		return nil, fmt.Errorf("encode k param: %w", err)
	}

	err = enc.Encode(f.counters)
	if err != nil {
		return nil, fmt.Errorf("encode counters param: %w", err)
	}

	return enc.Bytes(), nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (f *CountingFilter) MarshalBinary() ([]byte, error) {
	return f.ToBytes()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (f *CountingFilter) UnmarshalBinary(data []byte) error {
	restored, err := NewCountingFromBytes(data)
	if err != nil {
		return err
	}

	*f = *restored

	return nil
}

// NewCountingFromBytes deserializes a byte slice into a CountingFilter.
func NewCountingFromBytes(data []byte) (*CountingFilter, error) {
	if len(data) == 0 {
		return nil, ErrEmptyDump
	}

	dec := packer.NewBinaryDecoder(data)

	res := &CountingFilter{
		filter: &filter{
			h: fnv.New64(),
		},
	}

	err := dec.Decode(&res.m)
	if err != nil {
		return nil, fmt.Errorf("decode m param: %w", err)
	}

	err = dec.Decode(&res.k)
	if err != nil {
		return nil, fmt.Errorf("decode k param: %w", err)
	}

	err = dec.Decode(&res.counters)
	if err != nil {
		return nil, fmt.Errorf("decode counters param: %w", err)
	}

	return res, nil
}
