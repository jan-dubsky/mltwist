package repr

import "decomp/pkg/model"

type Instruction struct {
	model.Instruction

	Address model.Address
	Bytes   []byte
}
