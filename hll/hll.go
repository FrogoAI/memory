// Package hll provides a HyperLogLog probabilistic cardinality estimator.
//
// NOT safe for concurrent use. Callers must synchronize access externally.
package hll

import (
	"encoding"

	"github.com/segmentio/go-hll"
	"github.com/twmb/murmur3"
)

// Compile-time interface checks.
var (
	_ encoding.BinaryMarshaler   = (*HyperLogLog)(nil)
	_ encoding.BinaryUnmarshaler = (*HyperLogLog)(nil)
)

// HyperLogLog configuration constants.
const (
	DefaultLog2m    = 31
	DefaultRegwidth = 8
)

// DefaultSettings is the default HyperLogLog configuration used by New.
var DefaultSettings = hll.Settings{
	Log2m:             DefaultLog2m,
	Regwidth:          DefaultRegwidth,
	ExplicitThreshold: hll.AutoExplicitThreshold,
	SparseEnabled:     true,
}

// HyperLogLog is a probabilistic cardinality estimator.
type HyperLogLog struct {
	HLL hll.Hll
}

// New creates a new HyperLogLog with DefaultSettings.
func New() (*HyperLogLog, error) {
	h, err := hll.NewHll(DefaultSettings)
	if err != nil {
		return nil, err
	}

	return &HyperLogLog{
		HLL: h,
	}, nil
}

// FromBytes deserializes a HyperLogLog from its binary representation.
func FromBytes(raw []byte) (*HyperLogLog, error) {
	h, err := hll.FromBytes(raw)
	if err != nil {
		return nil, err
	}

	return &HyperLogLog{
		HLL: h,
	}, nil
}

// Union returns a new HyperLogLog that is the union of h and h2.
func (h *HyperLogLog) Union(h2 *HyperLogLog) (*HyperLogLog, error) {
	hc, err := hll.NewHll(DefaultSettings)
	if err != nil {
		return nil, err
	}

	hc.Union(h.HLL)
	hc.Union(h2.HLL)

	return &HyperLogLog{HLL: hc}, nil
}

// UnionCount returns the estimated cardinality of the union of h and h2.
func (h *HyperLogLog) UnionCount(h2 *HyperLogLog) (uint64, error) {
	hc, err := h.Union(h2)
	if err != nil {
		return 0, err
	}

	return hc.Count(), nil
}

// Add hashes the given byte slice and adds it to the HyperLogLog.
func (h *HyperLogLog) Add(data []byte) {
	h.HLL.AddRaw(murmur3.Sum64(data))
}

// AddAny gob-encodes the given value and adds it to the HyperLogLog.
func (h *HyperLogLog) AddAny(data any) error {
	v, err := GetBytes(data)
	if err != nil {
		return err
	}

	h.HLL.AddRaw(murmur3.Sum64(v))

	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (h *HyperLogLog) MarshalBinary() ([]byte, error) {
	return h.ToBytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (h *HyperLogLog) UnmarshalBinary(data []byte) error {
	restored, err := FromBytes(data)
	if err != nil {
		return err
	}

	*h = *restored

	return nil
}

// ToBytes serializes the HyperLogLog into its binary representation.
func (h *HyperLogLog) ToBytes() []byte {
	return h.HLL.ToBytes()
}

// Count returns the estimated number of distinct items as a uint64.
func (h *HyperLogLog) Count() uint64 {
	return h.HLL.Cardinality()
}

// Len returns the estimated number of distinct items in the set.
func (h *HyperLogLog) Len() int {
	//nolint:gosec // cardinality is a probabilistic estimate; overflow beyond MaxInt is acceptable
	return int(h.HLL.Cardinality())
}

// Clear resets the HyperLogLog to an empty state.
func (h *HyperLogLog) Clear() error {
	fresh, err := hll.NewHll(DefaultSettings)
	if err != nil {
		return err
	}

	h.HLL = fresh

	return nil
}

// IntersectionCount estimates the number of distinct items common to h and h2 using inclusion-exclusion.
func (h *HyperLogLog) IntersectionCount(h2 *HyperLogLog) (uint64, error) {
	hc, err := hll.NewHll(DefaultSettings)
	if err != nil {
		return 0, err
	}

	hc.Union(h.HLL)
	cardH := h.HLL.Cardinality()

	hc.Union(h2.HLL)
	cardH2 := h2.HLL.Cardinality()

	return cardH + cardH2 - hc.Cardinality(), nil
}
