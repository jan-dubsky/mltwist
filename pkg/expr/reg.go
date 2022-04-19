package expr

var (
	_ Expr   = RegLoad{}
	_ Effect = RegStore{}
)

// RegLoad represents value load from a generic key-value registry storage.
type RegLoad struct {
	key Key
	w   Width
}

func NewRegLoad(key Key, w Width) RegLoad {
	key.assertValid()

	return RegLoad{
		key: key,
		w:   w,
	}
}

// Key returns key identifying the value loaded.
func (l RegLoad) Key() Key { return l.key }

// Equal checks RegLoad equality.
func (l1 RegLoad) Equal(l2 RegLoad) bool {
	return l1.Key() == l2.Key() &&
		l1.Width() == l2.Width()
}

func (l RegLoad) Width() Width { return l.w }
func (RegLoad) internalExpr()  {}

// RegStore represents value store to a generic key-value registry storage.
type RegStore struct {
	value Expr
	key   Key
	w     Width
}

func NewRegStore(value Expr, key Key, w Width) RegStore {
	key.assertValid()

	return RegStore{
		value: value,
		key:   key,
		w:     w,
	}
}

// Value returns the value to store.
func (s RegStore) Value() Expr { return s.value }

// Key returns key identifying the value Stored.
func (s RegStore) Key() Key { return s.key }

// Equal checks RegStore equality.
func (s1 RegStore) Equal(s2 MemStore) bool {
	return s1.Key() == s2.Key() &&
		s1.Width() == s2.Width()
}

func (s RegStore) Width() Width  { return s.w }
func (RegStore) internalEffect() {}
