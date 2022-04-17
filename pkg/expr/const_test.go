package expr_test

import (
	"decomp/pkg/expr"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConst(t *testing.T) {
	r := require.New(t)

	bs := []byte{5, 6, 7, 8}
	c := expr.NewConst(bs, expr.Width32)
	r.Equal(bs, c.Bytes())
	r.Equal(expr.Width32, c.Width())
	bs[0] = 55
	r.Equal([]byte{5, 6, 7, 8}, c.Bytes())
}

func TestConstUint(t *testing.T) {
	r := require.New(t)

	r.Panics(func() {
		_ = expr.NewConstUint[uint16](999, expr.Width8)
	})

	c := expr.NewConstUint[uint8](9, expr.Width32)
	r.Equal([]byte{9, 0, 0, 0}, c.Bytes())
	r.Equal(expr.Width32, c.Width())

	const bigNum uint64 = 7733294320943090932
	expectedBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(expectedBytes, bigNum)

	c = expr.NewConstUint(bigNum, expr.Width64)
	r.Equal(expectedBytes, c.Bytes())
	r.Equal(expr.Width64, c.Width())

	r.Panics(func() {
		_ = expr.NewConstUint(bigNum, expr.Width32)
	})
}

func TestConstInt(t *testing.T) {
	r := require.New(t)

	r.Panics(func() {
		_ = expr.NewConstInt[int16](999, expr.Width8)
	})
	r.Panics(func() {
		// 11111101_11111111
		_ = expr.NewConstInt[int16](-513, expr.Width8)
	})
	r.Panics(func() {
		// 11111111_01010101
		_ = expr.NewConstInt[int16](-171, expr.Width8)
	})

	c := expr.NewConstInt[int8](-1, expr.Width32)
	r.Equal([]byte{0xff, 0xff, 0xff, 0xff}, c.Bytes())
	r.Equal(expr.Width32, c.Width())

	const bigNegNum int64 = -39293939384848383
	expectedBytes := make([]byte, 8)
	var zero uint64
	binary.LittleEndian.PutUint64(expectedBytes, zero-uint64(-bigNegNum))

	c = expr.NewConstInt(bigNegNum, expr.Width64)
	r.Equal(expectedBytes, c.Bytes())
	r.Equal(expr.Width64, c.Width())

	r.Panics(func() {
		_ = expr.NewConstInt(bigNegNum, expr.Width32)
	})

	const bigPosNum int64 = 33294320943090932
	binary.LittleEndian.PutUint64(expectedBytes, uint64(bigPosNum))

	c = expr.NewConstInt(bigPosNum, expr.Width64)
	r.Equal(expectedBytes, c.Bytes())
	r.Equal(expr.Width64, c.Width())

	r.Panics(func() {
		_ = expr.NewConstInt(bigPosNum, expr.Width32)
	})
}
