package btree

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/FrogoAI/memory/comparator"
)

func mustPut(t testing.TB, tree *Tree[int, any], key int, value any) {
	t.Helper()

	if err := tree.Put(key, value); err != nil {
		t.Fatalf("Put(%v, %v): %v", key, value, err)
	}
}

func mustRemove(t testing.TB, tree *Tree[int, any], key int) {
	t.Helper()

	if err := tree.Remove(key); err != nil {
		t.Fatalf("Remove(%v): %v", key, err)
	}
}

func mustGet(t testing.TB, tree *Tree[int, any], key int) (any, bool) {
	t.Helper()

	value, found, err := tree.Get(key)
	if err != nil {
		t.Fatalf("Get(%v): %v", key, err)
	}

	return value, found
}

func TestBTreeGet1(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 1, "a")
	mustPut(t, tree, 2, "b")
	mustPut(t, tree, 3, "c")
	mustPut(t, tree, 4, "d")
	mustPut(t, tree, 5, "e")
	mustPut(t, tree, 6, "f")
	mustPut(t, tree, 7, "g")

	v, e := mustGet(t, tree, 4)
	slog.Info("Test tree get", "v", v, "e", e)

	tests := [][]interface{}{
		{0, nil, false},
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, nil, false},
	}

	for _, test := range tests {
		if value, found := mustGet(t, tree, test[0].(int)); value != test[1] || found != test[2] {
			t.Errorf("Got %v,%v expected %v,%v", value, found, test[1], test[2])
		}
	}
}

func TestBTreeGet2(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 7, "g")
	mustPut(t, tree, 9, "i")
	mustPut(t, tree, 10, "j")
	mustPut(t, tree, 6, "f")
	mustPut(t, tree, 3, "c")
	mustPut(t, tree, 4, "d")
	mustPut(t, tree, 5, "e")
	mustPut(t, tree, 8, "h")
	mustPut(t, tree, 2, "b")
	mustPut(t, tree, 1, "a")

	tests := [][]interface{}{
		{0, nil, false},
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "h", true},
		{9, "i", true},
		{10, "j", true},
		{11, nil, false},
	}

	for _, test := range tests {
		if value, found := mustGet(t, tree, test[0].(int)); value != test[1] || found != test[2] {
			t.Errorf("Got %v,%v expected %v,%v", value, found, test[1], test[2])
		}
	}
}

func TestBTreePut1(t *testing.T) {
	// https://upload.wikimedia.org/wikipedia/commons/3/33/B_tree_insertion_example.png
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	assertValidTree(t, tree, 0)

	mustPut(t, tree, 1, 0)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root, 1, 0, []int{1}, false)

	mustPut(t, tree, 2, 1)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root, 2, 0, []int{1, 2}, false)

	mustPut(t, tree, 3, 2)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{3}, true)

	mustPut(t, tree, 4, 2)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 2, 0, []int{3, 4}, true)

	mustPut(t, tree, 5, 2)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{2, 4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 1, 0, []int{5}, true)

	mustPut(t, tree, 6, 2)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{2, 4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 2, 0, []int{5, 6}, true)

	mustPut(t, tree, 7, 2)
	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{6}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{7}, true)
}

func TestBTreePut2(t *testing.T) {
	tree, err := NewWithIntComparator(4)
	if err != nil {
		t.Fatal(err)
	}

	assertValidTree(t, tree, 0)

	mustPut(t, tree, 0, 0)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root, 1, 0, []int{0}, false)

	mustPut(t, tree, 2, 2)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root, 2, 0, []int{0, 2}, false)

	mustPut(t, tree, 1, 1)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 3, 0, []int{0, 1, 2}, false)

	mustPut(t, tree, 1, 1)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 3, 0, []int{0, 1, 2}, false)

	mustPut(t, tree, 3, 3)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{1}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 2, 0, []int{2, 3}, true)

	mustPut(t, tree, 4, 4)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{1}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 3, 0, []int{2, 3, 4}, true)

	mustPut(t, tree, 5, 5)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{1, 3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 2, 0, []int{4, 5}, true)
}

func TestBTreePut3(t *testing.T) {
	// http://www.geeksforgeeks.org/b-tree-set-1-insert-2/
	tree, err := NewWithIntComparator(6)
	if err != nil {
		t.Fatal(err)
	}

	assertValidTree(t, tree, 0)

	mustPut(t, tree, 10, 0)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root, 1, 0, []int{10}, false)

	mustPut(t, tree, 20, 1)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root, 2, 0, []int{10, 20}, false)

	mustPut(t, tree, 30, 2)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 3, 0, []int{10, 20, 30}, false)

	mustPut(t, tree, 40, 3)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 4, 0, []int{10, 20, 30, 40}, false)

	mustPut(t, tree, 50, 4)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root, 5, 0, []int{10, 20, 30, 40, 50}, false)

	mustPut(t, tree, 60, 5)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{30}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 3, 0, []int{40, 50, 60}, true)

	mustPut(t, tree, 70, 6)
	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{30}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 4, 0, []int{40, 50, 60, 70}, true)

	mustPut(t, tree, 80, 7)
	assertValidTree(t, tree, 8)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{30}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 5, 0, []int{40, 50, 60, 70, 80}, true)

	mustPut(t, tree, 90, 8)
	assertValidTree(t, tree, 9)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{30, 60}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{10, 20}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 2, 0, []int{40, 50}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 3, 0, []int{70, 80, 90}, true)
}

