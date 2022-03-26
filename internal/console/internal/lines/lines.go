package lines

import (
	"decomp/internal/deps"
	"fmt"
)

type Lines struct {
	offset int

	lines       []Line
	blockStarts []int

	p *deps.Program
}

func NewLines(p *deps.Program) *Lines {
	// Each block will have a header and will be delimited by a blank line.
	// Each instruction will be a single line.
	lns := make([]Line, 0, 2*p.Len()+p.NumInstr()+1)
	blockStarts := make([]int, p.Len())
	for i, b := range p.Blocks() {
		if i != 0 {
			lns = append(lns, newEmptyLine())
		}

		blockStarts[i] = len(lns)
		lns = append(lns, blockToLines(b)...)
	}

	lns = append(lns, newEmptyLine())

	return &Lines{
		offset:      0,
		lines:       lns,
		blockStarts: blockStarts,
		p:           p,
	}
}

func blockToLines(b *deps.Block) []Line {
	lines := make([]Line, 1, b.Len()+1)

	lines[0] = Line{
		mark:  "",
		value: fmt.Sprintf("Block %d: 0x%x", b.Idx()+1, b.Begin()),
		block: b.Idx(),
		instr: -1,
	}

	for j, ins := range b.Instructions() {
		lines = append(lines, Line{
			mark:  "",
			value: fmt.Sprintf("%s %s", blockIndent, ins.String()),
			block: b.Idx(),
			instr: j,
		})
	}

	return lines
}

func (l Lines) Len() int         { return len(l.lines) }
func (l Lines) Index(i int) Line { return l.lines[i] }

func (l *Lines) SetMark(lineIdx int, m Mark) { l.lines[lineIdx].setMark(m) }

func (l *Lines) UnmarkAll() {
	for i := range l.lines {
		l.lines[i].setMark("")
	}
}

func (l Lines) lineIndices(lineIdx int) (int, int) {
	line := l.lines[lineIdx]
	switch {
	case line.block < 0:
		return -1, -1
	case line.instr < 0:
		return line.block, -1
	default:
		return line.block, line.instr
	}
}

func (l *Lines) Reload(blockIdx int) {
	newBlock := blockToLines(l.p.Index(blockIdx))
	lines := l.lines[l.blockStarts[blockIdx]:]
	lines = lines[:len(newBlock)]
	copy(lines, newBlock)
}

func (l Lines) Block(lineIdx int) (*deps.Block, bool) {
	blockIdx := l.lines[lineIdx].block
	if blockIdx < 0 {
		return nil, false
	}
	return l.p.Index(blockIdx), true
}

func (l Lines) Instruction(lineIdx int) (deps.Instruction, bool) {
	line := l.lines[lineIdx]
	if line.instr < 0 {
		return deps.Instruction{}, false
	}
	return l.p.Index(line.block).Index(line.instr), true
}

func (l Lines) Line(block *deps.Block, ins deps.Instruction) int {
	blockLine := l.blockStarts[block.Idx()]
	zerothInsLine := blockLine + 1
	return zerothInsLine + ins.Idx()
}
