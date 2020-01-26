// Copyright 2020 Matthew Burr
//
// Package binary implements a binary tree.
package binary

import (
	"fmt"
	"strings"
)

const (
	// LT indicates a comparison resulted in Less Than
	LT = -1
	// EQ indicates a comparison resulted in Equal
	EQ = 0
	// GT indicates a comparison result in Greater Than
	GT = 1
)

// Interface describes the functions a type
// should support to work with the binary tree.
type Interface interface {
	// Compare compares this object to another and
	// returns an integer to indicate how the two
	// equate:
	// -1 = this object is less than the one passed in
	// 0 = this object is equal to the one passed in
	// 1 = this object is greater than the one passed in
	// This package provides the LT, EQ, and GT constants
	// to simplify this.
	Compare(interface{}) int
	// Value returns the underlying value of this interface.
	// This is the raw data that will be stored in the Node.
	Value() interface{}
	// The binary tree calls Update when an attempt is made
	// to insert a value that matches an already existing
	// value in the tree.
	// The binary tree calls update on the value that is
	// already in the tree and passes in the incoming value,
	// giving the existing value an opportunity to decide
	// what to do with the new value.
	Update(interface{})
}

// A Node is an element in a binary tree.
type Node struct {
	// left and right are pointers to the children of the Node.
	left, right *Node

	// Value contains the Value for this node.
	value Interface
}

// Compare compares the value in this Node to another value.
// Compare will use the Compare method of the Interface
// stored in this Node to perform the comparison, but is
// guaranteed to return a proper value - -1, 0, or 1 - even
// if the Interface's Compare function does not. Anything
// less than zero will be converted to -1, and anything
// greater than zero will be converted to 1.
func (n *Node) Compare(to interface{}) int {
	switch result := n.value.Compare(to); {
	case result < EQ:
		return LT
	case result > EQ:
		return GT
	default:
		return EQ
	}
}

// valud returns the value stored in this Node.
func (n *Node) Value() interface{} {
	return n.value.Value()
}

// Update updates the value stored in this Node.
func (n *Node) Update(with interface{}) {
	n.value.Update(with)
	return
}

// stringInterface is an implementation of Interface
// that works with strings.
type stringInterface struct {
	value      string
	ignoreCase bool
}

// String wraps a string in an Interface.
func String(value string, ignoreCase bool) stringInterface {
	return stringInterface{
		value:      value,
		ignoreCase: ignoreCase,
	}
}

// Compare compares this string to another value.
// Note that although this is a string comparison,
// the incoming value will be converted to a string
// before the comparison, so any type can be passed
// in.
func (s stringInterface) Compare(to interface{}) int {
	var value = s.value

	if s.ignoreCase {
		value = strings.ToLower(value)
	}

	var str = fmt.Sprintf("%v", to)

	if s.ignoreCase {
		str = strings.ToLower(str)
	}

	if value < str {
		return LT
	}

	if value > str {
		return GT
	}

	return EQ
}

// Value returns the string stored in this Interface.
func (s stringInterface) Value() interface{} {
	return s.value
}

// Update does nothing to the string stored in this
// Interface.
func (s stringInterface) Update(with interface{}) {
	return
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

	if Done == visitInOrder(with, subtreeRoot.left) {
		return Done
	}

	if Done == with(subtreeRoot.Value()) {
		return Done
	}

	return visitInOrder(with, subtreeRoot.right)
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

	if Done == visitInOrder(with, subtreeRoot.right) {
		return Done
	}

	if Done == with(subtreeRoot.Value()) {
		return Done
	}

	return visitInOrder(with, subtreeRoot.left)
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

// Contains searches the tree for a given item and returns true if it is
// found.
// Contains traverses the tree using a binary search pattern and compares
// the item to the Value of each node of the tree.
func (t *Tree) Contains(item Interface) bool {
	var cur, found = t.root, false

	for !found && cur != nil {
		switch result := item.Compare(cur.Value()); {
		case result < EQ:
			cur = cur.left
		case result > EQ:
			cur = cur.right
		default:
			found = true
		}
	}

	return found
}

// Get retrieves an item from the tree.
// If the item is not in the tree, it will return nil instead.
func (t *Tree) Get(item Interface) interface{} {
	var cur = t.root

	for cur != nil {
		switch result := item.Compare(cur.Value()); {
		case result < EQ:
			cur = cur.left
		case result > EQ:
			cur = cur.right
		default:
			return cur.Value()
		}
	}

	return nil
}
