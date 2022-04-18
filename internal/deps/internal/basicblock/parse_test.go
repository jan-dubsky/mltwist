package basicblock_test

import (
	"mltwist/internal/deps/internal/basicblock"
	"mltwist/internal/repr"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"testing"

	"github.com/stretchr/testify/require"
)

const instrLen = 4

type detail string

func (d detail) Name() string   { return string(d) }
func (d detail) String() string { return string(d) }

func ins(
	addr model.Address,
	tp model.Type,
	desc string,
	jmps ...model.Address,
) repr.Instruction {
	jmpExprs := make([]expr.Expr, len(jmps))
	for i, j := range jmps {
		jmpExprs[i] = model.AddressExpr(j)
	}

	return repr.Instruction{
		Address: addr,
		Instruction: model.Instruction{
			Type:    tp,
			ByteLen: instrLen,
			Details: detail(desc),
		},
		JumpTargets: jmpExprs,
	}
}

func TestParse_Succ(t *testing.T) {
	tests := []struct {
		name     string
		instrs   []repr.Instruction
		expected [][]repr.Instruction
	}{
		{
			name: "simple_loop",
			instrs: []repr.Instruction{
				ins(72, model.TypeAritm, "4"),
				ins(76, model.TypeAritm, "5"),
				ins(80, model.TypeCJump, "6", 72),
			},
			expected: [][]repr.Instruction{
				{
					ins(72, model.TypeAritm, "4"),
					ins(76, model.TypeAritm, "5"),
					ins(80, model.TypeCJump, "6", 72),
				},
			},
		},
		{
			name: "complex_example",
			instrs: []repr.Instruction{
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
				ins(104, model.TypeCJump, "14", 54),
			},
			expected: [][]repr.Instruction{
				{
					ins(50, model.TypeAritm, "1"),
				}, {
					ins(54, model.TypeAritm, "2"),
				}, {
					ins(59, model.TypeJump, "3"),
				}, {
					ins(63, model.TypeCJump, "4"),
				}, {
					ins(68, model.TypeAritm, "5"),
					ins(72, model.TypeJumpDyn, "6"),
				}, {
					ins(76, model.TypeAritm, "7"),
					ins(80, model.TypeAritm, "8"),
				}, {
					ins(84, model.TypeAritm, "9"),
					ins(88, model.TypeAritm, "10"),
					ins(92, model.TypeCJump, "11", 84),
				}, {
					ins(96, model.TypeAritm, "12"),
					ins(100, model.TypeCJump, "13", 96),
				}, {
					ins(104, model.TypeCJump, "14", 54),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			seqs, err := basicblock.Parse(tt.instrs)
			r.NoError(err)
			r.Len(seqs, len(tt.expected))

			for i, e := range tt.expected {
				r.Equal(e, seqs[i])
			}
		})
	}
}

func TestParse_Fail(t *testing.T) {
	tests := []struct {
		name   string
		instrs []repr.Instruction
	}{
		{
			name: "jump_in_between_blocks",
			instrs: []repr.Instruction{
				ins(50, model.TypeAritm, "1"),
				ins(54, model.TypeAritm, "2"),
				ins(58, model.TypeJump, "3"),
				ins(72, model.TypeAritm, "4"),
				ins(76, model.TypeAritm, "5"),
				ins(80, model.TypeCJump, "6", 64),
			},
		},
		{
			name: "jump_behind_all_blocks",
			instrs: []repr.Instruction{
				ins(72, model.TypeAritm, "4"),
				ins(76, model.TypeAritm, "5"),
				ins(80, model.TypeCJump, "6", 92),
			},
		},
		{
			name: "jump_into_instruction",
			instrs: []repr.Instruction{
				ins(72, model.TypeAritm, "4"),
				ins(76, model.TypeAritm, "5"),
				ins(80, model.TypeCJump, "6", 78),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			seqs, err := basicblock.Parse(tt.instrs)
			r.Error(err)
			r.Nil(seqs)
		})
	}
}
