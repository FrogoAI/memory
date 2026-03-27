package lru_test

import (
	"fmt"

	"github.com/FrogoAI/memory/lru"
)

func ExampleNewLRUCache() {
	cache := lru.NewLRUCache[string](3)

	cache.Put("a", "alpha")
	cache.Put("b", "beta")

	fmt.Println(cache.Len())
	// Output: 2
}

func ExampleCache_Put() {
	cache := lru.NewLRUCache[string](2)

	cache.Put("a", "alpha")
	cache.Put("b", "beta")
	cache.Put("c", "gamma") // evicts "a"

	_, found := cache.Get("a")
	fmt.Println("a found:", found)

	val, found := cache.Get("c")
	fmt.Println("c found:", found, "value:", val)
	// Output:
	// a found: false
	// c found: true value: gamma
}

func ExampleCache_Get() {
	cache := lru.NewLRUCache[int](3)

	cache.Put("x", 42)

	val, ok := cache.Get("x")
	fmt.Println(val, ok)

	val, ok = cache.Get("missing")
	fmt.Println(val, ok)
	// Output:
	// 42 true
	// 0 false
}

func ExampleCache_Iterator() {
	cache := lru.NewLRUCache[string](3)

	cache.Put("a", "alpha")
	cache.Put("b", "beta")
	cache.Put("c", "gamma")

	for key, val := range cache.Iterator() {
		fmt.Printf("%s=%s ", key, val)
	}

	fmt.Println()
	// Output: c=gamma b=beta a=alpha
}

func ExampleCache_Clear() {
	cache := lru.NewLRUCache[string](3)

	cache.Put("a", "alpha")
	cache.Put("b", "beta")
	fmt.Println(cache.Len())

	cache.Clear()
	fmt.Println(cache.Len())
	// Output:
	// 2
	// 0
}
