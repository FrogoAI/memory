package registry

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/FrogoAI/testutils"
)

const KeyTemporary = "test"

type MockEntity struct {
	counter int
	id      uint64
	log     []string
	x, y    int
}

func (m *MockEntity) GetID() uint64 {
	return m.id
}

func (m *MockEntity) Tick() {
	m.counter++
}

func (m *MockEntity) Construct() error {
	m.log = append(m.log, "calling construct")
	return nil
}

func (m *MockEntity) Destroy() error {
	m.log = append(m.log, "calling destroy")
	return nil
}

func (m *MockEntity) Coordinates() (int, int) {
	return m.x, m.y
}

func TestRegistry(t *testing.T) {
	r := NewRegistry[string, uint64, any]()
	testcases := map[string]struct {
		id     uint64
		key    string
		data   interface{}
		setErr error
		getErr error
		result interface{}
	}{
		"simple_usage": {
			id:     r.NextID(),
			key:    "string_key",
			data:   "simple string",
			setErr: nil,
			getErr: nil,
			result: "simple string",
		},
		"empty_key": {
			id:     r.NextID(),
			key:    "",
			data:   "simple string",
			setErr: nil,
			getErr: nil,
			result: "simple string",
		},
		"empty_id": {
			id:     0,
			key:    "string_key2",
			data:   "simple string",
			setErr: nil,
			getErr: nil,
			result: "simple string",
		},
		"struct_store": {
			id:  r.NextID(),
			key: "string_key3",
			data: &MockEntity{
				id: 2,
			},
			setErr: nil,
			getErr: nil,
			result: &MockEntity{
				id: 2,
			},
		},
	}

	for name, test := range testcases {
		test := test

		t.Run(name, func(t *testing.T) {
			err := r.Add(test.key, test.id, test.data)
			testutils.Equal(t, err, test.setErr)
			data, err = r.Get(test.key, test.id)
			testutils.Equal(t, err, test.getErr)
			testutils.Equal(t, data, test.data)
		})
	}
}

func TestConstructDestroy(t *testing.T) {
	r := NewRegistry[string, uint64, any]()
	d := &MockEntity{
		id: r.NextID(),
	}

	err := r.Add(KeyTemporary, d.id, d)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Add(KeyTemporary, d.id, d) // UPDATE
	if err != nil {
		t.Fatal(err)
	}

	err = r.Remove(KeyTemporary, d.id)
	if err != nil {
		t.Fatal(err)
	}

	result := strings.Join(d.log, "\n")
	if result != "calling construct\ncalling construct\ncalling destroy" {
		t.Fatalf("Unexpected result: %s", result)
	}
}

func TestAsyncGetNextID(t *testing.T) {
	r := NewRegistry[uint64, string, any]()
	generate := make(chan struct{}, 100)

	var wg sync.WaitGroup
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func() {
			for range generate {
				r.NextID()
			}

			wg.Done()
		}()
	}

	for i := 0; i < 100; i++ {
		generate <- struct{}{}
	}

	close(generate)
	wg.Wait()

	id := r.NextID()
	if id != 101 {
		t.Fatalf("Unexpected next id: %d", id)
	}
}

func TestAsyncRegistry(t *testing.T) {
	r := NewRegistry[string, uint64, any]()
	writer := make(chan *MockEntity, 100)

	gorutines := 100

	var wg sync.WaitGroup
	wg.Add(gorutines)

	for i := 0; i < gorutines; i++ {
		go func() {
			for entity := range writer {
				err := r.Add(KeyTemporary, entity.GetID(), entity)
				if err != nil {
					panic(err)
				}
			}

			wg.Done()
		}()
	}

	for i := 0; i < 100; i++ {
		e := &MockEntity{
			id: r.NextID(),
		}
		writer <- e
	}

	close(writer)
	wg.Wait()

	count := len(r.GetGroup(KeyTemporary).GetValues())
	if count != 100 {
		t.Fatalf("Unexpected count: %d", count)
	}

	e := r.SearchOne(KeyTemporary, func(_ interface{}, id interface{}, _ interface{}) bool {
		return id.(uint64) == 1
	})

	entity, ok := e.(*MockEntity)
	if !ok {
		t.Fatalf("Unexpected type: %+v", e)
	}

	result := strings.Join(entity.log, "\n")
	if result != "calling construct" {
		t.Fatalf("Unexpected result: %s", result)
	}

	if entity.GetID() != 1 {
		t.Fatalf("Unexpected id: %d", entity.GetID())
	}
}

