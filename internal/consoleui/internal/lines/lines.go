package lines

import (
	"fmt"
	"mltwist/internal/deps"
)

type Lines struct {
	lines []Line
	code  *deps.Code

	blockStarts []int
	marks       map[int]struct{}
}

func newLines(code *deps.Code) *Lines {
	// Each block will have a header and will be delimited by a blank line.
	// Each instruction will be a single line. In the end, there will be a
	// single empty line.
	lns := make([]Line, 0, 2*code.Len()+code.NumInstr()+1)
	blockStarts := make([]int, code.Len())
	for i, b := range code.Blocks() {
		if i != 0 {
			lns = append(lns, newEmptyLine())
		}

		blockStarts[i] = len(lns)
		lns = append(lns, blockToLines(b)...)
	}

	lns = append(lns, newEmptyLine())

	return &Lines{
		lines:       lns,
		code:        code,
		blockStarts: blockStarts,
		marks:       make(map[int]struct{}, 2),
	}
}

func blockToLines(b deps.Block) []Line {
	lines := make([]Line, 1, b.Num()+1)
	lines[0] = newBlockLine(b)

	for _, ins := range b.Instructions() {
		lines = append(lines, newInstrLine(b, ins))
	}

	return lines
}

func (l Lines) Len() int         { return len(l.lines) }
func (l Lines) Index(i int) Line { return l.lines[i] }

func (l *Lines) SetMark(lineIdx int, m Mark) {
	l.lines[lineIdx].setMark(m)
	l.marks[lineIdx] = struct{}{}
}

func (l *Lines) UnmarkAll() {
	for i := range l.marks {
		l.lines[i].setMark(MarkNone)
		delete(l.marks, i)
	}
}

func (l *Lines) Reload(blockIdx int) {
	newBlock := blockToLines(l.code.Index(blockIdx))
	lines := l.lines[l.blockStarts[blockIdx]:]
	lines = lines[:len(newBlock)]
	copy(lines, newBlock)
}

func (l *Lines) reloadRange(from int, to int) {
	if from > to {
		from, to = to, from
	}

	for i := from; i <= to; i++ {
		l.Reload(i)
	}
}

func (l *Lines) Move(fromLine int, toLine int) error {
	from, to := l.Index(fromLine), l.Index(toLine)

	fromBlock, fromBlockOK := from.Block()
	toBlock, toBlockOK := to.Block()
	if !fromBlockOK {
		return fmt.Errorf("from cannot be an empty line: %d", fromLine)
	}
	if !toBlockOK {
		return fmt.Errorf("to cannot be an empty line: %d", toLine)
	}

	fromIns, fromInsOK := from.Instruction()
	toIns, toInsOK := to.Instruction()
	if fromInsOK != toInsOK {
		return fmt.Errorf("cannot swap block and an instruction")
	}

	if !fromInsOK {
		err := l.code.Move(fromBlock, toBlock)
		if err != nil {
			return fmt.Errorf("block move failed: %w", err)
		}

		l.reloadRange(fromBlock, toBlock)
	} else {
		if fromBlock != toBlock {
			return fmt.Errorf("instructions cannot be moved among blocks")
		}

		err := l.code.Index(fromBlock).Move(fromIns, toIns)
		if err != nil {
			return fmt.Errorf("instruction move failed: %w", err)
		}

		l.Reload(fromBlock)
	}

	return nil
}

func (l Lines) Block(lineIdx int) (deps.Block, bool) {
	idx, ok := l.lines[lineIdx].Block()
	if !ok {
		return deps.Block{}, false
	}
	return l.code.Index(idx), true
}

func (l Lines) Line(block deps.Block, ins int) int {
	blockLine := l.blockStarts[block.Idx()]
	zerothInsLine := blockLine + 1
	return zerothInsLine + ins
}
