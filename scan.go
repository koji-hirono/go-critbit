package critbit

// Scanner provides controlled iteration over a crit-bit tree.
// It maintains a stack to perform in-order traversal of the tree
// without recursion, making it suitable for large trees.
//
// Scanner can traverse the tree in either forward (lexicographical)
// or reverse order depending on the reverse parameter passed to NewScanner.
type Scanner[V any] struct {
	// direction for traversal: 0 for forward, 1 for reverse
	dir int
	// stack of nodes for iterative traversal
	stack []Node[V]
}

// NewScanner creates a new Scanner for iterating over a crit-bit tree.
// The root parameter specifies the root node to start traversal from.
// If reverse is true, the scanner will traverse in reverse
// lexicographical order.
//
// The scanner performs an in-order traversal of the tree, visiting all
// leaf nodes in the specified order.
//
// Example:
//
//	// Forward traversal
//	scanner := NewScanner(tree.root, false)
//	for {
//	    leaf := scanner.Scan()
//	    if leaf == nil {
//	        break
//	    }
//	    fmt.Printf("Key: %v, Value: %v\n", leaf.Key, leaf.Value)
//	}
func NewScanner[V any](root Node[V], reverse bool) *Scanner[V] {
	s := new(Scanner[V])
	if reverse {
		s.dir = 1
	}
	s.push(root)
	return s
}

// Scan returns the next leaf node in the traversal order.
// Returns nil when there are no more nodes to visit.
//
// The method implements iterative in-order traversal:
//   - For forward traversal: visits left subtree, then node,
//     then right subtree
//   - For reverse traversal: visits right subtree, then node,
//     then left subtree
//
// Time complexity: Amortized O(1) per call, O(n) for complete traversal
// where n is the number of nodes in the tree.
func (s *Scanner[V]) Scan() *Leaf[V] {
	n := s.pop()
	for {
		inner := n.Inner
		if inner == nil {
			break
		}
		// Push the node we'll visit later (opposite direction)
		s.push(inner.child[s.dir^1])
		// Continue with the node we want to visit first
		n = inner.child[s.dir]
	}
	return n.Leaf
}

// push adds a node to the traversal stack.
// Internal method used by the scanning algorithm.
func (s *Scanner[V]) push(node Node[V]) {
	s.stack = append(s.stack, node)
}

// pop removes and returns the top node from the traversal stack.
// Returns a zero Node if the stack is empty.
// Internal method used by the scanning algorithm.
func (s *Scanner[V]) pop() Node[V] {
	n := len(s.stack)
	if n == 0 {
		return Node[V]{}
	}
	n--
	node := s.stack[n]
	s.stack = s.stack[:n]
	return node
}
