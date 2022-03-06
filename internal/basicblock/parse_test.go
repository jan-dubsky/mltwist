package basicblock_test

import (
	"decomp/internal/basicblock"
	"decomp/internal/repr"
	"decomp/pkg/model"
	"testing"

	"github.com/stretchr/testify/require"
)

const instrLen = 4

type detail string

func (d detail) String() string { return string(d) }

func ins(
	addr model.Address,
	tp model.Type,
	desc string,
	jmp ...model.Address,
) repr.Instruction {
	return repr.Instruction{
		Address: addr,
		Instruction: model.Instruction{
			Type:        tp,
			ByteLen:     instrLen,
			Details:     detail(desc),
			JumpTargets: jmp,
		},
	}
}

func TestParse(t *testing.T) {
	r := require.New(t)

	instrs := []repr.Instruction{
		ins(50, model.TypeAritm, "1"),
		ins(54, model.TypeAritm, "2"),
		ins(59, model.TypeJump, "3"),
		ins(63, model.TypeCJump, "4"),
		ins(68, model.TypeAritm, "5"),
		ins(72, model.TypeJumpDyn, "6"),
		ins(76, model.TypeAritm, "7"),
		ins(80, model.TypeAritm, "8"),
		ins(84, model.TypeAritm, "9"),
		ins(88, model.TypeAritm, "10"),
		ins(92, model.TypeCJump, "11", 84),
		ins(96, model.TypeAritm, "12"),
		ins(100, model.TypeCJump, "13", 96),
	}

	blocks, err := basicblock.Parse(instrs)
	r.NoError(err)
	r.Len(blocks, 7)

	r.Equal([]repr.Instruction{
		ins(50, model.TypeAritm, "1"),
		ins(54, model.TypeAritm, "2"),
	}, blocks[0].Seq)
	r.Equal([]repr.Instruction{
		ins(59, model.TypeJump, "3"),
	}, blocks[1].Seq)
	r.Equal([]repr.Instruction{
		ins(63, model.TypeCJump, "4"),
	}, blocks[2].Seq)
	r.Equal([]repr.Instruction{
		ins(68, model.TypeAritm, "5"),
		ins(72, model.TypeJumpDyn, "6"),
	}, blocks[3].Seq)
	r.Equal([]repr.Instruction{
		ins(76, model.TypeAritm, "7"),
		ins(80, model.TypeAritm, "8"),
	}, blocks[4].Seq)
	r.Equal([]repr.Instruction{
		ins(84, model.TypeAritm, "9"),
		ins(88, model.TypeAritm, "10"),
		ins(92, model.TypeCJump, "11", 84),
	}, blocks[5].Seq)
	r.Equal([]repr.Instruction{
		ins(96, model.TypeAritm, "12"),
		ins(100, model.TypeCJump, "13", 96),
	}, blocks[6].Seq)
}
