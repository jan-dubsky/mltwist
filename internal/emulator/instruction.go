package emulator

import (
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type Instruction interface {
	Effects() []expr.Effect
	Len() model.Addr
}
