package memview

import (
	"mltwist/internal/state/interval"
	"mltwist/pkg/model"
)

type memLine struct {
	addr   model.Addr
	ranges []interval.Interval[model.Addr]
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
			ranges: []interval.Interval[model.Addr]{interval.New(b, e)},
		})
	}

	return lines
}
