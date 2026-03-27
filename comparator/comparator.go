package comparator

import (
	"fmt"
	"math"
	"time"
)

// Comparator compares two values and returns:
//
//	negative , if a < b
//	zero     , if a == b
//	positive , if a > b
//
// Returns ErrInvalidType if a or b cannot be asserted to the expected type.
type Comparator func(a, b interface{}) (int, error)

func typeError(expected string, a, b interface{}) error {
	return fmt.Errorf("%w: expected %s, got a=%T b=%T", ErrInvalidType, expected, a, b)
}

// StringComparator provides a fast comparison on strings.
func StringComparator(a, b interface{}) (int, error) {
	s1, ok1 := a.(string)
	s2, ok2 := b.(string)

	if !ok1 || !ok2 {
		return 0, typeError("string", a, b)
	}

	minLen := len(s2)
	if len(s1) < len(s2) {
		minLen = len(s1)
	}

	diff := 0
	for i := 0; i < minLen && diff == 0; i++ {
		diff = int(s1[i]) - int(s2[i])
	}

	if diff == 0 {
		diff = len(s1) - len(s2)
	}

	if diff < 0 {
		return -1, nil
	}

	if diff > 0 {
		return 1, nil
	}

	return 0, nil
}

// IntComparator provides a basic comparison on int.
func IntComparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(int)
	bAsserted, ok2 := b.(int)

	if !ok1 || !ok2 {
		return 0, typeError("int", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// Int8Comparator provides a basic comparison on int8.
func Int8Comparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(int8)
	bAsserted, ok2 := b.(int8)

	if !ok1 || !ok2 {
		return 0, typeError("int8", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// Int16Comparator provides a basic comparison on int16.
func Int16Comparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(int16)
	bAsserted, ok2 := b.(int16)

	if !ok1 || !ok2 {
		return 0, typeError("int16", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// Int32Comparator provides a basic comparison on int32.
func Int32Comparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(int32)
	bAsserted, ok2 := b.(int32)

	if !ok1 || !ok2 {
		return 0, typeError("int32", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// Int64Comparator provides a basic comparison on int64.
func Int64Comparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(int64)
	bAsserted, ok2 := b.(int64)

	if !ok1 || !ok2 {
		return 0, typeError("int64", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// UIntComparator provides a basic comparison on uint.
func UIntComparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(uint)
	bAsserted, ok2 := b.(uint)

	if !ok1 || !ok2 {
		return 0, typeError("uint", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// UInt8Comparator provides a basic comparison on uint8.
func UInt8Comparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(uint8)
	bAsserted, ok2 := b.(uint8)

	if !ok1 || !ok2 {
		return 0, typeError("uint8", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// UInt16Comparator provides a basic comparison on uint16.
func UInt16Comparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(uint16)
	bAsserted, ok2 := b.(uint16)

	if !ok1 || !ok2 {
		return 0, typeError("uint16", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// UInt32Comparator provides a basic comparison on uint32.
func UInt32Comparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(uint32)
	bAsserted, ok2 := b.(uint32)

	if !ok1 || !ok2 {
		return 0, typeError("uint32", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// UInt64Comparator provides a basic comparison on uint64.
func UInt64Comparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(uint64)
	bAsserted, ok2 := b.(uint64)

	if !ok1 || !ok2 {
		return 0, typeError("uint64", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// Float32Comparator provides a basic comparison on float32.
func Float32Comparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(float32)
	bAsserted, ok2 := b.(float32)

	if !ok1 || !ok2 {
		return 0, typeError("float32", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// Float64Comparator provides a basic comparison on float64.
func Float64Comparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(float64)
	bAsserted, ok2 := b.(float64)

	if !ok1 || !ok2 {
		return 0, typeError("float64", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// Float64DiffComparator provides a basic comparison on float64
// using a tolerance of math.SmallestNonzeroFloat64.
func Float64DiffComparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(float64)
	bAsserted, ok2 := b.(float64)

	if !ok1 || !ok2 {
		return 0, typeError("float64", a, b)
	}

	diff := aAsserted - bAsserted

	switch {
	case diff <= math.SmallestNonzeroFloat64:
		return 0, nil
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// ByteComparator provides a basic comparison on byte.
func ByteComparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(byte)
	bAsserted, ok2 := b.(byte)

	if !ok1 || !ok2 {
		return 0, typeError("byte", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// RuneComparator provides a basic comparison on rune.
func RuneComparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(rune)
	bAsserted, ok2 := b.(rune)

	if !ok1 || !ok2 {
		return 0, typeError("rune", a, b)
	}

	switch {
	case aAsserted > bAsserted:
		return 1, nil
	case aAsserted < bAsserted:
		return -1, nil
	default:
		return 0, nil
	}
}

// TimeComparator provides a basic comparison on time.Time.
func TimeComparator(a, b interface{}) (int, error) {
	aAsserted, ok1 := a.(time.Time)
	bAsserted, ok2 := b.(time.Time)

	if !ok1 || !ok2 {
		return 0, typeError("time.Time", a, b)
	}

	switch {
	case aAsserted.After(bAsserted):
		return 1, nil
	case aAsserted.Before(bAsserted):
		return -1, nil
	default:
		return 0, nil
	}
}
