package expr

import "math"

// Width states number of bytes the operation is applied to. This value
// describes width of both input and output operands of an instruction.
//
// Every CPU architecture typically have some maximal width of a both input and
// output registers. All calculation are then performed only of this particular
// bit register width and all higher bits of the result are dropped.
//
// Width represents number of full bytes rather than bits for a few reasons.
// First of all, any reasonable architecture nowadays operated on bytes, so it
// doesn't make much sense to allow widths in more granular units then bytes.
// Another reason is that supporting a generic bit width would be more
// complicated the code, but more importantly, this very generic solution would
// lower the performance of any possible calculating logic significantly due to
// byte-ordered nature of CPUs this program will run on.
//
// These days, uint8 seems to be sufficient to describe any possible width as it
// support up to 2048 bit wide registers. As even the largest registers and SIMD
// operations can handle about 512 bits at most, so 2048 will be most likely
// sufficient for following few decades. Not mentioning the fact that the
// expression model we are using is not strong enough to describe SIMD
// instruction as a single operation and it would have to be composed out of
// many operations applied on smaller bits. This further decreases the need to
// describe operations wider than 2048 bits.
type Width uint8

const (
	// Width8 is width spanning single byte -> 8 bits.
	Width8 Width = 1 << iota
	// Width16 is width spanning two bytes -> 16 bits.
	Width16
	// Width32 is width spanning four bytes -> 32 bits.
	Width32
	// Width64 is width spanning eight bytes -> 64 bits.
	Width64
	// Width128 is width spanning sixteen bytes -> 128 bits.
	Width128
	// Width256 is width spanning thirty two bytes -> 256 bits.
	Width256
	// Width512 is width spanning sixty four bytes -> 512 bits.
	Width512
	// Width1024 is width spanning one hundred twenty eight bytes -> 1024
	// bits.
	Width1024
)

// MaxWidth is maximal allowed width value.
//
// This value is not supposed to be used as actual width as it's not power of
// two. This value should be used in condition expresions to check if number is
// reasonably small to be represented by width.
const MaxWidth Width = math.MaxUint8

// Bits returns maximal number of bits which can be stored in an expression of
// width w.
func (w Width) Bits() uint16 { return uint16(w) * 8 }
