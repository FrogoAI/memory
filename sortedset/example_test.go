package sortedset_test

import (
	"fmt"

	"github.com/FrogoAI/memory/comparator"
	"github.com/FrogoAI/memory/sortedset"
)

func ExampleNewSortedSet() {
	ss := sortedset.NewSortedSet[int, string](comparator.IntComparator)

	ss.Upsert(10, "alice")
	ss.Upsert(20, "bob")
	ss.Upsert(30, "charlie")

	fmt.Println(ss.Len())
	// Output: 3
}

func ExampleSortedSet_Upsert() {
	ss := sortedset.NewSortedSet[int, string](comparator.IntComparator)

	added := ss.Upsert(10, "alice")
	fmt.Println("added:", added)

	updated := ss.Upsert(20, "alice") // same value, new key — updates
	fmt.Println("added:", updated)
	// Output:
	// added: true
	// added: false
}

func ExampleSortedSet_GetByRank() {
	ss := sortedset.NewSortedSet[int, string](comparator.IntComparator)

	ss.Upsert(30, "charlie")
	ss.Upsert(10, "alice")
	ss.Upsert(20, "bob")

	node := ss.GetByRank(1, false) // lowest key
	fmt.Println(node.Key(), node.Value())

	node = ss.GetByRank(-1, false) // highest key
	fmt.Println(node.Key(), node.Value())
	// Output:
	// 10 alice
	// 30 charlie
}

func ExampleSortedSet_Contains() {
	ss := sortedset.NewSortedSet[int, string](comparator.IntComparator)

	ss.Upsert(10, "alice")

	fmt.Println(ss.Contains("alice"))
	fmt.Println(ss.Contains("bob"))
	// Output:
	// true
	// false
}

func ExampleSortedSet_Remove() {
	ss := sortedset.NewSortedSet[int, string](comparator.IntComparator)

	ss.Upsert(10, "alice")
	ss.Upsert(20, "bob")

	ss.Remove("alice")

	fmt.Println(ss.Len())
	fmt.Println(ss.Contains("alice"))
	// Output:
	// 1
	// false
}
