package memory

import (
	"fmt"
	"mltwist/internal/state/interval"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"sort"
)

var _ Memory = &Bytes{}

type ByteBlock interface {
	Begin() model.Addr
	Bytes() []byte
}

// Bytes is memory type representing memory as a sequence of bytes.
//
// This type of memory is in a way special as it's meant to represent mostly
// readonly data. A good example of such data can be data loaded from ELF file.
// For those data, we could use standard sparse memory with bytes rewritten an
// expressions. On the other hand such representation wouldn't be effective. For
// both memory and performance reasons, it makes a good sense to have those
// bytes stored as just array of bytes.
//
// Bytes memory supports writes as well. But as we can represent only expr.Const
// as array of bytes, Store call will panic of the expression stored is not
// expr.Const.
type Bytes struct {
	blocks []byteBlock
}

type byteBlock struct {
	begin model.Addr
	bytes []byte
}

func (b *byteBlock) end() model.Addr { return b.begin + model.Addr(len(b.bytes)) }
func (b *byteBlock) interval() interval.Interval[model.Addr] {
	return interval.New(b.begin, b.end())
}

// index returns index of a byte with address addr in the block.
func (b byteBlock) index(addr model.Addr) int { return int(addr - b.begin) }

func NewBytes(blocks []ByteBlock) (*Bytes, error) {
	bs := make([]byteBlock, len(blocks))
	for i, b := range blocks {
		bytes := b.Bytes()
		bytesCopy := make([]byte, len(bytes))
		copy(bytesCopy, b.Bytes())
		bs[i] = byteBlock{begin: b.Begin(), bytes: bytesCopy}
	}

	sort.Slice(bs, func(i, j int) bool { return bs[i].begin < bs[j].begin })

	bs, err := dedupBlocks(bs)
	if err != nil {
		return nil, err
	}

	return &Bytes{blocks: bs}, nil
}

func dedupBlocks(bs []byteBlock) ([]byteBlock, error) {
	if len(bs) < 2 {
		return bs, nil
	}

	j := 1
	for i := 1; i < len(bs); i, j = i+1, j+1 {
		if bs[j-1].end() > bs[i].begin {
			return nil, fmt.Errorf("blocks overlap: %v and %v",
				bs[j-1].end(), bs[i].begin)
		}

		if bs[j-1].end() == bs[i].begin {
			bs[j-1].bytes = append(bs[j-1].bytes, bs[i].bytes...)
			j--
			continue
		}

		bs[j] = bs[i]
	}

	return bs[:j], nil
}

func (b *Bytes) address(addr model.Addr) (int, bool) {
	idx := sort.Search(len(b.blocks), func(i int) bool {
		return addr < b.blocks[i].end()
	})

	if idx == len(b.blocks) || addr < b.blocks[idx].begin {
		return -1, false
	}

	return idx, true
}

func (b *Bytes) Load(addr model.Addr, w expr.Width) (expr.Expr, bool) {
	end := addr + model.Addr(w)

	blockIdx, ok := b.address(addr)
	if !ok || b.blocks[blockIdx].end() < end {
		return nil, false
	}

	block := b.blocks[blockIdx]
	baseIdx := block.index(addr)
	return expr.NewConst(block.bytes[baseIdx:baseIdx+int(w)], w), true
}

// Store is standard memory store but it panics of ex is not of type expr.Const.
// See Bytes doc-comment for thorough explanation.
func (b *Bytes) Store(addr model.Addr, ex expr.Expr, w expr.Width) {
	c, ok := ex.(expr.Const)
	if !ok {
		panic(fmt.Sprintf("memory Bytes allows only expr.Const writes: %T", ex))
	}

	bs := c.WithWidth(w).Bytes()
	for len(bs) > 0 {
		n := b.store(addr, bs)
		bs = bs[n:]
		addr += model.Addr(n)
	}

	var err error
	b.blocks, err = dedupBlocks(b.blocks)
	if err != nil {
		panic(fmt.Sprintf("bug: store resulted in byte overlap: %s", err.Error()))
	}
}

func (b *Bytes) store(addr model.Addr, bs []byte) int {
	blockIdx, ok := b.address(addr)
	if ok {
		block := b.blocks[blockIdx]
		beginIdx := block.index(addr)
		endIdx := beginIdx + len(bs)
		if endIdx > len(block.bytes) {
			endIdx = len(block.bytes)
		}

		copy(block.bytes[beginIdx:endIdx], bs)
		return endIdx - beginIdx
	}

	end := addr + model.Addr(len(bs))
	idx := sort.Search(len(b.blocks), func(i int) bool {
		return addr < b.blocks[i].begin
	})
	if idx != len(b.blocks) && b.blocks[idx].begin < end {
		end = b.blocks[idx].begin
	}

	b.blocks = append(b.blocks, byteBlock{})
	for i := idx; i < len(b.blocks)-1; i++ {
		b.blocks[i+1] = b.blocks[i]
	}

	b.blocks[idx] = byteBlock{
		begin: addr,
		bytes: bs[:end-addr],
	}

	return int(end - addr)
}

func (b *Bytes) intervalMap() interval.Map[model.Addr] {
	intervals := make([]interval.Interval[model.Addr], len(b.blocks))
	for i, b := range b.blocks {
		intervals[i] = b.interval()
	}

	return interval.NewMap(intervals...)
}

func (b *Bytes) Missing(addr model.Addr, w expr.Width) interval.Map[model.Addr] {
	intv := interval.New(addr, addr+model.Addr(w))
	return interval.MapComplement(interval.NewMap(intv), b.intervalMap())
}

func (b *Bytes) Blocks() interval.Map[model.Addr] {
	return b.intervalMap()
}