func TestConcurrentAddGetRemove(t *testing.T) {
	r := NewRegistry[string, uint64, string]()

	const (
		goroutines = 20
		ops        = 50
	)

	var wg sync.WaitGroup
	wg.Add(goroutines * 3)

	for i := range goroutines {
		go func(base int) {
			defer wg.Done()

			for j := range ops {
				id := uint64(base*ops + j + 1)
				if err := r.Add(KeyTemporary, id, "value"); err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		}(i)
	}

	for range goroutines {
		go func() {
			defer wg.Done()

			for j := range ops {
				if _, err := r.Get(KeyTemporary, uint64(j+1)); err != nil {
					// Concurrent removes may cause not-found errors; ignore them.
					_ = err
				}
			}
		}()
	}

	for i := range goroutines {
		go func(base int) {
			defer wg.Done()

			for j := range ops {
				id := uint64(base*ops + j + 1)
				if err := r.Remove(KeyTemporary, id); err != nil {
					// Concurrent removes may race; ignore not-found errors.
					_ = err
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestConcurrentGroupCreateDelete(t *testing.T) {
	r := NewRegistry[int, uint64, any]()

	const goroutines = 20

	var wg sync.WaitGroup
	wg.Add(goroutines * 3)

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()

			err := r.Add(i%5, uint64(i+1), &MockEntity{id: uint64(i + 1)})
			if err != nil {
				t.Errorf("add failed: %v", err)
			}
		}(i)
	}

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()

			r.DeleteGroup(i % 5)
		}(i)
	}

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()

			r.GetGroup(i % 5)
		}(i)
	}

	wg.Wait()
}

func TestConcurrentIndex(t *testing.T) {
	r := NewRegistry[string, uint64, any]()

	const (
		goroutines = 20
		ops        = 50
	)

	var wg sync.WaitGroup
	wg.Add(goroutines * 3)

	for i := range goroutines {
		go func(base int) {
			defer wg.Done()

			for j := range ops {
				r.AddIndex(uint64(base*ops+j), uint64(j))
			}
		}(i)
	}

	for range goroutines {
		go func() {
			defer wg.Done()

			for j := range ops {
				r.GetIndex(uint64(j))
			}
		}()
	}

	for range goroutines {
		go func() {
			defer wg.Done()

			for j := range ops {
				r.RemIndex(uint64(j))
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentGetKeysAndSize(t *testing.T) {
	r := NewRegistry[int, uint64, any]()

	const goroutines = 20

	var wg sync.WaitGroup
	wg.Add(goroutines * 3)

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()

			err := r.Add(i, uint64(i+1), &MockEntity{id: uint64(i + 1)})
			if err != nil {
				t.Errorf("add failed: %v", err)
			}
		}(i)
	}

	for range goroutines {
		go func() {
			defer wg.Done()

			r.GetKeys()
		}()
	}

	for range goroutines {
		go func() {
			defer wg.Done()

			r.Size()
		}()
	}

	wg.Wait()
}

func TestConcurrentIDOperations(t *testing.T) {
	r := NewRegistry[string, uint64, any]()

	const (
		goroutines = 20
		ops        = 100
	)

	var wg sync.WaitGroup
	wg.Add(goroutines * 3)

	for range goroutines {
		go func() {
			defer wg.Done()

			for range ops {
				r.NextID()
			}
		}()
	}

	for range goroutines {
		go func() {
			defer wg.Done()

			for range ops {
				r.LatestID()
			}
		}()
	}

	for range goroutines {
		go func() {
			defer wg.Done()

			for range ops {
				r.SetLatestID(42)
			}
		}()
	}

	wg.Wait()

	if r.LatestID() == 0 {
		t.Fatal("expected non-zero latest ID after concurrent operations")
	}
}

func TestConcurrentGetValues(t *testing.T) {
	r := NewRegistry[string, uint64, any]()

	for i := range 50 {
		err := r.Add(KeyTemporary, uint64(i+1), &MockEntity{id: uint64(i + 1)})
		if err != nil {
			t.Fatal(err)
		}
	}

	const goroutines = 20

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for range goroutines {
		go func() {
			defer wg.Done()

			r.GetValues(KeyTemporary)
		}()
	}

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()

			id := uint64(1000 + i + 1)

			err := r.Add(KeyTemporary, id, &MockEntity{id: id})
			if err != nil {
				t.Errorf("add failed: %v", err)
			}
		}(i)
	}

	wg.Wait()
}

