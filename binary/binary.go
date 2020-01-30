// Copyright 2020 Matthew Burr
//

// Package binary implements a binary tree.
package binary

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

// Value returns the value stored in this Node.
func (n *Node) Value() interface{} {
	return n.value.Value()
}

// Update updates the value stored in this Node.
func (n *Node) Update(with interface{}) {
	n.value.Update(with)
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
	// Done signals that tree traversal is done.
	Done = true
	// Continue signals that tree traversal is not done.
	Continue = false
)

// A VisitorFunc is an operation to perform on Nodes of a tree.
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
func (t *Tree) VisitInOrder(with VisitorFunc) *Tree {
	_ = visitInOrder(with, t.root)
	return t
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
func (t *Tree) VisitInReverse(with VisitorFunc) *Tree {
	_ = visitInReverse(with, t.root)
	return t
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

// findNodeAndParent is used internally to locate a node and its parent.
func (t *Tree) findNodeAndParent(item Interface) (found bool, node, parent *Node) {
	found, node, parent = false, t.root, nil
	for !found && node != nil {
		switch result := item.Compare(node.Value()); {
		case result < EQ:
			parent = node
			node = parent.left
		case result > EQ:
			parent = node
			node = parent.right
		default:
			found = true
		}
	}
	return
}

// Get retrieves an item from the tree.
// If the item is not in the tree, it will return nil instead.
func (t *Tree) Get(item Interface) interface{} {
	if found, cur, _ := t.findNodeAndParent(item); found {
		return cur.Value()
	}

	return nil
}

// Insert adds an item to the tree if it does
// not exist in the tree already.
// If there is already an item in the tree that
// matches the one you are adding, the tree
// will call Update on the existing item, passing
// in the value of the item you are trying to add.
// To facilitate easily inserting a chain of items,
// the method returns the Tree after having inserted
// an item.
func (t *Tree) Insert(item Interface) *Tree {
	var found, cur, parent = t.findNodeAndParent(item)

	if found {
		var before = cur.Value()
		cur.Update(item.Value())

		// If our update of the value changed the value, then
		// we remove and re-insert the node to insure it ends
		// up in the right place.
		if changed := cur.Compare(before) != EQ; changed {
			t.removeNode(cur, parent)
			// Because Node contains original Interface, it
			// can act as the interface.
			t.Insert(cur)
		}

		return t
	}

	// If we've reached this point, we didn't
	// find the item, so we insert a new Node.
	cur = &Node{value: item}
	if parent == nil { // Should only happen if there is an empty tree
		t.root = cur
	} else if item.Compare(parent.Value()) == LT {
		parent.left = cur
	} else {
		parent.right = cur
	}

	return t
}

func (t *Tree) removeNode(node, parent *Node) {
	// If the node we're deleting has two children, we need
	// to pick one of its subtrees to replace it. Below,
	// we're picking its right subtree
	if node.left != nil && node.right != nil {
		var succ = node.right
		parent = node
		// We're going to walk down the left branch of
		// its right subtree until we can't go any further.
		for succ.left != nil {
			parent = succ
			succ = succ.left
		}
		// Basically, instead of deleting the node we originally
		// found, we're going to replace its value with this
		// successor node's value, and then set up the successor
		// node for deletion.
		node.value = succ.value
		node = succ
	}

	// Now, we can just follow a pattern for deleting a node
	// with, at most, one child.
	var subtree = node.left

	if subtree == nil {
		subtree = node.right
	}

	if parent == nil {
		t.root = subtree
	} else if parent.left == node {
		parent.left = subtree
	} else {
		parent.right = subtree
	}
}

// Remove deletes the specified item from the tree (if it exists
// in the tree) and returns the tree.
func (t *Tree) Remove(item Interface) *Tree {
	// The algorithm used here is adapted from Chapter 12 of
	// "ADTs, Data Structures, and Problem Solving with C++",
	// by Larry Nyhoff (2nd Edition, 2005, Prentice Hall)

	// First, we need to find the node to delete (if it exists)
	// along with its parent.
	var found, node, parent = t.findNodeAndParent(item)

	if !found {
		return t
	}

	t.removeNode(node, parent)

	return t
}
