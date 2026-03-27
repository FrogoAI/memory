package bloom

import "errors"

// Bloom filter errors.
var (
	ErrEmptyDump    = errors.New("empty dump")
	ErrBitsetTooBig = errors.New("bitset overflow: parameters require more than 2^32 bits")
)
