package critbit

type Scanner[V any] struct {
	dir   int
	stack []Node[V]
}

func NewScanner[V any](root Node[V], reverse bool) *Scanner[V] {
	s := new(Scanner[V])
	if reverse {
		s.dir = 1
	}
	s.push(root)
	return s
}

func (s *Scanner[V]) Scan() *Leaf[V] {
	n := s.pop()
	for {
		inner := n.Inner
		if inner == nil {
			break
		}
		s.push(inner.child[s.dir^1])
		n = inner.child[s.dir]
	}
	return n.Leaf
}

func (s *Scanner[V]) push(node Node[V]) {
	s.stack = append(s.stack, node)
}

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
