package expr

type Store struct {
	value Expr
	addr  Expr
	width uint8
}

func NewStore(v Expr, a Expr, width uint8) Store {
	return Store{
		value: v,
		addr:  a,
		width: width,
	}
}
