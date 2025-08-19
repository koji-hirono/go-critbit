// Package critbit implements a crit-bit tree (also known as PATRICIA tree)
// data structure for efficient storage and retrieval of binary keys.
//
// Crit-bit trees provide O(k) operations where k is the key length,
// making them excellent for applications requiring fast lookups with
// ordered iteration and longest prefix matching capabilities.
package critbit

import (
	"iter"
)

// Leaf represents a terminal node in the crit-bit tree containing
// a key-value pair. All actual data is stored in leaf nodes.
type Leaf[V any] struct {
	Key   Key
	Value V
}

// Inner represents an internal node in the crit-bit tree.
// It contains a critical bit position and exactly two child nodes.
// The critical bit determines which child to follow during traversal.
type Inner[V any] struct {
	// bit encodes the critical bit position
	// where the two subtrees diverge
	bit int
	// child contains exactly two children:
	// [0] for left, [1] for right
	child [2]Node[V]
}

// Node represents either an internal node or a leaf node
// in the crit-bit tree.
// Exactly one of Leaf or Inner will be non-nil for any valid node.
type Node[V any] struct {
	Leaf  *Leaf[V]
	Inner *Inner[V]
}

// Longest performs longest prefix matching starting from this node.
// It finds the longest key in the subtree that is a prefix of the
// given key.
// This is particularly useful for applications like IP routing.
//
// Returns the value associated with the longest matching prefix and true,
// or the zero value and false if no prefix match is found.
func (n Node[V]) Longest(key Key) (V, bool) {
	if inner := n.Inner; inner != nil {
		// Internal node: choose direction and recurse
		dir := key.Direction(inner.bit)
		val, found := inner.child[dir].Longest(key)
		if found {
			return val, found
		}
		// If no match in preferred direction,
		// try the other direction
		if dir == 1 {
			return inner.child[0].Longest(key)
		}
	} else if leaf := n.Leaf; leaf != nil {
		// Leaf node: check if this key is a prefix of the
		// search key
		if key.HasPrefix(leaf.Key) {
			return leaf.Value, true
		}
	}
	var val V
	return val, false
}

// Tree represents a crit-bit tree that maps Keys to values of type V.
// The zero value of Tree is an empty tree ready for use.
//
// Tree is not safe for concurrent access. Use external synchronization
// if the tree needs to be accessed from multiple goroutines.
type Tree[V any] struct {
	nums int     // number of key-value pairs in the tree
	root Node[V] // root node of the tree
}

// Len returns the number of key-value pairs in the tree.
func (t *Tree[V]) Len() int {
	return t.nums
}

// Get retrieves the value associated with the given key.
// Returns the value and true if the key exists, or the zero value
// of V and false if the key is not found.
//
// Time complexity: O(k) where k is the length of the key in bits.
func (t *Tree[V]) Get(key Key) (V, bool) {
	var val V
	leaf := t.findLeaf(key)
	if leaf == nil {
		return val, false
	}
	if !leaf.Key.Equal(key) {
		return val, false
	}
	return leaf.Value, true
}

// Set inserts a key-value pair into the tree or updates the value
// if the key already exists.
//
// Time complexity: O(k) where k is the length of the key in bits.
func (t *Tree[V]) Set(key Key, value V) {
	leaf := t.findLeaf(key)
	if leaf == nil {
		// Tree is empty, create first leaf
		t.root.Leaf = &Leaf[V]{
			Key:   key,
			Value: value,
		}
		t.nums++
		return
	}

	bit := leaf.Key.Critbit(key)
	if bit == -1 {
		// Key already exists, replace value
		leaf.Value = value
		return
	}

	// Insert new internal node at the appropriate position
	n := t.findNode(key, bit)
	t.insertNode(n, &Leaf[V]{Key: key, Value: value}, bit)
}