func TestBTreePut4(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	assertValidTree(t, tree, 0)

	mustPut(t, tree, 6, nil)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root, 1, 0, []int{6}, false)

	mustPut(t, tree, 5, nil)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root, 2, 0, []int{5, 6}, false)

	mustPut(t, tree, 4, nil)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{5}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{6}, true)

	mustPut(t, tree, 3, nil)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{5}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{3, 4}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{6}, true)

	mustPut(t, tree, 2, nil)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{3, 5}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 1, 0, []int{6}, true)

	mustPut(t, tree, 1, nil)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{3, 5}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{1, 2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 1, 0, []int{6}, true)

	mustPut(t, tree, 0, nil)
	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{6}, true)

	mustPut(t, tree, -1, nil)
	assertValidTree(t, tree, 8)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 2, 0, []int{-1, 0}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{6}, true)

	mustPut(t, tree, -2, nil)
	assertValidTree(t, tree, 9)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 3, []int{-1, 1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{-2}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[2], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{6}, true)

	mustPut(t, tree, -3, nil)
	assertValidTree(t, tree, 10)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 3, []int{-1, 1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 2, 0, []int{-3, -2}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[2], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{6}, true)

	mustPut(t, tree, -4, nil)
	assertValidTree(t, tree, 11)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{-1, 3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{-3}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 1, 2, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{-4}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{-2}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[2].Children[0], 1, 0, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[2].Children[1], 1, 0, []int{6}, true)
}

func TestBTreeRemove1(t *testing.T) {
	// empty
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustRemove(t, tree, 1)
	assertValidTree(t, tree, 0)
}

func TestBTreeRemove2(t *testing.T) {
	// leaf node (no underflow)
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 1, nil)
	mustPut(t, tree, 2, nil)

	mustRemove(t, tree, 1)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root, 1, 0, []int{2}, false)

	mustRemove(t, tree, 2)
	assertValidTree(t, tree, 0)
}

func TestBTreeRemove3(t *testing.T) {
	// merge with right (underflow)
	{
		tree, err := NewWithIntComparator(3)
		if err != nil {
			t.Fatal(err)
		}

		mustPut(t, tree, 1, nil)
		mustPut(t, tree, 2, nil)
		mustPut(t, tree, 3, nil)

		mustRemove(t, tree, 1)
		assertValidTree(t, tree, 2)
		assertValidTreeNode(t, tree.Root, 2, 0, []int{2, 3}, false)
	}
	// merge with left (underflow)
	{
		tree, err := NewWithIntComparator(3)
		if err != nil {
			t.Fatal(err)
		}

		mustPut(t, tree, 1, nil)
		mustPut(t, tree, 2, nil)
		mustPut(t, tree, 3, nil)

		mustRemove(t, tree, 3)
		assertValidTree(t, tree, 2)
		assertValidTreeNode(t, tree.Root, 2, 0, []int{1, 2}, false)
	}
}

func TestBTreeRemove4(t *testing.T) {
	// rotate left (underflow)
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 1, nil)
	mustPut(t, tree, 2, nil)
	mustPut(t, tree, 3, nil)
	mustPut(t, tree, 4, nil)

	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 2, 0, []int{3, 4}, true)

	mustRemove(t, tree, 1)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{3}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{4}, true)
}

func TestBTreeRemove5(t *testing.T) {
	// rotate right (underflow)
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 1, nil)
	mustPut(t, tree, 2, nil)
	mustPut(t, tree, 3, nil)
	mustPut(t, tree, 0, nil)

	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{0, 1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{3}, true)

	mustRemove(t, tree, 3)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{1}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{0}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{2}, true)
}

func TestBTreeRemove6(t *testing.T) {
	// root height reduction after a series of underflows on right side
	// use simulator: https://www.cs.usfca.edu/~galles/visualization/BTree.html
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 1, nil)
	mustPut(t, tree, 2, nil)
	mustPut(t, tree, 3, nil)
	mustPut(t, tree, 4, nil)
	mustPut(t, tree, 5, nil)
	mustPut(t, tree, 6, nil)
	mustPut(t, tree, 7, nil)

	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{6}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{7}, true)

	mustRemove(t, tree, 7)
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{2, 4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 2, 0, []int{5, 6}, true)
}

