package expr

var _ Expr = Ref{}

type Ref struct {
	s string
	w Width
}

func NewRef(s string, w Width) Ref {
	return Ref{
		s: s,
		w: w,
	}
}

func (r Ref) Width() Width { return r.w }
func (Ref) internal()      {}
