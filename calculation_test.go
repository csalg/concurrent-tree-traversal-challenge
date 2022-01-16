package challenge

import "testing"

func TestEmpty(t *testing.T) {
	node := Node{1, 1, 0, []Node{}}
	indexingSize := CalculateIndexingSize(node)
	if indexingSize != node.IndexingSize {
		t.Fatalf(`Expected indexing size to be %d, but got %d instead`, node.IndexingSize, indexingSize)
	}
}

// Test cases:
// Root with no children should return size of root
// Should handle a tree structure with arbitrary fan-out and different heights