func TestBTreeRemove7(t *testing.T) {
	// root height reduction after a series of underflows on left side
	// use simulator: https://www.cs.usfca.edu/~galles/visualization/BTree.html
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 1, nil)
	mustPut(t, tree, 2, nil)
	mustPut(t, tree, 3, nil)
	mustPut(t, tree, 4, nil)
	mustPut(t, tree, 5, nil)
	mustPut(t, tree, 6, nil)
	mustPut(t, tree, 7, nil)

	assertValidTree(t, tree, 7)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{6}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{7}, true)

	mustRemove(t, tree, 1) // series of underflows
	assertValidTree(t, tree, 6)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{4, 6}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{2, 3}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 1, 0, []int{7}, true)

	// clear all remaining
	mustRemove(t, tree, 2)
	assertValidTree(t, tree, 5)
	assertValidTreeNode(t, tree.Root, 2, 3, []int{4, 6}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[2], 1, 0, []int{7}, true)

	mustRemove(t, tree, 3)
	assertValidTree(t, tree, 4)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{6}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 2, 0, []int{4, 5}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{7}, true)

	mustRemove(t, tree, 4)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{6}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 0, []int{7}, true)

	mustRemove(t, tree, 5)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.Root, 2, 0, []int{6, 7}, false)

	mustRemove(t, tree, 6)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.Root, 1, 0, []int{7}, false)

	mustRemove(t, tree, 7)
	assertValidTree(t, tree, 0)
}

func TestBTreeRemove8(t *testing.T) {
	// use simulator: https://www.cs.usfca.edu/~galles/visualization/BTree.html
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 1, nil)
	mustPut(t, tree, 2, nil)
	mustPut(t, tree, 3, nil)
	mustPut(t, tree, 4, nil)
	mustPut(t, tree, 5, nil)
	mustPut(t, tree, 6, nil)
	mustPut(t, tree, 7, nil)
	mustPut(t, tree, 8, nil)
	mustPut(t, tree, 9, nil)

	assertValidTree(t, tree, 9)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{4}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{2}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 2, 3, []int{6, 8}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 1, 0, []int{1}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{3}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{7}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[2], 1, 0, []int{9}, true)

	mustRemove(t, tree, 1)
	assertValidTree(t, tree, 8)
	assertValidTreeNode(t, tree.Root, 1, 2, []int{6}, false)
	assertValidTreeNode(t, tree.Root.Children[0], 1, 2, []int{4}, true)
	assertValidTreeNode(t, tree.Root.Children[1], 1, 2, []int{8}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[0], 2, 0, []int{2, 3}, true)
	assertValidTreeNode(t, tree.Root.Children[0].Children[1], 1, 0, []int{5}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[0], 1, 0, []int{7}, true)
	assertValidTreeNode(t, tree.Root.Children[1].Children[1], 1, 0, []int{9}, true)
}

func TestBTreeRemove9(t *testing.T) {
	const maxElm = 1000

	orders := []int{3, 4, 5, 6, 7, 8, 9, 10, 20, 100, 500, 1000, 5000, 10000}
	for _, order := range orders {
		tree, err := NewWithIntComparator(order)
		if err != nil {
			t.Fatal(err)
		}

		{
			for i := 1; i <= maxElm; i++ {
				mustPut(t, tree, i, i)
			}

			assertValidTree(t, tree, maxElm)

			for i := 1; i <= maxElm; i++ {
				if _, found := mustGet(t, tree, i); !found {
					t.Errorf("Not found %v", i)
				}
			}

			for i := 1; i <= maxElm; i++ {
				mustRemove(t, tree, i)
			}

			assertValidTree(t, tree, 0)
		}

		{
			for i := maxElm; i > 0; i-- {
				mustPut(t, tree, i, i)
			}

			assertValidTree(t, tree, maxElm)

			for i := maxElm; i > 0; i-- {
				if _, found := mustGet(t, tree, i); !found {
					t.Errorf("Not found %v", i)
				}
			}

			for i := maxElm; i > 0; i-- {
				mustRemove(t, tree, i)
			}

			assertValidTree(t, tree, 0)
		}
	}
}

