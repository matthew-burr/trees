package binary

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	// . "matthew-burr/trees/binary"
)

func SampleTree() *Tree {
	return &Tree{
		root: &Node{
			left: &Node{
				Value: "L",
			},
			right: &Node{
				Value: "R",
			},
			Value: "V",
		},
	}
}

func PrintNodeTo(w io.Writer) VisitorFunc {
	return func(value interface{}) bool {
		fmt.Fprintln(w, value)
		return Continue
	}
}

func TestVisitInOrder(t *testing.T) {
	var buf = new(bytes.Buffer)
	SampleTree().VisitInOrder(PrintNodeTo(buf))

	want := "L\nV\nR\n"
	got := buf.String()

	assert.Equal(t, want, got)
}

func TestVisitInOrder_StopsWhenToldTo(t *testing.T) {
	var want, got = 1, 0

	SampleTree().VisitInOrder(func(interface{}) bool {
		got++
		return Done
	})

	assert.Equal(t, want, got)
}

func TestVisitInReverse(t *testing.T) {
	var buf = new(bytes.Buffer)
	SampleTree().VisitInReverse(PrintNodeTo(buf))

	want := "R\nV\nL\n"
	got := buf.String()

	assert.Equal(t, want, got)
}

func TestVisitInReverse_StopsWhenToldTo(t *testing.T) {
	var want, got = 1, 0

	SampleTree().VisitInOrder(func(interface{}) bool {
		got++
		return Done
	})

	assert.Equal(t, want, got)
}
