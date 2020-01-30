package binary

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestInsert_AddsNewItemsToRightLocation(t *testing.T) {
	var buf = new(bytes.Buffer)
	new(Tree).
		Insert(String("M", false)).
		Insert(String("L", false)).
		Insert(String("R", false)).
		VisitInOrder(PrintNodeTo(buf))

	var want, got = "L\nM\nR\n", buf.String()
	assert.Equal(t, want, got)
}

func TestGenericInterface_ReturnsValue(t *testing.T) {
	want := 1
	got := Generic(1, func(interface{}, interface{}) int { return 0 }, nil).Value()
	assert.Equal(t, want, got)
}

func TestGenericInterface_Compares(t *testing.T) {
	want := -1
	got := Generic(1,
		func(interface{}, interface{}) int {
			return -1
		}, nil).Compare(2)
	assert.Equal(t, want, got)
}

func TestGeneric_PanicsIfNilComparer(t *testing.T) {
	assert.Panics(t, func() {
		_ = Generic(0, nil, nil)
	})
}

func TestGeneric_DoesNotPanicIfUpdaterIsNil(t *testing.T) {
	assert.NotPanics(t, func() {
		_ = Generic(0, func(interface{}, interface{}) int { return 0 }, nil)
	})
}

func TestGenericInterface_Updates(t *testing.T) {
	var want, got = 1, 0
	Generic(0,
		func(value, to interface{}) int {
			return 0
		},
		func(this, with interface{}) interface{} {
			got = 1
			return this
		}).Update(0)

	assert.Equal(t, want, got)
}

func TestGenericInterface_NoPanicIfUpdateWithNilUpdater(t *testing.T) {
	assert.NotPanics(t, func() {
		Generic(0, func(interface{}, interface{}) int { return 0 }, nil).Update(0)
	})
}

func TestGenericInterface_CompareReturnsLTGT(t *testing.T) {
	tt := []struct {
		name string
		arg  int
		want int
	}{
		{"LT", 2, LT},
		{"GT", -2, GT},
		{"EQ", 0, EQ},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var want = tc.want
			var got = Generic(0, func(value, to interface{}) int {
				return value.(int) - to.(int)
			}, nil).Compare(tc.arg)

			assert.Equal(t, want, got)
		})
	}
}

func TestInt_Value(t *testing.T) {
	var want = 1
	var got = Int(want).Value()
	assert.Equal(t, want, got)
}

func TestInt_Compare(t *testing.T) {
	tt := []struct {
		name string
		arg  int
		want int
	}{
		{"LT", 1, LT},
		{"GT", -1, GT},
		{"EQ", 0, EQ},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			want := tc.want
			got := Int(0).Compare(tc.arg)
			assert.Equal(t, want, got)
		})
	}
}

func TestTreeInsert_UpdatesExistingNode(t *testing.T) {
	var want, got = 1, 0
	new(Tree).
		Insert(
			Generic(
				0,
				func(this, to interface{}) int {
					return EQ
				},
				func(this, with interface{}) interface{} {
					got++
					return this
				},
			)).
		Insert(Int(0))

	assert.Equal(t, want, got)
}

func TestTreeRemove(t *testing.T) {
	tt := []struct {
		name string
		arg  string
		want string
	}{
		{"Left", "L", "M\nQ\nR\nT\n"},
		{"Right", "R", "L\nM\nQ\nT\n"},
		{"Root", "M", "L\nQ\nR\nT\n"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var buf = new(bytes.Buffer)
			new(Tree).
				Insert(String("L", false)).
				Insert(String("M", false)).
				Insert(String("R", false)).
				Insert(String("T", false)).
				Insert(String("Q", false)).
				Remove(String(tc.arg, false)).
				VisitInOrder(PrintNodeTo(buf))

			var want, got = tc.want, buf.String()
			assert.Equal(t, want, got)

		})
	}
}

func TestTreeInsert_LeavesNodesInProperOrderAfterUpdateChangesValues(t *testing.T) {
	var tree = new(Tree).
		Insert(Int(1)).
		Insert(Int(2)).
		Insert(Generic(
			3,
			func(this, to interface{}) int {
				return this.(int) - to.(int)
			},
			func(this, with interface{}) interface{} {
				return this.(int) * 0
			},
		))

	var buf = new(bytes.Buffer)
	tree.Insert(Int(3)).VisitInOrder(PrintNodeTo(buf))

	var want, got = "0\n1\n2\n", buf.String()

	assert.Equal(t, want, got)
}

func TestTreeInsert_ChainOfUpdates(t *testing.T) {
	var comparer = func(this, to interface{}) int {
		return this.(int) - to.(int)
	}

	var updater = func(this, with interface{}) interface{} {
		return this.(int) + 1
	}

	var item = func(i int) *InterfaceImpl {
		return Generic(i, comparer, updater)
	}

	var buf = new(bytes.Buffer)
	var tree = new(Tree).
		Insert(item(0)).
		Insert(item(1)).
		Insert(item(2)).
		VisitInOrder(PrintNodeTo(buf))

	var want, got = "0\n1\n2\n", buf.String()
	require.Equal(t, want, got)

	buf = new(bytes.Buffer)
	tree = tree.Insert(item(0)).VisitInOrder(PrintNodeTo(buf))

	want, got = "3\n", buf.String()
	assert.Equal(t, want, got)
}
