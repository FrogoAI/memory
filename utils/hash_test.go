//nolint:revive // package name is intentional
package utils

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestCRC32(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  uint32
	}{
		{name: "test1", input: "test1", want: 1409163093},
		{name: "test2", input: "test2", want: 1085205665},
		{name: "true", input: "true", want: 151551613},
		{name: "false", input: "false", want: 118305666},
		{name: "empty_string", input: "", want: 0},
		{name: "standard_check_value", input: "123456789", want: 0xE3069283},
		{name: "single_char", input: "a", want: 3251651376},
		{name: "space", input: " ", want: 1925242255},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, CRC32(tc.input), tc.want)
		})
	}
}

func TestCRC16(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  uint16
	}{
		{name: "test1", input: "test1", want: 4768},
		{name: "test2", input: "test2", want: 8899},
		{name: "true", input: "true", want: 62787},
		{name: "false", input: "false", want: 29756},
		{name: "empty_string", input: "", want: 0},
		{name: "standard_check_value", input: "123456789", want: 0x31C3},
		{name: "single_char_A", input: "A", want: 22757},
		{name: "single_char_a", input: "a", want: 31879},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, CRC16(tc.input), tc.want)
		})
	}
}