func TestBTreeHeight(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	if actualValue, expectedValue := tree.Height(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	mustPut(t, tree, 1, 0)

	if actualValue, expectedValue := tree.Height(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	mustPut(t, tree, 2, 1)

	if actualValue, expectedValue := tree.Height(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	mustPut(t, tree, 3, 2)

	if actualValue, expectedValue := tree.Height(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	mustPut(t, tree, 4, 2)

	if actualValue, expectedValue := tree.Height(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	mustPut(t, tree, 5, 2)

	if actualValue, expectedValue := tree.Height(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	mustPut(t, tree, 6, 2)

	if actualValue, expectedValue := tree.Height(), 2; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	mustPut(t, tree, 7, 2)

	if actualValue, expectedValue := tree.Height(), 3; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	mustRemove(t, tree, 1)
	mustRemove(t, tree, 2)
	mustRemove(t, tree, 3)
	mustRemove(t, tree, 4)
	mustRemove(t, tree, 5)
	mustRemove(t, tree, 6)
	mustRemove(t, tree, 7)

	if actualValue, expectedValue := tree.Height(), 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeLeftAndRight(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	if actualValue := tree.Left(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	if actualValue := tree.Right(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	mustPut(t, tree, 1, "a")
	mustPut(t, tree, 5, "e")
	mustPut(t, tree, 6, "f")
	mustPut(t, tree, 7, "g")
	mustPut(t, tree, 3, "c")
	mustPut(t, tree, 4, "d")
	mustPut(t, tree, 1, "x") // overwrite
	mustPut(t, tree, 2, "b")

	if actualValue, expectedValue := tree.LeftKey(), 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.LeftValue(), "x"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.RightKey(), 7; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.RightValue(), "g"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIteratorValuesAndKeys(t *testing.T) {
	tree, err := NewWithIntComparator(4)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 4, "d")
	mustPut(t, tree, 5, "e")
	mustPut(t, tree, 6, "f")
	mustPut(t, tree, 3, "c")
	mustPut(t, tree, 1, "a")
	mustPut(t, tree, 7, "g")
	mustPut(t, tree, 2, "b")
	mustPut(t, tree, 1, "x") // override
	slog.Info("Get keys", "keys", tree.Keys())

	if actualValue, expectedValue := fmt.Sprintf("%s%s%s%s%s%s%s", tree.Values()...), "xbcdefg"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue := tree.Len(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}
}

func TestBTreeIteratorNextOnEmpty(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	it := tree.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty tree")
	}
}

func TestBTreeIteratorPrevOnEmpty(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	it := tree.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty tree")
	}
}

func TestBTreeIterator1Next(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 5, "e")
	mustPut(t, tree, 6, "f")
	mustPut(t, tree, 7, "g")
	mustPut(t, tree, 3, "c")
	mustPut(t, tree, 4, "d")
	mustPut(t, tree, 1, "x")
	mustPut(t, tree, 2, "b")
	mustPut(t, tree, 1, "a") // overwrite
	it := tree.Iterator()

	count := 0
	for it.Next() {
		count++

		key := it.Key()
		switch key {
		case count:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}

	if actualValue, expectedValue := count, tree.Len(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator1Prev(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 5, "e")
	mustPut(t, tree, 6, "f")
	mustPut(t, tree, 7, "g")
	mustPut(t, tree, 3, "c")
	mustPut(t, tree, 4, "d")
	mustPut(t, tree, 1, "x")
	mustPut(t, tree, 2, "b")
	mustPut(t, tree, 1, "a") // overwrite

	it := tree.Iterator()
	for it.Next() {
	}

	countDown := tree.size

	for it.Prev() {
		key := it.Key()
		switch key {
		case countDown:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}

		countDown--
	}

	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator2Next(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 3, "c")
	mustPut(t, tree, 1, "a")
	mustPut(t, tree, 2, "b")
	it := tree.Iterator()

	count := 0
	for it.Next() {
		count++

		key := it.Key()
		switch key {
		case count:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}

	if actualValue, expectedValue := count, tree.Len(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator2Prev(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 3, "c")
	mustPut(t, tree, 1, "a")
	mustPut(t, tree, 2, "b")

	it := tree.Iterator()
	for it.Next() {
	}

	countDown := tree.size

	for it.Prev() {
		key := it.Key()
		switch key {
		case countDown:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}

		countDown--
	}

	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator3Next(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 1, "a")
	it := tree.Iterator()

	count := 0
	for it.Next() {
		count++

		key := it.Key()
		switch key {
		case count:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}

	if actualValue, expectedValue := count, tree.Len(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator3Prev(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 1, "a")

	it := tree.Iterator()
	for it.Next() {
	}

	countDown := tree.size

	for it.Prev() {
		key := it.Key()
		switch key {
		case countDown:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}

		countDown--
	}

	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator4Next(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 13, 5)
	mustPut(t, tree, 8, 3)
	mustPut(t, tree, 17, 7)
	mustPut(t, tree, 1, 1)
	mustPut(t, tree, 11, 4)
	mustPut(t, tree, 15, 6)
	mustPut(t, tree, 25, 9)
	mustPut(t, tree, 6, 2)
	mustPut(t, tree, 22, 8)
	mustPut(t, tree, 27, 10)
	it := tree.Iterator()

	count := 0
	for it.Next() {
		count++

		value := it.Value()
		switch value {
		case count:
			if actualValue, expectedValue := value, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := value, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}

	if actualValue, expectedValue := count, tree.Len(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIterator4Prev(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 13, 5)
	mustPut(t, tree, 8, 3)
	mustPut(t, tree, 17, 7)
	mustPut(t, tree, 1, 1)
	mustPut(t, tree, 11, 4)
	mustPut(t, tree, 15, 6)
	mustPut(t, tree, 25, 9)
	mustPut(t, tree, 6, 2)
	mustPut(t, tree, 22, 8)
	mustPut(t, tree, 27, 10)
	it := tree.Iterator()
	count := tree.Len()

	for it.Next() {
	}

	for it.Prev() {
		value := it.Value()
		switch value {
		case count:
			if actualValue, expectedValue := value, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := value, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}

		count--
	}

	if actualValue, expectedValue := count, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBTreeIteratorBegin(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 3, "c")
	mustPut(t, tree, 1, "a")
	mustPut(t, tree, 2, "b")
	it := tree.Iterator()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.Begin()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	for it.Next() {
	}

	it.Begin()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.Next()

	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
	}
}

func TestBTreeIteratorEnd(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	it := tree.Iterator()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.End()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	mustPut(t, tree, 3, "c")
	mustPut(t, tree, 1, "a")
	mustPut(t, tree, 2, "b")

	it.End()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.Prev()

	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}
}

func TestBTreeIteratorFirst(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 3, "c")
	mustPut(t, tree, 1, "a")
	mustPut(t, tree, 2, "b")

	it := tree.Iterator()
	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
	}
}

func TestBTreeIteratorLast(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 3, "c")
	mustPut(t, tree, 1, "a")
	mustPut(t, tree, 2, "b")

	it := tree.Iterator()
	if actualValue, expectedValue := it.Last(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}
}

func TestBTree_search(t *testing.T) {
	{
		tree, err := NewWithIntComparator(3)
		if err != nil {
			t.Fatal(err)
		}

		tree.Root = &Node[int, any]{Entries: []*Entry[int, any]{}, Children: make([]*Node[int, any], 0)}

		tests := [][]interface{}{
			{0, 0, false},
		}
		for _, test := range tests {
			index, found := tree.search(tree.Root, test[0])
			if actualValue, expectedValue := index, test[1]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}

			if actualValue, expectedValue := found, test[2]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
	{
		tree, err := NewWithIntComparator(3)
		if err != nil {
			t.Fatal(err)
		}

		tree.Root = &Node[int, any]{Entries: []*Entry[int, any]{{2, 0}, {4, 1}, {6, 2}}, Children: []*Node[int, any]{}}

		tests := [][]interface{}{
			{0, 0, false},
			{1, 0, false},
			{2, 0, true},
			{3, 1, false},
			{4, 1, true},
			{5, 2, false},
			{6, 2, true},
			{7, 3, false},
		}
		for _, test := range tests {
			index, found := tree.search(tree.Root, test[0])
			if actualValue, expectedValue := index, test[1]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}

			if actualValue, expectedValue := found, test[2]; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
}

func assertValidTree(t *testing.T, tree *Tree[int, any], expectedSize int) {
	if actualValue, expectedValue := tree.size, expectedSize; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for tree size", actualValue, expectedValue)
	}
}

func assertValidTreeNode(t *testing.T, node *Node[int, any], expectedEntries int, expectedChildren int, keys []int, hasParent bool) {
	if actualValue, expectedValue := node.Parent != nil, hasParent; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for hasParent", actualValue, expectedValue)
	}

	if actualValue, expectedValue := len(node.Entries), expectedEntries; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for entries size", actualValue, expectedValue)
	}

	if actualValue, expectedValue := len(node.Children), expectedChildren; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for children size", actualValue, expectedValue)
	}

	for i, key := range keys {
		if actualValue, expectedValue := node.Entries[i].Key, key; actualValue != expectedValue {
			t.Errorf("Got %v expected %v for key", actualValue, expectedValue)
		}
	}
}

func TestBTreeComparatorTypeMismatch(t *testing.T) {
	// Verify that Put, Get, Remove return errors (not panics) when the
	// comparator receives a key type it cannot assert.
	t.Run("Put returns error on type mismatch", func(t *testing.T) {
		tree, err := NewWith[string, any](3, comparator.IntComparator)
		if err != nil {
			t.Fatal(err)
		}

		err = tree.Put("hello", "world")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, comparator.ErrInvalidType) {
			t.Errorf("expected error wrapping comparator.ErrInvalidType, got: %v", err)
		}
	})

	t.Run("Get returns error on type mismatch", func(t *testing.T) {
		tree, err := NewWith[string, any](3, comparator.IntComparator)
		if err != nil {
			t.Fatal(err)
		}

		// Populate root directly so Get reaches the comparison path.
		tree.Root = &Node[string, any]{
			Entries:  []*Entry[string, any]{{Key: "a", Value: 1}},
			Children: []*Node[string, any]{},
		}
		tree.size = 1

		_, _, err = tree.Get("b")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, comparator.ErrInvalidType) {
			t.Errorf("expected error wrapping comparator.ErrInvalidType, got: %v", err)
		}
	})

	t.Run("Remove returns error on type mismatch", func(t *testing.T) {
		tree, err := NewWith[string, any](3, comparator.IntComparator)
		if err != nil {
			t.Fatal(err)
		}

		// Populate root directly so Remove reaches the comparison path.
		tree.Root = &Node[string, any]{
			Entries:  []*Entry[string, any]{{Key: "a", Value: 1}},
			Children: []*Node[string, any]{},
		}
		tree.size = 1

		err = tree.Remove("a")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, comparator.ErrInvalidType) {
			t.Errorf("expected error wrapping comparator.ErrInvalidType, got: %v", err)
		}
	})
}

func TestBTreeEdgeCases_EmptyTree(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "Empty returns true",
			fn: func(t *testing.T) {
				if !tree.Empty() {
					t.Error("expected Empty() == true")
				}
			},
		},
		{
			name: "Len returns zero",
			fn: func(t *testing.T) {
				if tree.Len() != 0 {
					t.Errorf("expected Len() == 0, got %d", tree.Len())
				}
			},
		},
		{
			name: "Height returns zero",
			fn: func(t *testing.T) {
				if tree.Height() != 0 {
					t.Errorf("expected Height() == 0, got %d", tree.Height())
				}
			},
		},
		{
			name: "Get on empty returns not found",
			fn: func(t *testing.T) {
				val, found, err := tree.Get(42)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if found {
					t.Errorf("expected found == false, got true with value %v", val)
				}
			},
		},
		{
			name: "Remove on empty is no-op",
			fn: func(t *testing.T) {
				if err := tree.Remove(42); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				assertValidTree(t, tree, 0)
			},
		},
		{
			name: "Left and Right return nil",
			fn: func(t *testing.T) {
				if tree.Left() != nil {
					t.Error("expected Left() == nil")
				}

				if tree.Right() != nil {
					t.Error("expected Right() == nil")
				}
			},
		},
		{
			name: "LeftKey LeftValue RightKey RightValue return nil",
			fn: func(t *testing.T) {
				if tree.LeftKey() != nil {
					t.Error("expected LeftKey() == nil")
				}

				if tree.LeftValue() != nil {
					t.Error("expected LeftValue() == nil")
				}

				if tree.RightKey() != nil {
					t.Error("expected RightKey() == nil")
				}

				if tree.RightValue() != nil {
					t.Error("expected RightValue() == nil")
				}
			},
		},
		{
			name: "Keys and Values return empty slices",
			fn: func(t *testing.T) {
				if keys := tree.Keys(); len(keys) != 0 {
					t.Errorf("expected empty Keys(), got %v", keys)
				}

				if vals := tree.Values(); len(vals) != 0 {
					t.Errorf("expected empty Values(), got %v", vals)
				}
			},
		},
		{
			name: "String on empty tree",
			fn: func(t *testing.T) {
				s := tree.String()
				if s != "BTree\n" {
					t.Errorf("expected %q, got %q", "BTree\n", s)
				}
			},
		},
		{
			name: "Clear on empty tree is no-op",
			fn: func(t *testing.T) {
				tree.Clear()
				assertValidTree(t, tree, 0)

				if tree.Root != nil {
					t.Error("expected Root == nil after Clear")
				}
			},
		},
		{
			name: "Iterator First and Last return false on empty",
			fn: func(t *testing.T) {
				it := tree.Iterator()
				if it.First() {
					t.Error("expected First() == false on empty tree")
				}

				it = tree.Iterator()
				if it.Last() {
					t.Error("expected Last() == false on empty tree")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, tc.fn)
	}
}

func TestBTreeEdgeCases_SingleElement(t *testing.T) {
	cases := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "boundary accessors on single element",
			fn: func(t *testing.T) {
				tree, err := NewWithIntComparator(3)
				if err != nil {
					t.Fatal(err)
				}

				mustPut(t, tree, 1, "a")

				if tree.LeftKey() != 1 {
					t.Errorf("expected LeftKey() == 1, got %v", tree.LeftKey())
				}

				if tree.LeftValue() != "a" {
					t.Errorf("expected LeftValue() == a, got %v", tree.LeftValue())
				}

				if tree.RightKey() != 1 {
					t.Errorf("expected RightKey() == 1, got %v", tree.RightKey())
				}

				if tree.RightValue() != "a" {
					t.Errorf("expected RightValue() == a, got %v", tree.RightValue())
				}

				if tree.Height() != 1 {
					t.Errorf("expected Height() == 1, got %d", tree.Height())
				}
			},
		},
		{
			name: "remove single element leaves tree empty",
			fn: func(t *testing.T) {
				tree, err := NewWithIntComparator(3)
				if err != nil {
					t.Fatal(err)
				}

				mustPut(t, tree, 5, "x")
				mustRemove(t, tree, 5)

				assertValidTree(t, tree, 0)

				if !tree.Empty() {
					t.Error("expected Empty() == true after removing sole element")
				}

				if tree.Root != nil {
					t.Error("expected Root == nil after removing sole element")
				}
			},
		},
		{
			name: "Clear on single element tree",
			fn: func(t *testing.T) {
				tree, err := NewWithIntComparator(3)
				if err != nil {
					t.Fatal(err)
				}

				mustPut(t, tree, 10, "z")
				tree.Clear()

				assertValidTree(t, tree, 0)

				if tree.Root != nil {
					t.Error("expected Root == nil after Clear")
				}

				if tree.LeftKey() != nil {
					t.Error("expected LeftKey() == nil after Clear")
				}

				if tree.RightKey() != nil {
					t.Error("expected RightKey() == nil after Clear")
				}
			},
		},
		{
			name: "String on single element tree",
			fn: func(t *testing.T) {
				tree, err := NewWithIntComparator(3)
				if err != nil {
					t.Fatal(err)
				}

				mustPut(t, tree, 42, "answer")

				s := tree.String()
				if s != "BTree\n42\n" {
					t.Errorf("expected %q, got %q", "BTree\n42\n", s)
				}
			},
		},
		{
			name: "iterator on single element",
			fn: func(t *testing.T) {
				tree, err := NewWithIntComparator(3)
				if err != nil {
					t.Fatal(err)
				}

				mustPut(t, tree, 7, "g")

				it := tree.Iterator()
				if !it.First() {
					t.Fatal("expected First() == true")
				}

				if it.Key() != 7 || it.Value() != "g" {
					t.Errorf("expected (7, g), got (%v, %v)", it.Key(), it.Value())
				}

				if it.Next() {
					t.Error("expected no more elements after First()")
				}

				it = tree.Iterator()
				if !it.Last() {
					t.Fatal("expected Last() == true")
				}

				if it.Key() != 7 || it.Value() != "g" {
					t.Errorf("expected (7, g), got (%v, %v)", it.Key(), it.Value())
				}

				if it.Prev() {
					t.Error("expected no more elements before Last()")
				}
			},
		},
		{
			name: "put then remove then put again",
			fn: func(t *testing.T) {
				tree, err := NewWithIntComparator(3)
				if err != nil {
					t.Fatal(err)
				}

				mustPut(t, tree, 1, "first")
				mustRemove(t, tree, 1)
				assertValidTree(t, tree, 0)

				mustPut(t, tree, 1, "second")
				assertValidTree(t, tree, 1)

				val, found := mustGet(t, tree, 1)
				if !found || val != "second" {
					t.Errorf("expected (second, true), got (%v, %v)", val, found)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, tc.fn)
	}
}

func TestBTreeEdgeCases_DuplicateKeys(t *testing.T) {
	cases := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "put same key updates value and keeps size",
			fn: func(t *testing.T) {
				tree, err := NewWithIntComparator(3)
				if err != nil {
					t.Fatal(err)
				}

				mustPut(t, tree, 1, "a")
				mustPut(t, tree, 1, "b")
				mustPut(t, tree, 1, "c")

				assertValidTree(t, tree, 1)

				val, found := mustGet(t, tree, 1)
				if !found || val != "c" {
					t.Errorf("expected (c, true), got (%v, %v)", val, found)
				}
			},
		},
		{
			name: "update key in leaf node",
			fn: func(t *testing.T) {
				tree, err := NewWithIntComparator(3)
				if err != nil {
					t.Fatal(err)
				}

				mustPut(t, tree, 1, "a")
				mustPut(t, tree, 2, "b")
				assertValidTree(t, tree, 2)

				mustPut(t, tree, 2, "updated")
				assertValidTree(t, tree, 2)

				val, found := mustGet(t, tree, 2)
				if !found || val != "updated" {
					t.Errorf("expected (updated, true), got (%v, %v)", val, found)
				}
			},
		},
		{
			name: "update key in internal node",
			fn: func(t *testing.T) {
				// Build a tree with order=3 where key 2 is in the root (internal node):
				// after inserting 1,2,3 the root holds [2] with children [1] and [3].
				tree, err := NewWithIntComparator(3)
				if err != nil {
					t.Fatal(err)
				}

				mustPut(t, tree, 1, "a")
				mustPut(t, tree, 2, "b")
				mustPut(t, tree, 3, "c")

				assertValidTree(t, tree, 3)
				assertValidTreeNode(t, tree.Root, 1, 2, []int{2}, false)

				// Update key 2 which lives in internal root node
				mustPut(t, tree, 2, "B-updated")
				assertValidTree(t, tree, 3)

				val, found := mustGet(t, tree, 2)
				if !found || val != "B-updated" {
					t.Errorf("expected (B-updated, true), got (%v, %v)", val, found)
				}
			},
		},
		{
			name: "update with nil value",
			fn: func(t *testing.T) {
				tree, err := NewWithIntComparator(3)
				if err != nil {
					t.Fatal(err)
				}

				mustPut(t, tree, 1, "a")
				mustPut(t, tree, 1, nil)
				assertValidTree(t, tree, 1)

				val, found := mustGet(t, tree, 1)
				if !found {
					t.Fatal("expected key 1 to be found")
				}

				if val != nil {
					t.Errorf("expected nil value, got %v", val)
				}
			},
		},
		{
			name: "many duplicate updates then delete",
			fn: func(t *testing.T) {
				tree, err := NewWithIntComparator(3)
				if err != nil {
					t.Fatal(err)
				}

				for i := 1; i <= 10; i++ {
					mustPut(t, tree, i, i)
				}

				assertValidTree(t, tree, 10)

				// Update all keys
				for i := 1; i <= 10; i++ {
					mustPut(t, tree, i, i*100)
				}

				assertValidTree(t, tree, 10)

				// Verify all values updated
				for i := 1; i <= 10; i++ {
					val, found := mustGet(t, tree, i)
					if !found || val != i*100 {
						t.Errorf("key %d: expected (%d, true), got (%v, %v)", i, i*100, val, found)
					}
				}

				// Remove all
				for i := 1; i <= 10; i++ {
					mustRemove(t, tree, i)
				}

				assertValidTree(t, tree, 0)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, tc.fn)
	}
}

func TestBTreeNewWithInvalidOrder(t *testing.T) {
	cases := []struct {
		name  string
		order int
	}{
		{name: "order 0", order: 0},
		{name: "order 1", order: 1},
		{name: "order 2", order: 2},
		{name: "order -1", order: -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewWithIntComparator(tc.order)
			if !errors.Is(err, ErrInvalidOrder) {
				t.Errorf("expected ErrInvalidOrder, got %v", err)
			}

			_, err = NewWithStringComparator(tc.order)
			if !errors.Is(err, ErrInvalidOrder) {
				t.Errorf("expected ErrInvalidOrder, got %v", err)
			}
		})
	}
}

func TestBTreeStringMultiLevel(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	mustPut(t, tree, 1, nil)
	mustPut(t, tree, 2, nil)
	mustPut(t, tree, 3, nil)

	// Root=[2], Children=[1],[3] — String should contain all keys
	s := tree.String()
	for _, key := range []string{"1", "2", "3"} {
		if !strings.Contains(s, key) {
			t.Errorf("expected String() to contain %q, got %q", key, s)
		}
	}
}

func TestBTreeEntryString(t *testing.T) {
	cases := []struct {
		name string
		key  int
		want string
	}{
		{name: "positive key", key: 42, want: "42"},
		{name: "zero key", key: 0, want: "0"},
		{name: "negative key", key: -1, want: "-1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			entry := &Entry[int, string]{Key: tc.key, Value: "v"}
			if got := entry.String(); got != tc.want {
				t.Errorf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestBTreeClearAfterInsertions(t *testing.T) {
	tree, err := NewWithIntComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 20; i++ {
		mustPut(t, tree, i, i)
	}

	assertValidTree(t, tree, 20)

	tree.Clear()
	assertValidTree(t, tree, 0)

	if tree.Root != nil {
		t.Error("expected Root == nil after Clear")
	}

	if tree.Height() != 0 {
		t.Errorf("expected Height() == 0 after Clear, got %d", tree.Height())
	}

	if !tree.Empty() {
		t.Error("expected Empty() == true after Clear")
	}

	// Verify tree is usable after Clear
	mustPut(t, tree, 99, "reuse")
	assertValidTree(t, tree, 1)

	val, found := mustGet(t, tree, 99)
	if !found || val != "reuse" {
		t.Errorf("expected (reuse, true), got (%v, %v)", val, found)
	}
}

func benchmarkSizes() []struct {
	name string
	size int
} {
	return []struct {
		name string
		size int
	}{
		{"n=100", 100},
		{"n=1000", 1000},
		{"n=10000", 10000},
		{"n=100000", 100000},
	}
}

func BenchmarkPut(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			tree, err := NewWithIntComparator(128)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				for n := 0; n < tc.size; n++ {
					mustPut(b, tree, n, struct{}{})
				}
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			tree, err := NewWithIntComparator(128)
			if err != nil {
				b.Fatal(err)
			}

			for n := 0; n < tc.size; n++ {
				mustPut(b, tree, n, struct{}{})
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				for n := 0; n < tc.size; n++ {
					mustGet(b, tree, n)
				}
			}
		})
	}
}

func BenchmarkRemove(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			tree, err := NewWithIntComparator(128)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				b.StopTimer()

				for n := 0; n < tc.size; n++ {
					mustPut(b, tree, n, struct{}{})
				}

				b.StartTimer()

				for n := 0; n < tc.size; n++ {
					mustRemove(b, tree, n)
				}
			}
		})
	}
}

