package memview

import (
	"fmt"
	"math"
	"mltwist/internal/consoleui/internal/cursor"
	"mltwist/internal/exprtransform"
	"mltwist/internal/state/interval"
	"mltwist/internal/state/memory"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"strings"
	"unsafe"
)

const bytesPerLine = 16
const bytesSpace = 8
const emptyByte = ".."

type memoryView struct {
	mem   memory.Memory
	lines []memLine
	c     *cursor.Cursor

	format          string
	emptyLineFormat string
}

func newMemoryView(mem memory.Memory) *memoryView {
	var lines []memLine
	if mem != nil {
		blocks := mem.Blocks()
		lines = memoryLines(blocks)
	}

	var c *cursor.Cursor
	if len(lines) > 0 {
		c = cursor.New(len(lines))
	}

	idFmt := fmt.Sprintf("%%%dd", numDigits(len(lines), 10))
	aChars := unsafe.Sizeof(model.Addr(0)) * 2
	format := fmt.Sprintf("%%1s %s  | 0x%%0%dx - 0x%%0%dx | %%s\n", idFmt,
		aChars, aChars)

	emptyLineFormat := fmt.Sprintf("%%1s %s  | ...\n", idFmt)

	return &memoryView{
		mem:   mem,
		lines: lines,
		c:     c,

		format:          format,
		emptyLineFormat: emptyLineFormat,
	}
}

func (v *memoryView) MinLines() int { return 5 }
func (v *memoryView) MaxLines() int { return -1 }

func (v *memoryView) Print(n int) error {
	if v.c == nil {
		fmt.Printf("\n\n")
		fmt.Printf("\tNO MEMORY TO SHOW\n")
		fmt.Printf("\n\n")
		return nil

	}

	cursorIndex := v.c.Value()

	// Golden ratio calculation.
	begin := cursorIndex - int(math.Floor(float64(n)/(math.Phi+1)))
	if begin < 0 {
		begin = 0
	}
	end := begin + n
	if end > len(v.lines) {
		end = len(v.lines)
	}

	for i := begin; i < end; i++ {
		var cursor string
		if i == cursorIndex {
			cursor = ">"
		}

		ln := v.lines[i]
		if len(ln.ranges) == 0 {
			fmt.Printf(v.emptyLineFormat, cursor, i)
		} else {
			end := ln.addr + bytesPerLine
			bytes := v.formatMemLine(ln)
			fmt.Printf(v.format, cursor, i, ln.addr, end, bytes)
		}
	}

	return nil
}

func (v *memoryView) formatMemLine(ln memLine) string {
	var sb strings.Builder

	i := 0
	for a := ln.addr; a < ln.addr+bytesPerLine; a++ {
		if a != ln.addr {
			sb.WriteByte(' ')
		}

		if diff := (a - ln.addr); diff != 0 && diff%bytesSpace == 0 {
			sb.WriteString("  ")
		}

		if i < len(ln.ranges) && a >= ln.ranges[i].end {
			i++
		}

		if i >= len(ln.ranges) || ln.ranges[i].begin > a {
			sb.WriteString(emptyByte)
			continue
		}

		ex, _ := v.mem.Load(a, expr.Width8)
		c, ok := exprtransform.ConstFold(ex).(expr.Const)
		if !ok {
			panic(fmt.Sprintf("bug: expected expr.Const at address 0x%x: %#v",
				a, ex))
		}

		sb.WriteString(fmt.Sprintf("%X", c.Bytes()[0]))
	}

	return sb.String()
}

type memRange struct {
	begin model.Addr
	end   model.Addr
}

type memLine struct {
	addr   model.Addr
	ranges []memRange
}

func memoryLines(blocks interval.Map[model.Addr]) []memLine {
	lines := make([]memLine, 0, blocks.Len())
	for _, b := range blocks.Intervals() {
		lines = append(lines, block2Lines(b)...)
	}

	for i, j := 0, 0; i < len(lines); i, j = i+1, j+1 {
		if j > 0 && lines[j-1].addr == lines[i].addr {
			lines[j-1].ranges = append(lines[j-1].ranges, lines[i].ranges...)
			j--
		} else {
			lines[j] = lines[i]
		}
	}

	return addEmptyLines(lines)
}

func addEmptyLines(lines []memLine) []memLine {
	lns := make([]memLine, 0, len(lines))
	for i, ln := range lines {
		if i == 0 && lines[i].addr != 0 {
			lns = append(lns, memLine{})
		} else if i > 0 && lines[i-1].addr+bytesPerLine < ln.addr {
			lns = append(lns, memLine{})
		}

		lns = append(lns, ln)
	}

	if len(lns) > 0 && lns[len(lns)-1].addr+bytesPerLine != model.MaxAddress {
		lns = append(lns, memLine{})
	}

	return lns
}

func block2Lines(block interval.Interval[model.Addr]) []memLine {
	begin := block.Begin() / bytesPerLine * bytesPerLine
	end := (block.End() + bytesPerLine - 1) / bytesPerLine * bytesPerLine

	lines := make([]memLine, 0, (end-begin)/bytesPerLine)
	for i := begin; i < end; i += bytesPerLine {
		b := i
		if b < block.Begin() {
			b = block.Begin()
		}

		e := i + bytesPerLine
		if e > block.End() {
			e = block.End()
		}

		lines = append(lines, memLine{
			addr:   b / bytesPerLine * bytesPerLine,
			ranges: []memRange{{begin: b, end: e}},
		})
	}

	return lines
}

func numDigits(num int, base int) int {
	if num == 0 {
		return 1
	}

	var cnt int

	// The minus (-) sign.
	if num < 0 {
		cnt++
	}

	for ; num != 0; num /= base {
		cnt++
	}

	return cnt
}
