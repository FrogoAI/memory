package comparator

import "errors"

// ErrInvalidType is returned when a comparator receives a value
// that cannot be type-asserted to the expected type.
var ErrInvalidType = errors.New("invalid type for comparator")
