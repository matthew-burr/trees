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

// A CompareFunc is a function that can be used to
// compare two values (presumably, of the same type).
type CompareFunc func(this, to interface{}) int

// An UpdateFunc is a function that can be used to
// update a value.
type UpdateFunc func(this, with interface{}) interface{}

// A InterfaceImpl is a generic implementation of
// interface that can be used to wrap values of
// any type.
type InterfaceImpl struct {
	value interface{}
	comp  CompareFunc
	upd   UpdateFunc
}

// Value returns the value of this Interface.
func (c *InterfaceImpl) Value() interface{} {
	return c.value
}

// Compare compare's this Interface's value to another value.
func (c *InterfaceImpl) Compare(to interface{}) int {
	switch result := c.comp(c.value, to); {
	case result < EQ:
		return LT
	case result > EQ:
		return GT
	default:
		return EQ
	}
}

// Update updates the value in the interface.
func (c *InterfaceImpl) Update(with interface{}) {
	if c.upd == nil {
		return
	}

	c.value = c.upd(c.value, with)
	return
}

// Generic returns an implementation of Interface.
// This implementation's behavior is defined by parameters
// passed to Generic at the time of its creation.
// The value parameter is the value stored in the interface.
// Comparer (required) is a function that can be used to compare the interface's
// value to another value.
// Updater (optional) is a function that can be used to update the value stored
// in the interface. If not supplied, then a call to Interface.Update is a no-op.
func Generic(value interface{}, comparer CompareFunc, updater UpdateFunc) *InterfaceImpl {
	if comparer == nil {
		panic("you must provide a comparer function")
	}

	return &InterfaceImpl{
		value: value,
		comp:  comparer,
		upd:   updater,
	}
}

// String wraps a string in an Interface.
func String(value string, ignoreCase bool) *InterfaceImpl {
	return Generic(
		value,
		func(this, to interface{}) int {
			value = this.(string)
			if ignoreCase {
				value = strings.ToLower(value)
			}

			var str = fmt.Sprintf("%v", to)

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
		},
		nil,
	)
}

// Int wraps an integer in an Interface.
func Int(value int) *InterfaceImpl {
	return Generic(
		value,
		func(this, to interface{}) int {
			return this.(int) - to.(int)
		},
		nil,
	)
}
