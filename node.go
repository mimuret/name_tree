package name_tree

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
	"github.com/mimuret/dnsutils"
)

var (
	ErrNotChild = fmt.Errorf("not child name")
	// ErrNotInDomain returns when arg node is in-domain.
	ErrNotInDomain = fmt.Errorf("name is not subdomain")
	// ErrRemoveItself by RemoveChildNode when remove itself.
	ErrRemoveItself = fmt.Errorf("can not remove itself")
)

type Node[T any] struct {
	value    *T
	name     string
	parent   *Node[T]
	children map[string]*Node[T]
}

func NewNode[T any](name string, value *T) *Node[T] {
	return &Node[T]{
		name:     dns.CanonicalName(name),
		children: make(map[string]*Node[T]),
		value:    value,
	}
}

func (n *Node[T]) insertNode(node *Node[T]) error {
	if dnsutils.Equals(n.GetName(), node.GetName()) {
		return ErrNotChild
	}
	if !dns.IsSubDomain(n.GetName(), node.GetName()) {
		return ErrNotChild
	}
	nlabels := dns.SplitDomainName(n.GetName())
	nodelabels := dns.SplitDomainName(node.GetName())
	if len(nodelabels)-len(nlabels) == 1 {
		// set child node
		node.parent = n
		n.children[nodelabels[0]] = node
		return nil
	}
	childLabel := nodelabels[len(nodelabels)-len(nlabels)-1]
	childNode, exist := n.children[childLabel]
	if !exist {
		// set ENT
		childNode = NewNode[T](strings.Join(nodelabels[len(nodelabels)-len(nlabels)-1:], "."), nil)
		childNode.parent = n
	}
	if err := childNode.insertNode(node); err != nil {
		return err
	}
	n.children[childLabel] = childNode
	return nil
}

func (n *Node[T]) getNode(name string) (node *Node[T], strict bool) {
	if !dns.IsSubDomain(n.GetName(), name) {
		return n, false
	}
	if dnsutils.Equals(n.GetName(), name) {
		return n, true
	}
	for _, child := range n.children {
		if dns.IsSubDomain(child.GetName(), name) {
			return child.getNode(name)
		}
	}
	return n, false
}

func (n *Node[T]) removeNode(name string) *Node[T] {
	name = dns.CanonicalName(name)
	child, exist := n.getNode(name)
	if !exist {
		// node not found
		return nil
	}
	if child.parent == nil {
		// name is root
		return nil
	}
	labels := dns.SplitDomainName(name)
	delete(child.parent.children, labels[0])
	return child
}

// ENT eval func
func IsENT[T any](v *T) bool {
	return v == nil
}

func (n *Node[T]) removeNodeWithENT(f func(*T) bool) {
	if len(n.children) > 0 {
		// has children
		return
	}
	if !f(n.value) {
		// value exist
		return
	}
	if n.parent == nil {
		// n is root
		return
	}
	labels := dns.SplitDomainName(n.GetName())
	delete(n.parent.children, labels[0])
	n.parent.removeNodeWithENT(f)
}

func (n *Node[T]) iterateNode(f func(nn *Node[T]) (bool, error)) error {
	ok, err := f(n)
	if err != nil {
		// stop iterate all nodes if err != nil
		return err
	}
	if !ok {
		// not iterate child nodes if ok == false
		return nil
	}
	for _, label := range n.getChildLabel() {
		err := n.children[label].iterateNode(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *Node[T]) getChildLabel() []string {
	var names = make([]string, 0, len(n.children))
	for name := range n.children {
		names = append(names, name)
	}
	dnsutils.SortNames(names)
	return names
}

// GetName returns canonical name
func (n *Node[T]) GetName() string {
	return n.name
}

func (n *Node[T]) SetValue(val *T) {
	n.value = val
}

func (n *Node[T]) ClearValue() {
	n.value = nil
}

func (n *Node[T]) GetVault() *T {
	return n.value
}
