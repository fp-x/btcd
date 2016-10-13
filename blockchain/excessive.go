package blockchain

// Some prototype objects to describe excessive block handling.

// mockChain is a simulated blockchain for prototyping.
type mockChain struct {
	// The current main chain.
	main *mockNode

	// An excessive block cannot be part of the main chain unless it
	// has a chain after it of at least this length.
	AcceptDepth uint32
}

func makeMockChain(acceptDepth uint32) *mockChain {
	return &mockChain{
		AcceptDepth: acceptDepth,
	}
}

// Add a new block.
func (ch *mockChain) addBlock(node *mockNode) {
	if node == nil {
		return
	}

	// The maximum block that could be accepted as the
	// head of the main chain.
	max := maxAcceptable(node, ch.AcceptDepth)

	h := height(max)
	maxh := height(ch.main)

	if h > maxh {
		ch.main = max
	}
}

// mockNode is a pretend block node.
type mockNode struct {
	parent *mockNode

	excessive bool
}

func height(node *mockNode) uint32 {
	if node == nil {
		return 0
	}

	return height(node.parent) + 1
}

// Is there an excessive block in this chain?
func excessiveChain(node *mockNode, depth uint32) bool {
	if node == nil || depth == 0 {
		return false
	}

	if node.excessive {
		return true
	}

	return excessiveChain(node.parent, depth-1)
}

// The maximum block that is acceptable.
func maxAcceptable(node *mockNode, acceptDepth uint32) *mockNode {
	if acceptDepth == 0 || node == nil {
		return node
	}

	maxNode := maxAcceptable(node.parent, acceptDepth-1)

	if maxNode == node.parent && !node.excessive {
		return node
	}

	return maxNode
}
