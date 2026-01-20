# memory

> A high-performance collection of generic data structures and algorithms for Go.

`memory` provides a suite of specialized, type-safe in-memory data structures optimized for efficiency and specific use cases. Unlike the standard library's container packages, this repository offers advanced algorithms like **Locality Sensitive Hashing (SimDict)**, **Counting Bloom Filters**, **Fuzzy Search**, and **Ordered Maps**, all built using Go 1.18+ generics.

## Features

### 🧠 Probabilistic Data Structures

- **Bloom Filter (`bloom`):** Space-efficient probabilistic data structure for checking if an element is in a set. Includes a **Counting Bloom Filter** that supports removals.
- **HyperLogLog (`hll`):** Efficient cardinality estimation for large datasets.
- **SimDict (`simdict`):** A similarity dictionary using **Locality Sensitive Hashing (LSH)** to cluster similar documents or strings into buckets.

### 🔍 Search & String Algo

- **Fuzzy Search (`fuzzysearch`):** Fast, lightweight fuzzy matching. Supports simple subsequence matching and **Levenshtein distance** ranking.

### 📦 Containers & Caches

- **LRU Cache (`lru`):** A generic, thread-safe Least Recently Used cache implementation.
- **Ordered Map (`orderedmap`):** A map that maintains insertion order, supporting iteration and index-based access.
- **Registry (`registry`):** A thread-safe structure for grouping items by ID and Category, useful for managing active sessions or grouped workers.
- **Set (`set`):** Generic `DataSet` and `OrderedDataSet` implementations for O(1) lookups.
- **Stack (`stack`):** A classic LIFO stack implementation with slice-based backing.
- **Sorted Set (`sortedset`):** (Redis-like) ZSET implementation for storing unique elements with scores.

## Installation

Bash

```
go get github.com/FrogoAI/memory
```

## Usage Examples

### 1. Counting Bloom Filter

Check for existence with a defined false-positive rate, with support for removing items.

Go

```
package main

import (
	"fmt"
	"github.com/FrogoAI/memory/bloom"
)

func main() {
	// Initialize for 1000 items with 0.01 (1%) false positive rate
	filter := bloom.NewCounting(1000, 0.01)

	data := []byte("user_123")

	filter.Add(data)

	if filter.Test(data) {
		fmt.Println("Item exists!")
	}

	filter.Remove(data)
}
```

### 2. Fuzzy Search & Ranking

Find strings that approximately match a target, ranked by Levenshtein distance.

Go

```
package main

import (
	"fmt"
	"github.com/FrogoAI/memory/fuzzysearch"
)

func main() {
	targets := []string{"cartwheel", "foobar", "wheel", "baz"}
	
	// Find simple matches (subsequence)
	matches := fuzzysearch.Find("whl", targets)
	fmt.Println(matches) // ["cartwheel", "wheel"]

	// Rank by Levenshtein distance
	ranks := fuzzysearch.RankFind("wheel", targets)
	for _, r := range ranks {
		fmt.Printf("Target: %s, Distance: %d\n", r.Target, r.Distance)
	}
}
```

### 3. LRU Cache

A type-safe cache that automatically evicts the least recently used items.

Go

```
package main

import (
	"fmt"
	"github.com/FrogoAI/memory/lru"
)

func main() {
	// Create a cache with capacity 2
	cache := lru.NewLRUCache[string](2)

	cache.Put("a", "alpha")
	cache.Put("b", "beta")
	
	val, found := cache.Get("a") // "a" is now most recently used
	
	cache.Put("c", "charlie") // Evicts "b" (least recently used)

	_, foundB := cache.Get("b") // false
	fmt.Println(foundB)
}
```

### 4. SimDict (Similarity Clustering)

Cluster documents into buckets based on similarity using LSH.

Go

```
package main

import (
	"fmt"
	"github.com/FrogoAI/memory/simdict"
)

func main() {
	manager := simdict.NewLSHManager()

	// Assign documents to buckets based on content similarity
	bucket1 := manager.ProcessAndAssign("The quick brown fox")
	bucket2 := manager.ProcessAndAssign("The quick brown fox jumps")
	bucket3 := manager.ProcessAndAssign("Completely different text")

	fmt.Println(bucket1 == bucket2) // true (Similar enough to share a bucket)
	fmt.Println(bucket1 == bucket3) // false
}
```

### 5. Registry

Manage grouped resources (e.g., connections per user) safely.

Go

```
package main

import (
	"github.com/FrogoAI/memory/registry"
)

func main() {
	// Key=String (GroupID), Index=Int (ItemID), Value=String (Data)
	reg := registry.NewRegistry[string, int, string]()

	// Add item 101 to group "admins"
	reg.Add("admins", 101, "User Alice")
	
	// Iterate over all admins
	for user := range reg.Iterator("admins") {
		println(user)
	}
}
```

## License

[MIT](https://www.google.com/search?q=LICENSE)