package utils

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestCRC32(t *testing.T) {
	testutils.Equal(t, CRC32("test1"), uint32(1409163093))
	testutils.Equal(t, CRC32("test2"), uint32(1085205665))
	testutils.Equal(t, CRC32("true"), uint32(151551613))
	testutils.Equal(t, CRC32("false"), uint32(118305666))
}

func TestCRC16(t *testing.T) {
	testutils.Equal(t, CRC16("test1"), uint16(4768))
	testutils.Equal(t, CRC16("test2"), uint16(8899))
	testutils.Equal(t, CRC16("true"), uint16(62787))
	testutils.Equal(t, CRC16("false"), uint16(29756))
}
