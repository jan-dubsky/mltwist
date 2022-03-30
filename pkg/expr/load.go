package expr

var _ Expr = Load{}

type Load struct {
	addr Expr
	w    Width
}

func NewLoad(addr Expr, w Width) Load {
	return Load{
		addr: addr,
		w:    w,
	}
}

func (l Load) Width() Width { return l.w }
func (Load) internal()      {}
