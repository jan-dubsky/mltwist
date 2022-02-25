package basicblock_test

import (
	"decomp/internal/addr"
	"decomp/internal/basicblock"
	"decomp/internal/instruction"
	"testing"

	"github.com/stretchr/testify/require"
)

const instrLen = 4

type detail string

func (d detail) String() string { return string(d) }

func ins(
	addr addr.Address,
	tp instruction.Type,
	desc string,
	jmp ...addr.Address,
) instruction.Instruction {
	return instruction.Instruction{
		Address:     addr,
		Type:        tp,
		ByteLen:     instrLen,
		Details:     detail(desc),
		JumpTargets: jmp,
	}
}

func TestParse(t *testing.T) {
	r := require.New(t)

	instrs := []instruction.Instruction{
		ins(50, instruction.TypeAritm, "1"),
		ins(54, instruction.TypeAritm, "2"),
		ins(59, instruction.TypeJump, "3"),
		ins(63, instruction.TypeCJump, "4"),
		ins(68, instruction.TypeAritm, "5"),
		ins(72, instruction.TypeJumpDyn, "6"),
		ins(76, instruction.TypeAritm, "7"),
		ins(80, instruction.TypeAritm, "8"),
		ins(84, instruction.TypeAritm, "9"),
		ins(88, instruction.TypeAritm, "10"),
		ins(92, instruction.TypeCJump, "11", 84),
		ins(96, instruction.TypeAritm, "12"),
		ins(100, instruction.TypeCJump, "13", 96),
	}

	blocks, err := basicblock.Parse(instrs)
	r.NoError(err)
	r.Len(blocks, 7)

	r.Equal([]instruction.Instruction{
		ins(50, instruction.TypeAritm, "1"),
		ins(54, instruction.TypeAritm, "2"),
	}, blocks[0].Seq)
	r.Equal([]instruction.Instruction{
		ins(59, instruction.TypeJump, "3"),
	}, blocks[1].Seq)
	r.Equal([]instruction.Instruction{
		ins(63, instruction.TypeCJump, "4"),
	}, blocks[2].Seq)
	r.Equal([]instruction.Instruction{
		ins(68, instruction.TypeAritm, "5"),
		ins(72, instruction.TypeJumpDyn, "6"),
	}, blocks[3].Seq)
	r.Equal([]instruction.Instruction{
		ins(76, instruction.TypeAritm, "7"),
		ins(80, instruction.TypeAritm, "8"),
	}, blocks[4].Seq)
	r.Equal([]instruction.Instruction{
		ins(84, instruction.TypeAritm, "9"),
		ins(88, instruction.TypeAritm, "10"),
		ins(92, instruction.TypeCJump, "11", 84),
	}, blocks[5].Seq)
	r.Equal([]instruction.Instruction{
		ins(96, instruction.TypeAritm, "12"),
		ins(100, instruction.TypeCJump, "13", 96),
	}, blocks[6].Seq)
}
