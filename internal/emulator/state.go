package emulator

import (
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type StateProvider interface {
	Register(key expr.Key) expr.Const
	Memory(key expr.Key, addr model.Addr, w expr.Width) expr.Const
}

var _ StateProvider = ZeroProvider{}

type ZeroProvider struct{}

func (ZeroProvider) Register(_ expr.Key) expr.Const { return expr.Zero }
func (ZeroProvider) Memory(_ expr.Key, _ model.Addr, w expr.Width) expr.Const {
	return expr.NewConstUint[uint8](0, w)
}
