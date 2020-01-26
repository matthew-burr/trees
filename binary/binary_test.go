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
				value: String("L", false),
			},
			right: &Node{
				value: String("R", false),
			},
			value: String("M", false),
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

	want := "L\nM\nR\n"
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

	want := "R\nM\nL\n"
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

func TestContains(t *testing.T) {
	type args struct {
		find       string
		ignoreCase bool
	}

	tt := []struct {
		name string
		args
		want bool
	}{
		{"Case Sensitive - Exists", args{"L", false}, true},
		{"Case Sensitive - Wrong case", args{"l", false}, false},
		{"Case Sensitive - Not exists", args{"Foo", false}, false},
		{"Ignore Case - Same Case", args{"R", true}, true},
		{"Ignore Case - Diff Case", args{"m", true}, true},
		{"Ignore Case - Not Exists", args{"Bar", true}, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var want = tc.want
			var got = SampleTree().Contains(String(tc.args.find, tc.args.ignoreCase))

			assert.Equal(t, want, got)
		})
	}
}

/*
func TestInt(t *testing.T) {
	type args struct {
		left, right int
	}

	tt := []struct {
		name string
		args
		want int
	}{
		{"LT-Negs", args{-2, -1}, LT},
		{"LT-Neg-0", args{-1, 0}, LT},
		{"LT-0-Pos", args{0, 1}, LT},
		{"LT-Poss", args{1, 2}, LT},
		{"GT-Negs", args{-1, -2}, GT},
		{"GT-Neg-0", args{0, -1}, GT},
		{"GT-0-Pos", args{1, 0}, GT},
		{"GT-Poss", args{2, 1}, GT},
		{"EQ-Negs", args{-1, -1}, EQ},
		{"EQ-0-0", args{0, 0}, EQ},
		{"EQ-Poss", args{1, 1}, EQ},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var want = tc.want
			var got = Int(tc.left).Compare(tc.right)
			assert.Equal(t, want, got)
		})
	}
}
*/

func TestGet(t *testing.T) {
	tt := []struct {
		name string
		arg  string
		want interface{}
	}{
		{"Exists", "R", "R"},
		{"Not Exists", "Foo", nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var want = tc.want
			var got = SampleTree().Get(String(tc.arg, true))
			assert.Equal(t, want, got)
		})
	}
}
