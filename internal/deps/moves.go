package deps

import "fmt"

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

// indexAware is any type which knows its position in an array and is able to
// update it.
type indexAware interface {
	setIndex(i int)
}

func move[T indexAware](arr []T, from int, to int) {
	if from == to {
		return
	}

	f := arr[from]
	if from < to {
		moveFwd(arr, from, to)
	} else {
		moveBack(arr, from, to)
	}
	arr[to] = f
	f.setIndex(to)
}

func moveFwd[T indexAware](arr []T, from int, to int) {
	for i := from; i < to; i++ {
		arr[i] = arr[i+1]
		arr[i].setIndex(i)
	}
}

func moveBack[T indexAware](arr []T, from int, to int) {
	for i := from; i > to; i-- {
		arr[i] = arr[i-1]
		arr[i].setIndex(i)
	}
}
