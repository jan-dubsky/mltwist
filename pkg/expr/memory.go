package expr

var (
	_ Expr   = MemLoad{}
	_ Effect = MemStore{}
)

// MemLoad implements memory load from a memory address space identified by key,
// from address adds. The width of load is w.
//
// Address of the memory store is not effected by w, but it keeps its original
// width.
type MemLoad struct {
	key  Key
	addr Expr
	w    Width
}

// NewMemLoad returns a new memory load reading w bytes from address addr in
// memory address space identified by key.
func NewMemLoad(key Key, addr Expr, w Width) MemLoad {
	key.assertValid(keyScopeMem, keyPermissionRead)

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

// Width returns width of l.
func (l MemLoad) Width() Width { return l.w }
func (MemLoad) internalExpr()  {}

// MemStore stores value of expression e to a memory address space identified by
// key to address adds. The width of store is w.
//
// Address of the memory store is not effected by w, it keeps its original
// width.
type MemStore struct {
	value Expr
	key   Key
	addr  Expr
	w     Width
}

// NewMemSotre returns new MemStore storing value to address addr in memory
// addres space identified by key. The width of write is w.
func NewMemStore(value Expr, key Key, addr Expr, w Width) MemStore {
	key.assertValid(keyScopeMem, keyPermissionWrite)

	return MemStore{
		value: value,
		key:   key,
		addr:  addr,
		w:     w,
	}
}

// Value returns the expression which value is stored to memory.
func (s MemStore) Value() Expr { return s.value }

// Key returns key identifying memory space.
func (s MemStore) Key() Key { return s.key }

// Addr returns address of memory store.
func (s MemStore) Addr() Expr { return s.addr }

// Width returns width of s.
func (s MemStore) Width() Width  { return s.w }
func (MemStore) internalEffect() {}
