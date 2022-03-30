package expr

var _ Expr = SignExtend{}

type SignExtend struct {
	e       Expr
	signBit Expr
	w       Width
}

func NewSignExtend(e Expr, signBit Expr, w Width) SignExtend {
	return SignExtend{
		e:       e,
		signBit: signBit,
		w:       w,
	}
}

func (e SignExtend) Width() Width { return e.w }
func (SignExtend) internal()      {}
