package utils

import (
	"hash/crc32"
)

func CRC32(field string) uint32 {
	return crc32.Checksum([]byte(field), crc32.MakeTable(crc32.Castagnoli))
}

func CRC16(field string) uint16 {
	return crc16XModem([]byte(field))
}

// crc16XModem computes a CRC-16 checksum using the XModem variant
// (polynomial 0x1021, initial value 0x0000, no bit reversal).
func crc16XModem(data []byte) uint16 {
	var crc uint16

	for _, b := range data {
		crc ^= uint16(b) << 8

		for range 8 {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}

	return crc
}
