package challenge

func CalculateIndexingSize(root Node) int {
	// NOTES
	// It seems that there isn't a popular queue implementation in the standard library that everyone uses.
	// Since the queue is not very large and it won't require persistence or anything like that I am actually going to
	// use a channel as a queue. Could otherwise use a library or write my own queue implementation but I suppose that's
	// not the point of the challenge.
	// Idea from: https://stackoverflow.com/a/39598511

	queue := make(chan Node, 100) // Actually, it doesn't have to be 100, just the maximum fan-out * maximum height
	queue <- root
	// 1. Calculate the indexing size of all nodes
	for {

	}
	// 2. Sum up all the indexing sizes
	return 0
}

type Node struct {
	Size int
	// Go doesn't seem to support nullable / optional values,
	// so I will use -1 as indexing size when it is not known
	IndexingSize int
	Parent       *Node
	Children     []*Node
}
