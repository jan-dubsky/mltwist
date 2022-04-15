package expr

var (
	_ Expr   = MemLoad{}
	_ Effect = MemStore{}
)

// MemLoad implements memory load from a memory address space identified by key,
// from address adds. The width of load is w.
type MemLoad struct {
	key  Key
	addr Expr
	w    Width
}

func NewMemLoad(key Key, addr Expr, w Width) MemLoad {
	return MemLoad{
		key:  key,
		addr: addr,
		w:    w,
	}
}

// Key returns key identifying memory space.
func (l MemLoad) Key() Key { return l.key }

// Addr returns address of memory load.
func (l MemLoad) Addr() Expr { return l.addr }

func (l MemLoad) Width() Width { return l.w }
func (MemLoad) internalExpr()  {}

// MemStore stores value of expression e to a memory address space identified by
// key to address adds. The width of store is w.
type MemStore struct {
	value Expr
	key   Key
	addr  Expr
	w     Width
}

func NewMemStore(value Expr, key Key, addr Expr, w Width) MemStore {
	return MemStore{
		value: value,
		key:   key,
		addr:  addr,
		w:     w,
	}
}

// Value returns the value to store.
func (s MemStore) Value() Expr { return s.value }

// Key returns key identifying memory space.
func (s MemStore) Key() Key { return s.key }

// Addr returns address of memory store.
func (s MemStore) Addr() Expr { return s.addr }

func (s MemStore) Width() Width  { return s.w }
func (MemStore) internalEffect() {}
