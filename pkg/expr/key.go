package expr

import "fmt"

const (
	// IPKey identifies instruction pointer registers. Writes to this
	// register will be interpreted as jumps.
	IPKey Key = "#r:ip"
)

// Key represents an arbitrary memory or register key used to identify memory
// address space or register.
//
// Value of key startinf with hash (#) are reserved to be defined by this
// package. Usage of keys starting with # which are not defined by this package
// might result in undefined behaviour. Empty key is invalid as well.
//
// For reserved keys starting with #, it's highly recommended to use constants
// defined by this package. Values of reserved strings might vary in between
// package minor versios, but constants are granted to be backward compatible.
type Key string

func NewKey(s string) Key { return Key(s) }

// validate checks that key value is valid (allowed) according to key
// definition.
func (k Key) validate() error {
	if k == "" {
		return fmt.Errorf("empty key is not allowed")
	}
	if k[0] != '#' {
		return nil
	}

	switch k {
	case IPKey:
		return nil
	default:
		return fmt.Errorf("unknown key starting with #: %s", k)
	}
}

// assertValid is a helper function which panics if key is not valid.
func (k Key) assertValid() {
	if err := k.validate(); err != nil {
		panic(fmt.Sprintf("invalid key %q: %s", k, err.Error()))
	}
}

// Reserved informs if key is reserved key or if it's standard (non-reserved)
// key.
//
// This function will panic for an empty key.
func (k Key) Reserved() bool { return k[0] == '#' }