func TestConcurrentRemoveIDEverywhere(t *testing.T) {
	r := NewRegistry[string, uint64, any]()
	groups := []string{"g1", "g2", "g3", "g4", "g5"}

	for _, g := range groups {
		for i := range 20 {
			err := r.Add(g, uint64(i+1), &MockEntity{id: uint64(i + 1)})
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	const goroutines = 10

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()

			if err := r.RemoveIDEverywhere(uint64(i + 1)); err != nil {
				// Concurrent removes may race; ignore not-found errors.
				_ = err
			}
		}(i)
	}

	for range goroutines {
		go func() {
			defer wg.Done()

			for _, g := range groups {
				r.GetValues(g)
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentClearGroup(t *testing.T) {
	r := NewRegistry[string, uint64, any]()

	const goroutines = 20

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()

			err := r.Add(KeyTemporary, uint64(i+1), &MockEntity{id: uint64(i + 1)})
			if err != nil {
				t.Errorf("add failed: %v", err)
			}
		}(i)
	}

	for range goroutines {
		go func() {
			defer wg.Done()

			r.ClearGroup(KeyTemporary)
		}()
	}

	wg.Wait()
}

func TestConcurrentGetGroups(t *testing.T) {
	r := NewRegistry[int, uint64, any]()

	const goroutines = 20

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()

			err := r.Add(i%5, uint64(i+1), &MockEntity{id: uint64(i + 1)})
			if err != nil {
				t.Errorf("add failed: %v", err)
			}
		}(i)
	}

	for range goroutines {
		go func() {
			defer wg.Done()

			r.GetGroups(0, 1, 2, 3, 4)
		}()
	}

	wg.Wait()
}

type MockConstructError struct{}

func (m *MockConstructError) Construct() error {
	return errors.New("construct failed")
}

func TestGetNotFound(t *testing.T) {
	cases := []struct {
		name    string
		key     string
		id      uint64
		wantErr error
	}{
		{name: "missing item in empty group", key: "empty", id: 999, wantErr: ErrNotFoundEntity},
		{name: "missing item in populated group", key: "populated", id: 999, wantErr: ErrNotFoundEntity},
	}

	r := NewRegistry[string, uint64, string]()

	err := r.Add("populated", 1, "exists")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := r.Get(tc.key, tc.id)
			testutils.Equal(t, err, tc.wantErr)
		})
	}
}

func TestConstructError(t *testing.T) {
	r := NewRegistry[string, uint64, any]()

	err := r.Add(KeyTemporary, 1, &MockConstructError{})
	if err == nil {
		t.Fatal("expected construct error, got nil")
	}

	if err.Error() != "construct failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIteratorSync(t *testing.T) {
	r := NewRegistry[string, uint64, string]()

	if err := r.Add("g1", 1, "a"); err != nil {
		t.Fatal(err)
	}

	if err := r.Add("g1", 2, "b"); err != nil {
		t.Fatal(err)
	}

	if err := r.Add("g2", 3, "c"); err != nil {
		t.Fatal(err)
	}

	collected := make(map[string]bool)
	for v := range r.Iterator("g1", "g2") {
		collected[v] = true
	}

	if len(collected) != 3 {
		t.Fatalf("Iterator() yielded %d items, want 3", len(collected))
	}

	for _, expected := range []string{"a", "b", "c"} {
		if !collected[expected] {
			t.Fatalf("Iterator() missing value %q", expected)
		}
	}
}

func TestAsyncIteratorMultiGroup(t *testing.T) {
	r := NewRegistry[string, uint64, string]()

	if err := r.Add("g1", 1, "a"); err != nil {
		t.Fatal(err)
	}

	if err := r.Add("g1", 2, "b"); err != nil {
		t.Fatal(err)
	}

	if err := r.Add("g2", 3, "c"); err != nil {
		t.Fatal(err)
	}

	if err := r.Add("g3", 4, "d"); err != nil {
		t.Fatal(err)
	}

	collected := make(map[string]bool)
	for v := range r.AsyncIterator("g1", "g2", "g3") {
		collected[v] = true
	}

	if len(collected) != 4 {
		t.Fatalf("AsyncIterator() yielded %d items, want 4", len(collected))
	}

	for _, expected := range []string{"a", "b", "c", "d"} {
		if !collected[expected] {
			t.Fatalf("AsyncIterator() missing value %q", expected)
		}
	}
}

func TestTickGroup(t *testing.T) {
	r := NewRegistry[string, uint64, any]()

	entities := []*MockEntity{{id: 1}, {id: 2}, {id: 3}}
	for _, e := range entities {
		if err := r.Add(KeyTemporary, e.id, e); err != nil {
			t.Fatal(err)
		}
	}

	r.TickGroup(KeyTemporary)

	for _, e := range entities {
		if e.counter != 1 {
			t.Fatalf("entity %d: counter = %d, want 1", e.id, e.counter)
		}
	}
}

func TestTickGroups(t *testing.T) {
	r := NewRegistry[string, uint64, any]()

	e1 := &MockEntity{id: 1}
	e2 := &MockEntity{id: 2}

	if err := r.Add("g1", e1.id, e1); err != nil {
		t.Fatal(err)
	}

	if err := r.Add("g2", e2.id, e2); err != nil {
		t.Fatal(err)
	}

	r.TickGroups("g1", "g2")

	if e1.counter != 1 {
		t.Fatalf("g1 entity: counter = %d, want 1", e1.counter)
	}

	if e2.counter != 1 {
		t.Fatalf("g2 entity: counter = %d, want 1", e2.counter)
	}
}

