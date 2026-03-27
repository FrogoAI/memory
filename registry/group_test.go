package registry

import (
	"errors"
	"testing"
)

func TestGroup(t *testing.T) {
	r := NewRegistry[string, uint64, any]()

	data := []*MockEntity{
		{
			id: r.NextID(),
		},
		{
			id: r.NextID(),
		},
		{
			id: r.NextID(),
		},
	}

	for _, e := range data {
		err := r.GetGroup(KeyTemporary).Add(r.NextID(), e)
		if err != nil {
			t.Fatal(err)
		}
	}

	result := r.GetValues(KeyTemporary)

	if len(data) != len(result) {
		t.Fatalf("Unexpected objects count: %d != %d", len(data), len(result))
	}

	for _, item := range result {
		var exists bool

		e := item.(*MockEntity)
		for _, i := range data {
			if i == e {
				exists = true
			}
		}

		if !exists {
			t.Fatalf("Not found expected object: %+v", e)
		}
	}

	r.DeleteGroup(KeyTemporary)

	result = r.GetValues(KeyTemporary)
	if len(result) != 0 {
		t.Fatalf("Unexpected objects count: %d", len(result))
	}

	for _, e := range data {
		err := r.GetGroup(KeyTemporary).Add(r.NextID(), e)
		if err != nil {
			t.Fatal(err)
		}
	}

	r.ClearGroup(KeyTemporary)

	result = r.GetValues(KeyTemporary)
	if len(result) != 0 {
		t.Fatalf("Unexpected objects count: %d", len(result))
	}

	res := r.SearchInGroup(KeyTemporary, func(_ interface{}, id interface{}, _ interface{}) bool {
		return id.(uint32) == 1
	})
	for e := range res {
		var exists bool

		entity := e.(*MockEntity)
		for _, i := range data {
			if i == entity {
				exists = true
			}
		}

		if !exists {
			t.Fatalf("Not found expected object: %+v", e)
		}
	}
}

func TestGroupSize(t *testing.T) {
	cases := []struct {
		name     string
		addCount int
		want     int
	}{
		{name: "empty group", addCount: 0, want: 0},
		{name: "one item", addCount: 1, want: 1},
		{name: "multiple items", addCount: 5, want: 5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGroup[uint64, string]()

			for i := range tc.addCount {
				if err := g.Add(uint64(i+1), "value"); err != nil {
					t.Fatal(err)
				}
			}

			if g.Size() != tc.want {
				t.Fatalf("Size() = %d, want %d", g.Size(), tc.want)
			}
		})
	}
}

func TestGroupGetKeys(t *testing.T) {
	g := NewGroup[string, int]()

	keys := []string{"alpha", "beta", "gamma"}
	for i, k := range keys {
		if err := g.Add(k, i); err != nil {
			t.Fatal(err)
		}
	}

	got := g.GetKeys()
	if len(got) != len(keys) {
		t.Fatalf("GetKeys() returned %d keys, want %d", len(got), len(keys))
	}

	keySet := make(map[string]bool)
	for _, k := range got {
		keySet[k] = true
	}

	for _, k := range keys {
		if !keySet[k] {
			t.Fatalf("GetKeys() missing key %q", k)
		}
	}
}

func TestGroupIterator(t *testing.T) {
	cases := []struct {
		name  string
		items map[uint64]string
	}{
		{name: "empty group", items: map[uint64]string{}},
		{name: "single item", items: map[uint64]string{1: "one"}},
		{name: "multiple items", items: map[uint64]string{1: "one", 2: "two", 3: "three"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGroup[uint64, string]()

			for id, val := range tc.items {
				if err := g.Add(id, val); err != nil {
					t.Fatal(err)
				}
			}

			var collected []string
			for v := range g.Iterator() {
				collected = append(collected, v)
			}

			if len(collected) != len(tc.items) {
				t.Fatalf("Iterator() yielded %d items, want %d", len(collected), len(tc.items))
			}
		})
	}
}

func TestGroupCallWithLock(t *testing.T) {
	cases := []struct {
		name    string
		items   map[string]int
		fn      func(d map[string]int) (int, error)
		want    int
		wantErr bool
	}{
		{
			name:  "sum values",
			items: map[string]int{"x": 10, "y": 20},
			fn: func(d map[string]int) (int, error) {
				sum := 0
				for _, v := range d {
					sum += v
				}

				return sum, nil
			},
			want: 30,
		},
		{
			name:  "callback error",
			items: map[string]int{},
			fn: func(_ map[string]int) (int, error) {
				return 0, errors.New("callback error")
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGroup[string, int]()

			for k, v := range tc.items {
				if err := g.Add(k, v); err != nil {
					t.Fatal(err)
				}
			}

			got, err := g.CallWithLock(tc.fn)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tc.wantErr && got != tc.want {
				t.Fatalf("CallWithLock() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestGroupTick(t *testing.T) {
	g := NewGroup[uint64, any]()

	entities := []*MockEntity{
		{id: 1},
		{id: 2},
		{id: 3},
	}

	for _, e := range entities {
		if err := g.Add(e.id, e); err != nil {
			t.Fatal(err)
		}
	}

	g.Tick()

	for _, e := range entities {
		if e.counter != 1 {
			t.Fatalf("entity %d: counter = %d, want 1", e.id, e.counter)
		}
	}

	g.Tick()

	for _, e := range entities {
		if e.counter != 2 {
			t.Fatalf("entity %d: counter = %d, want 2 after second tick", e.id, e.counter)
		}
	}
}

func TestGroupTickNonTicker(t *testing.T) {
	g := NewGroup[uint64, any]()

	if err := g.Add(1, "not a ticker"); err != nil {
		t.Fatal(err)
	}

	g.Tick()
}