func BenchmarkIterate(b *testing.B) {
	for _, tc := range benchmarkSizes() {
		b.Run(tc.name, func(b *testing.B) {
			tree, err := NewWithIntComparator(128)
			if err != nil {
				b.Fatal(err)
			}

			for n := 0; n < tc.size; n++ {
				mustPut(b, tree, n, struct{}{})
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				it := tree.Iterator()
				for it.Next() {
					_ = it.Key()
					_ = it.Value()
				}
			}
		})
	}
}

func TestCopy(t *testing.T) {
	cases := []struct {
		name string
		keys []int
	}{
		{
			name: "empty tree",
			keys: nil,
		},
		{
			name: "single element",
			keys: []int{1},
		},
		{
			name: "multiple elements causing splits",
			keys: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tree, err := NewWithIntComparator(3)
			if err != nil {
				t.Fatal(err)
			}

			for _, k := range tc.keys {
				mustPut(t, tree, k, fmt.Sprintf("v%d", k))
			}

			clone := tree.Copy()

			if clone.Len() != tree.Len() {
				t.Fatalf("clone.Len()=%d, want %d", clone.Len(), tree.Len())
			}

			// Verify all keys present in clone
			for _, k := range tc.keys {
				val, found := mustGet(t, clone, k)
				if !found {
					t.Fatalf("key %d not found in clone", k)
				}

				if val != fmt.Sprintf("v%d", k) {
					t.Fatalf("clone Get(%d)=%v, want v%d", k, val, k)
				}
			}

			// Verify independence: modify original, clone unaffected
			if len(tc.keys) > 0 {
				mustRemove(t, tree, tc.keys[0])

				_, found := mustGet(t, clone, tc.keys[0])
				if !found {
					t.Fatal("clone was affected by removing from original")
				}
			}
		})
	}
}

func TestCopyString(t *testing.T) {
	tree, err := NewWithStringComparator(3)
	if err != nil {
		t.Fatal(err)
	}

	if err := tree.Put("a", 1); err != nil {
		t.Fatal(err)
	}

	if err := tree.Put("b", 2); err != nil {
		t.Fatal(err)
	}

	clone := tree.Copy()

	val, found, err := clone.Get("a")
	if err != nil {
		t.Fatal(err)
	}

	if !found || val != 1 {
		t.Fatalf("clone Get(a)=%v found=%v, want 1 true", val, found)
	}

	if clone.String() != tree.String() {
		t.Fatal("clone String() does not match original")
	}
}