func TestAsyncTickGroups(t *testing.T) {
	r := NewRegistry[string, uint64, any]()

	e1 := &MockEntity{id: 1}
	e2 := &MockEntity{id: 2}

	if err := r.Add("g1", e1.id, e1); err != nil {
		t.Fatal(err)
	}

	if err := r.Add("g2", e2.id, e2); err != nil {
		t.Fatal(err)
	}

	r.AsyncTick("g1", "g2")

	if e1.counter != 1 {
		t.Fatalf("g1 entity: counter = %d, want 1", e1.counter)
	}

	if e2.counter != 1 {
		t.Fatalf("g2 entity: counter = %d, want 1", e2.counter)
	}
}

func TestTickNonTickerEntity(t *testing.T) {
	r := NewRegistry[string, uint64, any]()

	if err := r.Add(KeyTemporary, 1, "not a ticker"); err != nil {
		t.Fatal(err)
	}

	r.TickGroup(KeyTemporary)
}

/*
BenchmarkGetString-4            20000000                63.1 ns/op             0 B/op          0 allocs/op
BenchmarkGetEntity-4            20000000                64.5 ns/op             0 B/op          0 allocs/op
BenchmarkSetString-4             5000000               387 ns/op              83 B/op          0 allocs/op
BenchmarkSetEntity-4             1000000              1455 ns/op             184 B/op          2 allocs/op
BenchmarkSetDeleteString-4       5000000               283 ns/op               0 B/op          0 allocs/op
BenchmarkSetDeleteEntity-4       2000000               653 ns/op             112 B/op          3 allocs/op
*/

// Prevent optimization
var data interface{}

func BenchmarkGetString(b *testing.B) {
	b.StopTimer()

	r := NewRegistry[string, uint64, any]()

	err := r.Add(KeyTemporary, 1, "simple text insert benchmark")
	if err != nil {
		b.Fatal(err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		var err error

		data, err = r.Get(KeyTemporary, 1)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetEntity(b *testing.B) {
	b.StopTimer()

	r := NewRegistry[string, uint64, any]()

	err := r.Add(KeyTemporary, 1, &MockEntity{})
	if err != nil {
		b.Fatal(err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		var err error

		data, err = r.Get(KeyTemporary, 1)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSetString(b *testing.B) {
	r := NewRegistry[string, uint64, any]()
	for i := 0; i < b.N; i++ {
		id := r.NextID()

		err := r.Add(KeyTemporary, id, "simple text insert benchmark")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSetEntity(b *testing.B) {
	r := NewRegistry[string, uint64, any]()
	for i := 0; i < b.N; i++ {
		id := r.NextID()

		err := r.Add(KeyTemporary, id, &MockEntity{})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSetDeleteString(b *testing.B) {
	r := NewRegistry[string, uint64, any]()
	for i := 0; i < b.N; i++ {
		id := r.NextID()

		err := r.Add(KeyTemporary, id, "simple text insert benchmark")
		if err != nil {
			b.Fatal(err)
		}

		err = r.Remove(KeyTemporary, id)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIterator(b *testing.B) {
	b.StopTimer()

	r := NewRegistry[string, uint64, any]()
	for i := range 1000 {
		err := r.Add(KeyTemporary, uint64(i+1), &MockEntity{id: uint64(i + 1)})
		if err != nil {
			b.Fatal(err)
		}
	}

	b.StartTimer()

	for range b.N {
		for v := range r.Iterator(KeyTemporary) {
			data = v
		}
	}
}

func BenchmarkSearchOne(b *testing.B) {
	b.StopTimer()

	r := NewRegistry[string, uint64, any]()
	for i := range 1000 {
		err := r.Add(KeyTemporary, uint64(i+1), &MockEntity{id: uint64(i + 1)})
		if err != nil {
			b.Fatal(err)
		}
	}

	b.StartTimer()

	for i := range b.N {
		target := uint64(i%1000 + 1)
		data = r.SearchOne(KeyTemporary, func(_ interface{}, id interface{}, _ interface{}) bool {
			return id.(uint64) == target
		})
	}
}

func BenchmarkSetDeleteEntity(b *testing.B) {
	r := NewRegistry[string, uint64, any]()
	for i := 0; i < b.N; i++ {
		id := r.NextID()

		err := r.Add(KeyTemporary, id, &MockEntity{})
		if err != nil {
			b.Fatal(err)
		}

		err = r.Remove(KeyTemporary, id)
		if err != nil {
			b.Fatal(err)
		}
	}
}
