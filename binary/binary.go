// Copyright 2020 Matthew Burr
//
// Package binary implements a binary tree.
package binary

import (
	"fmt"
	"strings"
)

// A Node is an element in a binary tree.
type Node struct {
	// left and right are pointers to the children of the Node.
	left, right *Node

	// Value contains the Value for this node.
	Value interface{}
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

	if Done == with(subtreeRoot.Value) {
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

	if Done == with(subtreeRoot.Value) {
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

const (
	LT = -1
	EQ = 0
	GT = 1
)

// A Comparable is a type that supports the Compare function.
type Comparable interface {
	// Compare returns an int that indicates how the Comparable
	// compares in value to the input. Compare should return
	// a negative value if the Comparable is less than the input;
	// zero if the Comparable equals the input; and a positive value
	// if the Comparable is greater than the input.
	// Use the LT, EQ, and GT constants to simplify this.
	Compare(interface{}) int
}

// A CompareFunc is a function that supports the Comparable interface.
type CompareFunc func(interface{}) int

// Compare implements the Comparable interface for a CompareFunc.
func (c CompareFunc) Compare(item interface{}) int {
	return c(item)
}

// Some common comparison strategies
// String performs a string comparison.
// Optionally, the comparison can be case-insensitive by
// setting the second parameter to true.
func String(value string, ignoreCase bool) CompareFunc {
	if ignoreCase {
		value = strings.ToLower(value)
	}

	return func(item interface{}) int {
		var str = fmt.Sprintf("%v", item)
		if ignoreCase {
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
}

// Int performs an integer comparison between
// its value and the input items.
func Int(value int) CompareFunc {
	return func(item interface{}) int {
		switch i := item.(int); {
		case value < i:
			return LT
		case value > i:
			return GT
		default:
			return EQ
		}
	}
}

// Float64 performs a float64 comparison between
// its value and the input items.
func Float64(value float64) CompareFunc {
	return func(item interface{}) int {
		switch f := item.(float64); {
		case value < f:
			return LT
		case value > f:
			return GT
		default:
			return EQ
		}
	}
}

// Contains searches the tree for a given item and returns true if it is
// found.
// Contains traverses the tree using a binary search pattern and compares
// the item to the Value of each node of the tree.
func (t *Tree) Contains(item Comparable) bool {
	var cur, found = t.root, false

	for !found && cur != nil {
		switch result := item.Compare(cur.Value); {
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
func (t *Tree) Get(item Comparable) interface{} {
	var cur = t.root

	for cur != nil {
		switch result := item.Compare(cur.Value); {
		case result < EQ:
			cur = cur.left
		case result > EQ:
			cur = cur.right
		default:
			return cur.Value
		}
	}

	return nil
}
