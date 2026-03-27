package bloom_test

import (
	"fmt"

	"github.com/FrogoAI/memory/bloom"
)

func ExampleNewCounting() {
	bf, err := bloom.NewCounting(1000, 0.01)
	if err != nil {
		panic(err)
	}

	bf.Add([]byte("hello"))
	fmt.Println(bf.Test([]byte("hello")))
	// Output: true
}

func ExampleCountingFilter_Add() {
	bf, _ := bloom.NewCounting(1000, 0.01)

	bf.Add([]byte("apple"))
	bf.Add([]byte("banana"))
	bf.Add([]byte("cherry"))

	fmt.Println(bf.Test([]byte("apple")))
	fmt.Println(bf.Test([]byte("banana")))
	fmt.Println(bf.Test([]byte("cherry")))
	// Output:
	// true
	// true
	// true
}

func ExampleCountingFilter_Test() {
	bf, _ := bloom.NewCounting(1000, 0.01)

	bf.Add([]byte("apple"))

	fmt.Println(bf.Test([]byte("apple")))
	fmt.Println(bf.Test([]byte("banana")))
	// Output:
	// true
	// false
}

func ExampleCountingFilter_Remove() {
	bf, _ := bloom.NewCounting(1000, 0.01)

	bf.Add([]byte("apple"))
	fmt.Println(bf.Test([]byte("apple")))

	bf.Remove([]byte("apple"))
	fmt.Println(bf.Test([]byte("apple")))
	// Output:
	// true
	// false
}

func ExampleCountingFilter_Clear() {
	bf, _ := bloom.NewCounting(1000, 0.01)

	bf.Add([]byte("apple"))
	bf.Add([]byte("banana"))

	fmt.Println(bf.Test([]byte("apple")))

	bf.Clear()
	fmt.Println(bf.Test([]byte("apple")))
	// Output:
	// true
	// false
}
