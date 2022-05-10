package basicblock_test

import (
	"mltwist/internal/deps/internal/basicblock"
	"mltwist/internal/parser"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"testing"

	"github.com/stretchr/testify/require"
)

const instrLen = 4

type detail string

func (d detail) Name() string   { return string(d) }
func (d detail) String() string { return string(d) }

func jumpEffects(jumpAddrs []model.Addr) []expr.Effect {
	jumps := make([]expr.Effect, len(jumpAddrs))
	for i, j := range jumpAddrs {
		jumps[i] = expr.NewRegStore(model.AddrExpr(j), expr.IPKey, expr.Width32)
	}
	return jumps
}

func insInput(addr model.Addr, desc string, jumpAddrs ...model.Addr) parser.Instruction {
	return parser.Instruction{
		Addr:    addr,
		Bytes:   make([]byte, instrLen),
		Details: detail(desc),
		Effects: jumpEffects(jumpAddrs),
	}
}

func ins(addr model.Addr, desc string, jumpAddrs ...model.Addr) basicblock.Instruction {
	var jumps []expr.Expr
	for _, a := range jumpAddrs {
		jumps = append(jumps, model.AddrExpr(a))
	}

	return basicblock.Instruction{
		Addr:        addr,
		Bytes:       make([]byte, instrLen),
		Details:     detail(desc),
		Effects:     jumpEffects(jumpAddrs),
		JumpTargets: jumps,
	}
}

func TestParse_Succ(t *testing.T) {
	tests := []struct {
		name     string
		instrs   []parser.Instruction
		expected [][]basicblock.Instruction
	}{{
		name: "simple_loop",
		instrs: []parser.Instruction{
			insInput(72, "1"),
			insInput(76, "2"),
			insInput(80, "3", 72),
		},
		expected: [][]basicblock.Instruction{{
			ins(72, "1"),
			ins(76, "2"),
			ins(80, "3", 72),
		}},
	}, {
		name: "hole_in_between_instructions",
		instrs: []parser.Instruction{
			insInput(72, "1"),
			insInput(76, "2"),
			insInput(80, "3"),
			insInput(96, "4"),
			insInput(100, "5"),
		},
		expected: [][]basicblock.Instruction{{
			ins(72, "1"),
			ins(76, "2"),
			ins(80, "3"),
		}, {
			ins(96, "4"),
			ins(100, "5"),
		}},
	}, {
		name: "jump_instruction_splits_block",
		instrs: []parser.Instruction{
			insInput(72, "1"),
			insInput(76, "2"),
			insInput(80, "3", 72),
			insInput(84, "4"),
			insInput(88, "5"),
		},
		expected: [][]basicblock.Instruction{{
			ins(72, "1"),
			ins(76, "2"),
			ins(80, "3", 72),
		}, {
			ins(84, "4"),
			ins(88, "5"),
		}},
	}, {
		name: "jump_target_splits_block",
		instrs: []parser.Instruction{
			insInput(50, "1"),
			insInput(54, "2"),
			insInput(58, "3"),
			insInput(62, "4"),
			insInput(66, "5", 58, 70),
			insInput(70, "6"),
		},
		expected: [][]basicblock.Instruction{{
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
		instrs: []parser.Instruction{
			insInput(30, "1"),
			insInput(34, "2"),
			insInput(38, "3", 30, 42),
			insInput(42, "4", 38, 46),
			insInput(46, "5"),
		},
		expected: [][]basicblock.Instruction{{
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
		instrs: []parser.Instruction{
			insInput(50, "1"),
			insInput(54, "2"),
			insInput(58, "3", 50, 62),
			insInput(62, "4"),
			insInput(66, "5"),
			insInput(70, "6", 50, 62),
			insInput(78, "7"),
			insInput(82, "8"),
			insInput(86, "9", 78),
		},
		expected: [][]basicblock.Instruction{{
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

			seqs, err := basicblock.Parse(tt.instrs[0].Addr, tt.instrs)
			r.NoError(err)
			r.Equal(len(tt.expected), len(seqs), "Unexpected output length.")

			// Jump targets of instructions can differ as we remove
			// those targets which jump to the following
			// instruction.
			for i, e := range tt.expected {
				r.Equal(len(e), len(seqs[i]))
				for j, e := range e {
					r.Equal(e.Addr, seqs[i][j].Addr)
					r.Equal(e.Details, seqs[i][j].Details)
				}
			}
		})
	}
}

func TestParse_Fail(t *testing.T) {
	tests := []struct {
		name   string
		instrs []parser.Instruction
	}{{
		name: "jump_in_between_basic_blocks",
		instrs: []parser.Instruction{
			insInput(50, "1"),
			insInput(54, "2"),
			insInput(58, "3", 72),
			insInput(72, "4"),
			insInput(76, "5"),
			insInput(80, "6", 64),
		},
	}, {
		name: "jump_behind_all_blocks",
		instrs: []parser.Instruction{
			insInput(72, "4"),
			insInput(76, "5"),
			insInput(80, "6", 92),
		},
	}, {
		name: "jump_into_instruction",
		instrs: []parser.Instruction{
			insInput(72, "4"),
			insInput(76, "5"),
			insInput(80, "6", 78),
		},
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			seqs, err := basicblock.Parse(tt.instrs[0].Addr, tt.instrs)
			r.Error(err)
			r.Nil(seqs)
		})
	}
}
