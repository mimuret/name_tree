package name_tree

import (
	"github.com/miekg/dns"
)

type Tree[T any] struct {
	root *Node[T]
}

func NewTree[T any](name string, value *T) *Tree[T] {
	return &Tree[T]{
		root: NewNode(name, value),
	}
}

func (t *Tree[T]) InsertNode(node *Node[T]) error {
	return t.root.insertNode(node)
}

func (t *Tree[T]) GetNode(name string) (node *Node[T], strict bool) {
	if !dns.IsSubDomain(t.root.GetName(), name) {
		return nil, false
	}
	return t.root.getNode(name)
}

func (t *Tree[T]) RemoveNode(name string) {
	t.root.removeNode(name)
}

func (t *Tree[T]) RemoveNodeWithENT(name string, f func(*T) bool) {
	if f == nil {
		f = IsENT
	}
	child := t.root.removeNode(name)
	if child != nil {
		child.parent.removeNodeWithENT(f)
	}
}

func (t *Tree[T]) IterateNode(f func(nn *Node[T]) (bool, error)) error {
	return t.root.iterateNode(f)
}
