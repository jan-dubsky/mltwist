package expr

import (
	"fmt"
)

const (
	// IPKey identifies instruction pointer registers. Writes to this
	// register will be interpreted as jumps.
	//
	// Reads of this register might result in undefined behaviour. This is
	// true because internally the compiler is allowed to modify the
	// position of an instruction and change the meaning if instruction
	// pointer. Please prefer constant addressing based on the address of
	// the instruction rather than reading this register.
	//
	// As the effect signature doesn't allow to express "conditional write"
	// operation, a conditional jump has to unconditionally write this
	// register key. This is not an issue as jump to following instruction
	// address is equivalent to no jump at all, register write of a
	// conditional expression can be used to represent conditional jump.
	//
	// Please note that this design of jump expressions is not sufficient to
	// describe jump delay in any reasonable way. To support architectures
	// where jump don't take an immediate effect some sort of redesign of
	// this package or at least of this way of expressing jumps will be
	// necessary.
	IPKey Key = "#r:ip"
)

// Key represents an arbitrary memory or register key used to identify memory
// address space or register.
//
// Value of key starting with hash (#) are reserved to be defined by this
// package. Usage of keys starting with # which are not defined by this package
// might result in undefined behaviour. Empty key is invalid as well.
//
// For reserved keys starting with #, it's highly recommended to use constants
// defined by this package. Values of reserved strings might vary in between
// package minor versions, but constants are granted to be backward compatible.
type Key string

// NewKey creates a new key out of string s.
func NewKey(s string) Key { return Key(s) }

// allowedReserved checks if k is in a list of reserved keys defined by this
// package.
func (k Key) allowedReserved() bool {
	switch k {
	case IPKey:
	default:
		return false
	}
	return true
}

// validate checks that key value is valid (allowed) according to key
// definition in a scope scope.
func (k Key) validate(scope keyScope) error {
	if k == "" {
		return fmt.Errorf("empty key is not allowed")
	}
	if k[0] != '#' {
		return nil
	}
	if !k.allowedReserved() {
		return fmt.Errorf("unknown key starting with #: %s", k)
	}

	err := k.validateScope(scope)
	if err != nil {
		return fmt.Errorf("invalid scope: %w", err)
	}

	return nil
}

// validateScope checks that key k is allowed to be used in scope scope.
func (k Key) validateScope(scope keyScope) error {
	if l := len(k); l < 3 {
		return fmt.Errorf("key is to short (%d) for a reserved key: %s", l, k)
	}
	if k[2] != ':' {
		return fmt.Errorf(
			"reserved keys must start with single letter and colon: %s", k)
	}

	keyScope := keyScope(k[1])
	if err := keyScope.validate(); err != nil {
		return err
	}
	if !keyScope.allows(scope) {
		return fmt.Errorf("key scope (%c) is not allowed in scope %c: %s",
			keyScope, scope, k)
	}

	return nil
}

// assertValid is a helper function which panics if key is not valid.
func (k Key) assertValid(scope keyScope) {
	if err := k.validate(scope); err != nil {
		panic(fmt.Sprintf("invalid key %q: %s", k, err.Error()))
	}
}

// Reserved informs if key is reserved key or if it's standard (non-reserved)
// key.
//
// This function will panic for an empty key.
func (k Key) Reserved() bool { return k[0] == '#' }

// keyScope describes in which contexts is the reserved key valid.
//
// Reserved keys have a special meaning defined by this package. For this reason
// usage of reserved keys doesn't have to make sense in all contexts. In other
// words some keys might have well-defined meaning only in register operations
// and some only in memory operations.
//
// To differentiate those keys and contexts in which they are valid, we
// introduce key scopes. Key scope is a single byte information at the beginning
// of a reserved key (right after the starting hash), which is delimited from
// the rest of the key by colon. This single byte can be either k or m to
// represent either register or memory respectively. The scope b then indicates
// conjunction of r and m.
//
// As values of reserved keys are allowed to change in between package versions,
// the notation of scopes can change as well.
type keyScope byte

const (
	// keyScopeReg specifies that reserved key is valid only for register
	// operations.
	keyScopeReg = 'r'
	// keyScopeMem specifies that reserved key is valid only for memory
	// operations.
	keyScopeMem = 'm'
	// keyScopeRegMem specifies that reserved key is valid for both register
	// and memory operations.
	keyScopeRegMem = 'b'
)

// allows checks if scope s contains scope scope. In other words, this check
// asserts that scope is subset of (or equal to) s using set-related terms.
func (s keyScope) allows(scope keyScope) bool {
	if s == keyScopeRegMem {
		return true
	}

	return s == scope
}

// validate check is value of s is one of defined values and returns an error if
// it's not.
func (s keyScope) validate() error {
	switch s {
	case keyScopeReg, keyScopeMem, keyScopeRegMem:
		return nil
	default:
		return fmt.Errorf("unknown scope: %c", s)
	}
}
