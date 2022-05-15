package deps

import (
	"fmt"
	"mltwist/pkg/model"
)

// validateArrayIndex validates that variable called name with value val is
// valid index into an array with length l.
func validateArrayIndex(name string, val int, l int) error {
	if val < 0 {
		return fmt.Errorf("negative value of %q is not allowed: %d", name, val)
	}
	if val >= l {
		return fmt.Errorf("value of %q is above limit: %d >= %d", name, val, l)
	}

	return nil
}

func checkFromToIndex(from int, to int, l int) error {
	if err := validateArrayIndex("from", from, l); err != nil {
		return err
	} else if err := validateArrayIndex("to", to, l); err != nil {
		return err
	}

	return nil
}

// movable is any type which knows its position in an array and it's array. and
// which is also able to update those values.
type movable interface {
	setIndex(i int)

	Begin() model.Addr
	End() model.Addr
	setAddr(a model.Addr)
}

func move[T movable](arr []T, from int, to int) {
	if from == to {
		return
	}

	if from < to {
		moveFwd(arr, from, to)
	} else {
		moveBack(arr, from, to)
	}
}

func moveFwd[T movable](arr []T, from int, to int) {
	f := arr[from]
	a := f.Begin()

	for i := from; i < to; i++ {
		arr[i] = arr[i+1]
		arr[i].setIndex(i)

		arr[i].setAddr(a)
		a = arr[i].End()
	}

	f.setIndex(to)
	f.setAddr(a)
	arr[to] = f
}

func moveBack[T movable](arr []T, from int, to int) {
	f := arr[from]
	a := arr[to].Begin()

	for i := from; i > to; i-- {
		arr[i] = arr[i-1]
		arr[i].setIndex(i)
	}

	f.setIndex(to)
	arr[to] = f

	for i := to; i <= from; i++ {
		arr[i].setAddr(a)
		a = arr[i].End()
	}
}
