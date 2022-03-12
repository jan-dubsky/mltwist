package lines

import (
	"decomp/internal/deps"
	"fmt"
)

type Lines struct {
	offset int

	lines       []Line
	blockStarts []int

	m *deps.Model
}

func NewLines(m *deps.Model) *Lines {
	// Each block will have a header and will be delimited by a blank line.
	// Each instruction will be a single line.
	lns := make([]Line, 0, 2*m.Len()+m.NumInstr()+1)
	blockStarts := make([]int, m.Len())
	for i, b := range m.Blocks() {
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
		m:           m,
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

func (l Lines) Offset() int { return l.offset }
func (l Lines) Len() int    { return len(l.lines) }

func (l Lines) checkOffset(offset int) error {
	if offset < 0 {
		return fmt.Errorf("offset cannot be negative: %d", offset)
	}
	if l := len(l.lines); offset >= l {
		return fmt.Errorf("offset is too high: %d > %d", offset, l)
	}

	return nil
}

func (l *Lines) Shift(o int) error {
	newOffset := l.offset + o
	if err := l.checkOffset(newOffset); err != nil {
		return fmt.Errorf("new offset value is invalid: %w", err)
	}

	l.offset = newOffset
	return nil
}

func (l *Lines) Lines(n int) []Line {
	if n <= 0 {
		return nil
	}

	ret := make([]Line, 0, n)
	for i := 0; i < n; i++ {
		j := l.offset + int(i)
		if j >= len(l.lines) {
			break
		}

		ret = append(ret, l.lines[j])
	}

	return ret
}

func (l *Lines) SetMark(lineIdx int, mark string) { l.lines[lineIdx].setMark(mark) }

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
	newBlock := blockToLines(l.m.Index(blockIdx))
	lines := l.lines[l.blockStarts[blockIdx]:]
	lines = lines[:len(newBlock)]
	copy(lines, newBlock)
}

func (l Lines) Block(lineIdx int) (*deps.Block, bool) {
	blockIdx := l.lines[lineIdx].block
	if blockIdx < 0 {
		return nil, false
	}
	return l.m.Index(blockIdx), true
}

func (l Lines) Instruction(lineIdx int) (deps.Instruction, bool) {
	line := l.lines[lineIdx]
	if line.instr < 0 {
		return deps.Instruction{}, false
	}
	return l.m.Index(line.block).Index(line.instr), true
}

func (l Lines) Line(block *deps.Block, ins deps.Instruction) int {
	blockLine := l.blockStarts[block.Idx()]
	zerothInsLine := blockLine + 1
	return zerothInsLine + ins.Idx()
}
