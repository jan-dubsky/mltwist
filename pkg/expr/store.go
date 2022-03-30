package expr

var _ Expr = Store{}

type Store struct {
	value Expr
	addr  Expr
	w     Width
}

func NewStore(v Expr, a Expr, w Width) Store {
	return Store{
		value: v,
		addr:  a,
		w:     w,
	}
}

func (s Store) Width() Width { return s.w }
func (Store) internal()      {}
