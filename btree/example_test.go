package btree_test

import (
	"fmt"

	"github.com/FrogoAI/memory/btree"
)

func ExampleNewWithIntComparator() {
	tree, err := btree.NewWithIntComparator(3)
	if err != nil {
		panic(err)
	}

	_ = tree.Put(1, "one")
	_ = tree.Put(2, "two")
	_ = tree.Put(3, "three")

	fmt.Println(tree.Len())
	// Output: 3
}

func ExampleNewWithStringComparator() {
	tree, err := btree.NewWithStringComparator(3)
	if err != nil {
		panic(err)
	}

	_ = tree.Put("b", "banana")
	_ = tree.Put("a", "apple")
	_ = tree.Put("c", "cherry")

	fmt.Println(tree.Keys())
	// Output: [a b c]
}

func ExampleTree_Put() {
	tree, _ := btree.NewWithIntComparator(3)

	_ = tree.Put(3, "three")
	_ = tree.Put(1, "one")
	_ = tree.Put(2, "two")

	fmt.Println(tree.Keys())
	// Output: [1 2 3]
}

func ExampleTree_Get() {
	tree, _ := btree.NewWithIntComparator(3)

	_ = tree.Put(1, "one")

	value, found, _ := tree.Get(1)
	fmt.Println(value, found)

	value, found, _ = tree.Get(99)
	fmt.Println(value, found)
	// Output:
	// one true
	// <nil> false
}

func ExampleTree_Remove() {
	tree, _ := btree.NewWithIntComparator(3)

	_ = tree.Put(1, "one")
	_ = tree.Put(2, "two")
	_ = tree.Put(3, "three")

	_ = tree.Remove(2)

	fmt.Println(tree.Keys())
	// Output: [1 3]
}

func ExampleTree_Iterator() {
	tree, _ := btree.NewWithIntComparator(3)

	_ = tree.Put(3, "three")
	_ = tree.Put(1, "one")
	_ = tree.Put(2, "two")

	it := tree.Iterator()
	for it.Next() {
		fmt.Printf("%d=%v ", it.Key(), it.Value())
	}

	fmt.Println()
	// Output: 1=one 2=two 3=three
}
