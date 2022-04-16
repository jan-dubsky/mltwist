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
// sufficient for sollowing few decades. Not mentioning the fact that the
// expression model we are using is not strong enough to describe SIMD
// instruction as a single operation and it would have to be composed out of
// many operations applied on smaller bits. This further decreases the need to
// describe operations wider than 2048 bits.
type Width uint8

const (
	Width8 Width = 1 << iota
	Width16
	Width32
	Width64
	Width128
	width256
	width512
	width1024
)

// MaxWidth is maximal allowed width value.
const MaxWidth Width = math.MaxUint8

func (w Width) Bits() uint16 { return uint16(w) * 8 }
