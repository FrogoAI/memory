package fuzzysearch

import (
	"testing"

	"github.com/FrogoAI/set"
	"github.com/FrogoAI/testutils"
)

func TestSearch(t *testing.T) {
	s := NewIndex()
	d := NewDocument(1, "test string")
	d2 := NewDocument(2, "some name")
	d3 := NewDocument(3, "awesome name")
	s.Add(d, d2, d3)
	testutils.Equal(t, s.Search("test"), []int{1})
}

func TestSearchMultiToken(t *testing.T) {
	idx := NewIndex()
	idx.Add(
		NewDocument(1, "quick brown fox"),
		NewDocument(2, "quick red car"),
		NewDocument(3, "slow brown dog"),
	)

	t.Run("two token intersection", func(t *testing.T) {
		// "quick brown" should only match doc 1
		results := idx.Search("quick brown")
		testutils.Equal(t, results, []int{1})
	})

	t.Run("single token multiple results", func(t *testing.T) {
		results := idx.Search("quick")
		testutils.Equal(t, results, []int{1, 2})
	})

	t.Run("no match token", func(t *testing.T) {
		results := idx.Search("zebra")
		if len(results) != 0 {
			t.Errorf("expected no results, got %v", results)
		}
	})
}

func TestIndexRemove(t *testing.T) {
	idx := NewIndex()
	d1 := NewDocument(1, "hello world")
	d2 := NewDocument(2, "hello there")
	idx.Add(d1, d2)

	t.Run("remove existing document", func(t *testing.T) {
		idx.Remove(d1)
		results := idx.Search("hello")
		testutils.Equal(t, results, []int{2})
	})

	t.Run("remove last document for token", func(t *testing.T) {
		idx.Remove(d2)

		results := idx.Search("hello")
		if len(results) != 0 {
			t.Errorf("expected no results, got %v", results)
		}
	})

	t.Run("remove non-existent document", func(t *testing.T) {
		// Should not panic
		idx.Remove(NewDocument(99, "nonexistent"))
	})
}

func TestIntersection(t *testing.T) {
	cases := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{"both non-empty with overlap", []int{1, 2, 3, 5}, []int{2, 3, 4, 5}, []int{2, 3, 5}},
		{"no overlap", []int{1, 3, 5}, []int{2, 4, 6}, nil},
		{"identical", []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}},
		{"a empty", []int{}, []int{1, 2}, nil},
		{"b empty", []int{1, 2}, []int{}, nil},
		{"both empty", []int{}, []int{}, nil},
		{"single element match", []int{5}, []int{5}, []int{5}},
		{"b longer than a", []int{2}, []int{1, 2, 3, 4, 5}, []int{2}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Intersection(tc.a, tc.b)
			if len(tc.want) == 0 {
				if len(got) != 0 {
					t.Errorf("Intersection(%v, %v) = %v, want empty", tc.a, tc.b, got)
				}

				return
			}

			testutils.Equal(t, got, tc.want)
		})
	}
}

func TestDocumentTerms(t *testing.T) {
	d := NewDocument(1, "Hello World")
	terms := d.Terms()
	testutils.Equal(t, terms, []string{"hello", "world"})
}

func TestTokenize(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{"simple words", "hello world", []string{"hello", "world"}},
		{"with punctuation", "hello, world!", []string{"hello", "world"}},
		{"numbers", "abc123 def456", []string{"abc123", "def456"}},
		{"empty string", "", nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Tokenize(tc.input)
			if len(tc.want) == 0 {
				if len(got) != 0 {
					t.Errorf("Tokenize(%q) = %v, want empty", tc.input, got)
				}

				return
			}

			testutils.Equal(t, got, tc.want)
		})
	}
}

func TestStopwordFilter(t *testing.T) {
	stopwords := set.NewGenericDataSet("the", "a", "is")

	cases := []struct {
		name   string
		tokens []string
		want   []string
	}{
		{"filters stopwords", []string{"the", "cat", "is", "fast"}, []string{"cat", "fast"}},
		{"no stopwords", []string{"hello", "world"}, []string{"hello", "world"}},
		{"all stopwords", []string{"the", "a", "is"}, nil},
		{"empty input", []string{}, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := StopwordFilter(tc.tokens, stopwords)
			if len(tc.want) == 0 {
				if len(got) != 0 {
					t.Errorf("StopwordFilter(%v) = %v, want empty", tc.tokens, got)
				}

				return
			}

			testutils.Equal(t, got, tc.want)
		})
	}
}
