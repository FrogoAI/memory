# memory

> A high-performance collection of generic data structures and algorithms for Go.

`memory` provides specialized, type-safe in-memory data structures optimized for efficiency. Unlike the standard library's container packages, this library offers advanced algorithms like **Counting Bloom Filters**, **HyperLogLog**, **Fuzzy Search**, **B-trees**, and **Sorted Sets**, all built using Go generics.

Data structures are **not thread-safe by default** (matching Go stdlib convention). Thread-safe wrappers (`Safe*`) are provided where needed.

## Packages

### Probabilistic Data Structures

| Package | Description |
|---------|-------------|
| **bloom** | Counting Bloom Filter with Add/Remove/Test and binary serialization |
| **hll** | HyperLogLog cardinality estimator with Union and Intersection |

### Search & String Algorithms

| Package | Description |
|---------|-------------|
| **fuzzysearch** | Fuzzy string matching with Levenshtein distance ranking |

### Ordered Collections

| Package | Description |
|---------|-------------|
| **btree** | B-tree for ordered key-value storage with configurable order |
| **sortedset** | Redis-like ZSET (skip list) with scoring, ranking, and range queries |
| **linkedlist** | Doubly-linked list with ID-based O(1) lookup |
| **orderedmap** | Map that maintains insertion order |
| **stack** | Generic LIFO stack |

### Caches & Registries

| Package | Description |
|---------|-------------|
| **lru** | LRU cache with automatic eviction |
| **registry** | Thread-safe grouped item registry (sessions, workers, connections) |

### Infrastructure

| Package | Description |
|---------|-------------|
| **comparator** | Type-aware comparison functions (int, string, float, time, etc.) |
| **utils** | SafeMap, SafeList (thread-safe), hashing (CRC, Murmur3, SimHash), string helpers |

### Thread-Safe Wrappers

| Wrapper | Wraps | Use when |
|---------|-------|----------|
| `sortedset.SafeSortedSet` | `SortedSet` | Concurrent access to sorted set |
| `orderedmap.SafeOrderedMap` | `OrderedMap` | Concurrent access to ordered map |
| `utils.SafeMap` | `map[K]V` | Concurrent access to map |
| `utils.SafeList` | `[]V` | Concurrent access to slice |

## Installation

```bash
go get github.com/FrogoAI/memory
```

## Usage Examples

### Counting Bloom Filter

```go
filter, err := bloom.NewCounting(1000, 0.01) // 1000 items, 1% false-positive rate
if err != nil {
    panic(err)
}

filter.Add([]byte("user_123"))
filter.Test([]byte("user_123")) // true
filter.Remove([]byte("user_123"))
filter.Test([]byte("user_123")) // false
```

### B-Tree

```go
tree, err := btree.NewWithIntComparator(3) // order 3
if err != nil {
    panic(err)
}

tree.Put(5, "five")
tree.Put(3, "three")
tree.Put(7, "seven")

value, found, _ := tree.Get(5) // "five", true
keys := tree.Keys()            // [3, 5, 7] (sorted)
```

### Fuzzy Search

```go
targets := []string{"cartwheel", "foobar", "wheel", "baz"}

matches := fuzzysearch.Find("whl", targets) // ["cartwheel", "wheel"]

ranks := fuzzysearch.RankFind("wheel", targets)
for _, r := range ranks {
    fmt.Printf("%s (distance: %d)\n", r.Target, r.Distance)
}
```

### LRU Cache

```go
cache := lru.NewLRUCache[string](2) // capacity 2

cache.Put("a", "alpha")
cache.Put("b", "beta")
cache.Get("a")             // "alpha", true — promotes to front
cache.Put("c", "charlie")  // evicts "b"
cache.Get("b")             // "", false
```

### Sorted Set (Redis-like ZSET)

```go
ss := sortedset.NewSortedSet[int, string](comparator.IntComparator)

ss.Upsert(100, "alice")
ss.Upsert(200, "bob")
ss.Upsert(50, "charlie")

top := ss.GetTop(2, false) // bob (200), alice (100)
ss.FindRank("charlie")     // 1 (lowest score)
```

### Thread-Safe Sorted Set

```go
ss := sortedset.NewSafeSortedSet[int, string](comparator.IntComparator)

// Safe for concurrent use from multiple goroutines
ss.Upsert(100, "alice")
ss.GetByValue("alice") // protected by RWMutex
```

### Registry

```go
reg := registry.NewRegistry[string, uint64, string]()

reg.Add("admins", 1, "Alice")
reg.Add("admins", 2, "Bob")

for user := range reg.Iterator("admins") {
    println(user)
}
```

### HyperLogLog

```go
h, _ := hll.New()

h.Add([]byte("user1"))
h.Add([]byte("user2"))
h.Add([]byte("user1")) // duplicate

h.Count() // ~2
```

## Build & Test

```bash
go test ./...                            # run tests
go test -race ./...                      # with race detector
go test -bench=. -benchmem -run="^$" ./... # benchmarks
golangci-lint run ./...                  # lint
make ci                                  # lint + coverage
```

## License

[MIT](LICENSE)
