// Copyright 2020 Matthew Burr
//
// Package binary implements a binary tree.
package binary

// A Node is an element in a binary tree.
type Node struct {
	// left and right are pointers to the children of the Node.
	left, right *Node

	// Value contains the data for this node.
	Value interface{}
}

// Left returns this Node's left child or nil.
func (n *Node) Left() *Node {
	return n.left
}

// Right returns this Node's right child or nil.
func (n *Node) Right() *Node {
	return n.right
}

// A Tree represents a binary search tree.
// The Nodes of the tree are organized so that all of
// the Nodes in a given Node's left sub-tree have values
// that are less than the Nodes in its right sub-tree.
type Tree struct {
	root *Node
}

const (
	Done     = true
	Continue = false
)

// A VistorFunc is an operation to perform on Nodes of a tree.
// Visitor functions accept an interface{] as their pattern;
// the value of the interface{} is the Value of a Node in the tree.
// Visitors functions return a bool indicating whether to stop
// traversing the tree. A return value of "true" means we are done
// and should traverse the tree no futher. A value of "false" means
// we are not done and should continue.
type VisitorFunc func(interface{}) bool

// visitInOrder is used internally to visit all of the Nodes of a
// subtree in order, i.e. LVR.
func visitInOrder(with VisitorFunc, subtreeRoot *Node) bool {
	if subtreeRoot == nil {
		return Continue
	}

	if Done == visitInOrder(with, subtreeRoot.Left()) {
		return Done
	}

	if Done == with(subtreeRoot.Value) {
		return Done
	}

	return visitInOrder(with, subtreeRoot.Right())
}

// VisitInOrder visits the Nodes of a Tree in order.
// At each Node, it applies the visitorFunc to the
// value of the Node. VisitInOrder stops visiting
// when the VisitorFunc returns Done ("true") after
// visiting the Node.
func (t *Tree) VisitInOrder(with VisitorFunc) {
	_ = visitInOrder(with, t.root)
	return
}

// visitInReverse is used internally to visit all the Nodes of a
// subtree in reverse order, i.e. RVL
func visitInReverse(with VisitorFunc, subtreeRoot *Node) bool {
	if subtreeRoot == nil {
		return Continue
	}

	if Done == visitInOrder(with, subtreeRoot.Right()) {
		return Done
	}

	if Done == with(subtreeRoot.Value) {
		return Done
	}

	return visitInOrder(with, subtreeRoot.Left())
}

// VisitInReverse visits the Nodes of a tree in postfix or reverse order.
// At each Node, it applies the visitorFunc to the
// value of the Node. VisitInOrder stops visiting
// when the VisitorFunc returns Done ("true") after
// visiting the Node.
func (t *Tree) VisitInReverse(with VisitorFunc) {
	_ = visitInReverse(with, t.root)
	return
}
