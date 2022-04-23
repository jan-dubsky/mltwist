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
	addr model.Addr,
	desc string,
	jmps ...model.Addr,
) repr.Instruction {
	jmpExprs := make([]expr.Expr, len(jmps))
	for i, j := range jmps {
		jmpExprs[i] = model.AddrExpr(j)
	}

	return repr.Instruction{
		Address: addr,
		Instruction: model.Instruction{
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
	}{{
		name: "simple_loop",
		instrs: []repr.Instruction{
			ins(72, "1"),
			ins(76, "2"),
			ins(80, "3", 72),
		},
		expected: [][]repr.Instruction{{
			ins(72, "1"),
			ins(76, "2"),
			ins(80, "3", 72),
		}},
	}, {
		name: "hole_in_between_instructions",
		instrs: []repr.Instruction{
			ins(72, "1"),
			ins(76, "2"),
			ins(80, "3"),
			ins(96, "4"),
			ins(100, "5"),
		},
		expected: [][]repr.Instruction{{
			ins(72, "1"),
			ins(76, "2"),
			ins(80, "3"),
		}, {
			ins(96, "4"),
			ins(100, "5"),
		}},
	}, {
		name: "jump_instruction_splits_block",
		instrs: []repr.Instruction{
			ins(72, "1"),
			ins(76, "2"),
			ins(80, "3", 72),
			ins(84, "4"),
			ins(88, "5"),
		},
		expected: [][]repr.Instruction{{
			ins(72, "1"),
			ins(76, "2"),
			ins(80, "3", 72),
		}, {
			ins(84, "4"),
			ins(88, "5"),
		}},
	}, {
		name: "jump_target_splits_block",
		instrs: []repr.Instruction{
			ins(50, "1"),
			ins(54, "2"),
			ins(58, "3"),
			ins(62, "4"),
			ins(66, "5", 58, 70),
			ins(70, "6"),
		},
		expected: [][]repr.Instruction{{
			ins(50, "1"),
			ins(54, "2"),
		}, {
			ins(58, "3"),
			ins(62, "4"),
			ins(66, "5", 58, 70),
		}, {
			ins(70, "6"),
		}},
	}, {
		name: "split_by_jump_to_jump",
		instrs: []repr.Instruction{
			ins(30, "1"),
			ins(34, "2"),
			ins(38, "3", 30, 42),
			ins(42, "4", 38, 46),
			ins(46, "5"),
		},
		expected: [][]repr.Instruction{{
			ins(30, "1"),
			ins(34, "2"),
		}, {
			ins(38, "3", 30, 42),
		}, {
			ins(42, "4", 38, 46),
		}, {
			ins(46, "5"),
		}},
	}, {
		name: "no_double_splits",
		instrs: []repr.Instruction{
			ins(50, "1"),
			ins(54, "2"),
			ins(58, "3", 50, 62),
			ins(62, "4"),
			ins(66, "5"),
			ins(70, "6", 50, 62),
			ins(78, "7"),
			ins(82, "8"),
			ins(86, "9", 78),
		},
		expected: [][]repr.Instruction{{
			ins(50, "1"),
			ins(54, "2"),
			ins(58, "3", 50, 62),
		}, {
			ins(62, "4"),
			ins(66, "5"),
			ins(70, "6", 50, 62),
		}, {
			ins(78, "7"),
			ins(82, "8"),
			ins(86, "9", 78),
		}},
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			seqs, err := basicblock.Parse(tt.instrs)
			r.NoError(err)
			r.Equal(len(tt.expected), len(seqs), "Unexpected output length.")

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
	}{{
		name: "jump_in_between_basic_blocks",
		instrs: []repr.Instruction{
			ins(50, "1"),
			ins(54, "2"),
			ins(58, "3", 72),
			ins(72, "4"),
			ins(76, "5"),
			ins(80, "6", 64),
		},
	}, {
		name: "jump_behind_all_blocks",
		instrs: []repr.Instruction{
			ins(72, "4"),
			ins(76, "5"),
			ins(80, "6", 92),
		},
	}, {
		name: "jump_into_instruction",
		instrs: []repr.Instruction{
			ins(72, "4"),
			ins(76, "5"),
			ins(80, "6", 78),
		},
	}}

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
