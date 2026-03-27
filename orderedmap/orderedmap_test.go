package orderedmap

import (
	"strconv"
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestOrderedmap(t *testing.T) {
	o := &OrderedMap[string, string]{}
	o.Add("A", "1")
	o.Add("B", "2")
	o.Add("C", "3")
	o.Add("E", "4")
	o.Remove("C")
	testutils.Equal(t, o.Get("A"), "1")
	testutils.Equal(t, o.Get("C"), "")
	keys, values := o.GetAll()
	testutils.Equal(t, len(keys), 3)
	testutils.Equal(t, len(values), 3)
	testutils.Equal(t, keys, []string{"A", "B", "E"})
	testutils.Equal(t, values, []string{"1", "2", "4"})
}

func BenchmarkAdd(b *testing.B) {
	b.Run("new_key", func(b *testing.B) {
		o := &OrderedMap[string, int]{}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			o.Add(strconv.Itoa(i), i)
		}
	})

	b.Run("update_existing", func(b *testing.B) {
		o := &OrderedMap[string, int]{}
		for i := 0; i < 1000; i++ {
			o.Add(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			o.Add(strconv.Itoa(i%1000), i)
		}
	})
}

func BenchmarkGet(b *testing.B) {
	b.Run("hit", func(b *testing.B) {
		o := &OrderedMap[string, int]{}
		for i := 0; i < 1000; i++ {
			o.Add(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			o.Get(strconv.Itoa(i % 1000))
		}
	})

	b.Run("miss", func(b *testing.B) {
		o := &OrderedMap[string, int]{}
		for i := 0; i < 1000; i++ {
			o.Add(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			o.Get("miss" + strconv.Itoa(i))
		}
	})
}

func BenchmarkIterate(b *testing.B) {
	b.Run("get_all", func(b *testing.B) {
		o := &OrderedMap[string, int]{}
		for i := 0; i < 1000; i++ {
			o.Add(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			o.GetAll()
		}
	})

	b.Run("iterator", func(b *testing.B) {
		o := &OrderedMap[string, int]{}
		for i := 0; i < 1000; i++ {
			o.Add(strconv.Itoa(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			ch := o.Iterator(1000)
			for range ch {
			}
		}
	})
}
