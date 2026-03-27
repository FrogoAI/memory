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

func TestComparatorInvalidType(t *testing.T) {
	cases := []struct {
		name string
		cmp  Comparator
		a, b interface{}
	}{
		{"int with string", IntComparator, "a", "b"},
		{"string with int", StringComparator, 1, 2},
		{"int8 with int", Int8Comparator, 1, 2},
		{"float64 with string", Float64Comparator, "a", "b"},
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
