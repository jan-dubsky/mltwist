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
	IPKey Key = "#r:w:ip"
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
		return true
	default:
		return false
	}
}

// validate checks that key value is valid (allowed) according to key
// definition in a scope s and with permissions p.
func (k Key) validate(s keyScope, p keyPermission) error {
	if k == "" {
		return fmt.Errorf("empty key is not allowed")
	}
	if k[0] != '#' {
		return nil
	}

	if err := k.validateProps(s, p); err != nil {
		return fmt.Errorf("invalid reserved key properties: %w", err)
	}

	if !k.allowedReserved() {
		return fmt.Errorf("unknown reserved key: %s", k)
	}

	return nil
}

func (k Key) validateProps(s keyScope, p keyPermission) error {
	// Starting hash + scope + colon + permission + colon + at least one
	// character of register name => 6 is the minimum.
	if l := len(k); l < 6 {
		return fmt.Errorf("key is to short (%d) for a reserved key: %s", l, k)
	}
	if k[2] != ':' || k[4] != ':' {
		return fmt.Errorf("reserved key must have colons as index 2 and 4: %s", k)
	}

	keyScope, keyPermission := keyScope(k[1]), keyPermission(k[3])
	if err := validateKeyProp(keyScope, s); err != nil {
		return fmt.Errorf("reserved key scope not allowed: %w", err)
	}
	if err := validateKeyProp(keyPermission, p); err != nil {
		return fmt.Errorf("reserved key permission not allowed: %w", err)
	}

	return nil
}

// keyProp describes any sort of information encoded in reserved register name.
type keyProp[T any] interface {
	// Present just to make the fmt.Errorf lister happy about it being char.
	~byte
	allows(T) bool
	validate() error
}

// validateKeyProp validates whether both key and want are valid and then
// asserts of wanted property w is allowed under a property k of the string key.
func validateKeyProp[T keyProp[T]](k T, w T) error {
	if err := k.validate(); err != nil {
		return fmt.Errorf("invalid \"k\" property provides: %w", err)
	}
	if err := w.validate(); err != nil {
		return fmt.Errorf("invalid \"w\" property provides: %w", err)
	}
	if !k.allows(w) {
		return fmt.Errorf(
			"key property of value '%c' doesn't allow wanted property: %c",
			k, w)
	}

	return nil
}

// assertValid is a helper function which panics if key is not valid.
func (k Key) assertValid(s keyScope, p keyPermission) {
	if err := k.validate(s, p); err != nil {
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
// represent either register or memory respectively. The scope b (both) then
// indicates conjunction of r and m.
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

// validate check is value of s is one of defined values of scope and returns an
// error if it's not.
func (s keyScope) validate() error {
	switch s {
	case keyScopeReg, keyScopeMem, keyScopeRegMem:
		return nil
	default:
		return fmt.Errorf("unknown scope: %c", s)
	}
}

// keyPermission devotes which operations are valid for the key.
//
// Reserved keys have a special meaning defined by this package. For this reason
// usage of reserved keys don;t have to always support both reads and written.
// In other words some keys might have well-defined meaning only while read and
// some only when written.
//
// To differentiate which operations are allowed for a given key, we introduce
// key permissions. Key permission is a single byte information at the beginning
// of a reserved key (right after the starting hash followed by scope and
// colon), which is delimited from the rest of the key by (another) colon. This
// single byte can be either r or w to represent either read or write
// respectively. The permission b (both) then indicates conjunction of r and w.
//
// As values of reserved keys are allowed to change in between package versions,
// the notation of permissions can change as well.
type keyPermission byte

const (
	// keyPermissionRead specifies that the value can be only read.
	keyPermissionRead = 'r'
	// keyPermissionWrite specifies that the value can be only written.
	keyPermissionWrite = 'w'
	// keyPermissionReadWrite specifies that the value can be both read and
	// written.
	keyPermissionReadWrite = 'b'
)

// allows checks if scope p allows operation described by permission perm. In
// other words, this check asserts that perm is subset of (or equal to) p using
// set-related terms.
func (p keyPermission) allows(perm keyPermission) bool {
	if p == keyPermissionReadWrite {
		return true
	}

	return p == perm
}

// validate check is value of p is one of defined values of permission and
// returns an error if it's not.
func (p keyPermission) validate() error {
	switch p {
	case keyPermissionRead, keyPermissionWrite, keyPermissionReadWrite:
		return nil
	default:
		return fmt.Errorf("unknown permission string: %c", p)
	}
}
