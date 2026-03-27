// Package utils provides hashing, sorting, and string helpers (stateless, safe for concurrent use)
// as well as SafeMap and SafeList which are thread-safe via read-write mutexes.
package utils //nolint:revive // package name is intentional

import (
	"hash/crc32"
)

// CRC32 returns the CRC-32 checksum of field using the Castagnoli polynomial.
func CRC32(field string) uint32 {
	return crc32.Checksum([]byte(field), crc32.MakeTable(crc32.Castagnoli))
}

// CRC16 returns the CRC-16/XModem checksum of field.
func CRC16(field string) uint16 {
	return crc16XModem([]byte(field))
}

const (
	crc16BitsPerByte = 8      // number of bits processed per byte
	crc16Polynomial  = 0x1021 // CRC-16/XModem polynomial (x^16 + x^12 + x^5 + 1)
)

// crc16XModem computes a CRC-16 checksum using the XModem variant
// (polynomial 0x1021, initial value 0x0000, no bit reversal).
func crc16XModem(data []byte) uint16 {
	var crc uint16

	for _, b := range data {
		crc ^= uint16(b) << crc16BitsPerByte

		for range crc16BitsPerByte {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ crc16Polynomial
			} else {
				crc <<= 1
			}
		}
	}

	return crc
}
