package basicblock_test

import (
	"mltwist/internal/deps/internal/basicblock"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"testing"

	"github.com/stretchr/testify/require"
)

const instrLen = 4

type instruction struct {
	addr  model.Addr
	jumps []expr.Expr

	desc string
}

func (i instruction) Begin() model.Addr  { return i.addr }
func (i instruction) End() model.Addr    { return i.addr + instrLen }
func (i instruction) Jumps() []expr.Expr { return i.jumps }

func jumps(jumpAddrs []model.Addr) []expr.Expr {
	jumps := make([]expr.Expr, len(jumpAddrs))
	for i, j := range jumpAddrs {
		jumps[i] = expr.ConstFromUint(j)
	}
	return jumps
}

func insInput(addr model.Addr, desc string, jumpAddrs ...model.Addr) instruction {
	return instruction{
		addr:  addr,
		jumps: jumps(jumpAddrs),
		desc:  desc,
	}
}

func TestParse_Succ(t *testing.T) {
	tests := []struct {
		name     string
		instrs   []instruction
		expected [][]instruction
	}{{
		name: "simple_loop",
		instrs: []instruction{
			insInput(72, "1"),
			insInput(76, "2"),
			insInput(80, "3", 72),
		},
		expected: [][]instruction{{
			insInput(72, "1"),
			insInput(76, "2"),
			insInput(80, "3", 72),
		}},
	}, {
		name: "hole_in_between_instructions",
		instrs: []instruction{
			insInput(72, "1"),
			insInput(76, "2"),
			insInput(80, "3"),
			insInput(96, "4"),
			insInput(100, "5"),
		},
		expected: [][]instruction{{
			insInput(72, "1"),
			insInput(76, "2"),
			insInput(80, "3"),
		}, {
			insInput(96, "4"),
			insInput(100, "5"),
		}},
	}, {
		name: "jump_instruction_splits_block",
		instrs: []instruction{
			insInput(72, "1"),
			insInput(76, "2"),
			insInput(80, "3", 72),
			insInput(84, "4"),
			insInput(88, "5"),
		},
		expected: [][]instruction{{
			insInput(72, "1"),
			insInput(76, "2"),
			insInput(80, "3", 72),
		}, {
			insInput(84, "4"),
			insInput(88, "5"),
		}},
	}, {
		name: "jump_target_splits_block",
		instrs: []instruction{
			insInput(50, "1"),
			insInput(54, "2"),
			insInput(58, "3"),
			insInput(62, "4"),
			insInput(66, "5", 58, 70),
			insInput(70, "6"),
		},
		expected: [][]instruction{{
			insInput(50, "1"),
			insInput(54, "2"),
		}, {
			insInput(58, "3"),
			insInput(62, "4"),
			insInput(66, "5", 58, 70),
		}, {
			insInput(70, "6"),
		}},
	}, {
		name: "split_by_jump_to_jump",
		instrs: []instruction{
			insInput(30, "1"),
			insInput(34, "2"),
			insInput(38, "3", 30, 42),
			insInput(42, "4", 38, 46),
			insInput(46, "5"),
		},
		expected: [][]instruction{{
			insInput(30, "1"),
			insInput(34, "2"),
		}, {
			insInput(38, "3", 30, 42),
		}, {
			insInput(42, "4", 38, 46),
		}, {
			insInput(46, "5"),
		}},
	}, {
		name: "no_double_splits",
		instrs: []instruction{
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
		expected: [][]instruction{{
			insInput(50, "1"),
			insInput(54, "2"),
			insInput(58, "3", 50, 62),
		}, {
			insInput(62, "4"),
			insInput(66, "5"),
			insInput(70, "6", 50, 62),
		}, {
			insInput(78, "7"),
			insInput(82, "8"),
			insInput(86, "9", 78),
		}},
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			seqs, err := basicblock.Parse(tt.instrs[0].addr, tt.instrs)
			r.NoError(err)
			r.Equal(len(tt.expected), len(seqs), "Unexpected output length.")

			// Jump targets of instructions can differ as we remove
			// those targets which jump to the following
			// instruction.
			for i, e := range tt.expected {
				r.Equal(len(e), len(seqs[i]))
				for j, e := range e {
					r.Equal(e.addr, seqs[i][j].addr)
					r.Equal(e.desc, seqs[i][j].desc)
				}
			}
		})
	}
}

func TestParse_Fail(t *testing.T) {
	tests := []struct {
		name   string
		instrs []instruction
	}{{
		name: "jump_in_between_basic_blocks",
		instrs: []instruction{
			insInput(50, "1"),
			insInput(54, "2"),
			insInput(58, "3", 72),
			insInput(72, "4"),
			insInput(76, "5"),
			insInput(80, "6", 64),
		},
	}, {
		name: "jump_behind_all_blocks",
		instrs: []instruction{
			insInput(72, "4"),
			insInput(76, "5"),
			insInput(80, "6", 92),
		},
	}, {
		name: "jump_into_instruction",
		instrs: []instruction{
			insInput(72, "4"),
			insInput(76, "5"),
			insInput(80, "6", 78),
		},
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			seqs, err := basicblock.Parse(tt.instrs[0].addr, tt.instrs)
			r.Error(err)
			r.Nil(seqs)
		})
	}
}
