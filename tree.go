package critbit

import (
	"iter"
)

type Leaf[V any] struct {
	Key   Key
	Value V
}

type Inner[V any] struct {
	bit   int
	child [2]Node[V]
}

type Node[V any] struct {
	Leaf  *Leaf[V]
	Inner *Inner[V]
}

func (n Node[V]) Longest(key Key) (V, bool) {
	if inner := n.Inner; inner != nil {
		dir := key.Direction(inner.bit)
		val, found := inner.child[dir].Longest(key)
		if found {
			return val, found
		}
		if dir == 1 {
			return inner.child[0].Longest(key)
		}
	} else if leaf := n.Leaf; leaf != nil {
		if key.HasPrefix(leaf.Key) {
			return leaf.Value, true
		}
	}
	var val V
	return val, false
}

type Tree[V any] struct {
	nums int
	root Node[V]
}

func (t *Tree[V]) Len() int {
	return t.nums
}

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

func (t *Tree[V]) Set(key Key, value V) {
	leaf := t.findLeaf(key)
	if leaf == nil {
		t.root.Leaf = &Leaf[V]{
			Key:   key,
			Value: value,
		}
		t.nums++
		return
	}
	bit := leaf.Key.Critbit(key)
	if bit == -1 {
		// replace
		leaf.Value = value
		return
	}
	n := t.findNode(key, bit)
	t.insertNode(n, &Leaf[V]{Key: key, Value: value}, bit)
}

func (t *Tree[V]) Delete(key Key) {
	var p *Node[V]
	var dir int
	n := &t.root
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
		return
	}
	if !leaf.Key.Equal(key) {
		return
	}
	if p == nil {
		t.root = Node[V]{}
	} else {
		*p = p.Inner.child[dir^1]
	}
	t.nums--
}

func (t *Tree[V]) Longest(key Key) (V, bool) {
	return t.root.Longest(key)
}

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

func (t *Tree[V]) findNode(key Key, bit int) *Node[V] {
	n := &t.root
	for {
		inner := n.Inner
		if inner == nil {
			break
		}
		if inner.bit > bit {
			break
		}
		dir := key.Direction(inner.bit)
		n = &inner.child[dir]
	}
	return n
}

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

func (t *Tree[V]) insertNode(n *Node[V], leaf *Leaf[V], bit int) {
	dir := leaf.Key.Direction(bit)
	inner := new(Inner[V])
	inner.child[dir].Leaf = leaf
	inner.child[dir^1] = *n
	inner.bit = bit
	*n = Node[V]{Inner: inner}
	t.nums++
}
