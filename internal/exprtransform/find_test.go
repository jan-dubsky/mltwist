package exprtransform_test

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	r := require.New(t)

	binary1 := expr.NewBinary(expr.Rsh,
		expr.One,
		expr.NewRegLoad("r2", expr.Width64),
		expr.Width64,
	)
	binary2 := expr.NewBinary(expr.Div,
		expr.NewRegLoad("r2", expr.Width16),
		expr.NewRegLoad("r1", expr.Width64),
		expr.Width32,
	)

	load1 := expr.NewMemLoad("mem1",
		expr.NewRegLoad("r1", expr.Width32),
		expr.Width64,
	)
	load2 := expr.NewMemLoad("mem2",
		expr.ConstFromUint[uint16](0xffbc),
		expr.Width32,
	)

	cond1 := expr.NewCond(expr.Ltu,
		expr.NewRegLoad("r1", expr.Width64),
		expr.ConstFromUint[uint16](6789),
		binary1,
		binary2,
		expr.Width16,
	)
	cond2 := expr.NewCond(expr.Eq,
		expr.Zero,
		expr.NewRegLoad("r3", expr.Width16),
		load1,
		expr.NewRegLoad("r3", expr.Width32),
		expr.Width32,
	)

	binary3 := expr.NewBinary(expr.Div, load2, cond2, expr.Width64)
	load3 := expr.NewMemLoad("mem01", binary3, expr.Width64)

	binary4 := expr.NewBinary(expr.Sub,
		expr.NewRegLoad("r2", expr.Width16),
		load3,
		expr.Width8,
	)

	e := expr.NewBinary(expr.Mul, cond1, binary4, expr.Width32)

	consts := []expr.Const{
		expr.Zero,
		expr.One,
		expr.ConstFromUint[uint16](6789),
		expr.ConstFromUint[uint16](0xffbc),
	}
	r.ElementsMatch(consts, exprtransform.FindAll[expr.Const](e))

	regs := []expr.RegLoad{
		expr.NewRegLoad("r2", expr.Width64),
		expr.NewRegLoad("r2", expr.Width16),
		expr.NewRegLoad("r1", expr.Width64),
		expr.NewRegLoad("r1", expr.Width64),
		expr.NewRegLoad("r3", expr.Width16),
		expr.NewRegLoad("r1", expr.Width32),
		expr.NewRegLoad("r3", expr.Width32),
		expr.NewRegLoad("r2", expr.Width16),
	}
	r.ElementsMatch(regs, exprtransform.FindAll[expr.RegLoad](e))

	binary := []expr.Binary{
		e,
		binary1,
		binary2,
		binary3,
		binary4,
	}
	r.ElementsMatch(binary, exprtransform.FindAll[expr.Binary](e))

	conds := []expr.Cond{
		cond1,
		cond2,
	}
	r.ElementsMatch(conds, exprtransform.FindAll[expr.Cond](e))

	loads := []expr.MemLoad{
		load1,
		load2,
		load3,
	}
	r.ElementsMatch(loads, exprtransform.FindAll[expr.MemLoad](e))
}
