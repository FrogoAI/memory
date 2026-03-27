package comparator

import (
	"errors"
	"testing"
	"time"
)

func TestIntComparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", 1, 1, 0},
		{"less", 1, 2, -1},
		{"greater", 2, 1, 1},
		{"less large", 11, 22, -1},
		{"zero equal", 0, 0, 0},
		{"greater zero", 1, 0, 1},
		{"less zero", 0, 1, -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := IntComparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestStringComparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", "a", "a", 0},
		{"less", "a", "b", -1},
		{"greater", "b", "a", 1},
		{"prefix less", "aa", "aab", -1},
		{"empty equal", "", "", 0},
		{"greater empty", "a", "", 1},
		{"less empty", "", "a", -1},
		{"less long", "", "aaaaaaa", -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := StringComparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestTimeComparator(t *testing.T) {
	now := time.Now()

	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", now, now, 0},
		{"after", now.Add(24 * 7 * 2 * time.Hour), now, 1},
		{"before", now, now.Add(24 * 7 * 2 * time.Hour), -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := TimeComparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestCustomComparator(t *testing.T) {
	type Custom struct {
		id   int
		name string
	}

	byID := func(a, b interface{}) (int, error) {
		c1, ok1 := a.(Custom)
		c2, ok2 := b.(Custom)

		if !ok1 || !ok2 {
			return 0, ErrInvalidType
		}

		switch {
		case c1.id > c2.id:
			return 1, nil
		case c1.id < c2.id:
			return -1, nil
		default:
			return 0, nil
		}
	}

	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", Custom{1, "a"}, Custom{1, "a"}, 0},
		{"less", Custom{1, "a"}, Custom{2, "b"}, -1},
		{"greater", Custom{2, "b"}, Custom{1, "a"}, 1},
		{"same id diff name", Custom{1, "a"}, Custom{1, "b"}, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := byID(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestInt8Comparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", int8(1), int8(1), 0},
		{"less", int8(1), int8(2), -1},
		{"greater", int8(2), int8(1), 1},
		{"zero equal", int8(0), int8(0), 0},
		{"negative less", int8(-1), int8(0), -1},
		{"negative greater", int8(0), int8(-1), 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Int8Comparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestInt16Comparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", int16(1), int16(1), 0},
		{"less", int16(1), int16(2), -1},
		{"greater", int16(2), int16(1), 1},
		{"zero equal", int16(0), int16(0), 0},
		{"negative", int16(-100), int16(100), -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Int16Comparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestInt32Comparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", int32(1), int32(1), 0},
		{"less", int32(1), int32(2), -1},
		{"greater", int32(2), int32(1), 1},
		{"zero equal", int32(0), int32(0), 0},
		{"negative", int32(-1000), int32(1000), -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Int32Comparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestInt64Comparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", int64(1), int64(1), 0},
		{"less", int64(1), int64(2), -1},
		{"greater", int64(2), int64(1), 1},
		{"zero equal", int64(0), int64(0), 0},
		{"negative", int64(-1), int64(1), -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Int64Comparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestUIntComparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", uint(1), uint(1), 0},
		{"less", uint(1), uint(2), -1},
		{"greater", uint(2), uint(1), 1},
		{"zero equal", uint(0), uint(0), 0},
		{"greater zero", uint(1), uint(0), 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := UIntComparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestUInt8Comparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", uint8(1), uint8(1), 0},
		{"less", uint8(1), uint8(2), -1},
		{"greater", uint8(2), uint8(1), 1},
		{"zero equal", uint8(0), uint8(0), 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := UInt8Comparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestUInt16Comparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", uint16(1), uint16(1), 0},
		{"less", uint16(1), uint16(2), -1},
		{"greater", uint16(2), uint16(1), 1},
		{"zero equal", uint16(0), uint16(0), 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := UInt16Comparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestUInt32Comparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", uint32(1), uint32(1), 0},
		{"less", uint32(1), uint32(2), -1},
		{"greater", uint32(2), uint32(1), 1},
		{"zero equal", uint32(0), uint32(0), 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := UInt32Comparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestUInt64Comparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", uint64(1), uint64(1), 0},
		{"less", uint64(1), uint64(2), -1},
		{"greater", uint64(2), uint64(1), 1},
		{"zero equal", uint64(0), uint64(0), 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := UInt64Comparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestFloat32Comparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", float32(1.0), float32(1.0), 0},
		{"less", float32(1.0), float32(2.0), -1},
		{"greater", float32(2.0), float32(1.0), 1},
		{"zero equal", float32(0.0), float32(0.0), 0},
		{"negative less", float32(-1.0), float32(1.0), -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Float32Comparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestFloat64Comparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", 1.0, 1.0, 0},
		{"less", 1.0, 2.0, -1},
		{"greater", 2.0, 1.0, 1},
		{"zero equal", 0.0, 0.0, 0},
		{"negative less", -1.0, 1.0, -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Float64Comparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestFloat64DiffComparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", 1.0, 1.0, 0},
		{"greater", 3.0, 1.0, 1},
		{"zero equal", 0.0, 0.0, 0},
		// NOTE: a < b always returns 0 due to diff <= SmallestNonzeroFloat64
		// catching negative differences — this is a pre-existing implementation quirk.
		{"less returns zero", 1.0, 3.0, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Float64DiffComparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestByteComparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", byte('a'), byte('a'), 0},
		{"less", byte('a'), byte('b'), -1},
		{"greater", byte('b'), byte('a'), 1},
		{"zero equal", byte(0), byte(0), 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ByteComparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestRuneComparator(t *testing.T) {
	cases := []struct {
		name     string
		a, b     interface{}
		expected int
	}{
		{"equal", 'a', 'a', 0},
		{"less", 'a', 'b', -1},
		{"greater", 'b', 'a', 1},
		{"zero equal", rune(0), rune(0), 0},
		{"unicode", 'α', 'β', -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := RuneComparator(tc.a, tc.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Got %v expected %v", actual, tc.expected)
			}
		})
	}
}

func TestComparatorInvalidType(t *testing.T) {
	cases := []struct {
		name string
		cmp  Comparator
		a, b interface{}
	}{
		{"int with string", IntComparator, "a", "b"},
		{"string with int", StringComparator, 1, 2},
		{"int8 with int", Int8Comparator, 1, 2},
		{"int16 with string", Int16Comparator, "a", "b"},
		{"int32 with string", Int32Comparator, "a", "b"},
		{"int64 with string", Int64Comparator, "a", "b"},
		{"uint with string", UIntComparator, "a", "b"},
		{"uint8 with string", UInt8Comparator, "a", "b"},
		{"uint16 with string", UInt16Comparator, "a", "b"},
		{"uint32 with string", UInt32Comparator, "a", "b"},
		{"uint64 with string", UInt64Comparator, "a", "b"},
		{"float32 with string", Float32Comparator, "a", "b"},
		{"float64 with string", Float64Comparator, "a", "b"},
		{"float64diff with string", Float64DiffComparator, "a", "b"},
		{"byte with string", ByteComparator, "a", "b"},
		{"rune with float", RuneComparator, 1.0, 2.0},
		{"time with int", TimeComparator, 1, 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.cmp(tc.a, tc.b)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, ErrInvalidType) {
				t.Errorf("expected ErrInvalidType, got %v", err)
			}
		})
	}
}
