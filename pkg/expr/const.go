package expr

import "decomp/pkg/expr/internal/value"

var _ Expr = Const{}

var (
	Zero = NewConst([]byte{0})
	One  = NewConst([]byte{1})
)

type Const struct {
	v value.Value
}

func NewConst(b []byte) Const {
	bCopy := make([]byte, len(b))
	copy(bCopy, b)

	return Const{
		v: value.New(bCopy),
	}
}

func (c Const) value() (value.Value, bool) { return value.Value{}, true }
func (c Const) Value() []byte              { return c.v.Bytes() }
