package fuzzysearch_test

import (
	"fmt"

	"github.com/FrogoAI/memory/fuzzysearch"
)

func ExampleMatchFold() {
	fmt.Println(fuzzysearch.MatchFold("WHL", "cartwheel"))
	// Output: true
}

func ExampleMatchNormalizedFold() {
	fmt.Println(fuzzysearch.MatchNormalizedFold("cafe", "Café"))
	// Output: true
}

func ExampleFindFold() {
	results := fuzzysearch.FindFold("WHL", []string{"cartwheel", "foobar", "Wheel"})
	fmt.Println(results)
	// Output: [cartwheel Wheel]
}

func ExampleRankMatchNormalized() {
	distance := fuzzysearch.RankMatchNormalized("cafe", "café")
	fmt.Println(distance)
	// Output: 0
}

func ExampleRankFindFold() {
	ranks := fuzzysearch.RankFindFold("WHL", []string{"cartwheel", "foobar", "Wheel"})
	for _, r := range ranks {
		fmt.Printf("%s: %d\n", r.Target, r.Distance)
	}
	// Output:
	// cartwheel: 9
	// Wheel: 4
}

func ExampleLevenshteinDistance() {
	fmt.Println(fuzzysearch.LevenshteinDistance("kitten", "sitting"))
	// Output: 3
}
