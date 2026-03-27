package hll

import (
	"bytes"
	"encoding/gob"
)

// GetBytes encodes the given value into a byte slice using gob encoding.
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
