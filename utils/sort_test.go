package utils

import (
	"math/rand"
	"testing"

	"github.com/FrogoAI/memory/comparator"
)

func TestSortInts(t *testing.T) {
	ints := []interface{}{}
	ints = append(ints, 4)
	ints = append(ints, 1)
	ints = append(ints, 2)
	ints = append(ints, 3)

	if err := Sort(ints, comparator.IntComparator); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 1; i < len(ints); i++ {
		if ints[i-1].(int) > ints[i].(int) {
			t.Errorf("Not sorted!")
		}
	}
}

func TestSortStrings(t *testing.T) {
	strings := []interface{}{}
	strings = append(strings, "d")
	strings = append(strings, "a")
	strings = append(strings, "b")
	strings = append(strings, "c")

	if err := Sort(strings, comparator.StringComparator); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 1; i < len(strings); i++ {
		if strings[i-1].(string) > strings[i].(string) {
			t.Errorf("Not sorted!")
		}
	}
}

func TestSortStructs(t *testing.T) {
	type User struct {
		id   int
		name string
	}

	byID := func(a, b interface{}) (int, error) {
		c1 := a.(User)
		c2 := b.(User)

		switch {
		case c1.id > c2.id:
			return 1, nil
		case c1.id < c2.id:
			return -1, nil
		default:
			return 0, nil
		}
	}

	users := []interface{}{
		User{4, "d"},
		User{1, "a"},
		User{3, "c"},
		User{2, "b"},
	}

	if err := Sort(users, byID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 1; i < len(users); i++ {
		if users[i-1].(User).id > users[i].(User).id {
			t.Errorf("Not sorted!")
		}
	}
}

func TestSortRandom(t *testing.T) {
	ints := []interface{}{}
	for i := 0; i < 10000; i++ {
		ints = append(ints, rand.Int()) //nolint:gosec
	}

	if err := Sort(ints, comparator.IntComparator); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 1; i < len(ints); i++ {
		if ints[i-1].(int) > ints[i].(int) {
			t.Errorf("Not sorted!")
		}
	}
}

func BenchmarkGoSortRandom(b *testing.B) {
	b.StopTimer()

	ints := []interface{}{}
	for i := 0; i < 100000; i++ {
		ints = append(ints, rand.Int()) //nolint:gosec
	}

	b.StartTimer()

	_ = Sort(ints, comparator.IntComparator)

	b.StopTimer()
}
