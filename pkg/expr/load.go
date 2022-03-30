package expr

type Load struct {
	addr  Expr
	width uint8
}

func NewLoad(addr Expr, width uint8) Load {
	return Load{
		addr:  addr,
		width: width,
	}
}