// Delete removes the key-value pair with the given key from the tree.
// If the key does not exist, Delete is a no-op.
//
// Time complexity: O(k) where k is the length of the key in bits.
func (t *Tree[V]) Delete(key Key) {
	var p *Node[V] // parent of current node
	var dir int    // direction taken from parent
	n := &t.root

	// Find the leaf node and its parent
	for {
		inner := n.Inner
		if inner == nil {
			break
		}
		p = n
		dir = key.Direction(inner.bit)
		n = &inner.child[dir]
	}

	leaf := n.Leaf
	if leaf == nil {
		return // key not found
	}
	if !leaf.Key.Equal(key) {
		return // key not found
	}

	// Remove the node
	if p == nil {
		// Removing the only node in the tree
		t.root = Node[V]{}
	} else {
		// Replace parent with sibling
		*p = p.Inner.child[dir^1]
	}
	t.nums--
}

// Longest performs longest prefix matching on the entire tree.
// It finds the longest key in the tree that is a prefix of the given key.
//
// This is particularly useful for applications like IP routing where
// you need to find the most specific route that matches a destination.
//
// Returns the value associated with the longest matching prefix and true,
// or the zero value and false if no prefix match is found.
func (t *Tree[V]) Longest(key Key) (V, bool) {
	return t.root.Longest(key)
}

// Keys returns an iterator over all keys in the tree
// in lexicographical order.
// The iterator follows Go 1.23+ iterator conventions and can be used with
// range loops.
//
// Example:
//
//	for key := range tree.Keys() {
//	    fmt.Println("Key:", key)
//	}
func (t *Tree[V]) Keys() iter.Seq[Key] {
	return func(yield func(Key) bool) {
		s := NewScanner(t.root, false)
		for {
			leaf := s.Scan()
			if leaf == nil {
				break
			}
			if !yield(leaf.Key) {
				break
			}
		}
	}
}

// Values returns an iterator over all values in the tree in the order
// corresponding to their keys' lexicographical order.
// The iterator follows Go 1.23+ iterator conventions.
//
// Example:
//
//	for value := range tree.Values() {
//	    fmt.Println("Value:", value)
//	}
func (t *Tree[V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		s := NewScanner(t.root, false)
		for {
			leaf := s.Scan()
			if leaf == nil {
				break
			}
			if !yield(leaf.Value) {
				break
			}
		}
	}
}

// All returns an iterator over all key-value pairs in the tree
// in lexicographical order of keys.
// The iterator follows Go 1.23+ iterator conventions.
//
// Example:
//
//	for key, value := range tree.All() {
//	    fmt.Printf("Key: %v, Value: %v\n", key, value)
//	}
func (t *Tree[V]) All() iter.Seq2[Key, V] {
	return func(yield func(Key, V) bool) {
		s := NewScanner(t.root, false)
		for {
			leaf := s.Scan()
			if leaf == nil {
				break
			}
			if !yield(leaf.Key, leaf.Value) {
				break
			}
		}
	}
}

// findNode locates the position where a new internal node with the given
// critical bit should be inserted. It returns a pointer to the node
// that should become a child of the new internal node.
func (t *Tree[V]) findNode(key Key, bit int) *Node[V] {
	n := &t.root
	for {
		inner := n.Inner
		if inner == nil {
			break
		}
		// Stop when we find a node with a larger critical bit
		// (meaning we've found the insertion point)
		if inner.bit > bit {
			break
		}
		dir := key.Direction(inner.bit)
		n = &inner.child[dir]
	}
	return n
}

// findLeaf follows the path through the tree according to the given key
// and returns the leaf node that would contain the key if it exists.
// Returns nil if the tree is empty.
func (t *Tree[V]) findLeaf(key Key) *Leaf[V] {
	n := &t.root
	for {
		inner := n.Inner
		if inner == nil {
			break
		}
		dir := key.Direction(inner.bit)
		n = &inner.child[dir]
	}
	return n.Leaf
}

// insertNode creates a new internal node with the given critical bit
// and inserts it at the specified position in the tree.
// The existing node becomes one child,
// and the new leaf becomes the other child.
func (t *Tree[V]) insertNode(n *Node[V], leaf *Leaf[V], bit int) {
	dir := leaf.Key.Direction(bit)
	inner := new(Inner[V])
	inner.child[dir].Leaf = leaf
	inner.child[dir^1] = *n
	inner.bit = bit
	*n = Node[V]{Inner: inner}
	t.nums++
}
